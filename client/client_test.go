package client

import (
	"encoding/json"
	"fmt"
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
		t.Errorf("Error unmarshalling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetDevices()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestGetContacts(t *testing.T) {
	body :=
		`
		{
		  "contacts": [
		  {
		    "iden": "ubdcjAfszs0Smi",
		    "name": "Ryan Oldenburg",
		    "created": 1399011660.4298899,
		    "modified": 1399011660.42976,
		    "email": "ryanjoldenburg@gmail.com",
		    "email_normalized": "ryanjoldenburg@gmail.com",
		    "active": true
		  }
		  ]
		}
		`
	var expected Contacts
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		fmt.Println(err)
		t.Errorf("Error unmarshalling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetContacts()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}

func TestGetPushes(t *testing.T) {
	body := `
	{
	"pushes": [
	    {
	      "iden": "ubdprOsjAhOzf0XYq",
	      "type": "link",
	      "title": "Pushbullet",
	      "body": "Documenting our API",
	      "url": "http://docs.pushbullet.com",
	      "created": 1411595135.9685705,
	      "modified": 1411595135.9686127,
	      "active": true,
	      "dismissed": false,
	      "sender_iden": "ubd",
	      "sender_email": "ryan@pushbullet.com",
	      "sender_email_normalized": "ryan@pushbullet.com",
	      "receiver_iden": "ubd",
	      "receiver_email": "ryan@pushbullet.com",
	      "receiver_email_normalized": "ryan@pushbullet.com"
	    }
	    ]
	}
	`
	var expected Pushes
	err := json.Unmarshal([]byte(body), &expected)
	if err != nil {
		fmt.Println(err)
		t.Errorf("Error unmarshalling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetPushes()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
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
		fmt.Println(err)
		t.Errorf("Error unmarshalling JSON")
	}
	fakeRT := &FakeRoundTripper{message: body, status: http.StatusOK}
	client := newTestClient(fakeRT)
	got, _ := client.GetMe()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Error, expected %#v, got %#v", expected, got)
	}
}
