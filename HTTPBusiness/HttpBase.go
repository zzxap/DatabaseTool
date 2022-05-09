package HTTPBusiness

import (
	//"fmt"
	"bytes"
	"crypto/tls"
	"DataBaseManage/public"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var transport *http.Transport
var client *http.Client

func HttpRequest(urls string, paras map[string][]string) (body []byte, err error) {

	if transport == nil {
		if strings.Contains(urls, "https") {

			transport = &http.Transport{
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
				DisableKeepAlives: true,
			}
		} else {
			transport = &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			}
		}
		client = &http.Client{
			Transport: transport,
			Timeout:   time.Second * 10,
		}
	}
	//json := `{"key":"value"}`
	//b := strings.NewReader(json)
	//w, errr := client.Post(url, "application/x-www-form-urlencoded", b)
	w, errr := client.PostForm(urls, paras) //strings.NewReader("name=cjb")
	if errr != nil {
		public.Log(errr)
		public.Log(urls)
		return nil, errr
	}

	bodyy, err := ioutil.ReadAll(w.Body)
	if err != nil {
		public.Log(err)
		return nil, err
	}
	w.Body.Close()
	return bodyy, err

}

func myhttpRequest(url string, params url.Values) (body []byte, err error) {

	if transport == nil {
		if strings.Contains(url, "https") {

			transport = &http.Transport{
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
				DisableKeepAlives: true,
			}
		} else {
			transport = &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			}
		}
		client = &http.Client{
			Transport: transport,
			Timeout:   time.Second * 10,
		}

	}

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(params.Encode()))
	if err != nil {
		public.Log("Error Occured. %+v", err)
		return nil, err
	}
	//("Authorization", " Bearer " + authorization);
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, errr := client.Do(req)
	//res, errr := client.Post(url, "application/x-www-form-urlencoded", nil) //strings.NewReader("name=cjb")
	if errr != nil {
		public.Log("client.Post error")
		public.Log(errr)
		public.Log(url)
		return nil, errr
	}

	bodyy, err := ioutil.ReadAll(res.Body)
	if err != nil {
		public.Log("client.Post read error")
		public.Log(err)
		return nil, err
	}
	res.Body.Close()
	return bodyy, err

}
