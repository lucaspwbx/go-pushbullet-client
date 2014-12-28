package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type FakeRoundTripper struct {
	message  string
	status   int
	header   map[string]string
	requests []*http.Request
}

func newTestClient(rt *FakeRoundTripper) *Client {
	client := &Client{
		token:      "foobar",
		HttpClient: &http.Client{Transport: rt},
	}
	return client
}

func (rt *FakeRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	body := strings.NewReader(rt.message)
	rt.requests = append(rt.requests, r)
	res := &http.Response{
		StatusCode: rt.status,
		Body:       ioutil.NopCloser(body),
		Header:     make(http.Header),
	}
	for k, v := range rt.header {
		res.Header.Set(k, v)
	}
	return res, nil
}

func (rt *FakeRoundTripper) Reset() {
	rt.requests = nil
}

func TestGetMe(t *testing.T) {
	body := `
	{
	  "iden": "ubdpjxxxOK0sKG",
	  "email": "ryan@pushbullet.com",
	  "email_normalized": "ryan@pushbullet.com",
	  "created": 1357941753.8287899,
	  "modified": 1399325992.1842301,
	  "name": "Ryan Oldenburg",
	  "image_url": "https://lh4.googleusercontent.com/-YGdcF2MteeI/AAAAAAAAAAI/AAAAAAAADPU/uo9z33FoEYs/photo.jpg",
	  "preferences": {
	    "onboarding": {
	      "app": false,
	      "friends": false,
	      "extension": false
	    },
	    "social": false
	  }
	}
    `
	var expected User
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetMe()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

// TODO
func TestUpdateMe(t *testing.T) {
}

func TestSubscribeNoChannelTag(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	_, err := client.Subscribe(Params{})
	if err != noChannelTagError {
		t.Errorf("Error, expected %#v, got %#v", noChannelTagError.Error(), err)
	}
}

func TestSubscribeChannel(t *testing.T) {
	body := `
	{
	  "iden": "udprOsjAoRtnM0jc",
	  "created": 1412047948.579029,
	  "modified": 1412047948.5790315,
	  "active": true,
	  "channel": {
	    "iden": "ujxPklLhvyKsjAvkMyTVh6",
	    "tag": "jblow",
	    "name": "Jonathan Blow",
	    "description": "New comments on the web by Jonathan Blow.",
	    "image_url": "https://pushbullet.imgix.net/ujxPklLhvyK-6fXf4O2JQ1dBKQedhypIKwPX0lyFfwXW/jonathan-blow.png"
	  }
	}
	`
	var expected Subscription
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.Subscribe(Params{"channel_tag": "jblow"})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestGetChannelNoChannel(t *testing.T) {
	fakeRT := &FakeRoundTripper{message: "", status: http.StatusOK}
	client := newTestClient(fakeRT)
	_, err := client.GetChannel(Params{})
	if err != noChannelTagError {
		t.Errorf("Error, expected %#v, got %#v", noChannelTagError.Error(), err)
	}
}

func TestGetChannel(t *testing.T) {
	body := `
	{
	  "iden": "ujxPklLhvyKsjAvkMyTVh6",
	  "tag": "jblow",
	  "name": "Jonathan Blow",
	  "description": "New comments on the web by Jonathan Blow.",
	  "image_url": "https://pushbullet.imgix.net/ujxPklLhvyK-6fXf4O2JQ1dBKQedhypIKwPX0lyFfwXW/jonathan-blow.png"
	}
	`
	var expected Channel
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetChannel(Params{"tag": "jblow"})
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestGetDevices(t *testing.T) {
	body :=
		`
		{
		  "devices": [
		  {
		    "iden": "u1qSJddxeKwOGuGW",
		    "push_token": "u1qSJddxeKwOGuGWu1qdxeKwOGuGWu1qSJddxeK",
		    "app_version": 74,
		    "fingerprint": "<json_string>",
		    "active": true,
		    "nickname": "Galaxy S4",
		    "manufacturer": "samsung",
		    "type": "android",
		    "created": 1394748080.0139201,
		    "modified": 1399008037.8487799,
		    "model": "SCH-I545",
		    "pushable": true
		  }
		  ]
		}
		`
	var expected Devices
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetDevices()
	if !reflect.DeepEqual(got, expected.Devices) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}
