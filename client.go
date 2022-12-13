package go_aptos_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net/http"
	"reflect"
)

type AleoClient struct {
	imp            *http.Client
	endpoint       string
	debug          bool
	faucetEndpoint string
}

func NewAleoClient(endpoint string, debug bool) (*AleoClient, error) {
	client := &AleoClient{
		endpoint: endpoint,
		imp:      http.DefaultClient,
		debug:    debug,
	}
	return client, nil
}

func (ap *AleoClient) GetLatestHeight() (uint64, error) {
	var result uint64
	err := ap.get("testnet3/latest/height", nil, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (ap *AleoClient) GetTransactionsByHeight(height uint64) (uint64, error) {
	var result uint64
	err := ap.post(fmt.Sprintf("testnet3/transactions/%v", height), nil, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (ap *AleoClient) Transaction(txId string) (uint64, error) {
	var result uint64
	err := ap.get(fmt.Sprintf("testnet3/transaction/%v", txId), nil, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (ap *AleoClient) post(method string, param Params, value interface{}, options ...Option) error {
	return ap.httpReq(http.MethodPost, method, param, value, options...)
}

func (ap *AleoClient) put(method string, param Params, value interface{}, options ...Option) error {
	return ap.httpReq(http.MethodPut, method, param, value, options...)
}

func (ap *AleoClient) delete(method string, param Params, value interface{}, options ...Option) error {
	return ap.httpReq(http.MethodDelete, method, param, value, options...)
}

func (ap *AleoClient) get(path string, params Params, value interface{}, options ...Option) error {
	for _, opt := range options {
		if params == nil {
			break
		}
		params.SetValue(opt.Key, opt.Value)
	}
	return ap.httpReq(http.MethodGet, fmt.Sprintf("%v?%v", path, params.Encode()), nil, value, []Option{}...)

}

func (ap *AleoClient) newRequest(httpMethod, url string, reqData []byte) (*http.Request, error) {
	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer MWRkMWMxMWEtMzUzMC00YTRmLTg5NDQtZjdkZDMwN2YwMjIy")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (ap *AleoClient) http(httpMethod, url string, reqData []byte) ([]byte, error) {
	request, err := ap.newRequest(httpMethod, url, reqData)
	if err != nil {
		return nil, err
	}
	response, err := ap.imp.Do(request)
	if err != nil {
		panic(err)
	}
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("%v %v %v", response.StatusCode, response.Status, string(data))
	}
	return data, nil
}

func (ap *AleoClient) httpReq(httpMethod, path string, param Params, value interface{}, options ...Option) (err error) {
	vi := reflect.ValueOf(value)
	if vi.Kind() != reflect.Ptr {
		return fmt.Errorf("value must be pointer")
	}
	if param != nil && len(options) > 0 {
		for _, opt := range options {
			param.SetValue(opt.Key, opt.Value)
		}
	}
	var requestData []byte
	if param != nil {
		requestData, err = json.Marshal(param)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	url := fmt.Sprintf("%s/%s", ap.endpoint, path)
	if ap.debug {
		fmt.Printf("request: %v  %v \n", url, string(requestData))
	}

	req, err := ap.newRequest(httpMethod, url, requestData)
	if err != nil {
		return err
	}
	resp, err := ap.imp.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		return err
	}
	if resp == nil || resp.StatusCode < http.StatusOK || resp.StatusCode > 300 {
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("response err: %v %v %v \n", resp.StatusCode, resp.Status, string(data))
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if ap.debug {
		fmt.Printf("reqponse: %v  %v \n", url, string(data))
	}
	if len(data) != 0 {
		err = json.Unmarshal(data, value)
		if err != nil {
			return fmt.Errorf("%v %s %s ", err, path, string(data))
		}
	}
	return nil

}
