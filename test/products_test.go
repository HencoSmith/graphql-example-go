package productstest

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	source "github.com/HencoSmith/graphql-example-go/source"
)

func TestProductList(t *testing.T) {
	config := source.GetConfig("..")
	res, err := http.Get("http://localhost:" + config.Server.Port + "/graphql?query={list{id,name,info,price}}")
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, body)

	testIDData := [3]int64{1, 2, 3}
	testNameData := [3]string{
		"Chicha Morada",
		"Chicha de jora",
		"Pisco",
	}
	testInfoData := [3]string{
		"Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		"Chicha de jora is a corn beer chicha prepared by germinating maize, extracting the malt sugars, boiling the wort, and fermenting it in large vessels (traditionally huge earthenware vats) for several days (wiki)",
		"Pisco is a colorless or yellowish-to-amber colored brandy produced in winemaking regions of Peru and Chile (wiki)",
	}
	testPriceData := [3]float64{7.99, 5.95, 9.95}

	values := gjson.Get(string(body), "data.list")

	assert.Equal(t, len(values.Array()), 3, "Item list must be of length 3")

	for count, item := range values.Array() {
		itemStr := item.String()
		id := gjson.Get(itemStr, "id").Int()
		name := gjson.Get(itemStr, "name").String()
		info := gjson.Get(itemStr, "info").String()
		price := gjson.Get(itemStr, "price").Float()

		assert.Equal(t, testIDData[count], id, "IDs should be equal")
		assert.Equal(t, testNameData[count], name, "Names should be equal")
		assert.Equal(t, testInfoData[count], info, "Info should be equal")
		assert.Equal(t, testPriceData[count], price, "Prices should be equal")
	}
}

func TestGetProduct(t *testing.T) {
	config := source.GetConfig("..")
	res, err := http.Get("http://localhost:" + config.Server.Port + "/graphql?query={product(id:1){id,name,info,price}}")
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, body)

	bodyStr := string(body)
	id := gjson.Get(bodyStr, "data.product.id").Int()
	name := gjson.Get(bodyStr, "data.product.name").String()
	info := gjson.Get(bodyStr, "data.product.info").String()
	price := gjson.Get(bodyStr, "data.product.price").Float()

	assert.Equal(t, int64(1), id, "IDs should be equal")
	assert.Equal(t, "Chicha Morada", name, "Names should be equal")
	assert.Equal(t, "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)", info, "Info should be equal")
	assert.Equal(t, float64(7.99), price, "Prices should be equal")
}

func TestCreateProduct(t *testing.T) {

}

func TestUpdateProduct(t *testing.T) {

}

func TestDeleteProduct(t *testing.T) {

}
