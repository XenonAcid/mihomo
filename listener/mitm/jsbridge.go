package mitm

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type jsRequest struct {
	Method string      `json:"method"`
	URL    string      `json:"url"`
	Header http.Header `json:"header"`
	Body   string      `json:"body"`
}

type jsResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Header     http.Header `json:"header"`
	Body       string      `json:"body"`
}

// RequestToJSON serializes *http.Request into JSON for JavaScript runtime.
func RequestToJSON(r *http.Request) ([]byte, error) {
	if r == nil {
		return nil, nil
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	jsReq := jsRequest{
		Method: r.Method,
		URL:    r.URL.String(),
		Header: r.Header,
		Body:   string(bodyBytes),
	}
	return json.Marshal(jsReq)
}

// ResponseToJSON serializes *http.Response into JSON for JavaScript runtime.
func ResponseToJSON(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, nil
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	jsResp := jsResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       string(bodyBytes),
	}
	return json.Marshal(jsResp)
}

// JSONToRequest parses JSON returned from JavaScript into an http.Request.
func JSONToRequest(data []byte) (*http.Request, error) {
	if len(data) == 0 {
		return nil, nil
	}
	var jsReq jsRequest
	if err := json.Unmarshal(data, &jsReq); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(jsReq.Method, jsReq.URL, io.NopCloser(bytes.NewReader([]byte(jsReq.Body))))
	if err != nil {
		return nil, err
	}
	req.Header = jsReq.Header
	return req, nil
}

// JSONToResponse parses JSON returned from JavaScript into an http.Response.
func JSONToResponse(data []byte, req *http.Request) (*http.Response, error) {
	if len(data) == 0 {
		return nil, nil
	}
	var jsResp jsResponse
	if err := json.Unmarshal(data, &jsResp); err != nil {
		return nil, err
	}
	res := &http.Response{
		Status:        jsResp.Status,
		StatusCode:    jsResp.StatusCode,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        jsResp.Header,
		Body:          io.NopCloser(bytes.NewReader([]byte(jsResp.Body))),
		ContentLength: int64(len(jsResp.Body)),
		Request:       req,
	}
	return res, nil
}
