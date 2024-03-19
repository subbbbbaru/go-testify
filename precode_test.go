package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMainHandler(t *testing.T) {
	type args struct {
		url                string
		expectedStatusCode int
		expectedCount      int
		expectedBody       string
	}

	reqOK := args{
		url:                "/cafe?count=4&city=moscow",
		expectedStatusCode: http.StatusOK,
		expectedCount:      2,
		expectedBody:       "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент",
	}

	badCity := args{
		url:                "/cafe?count=10&city=paris",
		expectedStatusCode: http.StatusBadRequest,
		expectedCount:      0,
		expectedBody:       "wrong city value",
	}

	countMoreThanTotal := args{
		url:                "/cafe?count=5&city=moscow",
		expectedStatusCode: http.StatusOK,
		expectedCount:      4,
		expectedBody:       "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент",
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Good request",
			args: reqOK,
		},
		{
			name: "Wrong city value",
			args: badCity,
		},
		{
			name: "Count more than total",
			args: countMoreThanTotal,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := getResponse(test.args.url)
			require.Equal(t, test.args.expectedStatusCode, resp.Code,
				"wrong status code, want %d, got %d", test.args.expectedStatusCode, resp.Code)

			require.Equal(t, test.args.expectedBody, resp.Body.String(),
				"wrong body, want %s, got %s", test.args.expectedBody, resp.Body.String())
		})
	}
}

func getResponse(url string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	return responseRecorder
}
