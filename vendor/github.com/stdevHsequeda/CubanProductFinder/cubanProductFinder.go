package storeClient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	httpClient "github.com/stdevHsequeda/CubanProductFinder/client"
	html "github.com/zlepper/encoding-html"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var sectionList []Section

type StoreClient struct {
	pool   *Pool // Pool of workers
	stores []Store
	client *httpClient.Client // Client
	cache  *Cache             // Cache data
}

func NewStoreClient() *StoreClient {
	// Init all attribs
	httpClient.MaxRetry = 5
	return &StoreClient{client: httpClient.NewClient(), pool: NewPool(runtime.NumCPU()), cache: NewCache()}
}

func (sc *StoreClient) Start() {
	go sc.start()
}

func (sc *StoreClient) start() {
	logrus.Info("Starting client...")
	defer sc.pool.Shutdown()

	storeList, err := sc.getStoreList()
	if err != nil {
		logrus.Fatal(err)
	}

	sectionList = make([]Section, 0)
	for i := range storeList {
		if storeList[i].Online {
			tuEnvioSectionList, err := sc.getSectionsFromTuEnvioStore(storeList[i])
			if err != nil {
				logrus.Fatal(err)
			}
			sectionList = append(sectionList, tuEnvioSectionList...)
		}
	}

	quintaY42SectionList, err := sc.getSectionsFrom5tay42()
	if err != nil {
		logrus.Warn(err)
	}

	sectionList = append(sectionList, quintaY42SectionList...)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(sectionListInternal []Section, sc *StoreClient, ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for i, section := range sectionListInternal {
					if sectionListInternal[i].GetReadyTime().Before(time.Now()) {
						sectionListInternal[i].SetReadyTime(time.Now().Add(1 * time.Minute))
						sc.pool.Run(
							&W{
								ctx: context.WithValue(context.WithValue(ctx, "sc", sc), "section", section),
							},
						)
					}
				}
			}
		}
	}(sectionList, sc, ctx)
	<-ctx.Done()
}

func (sc *StoreClient) SearchProduct(pattern string) ([]Product, error) {
	_, list, err := sc.cache.SearchProducts(pattern)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (sc *StoreClient) getStoreList() ([]Store, error) {
	logrus.Info("Getting list of stores")

	req, err := http.NewRequest(http.MethodGet, "https://www.tuenvio.cu/stores.json", nil)
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	resp, err := sc.client.CallRetryable(req)
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	var storeList = make([]Store, 0)
	err = json.NewDecoder(resp).Decode(&storeList)
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}
	return storeList, nil
}

func (sc *StoreClient) getSectionsFrom5tay42() ([]Section, error) {
	logrus.Info("Getting sections from store: 5taY42")
	req, err := http.NewRequest("GET", "https://5tay42.xetid.cu/", nil)
	if err != nil {
		return nil, err
	}

	resp, err := sc.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var htmlContent = struct {
		Parents []struct {
			Name     string             `css:".dropdown-toggle"`
			Sections []QuintaY42Section `css:".block .level-2"`
		} `css:".navbar-nav .level-1"`
	}{}

	err = html.NewDecoder(resp).Decode(&htmlContent)
	if err != nil {
		return nil, err
	}

	var result = make([]Section, 0)
	for _, parent := range htmlContent.Parents {
		for _, section := range parent.Sections {
			result = append(result, &QuintaY42Section{
				Name:   section.Name,
				Url:    section.Url,
				Parent: parent.Name,
				Store: &Store{
					Id:       100,
					Name:     "5taY42",
					Province: "La Habana",
					Online:   true,
					Url:      "https://5tay42.xetid.cu",
				},
				ReadyTime: time.Now(),
			})
		}
	}

	return result, nil
}

func (sc *StoreClient) getSectionsFromTuEnvioStore(store Store) ([]Section, error) {
	logrus.Infof("Getting sections from store: %s", store.Name)

	req, err := http.NewRequest("GET", store.Url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := sc.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var htmlContent = struct {
		Content []TuEnvioSection `css:".nav li"`
	}{}

	err = html.NewDecoder(resp).Decode(&htmlContent)
	if err != nil {
		return nil, err
	}

	var result = make([]Section, 0)

	var currentParent string
	for _, section := range htmlContent.Content {
		switch section.Url {
		case "default.aspx":
			continue
		case "#":
			currentParent = section.Name
			continue
		default:
			result = append(result, &TuEnvioSection{
				Name:      section.Name,
				Url:       fmt.Sprintf("%s/%s", strings.TrimSpace(store.Url), strings.TrimSpace(section.Url)),
				Parent:    currentParent,
				Store:     &store,
				ReadyTime: time.Now(),
			})
		}
	}

	return result, nil
}

func (sc *StoreClient) getProductsFromSection(section Section) ([]Product, error) {
	logrus.Infof("Getting products from %s in %s", section.GetName(), section.GetStore().Name)
	req, err := http.NewRequest("GET", section.GetUrl(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := sc.client.CallRetryable(req)
	if err != nil {
		return nil, err
	}

	var result = make([]Product, 0)
	if section.GetStore().Name == "5taY42" {
		var list = struct {
			Content []QuintaY42Product `css:"#listado-prod li"`
		}{}
		err = html.NewDecoder(resp).Decode(&list)
		for _, product := range list.Content {
			if product.Available != "" {
				result = append(result, &TuEnvioProduct{
					Name:    strings.TrimSpace(product.Name),
					Price:   strings.TrimSpace(product.Price),
					Link:    strings.TrimSpace(product.Link),
					Section: section,
				})
			}
		}

	} else {
		var list = struct {
			Content []TuEnvioProduct `css:".hProductItems .clearfix"`
		}{}
		err = html.NewDecoder(resp).Decode(&list)

		for _, product := range list.Content {
			result = append(result, &TuEnvioProduct{
				Name:    strings.TrimSpace(product.Name),
				Price:   strings.TrimSpace(product.Price),
				Link:    fmt.Sprintf("%s/%s", section.GetStore().Url, strings.TrimSpace(product.Link)),
				Section: section,
			})
		}
	}

	return result, nil
}
