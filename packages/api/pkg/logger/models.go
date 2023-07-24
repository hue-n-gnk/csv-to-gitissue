package logger

import (
	"bytes"
	"net/http"
	"net/url"
	"time"
)

type apiLogger struct {
	srv http.Handler
	log Logger
}

type respWriter struct {
	w          http.ResponseWriter
	r          *http.Request
	buf        *bytes.Buffer
	statusCode int
	log        Logger
}

type LogRequest struct {
	Id            string
	StartTime     time.Time
	RequestMethod string
	Url           *url.URL

	RemoteAddr  string
	Referer     string
	UserAgent   string
	Protocol    string
	RequestSize string
	Header      http.Header
	Body        any
}

type LogResponse struct {
	Id            string
	StartTime     time.Time
	RequestMethod string
	Url           *url.URL
	StatusCode    int
	Header        http.Header
	Body          any
}

type GraphQLRequest struct {
	Id            string
	TimeStart     time.Time
	OperationName string
	RawQuery      string
	Variable      map[string]interface{}
}
type GraphQLResponse struct {
	Id         string
	TimeStart  time.Time
	Errors     string
	Data       any
	Extensions map[string]interface{}
}
