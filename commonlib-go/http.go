package lib

import (
	"time"
	"encoding/json"
	"net/http"
	"errors"
	"fmt"
	"strings"
	"net/url"
	"io"
)

type Http interface{
	Url(uri string) string
	GetHeader() Header
}

const Timeout = 30

func (h Header) SetHeader(key, value string) Header{
	h[key] = value
	return h
}

// 响应参数
type Response struct {
	Code 		int 		`json:"code"`
	Data 		interface{} `json:"data"`
	Message 	string		`json:"message"`
}

type ResponseMessage struct {
	Message 	string
	Level 		string
}

func (r *ResponseMessage) Error() string {
	Error("response message error", r.Message, "http")
	return fmt.Sprintf("[%s] response message: [%s]", r.Level, r.Message)
}

const (
	POST = "POST"
	GET  = "GET"
)

var (
	errRespCodeException = errors.New("response code not correct")
	errResMethodNotFound = errors.New("request method not found")
	errGetParamException = errors.New("get params must be map[string]string")
	errHttpUnexpected    = errors.New("unexpected http error")
)

func (q *Response) Error() string {
	return fmt.Sprintf("response code: [%d], data: [%s], message: [%s]", q.Code, q.Data, q.Message)
}

func PostJson(h Http, action string, params interface{}, timeout int)(io.ReadCloser, error){
	header := h.GetHeader().SetHeader("Content-type", "application/json")
	return sendRequest(h, POST, action, params, timeout, header)
}

func Post(h Http, action string, params interface{}, timeout int)(io.ReadCloser, error){
	return sendRequest(h, POST, action, params, timeout, h.GetHeader())
}

func Get(h Http, action string, params interface{}, timeout int)(io.ReadCloser, error){
	return sendRequest(h, GET, action, params, timeout, h.GetHeader())
}

func sendRequest(h Http, method, action string, params interface{}, timeout int, header Header) (io.ReadCloser, error){
	u := h.Url(action)
	client := new(http.Client)
	b, err := getBody(method, params, true)
	if err != nil {
		return nil, err
	}
	// 获取 http.Request
	req, err := getReq(method, u, b, header)
	if err != nil {
		return nil, err
	}
	client.Timeout = time.Duration(timeout) * time.Second
	resp, err := client.Do(req)
	Info("http request", map[string]interface{}{"param": params, "url": u, "header": header}, "http")
	if err != nil {
		return nil, err
	}
	Info("http request success", map[string]interface{}{"param": params, "url": u, "header": header, "status": resp.Status, "resp": resp.StatusCode}, "http")
	if resp.StatusCode != 200 {
		return nil, errRespCodeException
	}


	return resp.Body, nil
}

func getBody(method string, params interface{}, isJson bool) (string, error){
	defer func() (string, error){
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				return "", err
			}
		}

		return "", nil
	}()
	switch method {
	case GET:
		v := url.Values{}
		if b, ok := params.(map[string]string); ok {
			for key, value := range b {
				v.Add(key, value)
			}
			return v.Encode(), nil
		}

		return "", errGetParamException
	case POST:
		if isJson {
			b, err := json.Marshal(params)
			return string(b), err
		} else {
			return "", errors.New("当前内容未补充完整，用到的大佬自己补充代码")
		}
	default:
		return "", errResMethodNotFound
	}

	panic(errHttpUnexpected)
}

func getReq(method, url string, b string, header Header) (*http.Request, error){
	var err error
	var req *http.Request
	switch method {
	case GET:
		req, err = http.NewRequest(method, url, nil)
		req.URL.RawQuery = b
	case POST:
		req, err = http.NewRequest(method, url, strings.NewReader(b))
	}

	if err != nil {
		return nil, err
	}
	for key, value := range header {
		req.Header.Add(key, value)
	}
	return req, nil
}