package productstest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"text/template"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	source "github.com/HencoSmith/graphql-example-go/source"
)

type TestProduct struct {
	Name, Info string
	Price      float64
}
type TestProductUpdate struct {
	Name, Info string
	Price      float64
	ID         int64
}

func CreateProduct(input TestProduct) (body []byte, err error) {
	config := source.GetConfig("..")

	// Parse the base URL
	URL, errParse := url.Parse("http://localhost:" + config.Server.Port + "/graphql?")
	if errParse != nil {
		return nil, errParse
	}

	query := `mutation{create(name:"{{.Name}}",info:"{{.Info}}",price:{{.Price}}){id,name,info,price}}`
	queryTemplate := template.Must(template.New("query").Parse(query))
	var queryParsed bytes.Buffer
	if errExecute := queryTemplate.Execute(&queryParsed, input); errExecute != nil {
		return nil, errExecute
	}

	// Parse the query parameters
	v := url.Values{}
	v.Add("query", queryParsed.String())

	// Add the encoded query parameters to the base URL and format as a String
	strURL := URL.String() + v.Encode()

	// Make the GraphQL mutation request
	res, errPost := http.Post(strURL, "application/json", bytes.NewBuffer([]byte("{}")))
	if errPost != nil {
		return nil, errPost
	}
	defer res.Body.Close()

	// Handle the response
	buff, errRead := ioutil.ReadAll(res.Body)
	if errRead != nil {
		return nil, errRead
	}

	return buff, nil
}

func DeleteProduct(ID int64) (body []byte, err error) {
	config := source.GetConfig("..")

	// Parse the base URL
	URL, errParse := url.Parse("http://localhost:" + config.Server.Port + "/graphql?")
	if errParse != nil {
		return nil, errParse
	}

	query := "mutation{delete(id:" + strconv.FormatInt(ID, 10) + "){id,name,info,price}}"

	// Parse the query parameters
	v := url.Values{}
	v.Add("query", query)

	// Add the encoded query parameters to the base URL and format as a String
	strURL := URL.String() + v.Encode()

	// Make the GraphQL mutation request
	res, errPost := http.Post(strURL, "application/json", bytes.NewBuffer([]byte("{}")))
	if errPost != nil {
		return nil, errPost
	}
	defer res.Body.Close()

	// Handle the response
	buff, errRead := ioutil.ReadAll(res.Body)
	if errRead != nil {
		return nil, errRead
	}

	return buff, nil
}

func UpdateProduct(input TestProductUpdate) (body []byte, err error) {
	config := source.GetConfig("..")

	// Parse the base URL
	URL, errParse := url.Parse("http://localhost:" + config.Server.Port + "/graphql?")
	if errParse != nil {
		return nil, errParse
	}

	query := `mutation{update(id:{{.ID}},name:"{{.Name}}",info:"{{.Info}}",price:{{.Price}}){id,name,info,price}}`
	queryTemplate := template.Must(template.New("query").Parse(query))
	var queryParsed bytes.Buffer
	if errExecute := queryTemplate.Execute(&queryParsed, input); errExecute != nil {
		return nil, errExecute
	}

	// Parse the query parameters
	v := url.Values{}
	v.Add("query", queryParsed.String())

	// Add the encoded query parameters to the base URL and format as a String
	strURL := URL.String() + v.Encode()

	// Make the GraphQL mutation request
	res, errPost := http.Post(strURL, "application/json", bytes.NewBuffer([]byte("{}")))
	if errPost != nil {
		return nil, errPost
	}
	defer res.Body.Close()

	// Handle the response
	buff, errRead := ioutil.ReadAll(res.Body)
	if errRead != nil {
		return nil, errRead
	}

	return buff, nil
}

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

func TestCreateAndDeleteProduct(t *testing.T) {
	// Create a new product
	input := TestProduct{
		Name:  "Inca Kola",
		Info:  "Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",
		Price: 1.99,
	}

	createBody, errCreate := CreateProduct(input)
	if errCreate != nil {
		t.Fatal(errCreate)
	}

	assert.NotNil(t, createBody)

	strCreateBody := string(createBody)

	id := gjson.Get(strCreateBody, "data.create.id").Int()
	name := gjson.Get(strCreateBody, "data.create.name").String()
	info := gjson.Get(strCreateBody, "data.create.info").String()
	price := gjson.Get(strCreateBody, "data.create.price").Float()

	assert.Equal(t, reflect.TypeOf(id).String(), "int64", "IDs should be equal")
	assert.Equal(t, input.Name, name, "Names should be equal")
	assert.Equal(t, input.Info, info, "Info should be equal")
	assert.Equal(t, float64(input.Price), price, "Prices should be equal")

	// Remove the created product
	deleteBody, errDelete := DeleteProduct(id)
	if errDelete != nil {
		t.Fatal(errDelete)
	}

	assert.NotNil(t, deleteBody)

	strDeleteBody := string(deleteBody)

	idDelete := gjson.Get(strDeleteBody, "data.delete.id").Int()
	nameDelete := gjson.Get(strDeleteBody, "data.delete.name").String()
	infoDelete := gjson.Get(strDeleteBody, "data.delete.info").String()
	priceDelete := gjson.Get(strDeleteBody, "data.delete.price").Float()

	assert.Equal(t, id, idDelete, "IDs should be equal")
	assert.Equal(t, nameDelete, name, "Names should be equal")
	assert.Equal(t, infoDelete, info, "Info should be equal")
	assert.Equal(t, priceDelete, price, "Prices should be equal")
}

func TestUpdateProduct(t *testing.T) {
	input := TestProductUpdate{
		ID:    1,
		Name:  "Chicha Morada (New)",
		Info:  "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki). [unknown source]",
		Price: 6.99,
	}

	buff, errUpdate := UpdateProduct(input)
	if errUpdate != nil {
		t.Fatal(errUpdate)
	}

	assert.NotNil(t, buff)

	strUpdateBody := string(buff)

	idUpdate := gjson.Get(strUpdateBody, "data.update.id").Int()
	nameUpdate := gjson.Get(strUpdateBody, "data.update.name").String()
	infoUpdate := gjson.Get(strUpdateBody, "data.update.info").String()
	priceUpdate := gjson.Get(strUpdateBody, "data.update.price").Float()

	assert.Equal(t, idUpdate, input.ID, "IDs should be equal")
	assert.Equal(t, nameUpdate, input.Name, "Names should be equal")
	assert.Equal(t, infoUpdate, input.Info, "Info should be equal")
	assert.Equal(t, priceUpdate, input.Price, "Prices should be equal")

	// Reverse the changes
	inputReverse := TestProductUpdate{
		ID:    1,
		Name:  "Chicha Morada",
		Info:  "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price: 7.99,
	}

	buffReverse, errUpdateReverse := UpdateProduct(inputReverse)
	if errUpdateReverse != nil {
		t.Fatal(errUpdateReverse)
	}

	assert.NotNil(t, buffReverse)

	strReverseBody := string(buffReverse)

	idReverse := gjson.Get(strReverseBody, "data.update.id").Int()
	nameReverse := gjson.Get(strReverseBody, "data.update.name").String()
	infoReverse := gjson.Get(strReverseBody, "data.update.info").String()
	priceReverse := gjson.Get(strReverseBody, "data.update.price").Float()

	assert.Equal(t, idReverse, inputReverse.ID, "IDs should be equal")
	assert.Equal(t, nameReverse, inputReverse.Name, "Names should be equal")
	assert.Equal(t, infoReverse, inputReverse.Info, "Info should be equal")
	assert.Equal(t, priceReverse, inputReverse.Price, "Prices should be equal")
}
