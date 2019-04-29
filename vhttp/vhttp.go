package vhttp

import (
	"bytes"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	defaultQyHttp *VClient
)

type VClient struct{}

func Do() *VClient {
	if defaultQyHttp == nil {
		defaultQyHttp = new(VClient)
	}
	return defaultQyHttp
}

func (vc *VClient) Get(reqUrl string, params map[string]string) (reqBody []byte, err error) {
	if params != nil {
		reqUrl = reqUrl + ParseHttpParamForGet(params)
	}
	client := &http.Client{}
	resp, err := client.Get(reqUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	reqBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func (vc *VClient) PostForJson(reqUrl string, params interface{}, head map[string]string) (reqBody []byte, err error) {
	bts, err := jsoniter.Marshal(params)
	if err != nil {
		return
	}
	// TODO 这里可以改为 sync.Pool 对象池
	client := &http.Client{}
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(bts))
	req.Close = true
	req.Header.Add("Content-Type", "application/json")

	if head != nil {
		for k, v := range head {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if reqBody, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	return
}

func (vc *VClient) PostXForm(reqUrl string, params, head map[string]string) (reqBody []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(ParseHttpParamForGet(params)))
	req.Close = true
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if head != nil {
		for k, v := range head {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("%d", resp.StatusCode))
	}
	reqBody, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	return
}

func (vc *VClient) PostForJsonWithBaseAuth(reqUrl string, params interface{}, head map[string]string, userName, password string) (reqBody []byte, err error) {
	bts, err := jsoniter.Marshal(params)
	if err != nil {
		return
	}
	// TODO 这里可以改为 sync.Pool 对象池
	client := &http.Client{}
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(bts))
	req.Close = true
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(userName, password)
	if head != nil {
		for k, v := range head {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if reqBody, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	return
}
