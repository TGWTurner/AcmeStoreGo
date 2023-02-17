package test

import (
	"bjssStoreGo/backend/layers/api"
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
)

func NewApiRequester(w *api.Wiring) *ApiRequester {
	return &ApiRequester{
		apiBase: "http://localhost:4001",
		wiring:  w,
	}
}

func (ar *ApiRequester) Get(path string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, ar.apiBase+path, nil)

	if ar.session != "" {
		request.Header.Set("Cookie", ar.session)
	}

	response := ar.executeRequest(
		request,
		ar.wiring,
	)

	return response
}

func (ar *ApiRequester) Post(path string, body []byte) *httptest.ResponseRecorder {
	requestBody := bytes.NewBuffer(body)
	request := httptest.NewRequest(http.MethodPost, ar.apiBase+path, requestBody)
	request.Header.Set("Content-Type", "application/json")

	if ar.session != "" {
		request.Header.Set("Cookie", ar.session)
	}

	response := ar.executeRequest(
		request,
		ar.wiring,
	)

	cookie := response.Result().Header.Get("set-cookie")

	if len(cookie) > 0 {
		ar.session = strings.Split(cookie, ";")[0]
	}

	return response
}

func (ar *ApiRequester) executeRequest(req *http.Request, wiring *api.Wiring) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	wiring.Router.ServeHTTP(rr, req)

	return rr
}

type ApiRequester struct {
	apiBase string
	wiring  *api.Wiring
	session string
}
