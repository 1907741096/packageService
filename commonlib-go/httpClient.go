package lib

// @todo 该包会逐渐移除， 请使用和维护http.go

import (
	"strings"
	"time"
	"io/ioutil"
	"net/url"
	"encoding/json"
	"net/http"
)

type BaseInterface interface {
	SetDomain() string
	GetClient() BaseHttp
}

type BaseHttp struct {
	Domain string
	Client http.Client
	Status int
	Params interface{}
}

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CommonResponse struct {
	BaseResponse
	Data map[string]interface{} `json:"data"`
}

type Header map[string]string

var Domain string

const category = "API_REQUEST"

func (base *BaseHttp) Get(uri string, body map[string]string, headers map[string]string) ([]byte, error) {
	httpUrl := base.getUrl(uri)
	base.Params = body

	data := getForm(body)

	request, err := http.NewRequest("GET", httpUrl, nil)
	request.URL.RawQuery = data

	if err != nil {
		Error("Http Exception", err.Error(), category)
		return nil, err
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	return base.request(request)
}

func (base *BaseHttp) PostJson(uri string, body map[string]string, headers Header) ([]byte, error) {
	responseUrl := base.getUrl(uri)
	base.Params = body

	data := getJson(body)
	request, err := http.NewRequest("POST", responseUrl, strings.NewReader(data))
	if err != nil {
		Error("Http Exception", err.Error(), category)
		return nil, err
	}

	request.Header.Add("Content-type", "application/json")

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	return base.request(request)
}

func (base *BaseHttp) request(request *http.Request) ([]byte, error) {
	// 设置超时时间
	base.Client.Timeout = time.Duration(time.Second * 5)
	response, err := base.Client.Do(request)
	if err != nil {
		Error("Http Exception", err.Error(), category)
		return nil, err
	}

	result, err := ioutil.ReadAll(response.Body)
	Info("外部请求日志", map[string]interface{}{
		"uri":    request.URL,
		"params": base.Params,
		"result": string(result),
	}, "http")
	if err != nil {
		Error("Http Exception", err.Error(), category)
		return nil, err
	}
	return result, nil
}

func (base *BaseHttp) getUrl(uri string) string {
	return base.Domain + uri
}

func getForm(body map[string]string) string {
	var params = url.Values{}

	for key, value := range body {
		params.Add(key, value)
	}

	data := params.Encode()

	return data
}

func getJson(body map[string]string) string {
	requestJson, _ := json.Marshal(body)
	return string(requestJson)
}

type RequestMessage struct {
	Url               string
	Request, Response interface{}
}

func HttpInfo(url string, request, response interface{}) {
	message := RequestMessage{
		Url:      url,
		Request:  request,
		Response: response,
	}
	Error("http request", message, category)
}

func HttpError(url string, request, response interface{}) {
	message := RequestMessage{
		Url:      url,
		Request:  request,
		Response: response,
	}
	Error("http request error", message, category)
}
