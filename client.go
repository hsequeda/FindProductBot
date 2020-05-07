package main

import (
	"crypto/tls"
	"findTuEnvioBot/products"
	"fmt"
	html "github.com/zlepper/encoding-html"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const tuEnvioUrl = "https://www.tuenvio.cu"
const quintaY42 = "http://5tay42.xetid.cu"

var client *http.Client
var once sync.Once

func NewClient() *http.Client {
	once.Do(func() {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	})

	return client
}

/*
 GetProductsByPattern return a Product list by a store and a pattern.
 Example:

	var store = "carlos3"
	var pattern = "ron"
	list, err := GetProductByPattern( store, pattern)
	if err!=nil{
		panic(err)
	}

	for product := range list{
		fmt.Printf("%#v \n", product)
	}
*/
func GetProductsByPattern(store, pattern string) (result []products.Product, err error) {
	_ = NewClient()
	req, err := getRequest(store, pattern)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if store == "5taY42" {
		return decodeProduct5tay42(resp)
	} else {
		return decodeProductTuEnvio(resp, store)
	}
}

func decodeProduct5tay42(response *http.Response) ([]products.Product, error) {
	var list struct {
		Content []products.QuintaY42Product `css:"#listado-prod li"`
	}

	err := html.NewDecoder(response.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	var productList = make([]products.Product, 0)
	for _, prod := range list.Content {
		prod.Name = strings.TrimSpace(prod.Name)
		prod.Price = strings.TrimSpace(prod.Price) + "CUP"
		prod.Link = strings.TrimSpace(prod.Link)
		prod.Store = "5taY42"
		productList = append(productList, prod)
	}

	return productList, nil
}

func decodeProductTuEnvio(response *http.Response, store string) ([]products.Product, error) {
	var list struct {
		Content []products.TuEnvioProduct `css:".hProductItems .clearfix"`
	}

	err := html.NewDecoder(response.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	var productList = make([]products.Product, 0)
	for _, prod := range list.Content {
		prod.Name = strings.TrimSpace(prod.Name)
		prod.Price = strings.TrimSpace(prod.Price)
		prod.Link = strings.TrimSpace(fmt.Sprintf("%s/%s/%s", tuEnvioUrl, store, prod.Link))
		prod.Store = store
		productList = append(productList, prod)
	}
	return productList, nil
}

func getRequest(store, pattern string) (*http.Request, error) {
	if store == "5taY42" {
		return getRequest5tay42(pattern)
	} else {
		return getRequestTuEnvio(store, pattern)
	}
}

func getRequest5tay42(pattern string) (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/module/categorysearch", quintaY42), nil)
	if err != nil {
		return nil, err
	}

	qr, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	qr.Add("fc", "module")
	qr.Add("module", "categorysearch")
	qr.Add("controller", "catesearch")
	qr.Add("orderby", "position")
	qr.Add("orderway", "desc")
	qr.Add("search_category", "all")
	qr.Add("search_query", pattern)
	qr.Add("submit_search", "")

	req.URL.RawQuery = qr.Encode()

	return req, nil
}

func getRequestTuEnvio(store, pattern string) (*http.Request, error) {
	req, err := http.NewRequest("GET", tuEnvioUrl+fmt.Sprintf("/%s/Search.aspx", store), nil)
	if err != nil {
		return nil, err
	}

	qr, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	qr.Add("depPid", "0")
	qr.Add("keywords", pattern)
	req.URL.RawQuery = qr.Encode()

	return req, nil
}
