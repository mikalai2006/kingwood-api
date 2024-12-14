package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

var TestLanuageData = domain.Language{
	Name:      "testLang",
	Code:      "ru",
	Flag:      "flag",
	Publish:   true,
	Locale:    "ru-Ru",
	SortOrder: 1,
}

var langPath = "/api/v1/lang"

func (s *TestSuite) TestCreateLangNotAuth() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	dataJSON, err := json.Marshal(TestLanuageData)
	s.NoError(err)

	// test invalid header
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		langPath,
		bytes.NewBuffer(dataJSON),
	)
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer")
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	// test empty token
	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		langPath,
		bytes.NewBuffer(dataJSON),
	)
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer ")
	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()

	r.Equal(http.StatusUnauthorized, response.StatusCode)
}

func (s *TestSuite) TestCreateLangAuth() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))

	auth, err := s.Auth(router)
	s.NoError(err)

	coll := s.db.Collection(repository.TblLanguage)
	err = coll.Drop(context.Background())
	s.NoError(err)

	r := s.Require()

	dataJSON, err := json.Marshal(TestLanuageData)
	s.NoError(err)

	// test invalid auth token.
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		langPath,
		bytes.NewBuffer(dataJSON),
	)
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer 12345")
	req.Close = true
	s.NoError(err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	response := w.Result()
	defer response.Body.Close()
	r.Equal(http.StatusUnauthorized, response.StatusCode)

	// create with auth user.
	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		langPath,
		bytes.NewBuffer(dataJSON),
	)
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()

	var re domain.Language
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	r.Equal(re.Code, TestLanuageData.Code)
	r.Equal(re.Flag, TestLanuageData.Flag)
	r.Equal(re.Name, TestLanuageData.Name)
	r.Equal(re.Publish, TestLanuageData.Publish)
	r.Equal(re.SortOrder, TestLanuageData.SortOrder)
	r.Equal(re.Locale, TestLanuageData.Locale)

	r.Equal(http.StatusOK, response.StatusCode)
}

func (s *TestSuite) TestFindLangByLimitOne() {
	limit := 1

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, langPath, nil)
	q := req.URL.Query()
	q.Add("$limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	req.Close = true
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var re domain.Response[domain.Language]
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	r.Equal(re.Limit, limit)
	r.Equal(len(re.Data), limit)

	r.Equal(http.StatusOK, response.StatusCode)
}
func (s *TestSuite) TestFindLangByLimitBig() {
	limit := 1000

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, langPath, nil)
	q := req.URL.Query()
	q.Add("$limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	req.Close = true
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var re domain.Response[domain.Language]
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	r.Equal(10, re.Limit)
	r.Equal(re.Total, len(re.Data))

	r.Equal(http.StatusOK, response.StatusCode)
}

func (s *TestSuite) TestFindLangBySkip() {
	skip := 1

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, langPath, nil)
	q := req.URL.Query()
	q.Add("$skip", fmt.Sprintf("%d", skip))
	req.URL.RawQuery = q.Encode()
	req.Close = true
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var re domain.Response[domain.Language]
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	r.Equal(re.Skip, skip)
	r.Equal(len(re.Data), re.Total)

	r.Equal(http.StatusOK, response.StatusCode)
}

func (s *TestSuite) TestFindLangBySort() {
	sort := -1

	router := gin.New()
	s.handler.Init(router.Group("/api"))

	// create two item.
	auth, err := s.Auth(router)
	s.NoError(err)

	testLanuageDataTwo := domain.Language{
		Name:      "Two",
		SortOrder: 10,
		Publish:   true,
	}
	dataJSON, err := json.Marshal(testLanuageDataTwo)
	s.NoError(err)
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		langPath,
		bytes.NewBuffer(dataJSON),
	)
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var twoItem domain.Language
	err = json.NewDecoder(response.Body).Decode(&twoItem)
	s.NoError(err)

	r := s.Require()

	req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, langPath, nil)
	q := url.Values{}
	q.Add("$sort[sort_order]", fmt.Sprintf("%v", sort))
	req.URL.RawQuery = q.Encode()
	req.Close = true
	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()

	var re domain.Response[domain.Language]
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	r.Equal(re.Data[0].SortOrder, testLanuageDataTwo.SortOrder)
	// r.Equal(len(re.Data), re.Total)

	r.Equal(http.StatusOK, response.StatusCode)
}

func (s *TestSuite) TestPatchLang() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))

	// create item.
	auth, err := s.Auth(router)
	s.NoError(err)

	testData := domain.Language{
		Name:      "For patch",
		SortOrder: 2,
		Publish:   true,
	}
	dataJSON, err := json.Marshal(testData)
	patchData := domain.Language{
		Name:      "For patch(patch)",
		SortOrder: 3,
		Publish:   false,
	}
	dataPatchJSON, err := json.Marshal(patchData)

	s.NoError(err)
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		langPath,
		bytes.NewBuffer(dataJSON),
	)
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	s.NoError(err)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	response := w.Result()
	defer response.Body.Close()
	var twoItem domain.Language
	err = json.NewDecoder(response.Body).Decode(&twoItem)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
	r := s.Require()

	// test empty id.
	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodPatch,
		fmt.Sprintf("%s/%s", langPath, ""),
		nil,
	)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Close = true
	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()
	r.Equal(http.StatusNotFound, response.StatusCode)

	// test  empty patch data.
	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodPatch,
		fmt.Sprintf("%s/%s", langPath, "12345"),
		nil,
	)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Close = true
	s.NoError(err)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	response = w.Result()
	defer response.Body.Close()
	r.Equal(http.StatusBadRequest, response.StatusCode)

	// test not objectId as id.
	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodPatch,
		fmt.Sprintf("%s/%s", langPath, "12345"),
		bytes.NewBuffer(dataPatchJSON),
	)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Close = true
	s.NoError(err)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	response = w.Result()
	defer response.Body.Close()
	r.Equal(http.StatusInternalServerError, response.StatusCode)

	// test with id.
	idForPatch := twoItem.ID.Hex()
	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodPatch,
		fmt.Sprintf("%s/%s", langPath, idForPatch),
		bytes.NewBuffer(dataPatchJSON),
	)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Close = true
	s.NoError(err)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	response = w.Result()
	defer response.Body.Close()
	var re domain.Language
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	r.Equal(re.ID.Hex(), twoItem.ID.Hex())
	// r.Equal(len(re.Data), re.Total)

	r.Equal(http.StatusOK, response.StatusCode)
}

func (s *TestSuite) TestDeleteLang() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))

	// create item.
	auth, err := s.Auth(router)
	s.NoError(err)

	testData := domain.Language{
		Name:      "For remove",
		SortOrder: 1,
		Publish:   false,
	}
	dataJSON, err := json.Marshal(testData)
	s.NoError(err)
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		langPath,
		bytes.NewBuffer(dataJSON),
	)
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var twoItem domain.Language
	err = json.NewDecoder(response.Body).Decode(&twoItem)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	r := s.Require()

	// test empty id.
	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodDelete,
		fmt.Sprintf("%s/%s", langPath, ""),
		nil,
	)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Close = true
	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()
	r.Equal(http.StatusNotFound, response.StatusCode)

	// // test empty id.
	// req, err = http.NewRequestWithContext(
	// 	context.Background(),
	// 	http.MethodDelete,
	// 	fmt.Sprintf("%s/%s", langPath, ""),
	// 	nil,
	// )
	// req.Header.Set("Content-type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	// req.Close = true
	// s.NoError(err)
	// w = httptest.NewRecorder()
	// router.ServeHTTP(w, req)
	// response = w.Result()
	// defer response.Body.Close()
	// r.Equal(http.StatusBadRequest, response.StatusCode)

	// test not objectId as id.
	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodDelete,
		fmt.Sprintf("%s/%s", langPath, "12345"),
		nil,
	)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Close = true
	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()
	r.Equal(http.StatusBadRequest, response.StatusCode)

	// test with id.
	idForRemove := twoItem.ID.Hex()
	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodDelete,
		fmt.Sprintf("%s/%s", langPath, idForRemove),
		nil,
	)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Close = true

	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()

	var re domain.Language
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	// r.Equal(re.ID.Hex(), twoItem.ID.Hex())
	// r.Equal(len(re.Data), re.Total)

	r.Equal(http.StatusOK, response.StatusCode)
}
