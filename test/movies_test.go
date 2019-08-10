package moviestest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"text/template"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	source "github.com/HencoSmith/graphql-example-go/source"
)

type TestMovie struct {
	Name, Description string
	ReleaseYear       int64
}
type TestMovieUpdate struct {
	ID, Name, Description string
	ReleaseYear           int64
}

func CreateMovie(input TestMovie) (body []byte, err error) {
	config := source.GetConfig("..")

	// Parse the base URL
	URL, errParse := url.Parse("http://localhost:" + config.Server.Port + "/graphql?")
	if errParse != nil {
		return nil, errParse
	}

	query := `mutation{create(name:"{{.Name}}",description:"{{.Description}}",releaseYear:{{.ReleaseYear}}){id,name,description,release_year}}`
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

func DeleteMovie(ID string) (body []byte, err error) {
	config := source.GetConfig("..")

	// Parse the base URL
	URL, errParse := url.Parse("http://localhost:" + config.Server.Port + "/graphql?")
	if errParse != nil {
		return nil, errParse
	}

	query := "mutation{delete(id:" + ID + "){id,name,description,release_year}}"

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

func UpdateMovie(input TestMovieUpdate) (body []byte, err error) {
	config := source.GetConfig("..")

	// Parse the base URL
	URL, errParse := url.Parse("http://localhost:" + config.Server.Port + "/graphql?")
	if errParse != nil {
		return nil, errParse
	}

	query := `mutation{update(id:"{{.ID}}",name:"{{.Name}}",description:"{{.Description}}",releaseYear:{{.ReleaseYear}}){id,name,description,release_year}}`
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

func TestMovieList(t *testing.T) {
	config := source.GetConfig("..")
	res, err := http.Get("http://localhost:" + config.Server.Port + "/graphql?query={list{id,name,release_year,description,rating,review_count}}")
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, body)

	testIDData := [3]string{
		"13cbd25a-4a9d-4e71-9c39-4fc515083c95",
		"77034dd5-d3e4-4a44-a7fa-c2730dfe5370",
		"a774e5ff-a5f9-4643-832d-27d131344fe3",
	}
	testNameData := [3]string{
		"Scary Stories to Tell in the Dark",
		"Dora and the Lost City of Gold",
		"The Art of Racing in the Rain",
	}
	testDescriptionData := [3]string{
		"A group of teens face their fears in order to save their lives.",
		"Dora, a teenage explorer, leads her friends on an adventure to save her parents and solve the mystery behind a lost city of gold.",
		"Through his bond with his owner, aspiring Formula One race car driver Denny, golden retriever Enzo learns that the techniques needed on the racetrack can also be used to successfully navigate the journey of life.",
	}
	testReleaseYearData := [3]int64{
		2019,
		2019,
		2019,
	}

	values := gjson.Get(string(body), "data.list")

	assert.Equal(t, len(values.Array()), 3, "Item list must be of length 3")

	for count, item := range values.Array() {
		itemStr := item.String()
		id := gjson.Get(itemStr, "id").String()
		name := gjson.Get(itemStr, "name").String()
		description := gjson.Get(itemStr, "description").String()
		releaseYear := gjson.Get(itemStr, "release_year").Int()

		assert.Equal(t, testIDData[count], id, "IDs should be equal")
		assert.Equal(t, testNameData[count], name, "Names should be equal")
		assert.Equal(t, testDescriptionData[count], description, "Description should be equal")
		assert.Equal(t, testReleaseYearData[count], releaseYear, "Release Years should be equal")
	}
}

func TestGetMovie(t *testing.T) {
	config := source.GetConfig("..")
	res, err := http.Get("http://localhost:" + config.Server.Port + "/graphql?query={movie(id:\"13cbd25a-4a9d-4e71-9c39-4fc515083c95\"){id,name,release_year,description,rating,review_count}}")
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
	id := gjson.Get(bodyStr, "data.movie.id").String()
	name := gjson.Get(bodyStr, "data.movie.name").String()
	description := gjson.Get(bodyStr, "data.movie.description").String()
	releaseYear := gjson.Get(bodyStr, "data.movie.release_year").Int()

	assert.Equal(t, "13cbd25a-4a9d-4e71-9c39-4fc515083c95", id, "IDs should be equal")
	assert.Equal(t, "Scary Stories to Tell in the Dark", name, "Names should be equal")
	assert.Equal(t, "A group of teens face their fears in order to save their lives.", description, "Description should be equal")
	assert.Equal(t, int64(2019), releaseYear, "Release Years should be equal")
}

func TestCreateAndDeleteMovie(t *testing.T) {
	// Create a new movie
	input := TestMovie{
		Name:        "New Untitled Movie",
		Description: "Some very long description...",
		ReleaseYear: 2018,
	}

	createBody, errCreate := CreateMovie(input)
	if errCreate != nil {
		t.Fatal(errCreate)
	}

	assert.NotNil(t, createBody)

	strCreateBody := string(createBody)

	id := gjson.Get(strCreateBody, "data.create.id").String()
	name := gjson.Get(strCreateBody, "data.create.name").String()
	description := gjson.Get(strCreateBody, "data.create.description").String()
	releaseYear := gjson.Get(strCreateBody, "data.create.release_year").Int()

	assert.Equal(t, reflect.TypeOf(id), "String", "IDs should be equal")
	assert.Equal(t, input.Name, name, "Names should be equal")
	assert.Equal(t, input.Description, description, "Description should be equal")
	assert.Equal(t, int64(input.ReleaseYear), releaseYear, "Release Years should be equal")

	// Remove the created movie
	deleteBody, errDelete := DeleteMovie(id)
	if errDelete != nil {
		t.Fatal(errDelete)
	}

	assert.NotNil(t, deleteBody)

	strDeleteBody := string(deleteBody)

	idDelete := gjson.Get(strDeleteBody, "data.delete.id").Int()
	nameDelete := gjson.Get(strDeleteBody, "data.delete.name").String()
	descriptionDelete := gjson.Get(strDeleteBody, "data.delete.description").String()
	releaseYearDelete := gjson.Get(strDeleteBody, "data.delete.release_year").Int()

	assert.Equal(t, id, idDelete, "IDs should be equal")
	assert.Equal(t, nameDelete, name, "Names should be equal")
	assert.Equal(t, descriptionDelete, descriptionDelete, "Description should be equal")
	assert.Equal(t, releaseYearDelete, releaseYearDelete, "Release Years should be equal")
}

func TestUpdateMovie(t *testing.T) {
	input := TestMovieUpdate{
		ID:          "13cbd25a-4a9d-4e71-9c39-4fc515083c95",
		Name:        "Scary Stories to Tell in the Dark (New)",
		Description: "A group of teens face their fears in order to save their lives. [unknown source]",
		ReleaseYear: 2017,
	}

	buff, errUpdate := UpdateMovie(input)
	if errUpdate != nil {
		t.Fatal(errUpdate)
	}

	assert.NotNil(t, buff)

	strUpdateBody := string(buff)

	idUpdate := gjson.Get(strUpdateBody, "data.update.id").Int()
	nameUpdate := gjson.Get(strUpdateBody, "data.update.name").String()
	descriptionUpdate := gjson.Get(strUpdateBody, "data.update.description").String()
	releaseYearUpdate := gjson.Get(strUpdateBody, "data.update.release_year").Int()

	assert.Equal(t, idUpdate, input.ID, "IDs should be equal")
	assert.Equal(t, nameUpdate, input.Name, "Names should be equal")
	assert.Equal(t, descriptionUpdate, input.Description, "Description should be equal")
	assert.Equal(t, releaseYearUpdate, input.ReleaseYear, "Release Years should be equal")

	// Reverse the changes
	inputReverse := TestMovieUpdate{
		ID:          "13cbd25a-4a9d-4e71-9c39-4fc515083c95",
		Name:        "Scary Stories to Tell in the Dark",
		Description: "A group of teens face their fears in order to save their lives.",
		ReleaseYear: 2018,
	}

	buffReverse, errUpdateReverse := UpdateMovie(inputReverse)
	if errUpdateReverse != nil {
		t.Fatal(errUpdateReverse)
	}

	assert.NotNil(t, buffReverse)

	strReverseBody := string(buffReverse)

	idReverse := gjson.Get(strReverseBody, "data.update.id").Int()
	nameReverse := gjson.Get(strReverseBody, "data.update.name").String()
	descriptionReverse := gjson.Get(strReverseBody, "data.update.description").String()
	releaseYearReverse := gjson.Get(strReverseBody, "data.update.release_year").Int()

	assert.Equal(t, idReverse, inputReverse.ID, "IDs should be equal")
	assert.Equal(t, nameReverse, inputReverse.Name, "Names should be equal")
	assert.Equal(t, descriptionReverse, inputReverse.Description, "Description should be equal")
	assert.Equal(t, releaseYearReverse, inputReverse.ReleaseYear, "Release Years should be equal")
}
