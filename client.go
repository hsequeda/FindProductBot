package main

import (
	"fmt"
	html "github.com/zlepper/encoding-html"
	"net/http"
	"net/url"
)

type Product struct {
	Name    string `css:".thumbTitle"`
	Price   string `css:".thumbPrice"`
	Picture string `css:".thumbnail img" extract:"attr" attr:"src"`
	Link    string `css:".thumbnail a" extract:"attr" attr:"href"`
	Store   string
}

type ProductList struct {
	Content []Product `css:".hProductItems .clearfix"`
}

const tuEnvioUrl = "https://www.tuenvio.cu"

// GetProductsByPattern return a Product list by a store and a pattern.
// Example:
/*
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
func GetProductsByPattern(store, pattern string) (*ProductList, error) {
	req, err := getRequest(store, pattern)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var list ProductList
	err = html.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		return nil, err
	}

	for i := range list.Content {
		list.Content[i].Link = fmt.Sprintf("%s/%s/%s", tuEnvioUrl, store, list.Content[i].Link)
	}

	return &list, err
}

func getRequest(store, pattern string) (*http.Request, error) {
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
