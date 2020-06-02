package storeClient

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Cache struct {
	pool *redis.Pool
}

func NewCache(network, address string) (*Cache, error) {
	pool := redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 60 * time.Second,
      Wait: true,
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial(network,address)
		},
	}

	conn := pool.Get()

	defer conn.Close()

	err := conn.Send("MULTI")
	if err != nil {
      return nil, err
	}

	_, err = conn.Do("FLUSHALL")
	if err != nil {
      return nil, err
	}

	err = conn.Send(
		"FT.CREATE", "products", "SCHEMA",
		"name", "TEXT", "SORTABLE",
		"price", "TEXT", "NOINDEX",
		"link", "TEXT", "NOINDEX",
		"store", "TEXT", "NOINDEX",
		"timestamp", "NUMERIC", "NOINDEX")
	if err != nil {
      return nil, err
	}

	err = conn.Send("EXEC")
	if err != nil {
		logrus.Fatal(err)
	}

	return &Cache{pool: &pool}, nil
}

func (c Cache) AddProduct(product Product) error {
	conn := c.pool.Get()

	defer conn.Close()

	_, err := conn.Do(
		"FT.ADD", "products",
		fmt.Sprintf("%s:%s", product.GetName(), product.GetSection().GetStore().Name), "1", "REPLACE",
		"FIELDS",
		"name", product.GetName(),
		"price", product.GetPrice(),
		"link", product.GetLink(),
		"store", product.GetSection().GetStore().Name,
		"timestamp", time.Now().Add(2*time.Hour).Unix(),
	)

	return err
}

func (c Cache) SearchProducts(pattern string) (int, []Product, error) {

	conn := c.pool.Get()

	defer conn.Close()

	reply, err := conn.Do("FT.SEARCH", "products", pattern)
	if err != nil {
		return -1, nil, err
	}
	rawData, ok := reply.([]interface{})
	if !ok {
		return -1, nil, errors.New("'reply' cannot be treated as '[]interface{}'")
	}

	total, err := redis.Int(rawData[0], nil)
	if err != nil {
		fmt.Println(err)
		return -1, nil, err
	}

	rawData = rawData[1:]

	var productList = make([]Product, 0)
	for i := 1; i < len(rawData); i += 2 {
		mapProd, err := redis.StringMap(rawData[i], nil)
		if err != nil {
			return total, nil, err
		}

		timestamp, err := strconv.ParseInt(mapProd["timestamp"], 10, 64)
		if err != nil {
			return total, nil, err
		}
		if time.Now().Unix() > timestamp {
			key, err := redis.String(rawData[i-1], nil)
			if err != nil {
				return total, nil, err
			}

			_, err = conn.Do("FT.DEL", "products", key, "DD")
			if err != nil {
				return total, nil, err
			}
			// err=conn.Send("EXEC")
			continue
		}

		productList = append(productList, &GenericProduct{
			Name:  mapProd["name"],
			Price: mapProd["price"],
			Link:  mapProd["link"],
			Section: &GenericSection{
				Store: &Store{
					Name: mapProd["store"],
				},
				ReadyTime: time.Time{},
			},
		})
	}

	return total, productList, nil
}
