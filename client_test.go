package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func Test_GetProductsByPattern(t *testing.T) {
	testCases := []struct {
		name    string
		store   string
		pattern string
	}{{
		name:    "ron-4caminos",
		store:   "4caminos",
		pattern: "ron",
	},
		{
			name:    "pollo-5taY42",
			store:   "5taY42",
			pattern: "pollo",
		}, {
			name:    "aceite-carlos3",
			store:   "carlos3",
			pattern: "aceite",
		},
		{
			name:    "refresco-5taY42",
			store:   "5taY42",
			pattern: "refresco",
		},
	}
	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {

			resp, err := GetProductsByPattern(testCases[i].store, testCases[i].pattern)
			require.NoError(t, err)

			for i := range resp {
				t.Logf("Nombre:(%s), Precio(%s), Link(%s), Tienda(%s),Disponivilidad(%t)", resp[i].GetName(), resp[i].GetPrice(),
					resp[i].GetLink(), resp[i].GetStore(), resp[i].IsAvailable())
			}
		})
	}
	// resp, err := getRequest("4caminos", "ron")
	// require.NoError(t, err)
	//
	// var list ProductList
	// err = html.NewDecoder(resp.Body).Decode(&list)
	// require.NoError(t, err)
	//
	// t.Log(list)
}

func TestEnzona(t *testing.T) {
	_ = NewClient()
	form := url.Values{}
	form.Add("action", "data_products")
	form.Add("provincia", "Cuba")
	form.Add("filtro", "ron")
	form.Add("minPrice", "0")
	form.Add("maxPrice", "5000")
	form.Add("fav", "0")
	form.Add("tienda", "Todas")
	form.Add("disp", "0")
	req, err := http.NewRequest(http.MethodPost, "https://buscador.enzona.net/action", strings.NewReader(form.Encode()))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log(string(b))
}
