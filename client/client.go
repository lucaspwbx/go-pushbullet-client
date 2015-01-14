package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var (
	v2Api        = "https://api.pushbullet.com/v2/"
	apiEndpoints = Endpoint{
		"contacts":       v2Api + "contacts",
		"pushes":         v2Api + "pushes",
		"devices":        v2Api + "devices",
		"me":             v2Api + "users/me",
		"subscriptions":  v2Api + "subscriptions",
		"channels":       v2Api + "channel-info",
		"upload_request": v2Api + "upload-request",
	}
	noChannelTagError   = errors.New("No channel tag parameter")
	noIdenError         = errors.New("No iden parameter")
	noFileNameError     = errors.New("No file name")
	noFileTypeError     = errors.New("No file type")
	pushNoLinkError     = errors.New("No link for push of type link")
	pushNoAddressError  = errors.New("No address for push of type address")
	pushNoItemsError    = errors.New("No items for push of type checklist")
	pushNoUrlError      = errors.New("No url for push of type file")
	pushNoTypeError     = errors.New("No type error")
	pushNoFileNameError = errors.New("No filename for push of type file")
	pushNoFileTypeError = errors.New("No filetype for push of type file")
)

// HttpError encapsulates HTTP request errors.
type HttpError struct {
	Status  int
	Message string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Status: %d, Message: %s", e.Status, e.Message)
}

// Used for HTTP requests.
func (c *Client) do(method, endpoint string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		return data, nil
	} else {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, &HttpError{Status: resp.StatusCode, Message: "Bad Request"}
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, &HttpError{Status: resp.StatusCode, Message: "Unauthorized"}
		}
		if resp.StatusCode == http.StatusForbidden {
			return nil, &HttpError{Status: resp.StatusCode, Message: "Forbidden"}
		}
		if resp.StatusCode == http.StatusNotFound {
			return nil, &HttpError{Status: resp.StatusCode, Message: "StatusNotFound"}
		}
		if resp.StatusCode == http.StatusInternalServerError {
			return nil, &HttpError{Status: resp.StatusCode, Message: "Internal Server Error"}
		}
	}
	return nil, nil
}

// Get information about user.
// See: https://api.pushbullet.com/v2/users/me
//
// Usage:
//   user, err := client.GetMe()
func (c *Client) GetMe() (User, error) {
	body, err := c.do("GET", apiEndpoints["me"], nil)
	if err != nil {
		return User{}, err
	}
	var user User
	if err = json.Unmarshal(body, &user); err != nil {
		return User{}, err
	}
	return user, nil
}

// Update information about user.
// See: https://api.pushbullet.com/v2/users/me
//
// Usage:
//   obj := make(map[string]client.Preferences)
//   obj["preferences"] = client.Preferences{Social: false}
//   user, err := client.UpdateMe(obj)

// TODO: improve implementation
func (c *Client) UpdateMe(params map[string]Preferences) (User, error) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return User{}, err
	}
	body, err := c.do("POST", apiEndpoints["me"], bytes.NewBuffer(jsonParams))
	if err != nil {
		return User{}, err
	}

	var user User
	if err = json.Unmarshal(body, &user); err != nil {
		return User{}, err
	}
	return user, nil
}

// Subscribe to a channel.
// See: https://api.pushbullet.com/v2/subscriptions

// Usage:
//   client.Subscribe(client.Params{
//     "channel_tag": "jblow"
//   })
//
// If no channel tag is passed a noChannelTagError will be returned.
func (c *Client) Subscribe(params Params) (Subscription, error) {
	if _, ok := params["channel_tag"]; !ok {
		return Subscription{}, noChannelTagError
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Subscription{}, err
	}
	body, err := c.do("POST", apiEndpoints["subscriptions"], bytes.NewBuffer(jsonParams))
	if err != nil {
		return Subscription{}, err
	}

	var subscription Subscription
	if err = json.Unmarshal(body, &subscription); err != nil {
		return Subscription{}, err
	}
	return subscription, nil
}

// Get all subscriptions.
// See: https://api.pushbullet.com/v2/subscriptions
//
// Usage:
//   subscriptions, err := client.Subscribtions()
func (c *Client) Subscriptions() ([]Subscription, error) {
	body, err := c.do("GET", apiEndpoints["subscriptions"], nil)
	if err != nil {
		return nil, err
	}

	var resultSet Subscriptions
	if err = json.Unmarshal(body, &resultSet); err != nil {
		return nil, err
	}
	return resultSet.Subscriptions, nil
}

// Get information about a channel.
// See: https://docs.pushbullet.com/v2/subscriptions/
//
// Usage:
//   channel, err := client.GetChannel(client.Params{"tag": "jblow"})
//
// If no channel tag is passed, a noChannelTagError will be returned.
func (c *Client) GetChannel(params Params) (Channel, error) {
	tag, ok := params["tag"]
	if !ok {
		return Channel{}, noChannelTagError
	}
	endpoint := fmt.Sprintf(apiEndpoints["channels"]+"?tag=%s", tag)
	body, err := c.do("GET", endpoint, nil)
	if err != nil {
		return Channel{}, err
	}

	var channel Channel
	if err = json.Unmarshal(body, &channel); err != nil {
		return Channel{}, err
	}
	return channel, nil
}

// Unsubscribe from a channel.
// See: https://api.pushbullet.com/v2/subscriptions
//
// Usage:
//   err := client.Subscribe(client.Params{"iden": "0xbababcdk"})
//
// If no iden is passed a noIdenError will be returned.
func (c *Client) Unsubscribe(params Params) error {
	id, ok := params["iden"]
	if !ok {
		return noIdenError
	}
	endpoint := fmt.Sprintf(apiEndpoints["subscriptions"]+"/%s", id)
	_, err := c.do("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	return nil
}

//UPDATED - 12/2014 - need new tests and review of active/non active contacts
// Get contacts.
// See: https://docs.pushbullet.com/v2/contacts/
//
// Usage:
//   contacts, err := client.GetContacts()
func (c *Client) GetContacts() ([]Contact, error) {
	body, err := c.do("GET", apiEndpoints["contacts"], nil)
	if err != nil {
		return nil, err
	}

	var resultSet Contacts
	if err = json.Unmarshal(body, &resultSet); err != nil {
		return nil, err
	}
	return resultSet.Contacts, nil
}

// Create contact.
// See: https://docs.pushbullet.com/v2/contacts/
//
// Usage:
//   contact, err := client.CreateContact(client.Params{"name": "foo", "email": "bar"})
func (c *Client) CreateContact(params Params) (Contact, error) {
	if _, ok := params["name"]; !ok {
		return Contact{}, errors.New("no name has been given")
	}
	if _, ok := params["email"]; !ok {
		return Contact{}, errors.New("no email has been given")
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Contact{}, err
	}
	body, err := c.do("POST", apiEndpoints["contacts"], bytes.NewBuffer(jsonParams))
	if err != nil {
		return Contact{}, err
	}

	var contact Contact
	if err = json.Unmarshal(body, &contact); err != nil {
		return Contact{}, err
	}
	return contact, nil
}

// Update contact.
// See: https://docs.pushbullet.com/v2/contacts/
//
// Usage:
//   contact, err := client.UpdateContact(client.Params{"iden": "0xyz", "name": "foo"})
//
// If no iden is passed a noIdenError is returned.
func (c *Client) UpdateContact(params Params) (Contact, error) {
	id, ok := params["iden"]
	if !ok {
		return Contact{}, noIdenError
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["contacts"]+"/%s", id)

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Contact{}, err
	}
	body, err := c.do("POST", endpoint, bytes.NewBuffer(jsonParams))
	if err != nil {
		return Contact{}, err
	}

	var contact Contact
	if err = json.Unmarshal(body, &contact); err != nil {
		return Contact{}, err
	}
	return contact, nil
}

// Delete contact.
// See: https://docs.pushbullet.com/v2/contacts/
//
// Usage:
//   contact, err := client.DeleteContact(client.Params{"iden": "0xyz")
//
// If no iden is passed a noIdenError is returned.
func (c *Client) DeleteContact(params Params) error {
	id, ok := params["iden"]
	if !ok {
		return noIdenError
	}
	endpoint := fmt.Sprintf(apiEndpoints["contacts"]+"/%s", id)
	_, err := c.do("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	return nil
}

// Get all devices.
// See: https://docs.pushbullet.com/v2/devices/
//
// Usage:
//   devices, err := client.GetDevices()
func (c *Client) GetDevices() ([]Device, error) {
	body, err := c.do("GET", apiEndpoints["devices"], nil)
	if err != nil {
		return nil, err
	}

	var resultSet Devices
	if err = json.Unmarshal(body, &resultSet); err != nil {
		return nil, err
	}
	return resultSet.Devices, nil
}

// Create device.
// See: https://docs.pushbullet.com/v2/devices/
//
// Usage:
//   device, err := client.CreateDevice(client.Params{"nickname": "foo", "type": "stream"})
func (c *Client) CreateDevice(params Params) (Device, error) {
	if _, ok := params["nickname"]; !ok {
		return Device{}, errors.New("no nickname has been given")
	}
	if _, ok := params["type"]; !ok {
		return Device{}, errors.New("no type has been given")
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Device{}, err
	}
	body, err := c.do("POST", apiEndpoints["devices"], bytes.NewBuffer(jsonParams))
	if err != nil {
		return Device{}, err
	}

	var device Device
	if err = json.Unmarshal(body, &device); err != nil {
		return Device{}, err
	}
	return device, nil
}

// Update device.
// See: https://docs.pushbullet.com/v2/devices/
//
// Usage:
//   device, err := client.UpdateDevice(client.Params{"iden": "0xyz", "nickname": "foo"})
func (c *Client) UpdateDevice(params Params) (Device, error) {
	id, ok := params["iden"]
	if !ok {
		return Device{}, noIdenError
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["devices"]+"/%s", id)

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Device{}, err
	}
	body, err := c.do("POST", endpoint, bytes.NewBuffer(jsonParams))
	if err != nil {
		return Device{}, err
	}

	var device Device
	if err = json.Unmarshal(body, &device); err != nil {
		return Device{}, err
	}
	return device, nil
}

// Delete device.
// See: https://docs.pushbullet.com/v2/devices/
//
// Usage:
//   err := client.DeleteDevice(client.Params{"iden": "0xyz"})

// If no iden is provided a noIdenError is returned.
func (c *Client) DeleteDevice(params Params) error {
	id, ok := params["iden"]
	if !ok {
		return noIdenError
	}
	endpoint := fmt.Sprintf(apiEndpoints["devices"]+"/%s", id)
	_, err := c.do("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	return nil
}

// Get pushes.
// See: https://docs.pushbullet.com/v2/pushes/
//
// Usage:
//   pushes, err := client.GetPushes()
func (c *Client) GetPushes() ([]Push, error) {
	//TODO add params and allow modified_after
	body, err := c.do("GET", apiEndpoints["pushes"], nil)
	if err != nil {
		return nil, err
	}

	var resultSet Pushes
	if err = json.Unmarshal(body, &resultSet); err != nil {
		return nil, err
	}
	return resultSet.Pushes, nil
}

// Create push.
// See: https://docs.pushbullet.com/v2/pushes/
//
// Usage:
//   push, err := client.CreatePush(client.Params{"type": "link", "title": "baz"})
//   push, err := client.CreatePush(client.Params{"type": "address", "address": "baz"})
//   push, err := client.CreatePush(client.Params{"type": "list", "title": "titulo", "items": []string{"foo", "bar"}})
//   push, err := client.CreatePush(client.Params{"type": "file", "file_name": "foo.txt", "file_type": "text/plain"})
func (c *Client) CreatePush(params Params) (Push, error) {
	if _, ok := params["type"]; !ok {
		return Push{}, pushNoTypeError
	}
	switch params["type"] {
	case "link":
		if _, ok := params["link"]; !ok {
			return Push{}, pushNoLinkError
		}
	case "address":
		if _, ok := params["address"]; !ok {
			return Push{}, pushNoAddressError
		}
	case "list":
		if _, ok := params["items"]; !ok {
			return Push{}, pushNoItemsError
		}
	case "file":
		if _, ok := params["file_name"]; !ok {
			return Push{}, pushNoFileNameError
		}
		if _, ok := params["file_type"]; !ok {
			return Push{}, pushNoFileTypeError
		}
	}
	if params["type"] == "file" {
		filename := params["file_name"].(string)
		filetype := params["file_type"].(string)
		url, err := c.PushFile(filename, filetype, filename)
		if err != nil {
			return Push{}, err
		}
		params["file_url"] = url
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Push{}, err
	}
	body, err := c.do("POST", apiEndpoints["pushes"], bytes.NewBuffer(jsonParams))
	if err != nil {
		return Push{}, err
	}

	var push Push
	if err = json.Unmarshal(body, &push); err != nil {
		return Push{}, err
	}
	return push, nil
}

// Update push.
// See: https://docs.pushbullet.com/v2/pushes/
//
// Usage:
//   push, err := client.UpdatePush(client.Params{"iden": "0xyz", "title": "foobaz"})
//
// If no iden is provided a noIdenError is returned.
func (c *Client) UpdatePush(params Params) (Push, error) {
	id, ok := params["iden"]
	if !ok {
		return Push{}, noIdenError
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["pushes"]+"/%s", id)

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Push{}, err
	}
	body, err := c.do("POST", endpoint, bytes.NewBuffer(jsonParams))
	if err != nil {
		return Push{}, err
	}

	var push Push
	if err = json.Unmarshal(body, &push); err != nil {
		return Push{}, err
	}
	return push, nil
}

// Delete push.
// See: https://docs.pushbullet.com/v2/pushes/
//
// Usage:
//   push, err := client.DeletePush(client.Params{"iden": "0xyz"})
//
// If no iden is provided a noIdenError is returned.
func (c *Client) DeletePush(params Params) error {
	id, ok := params["iden"]
	if !ok {
		return noIdenError
	}
	endpoint := fmt.Sprintf(apiEndpoints["pushes"]+"/%s", id)
	_, err := c.do("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	return nil
}

// Upload request.
// See: https://docs.pushbullet.com/v2/upload-request/
//
// Usage:
//   req, err := client.UploadRequest(client.Params{"file_name": "foo", "file_type": "text"})
func (c *Client) UploadRequest(params Params) (UploadRequest, error) {
	if _, ok := params["file_name"]; !ok {
		return UploadRequest{}, noFileNameError
	}
	if _, ok := params["file_type"]; !ok {
		return UploadRequest{}, noFileTypeError
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return UploadRequest{}, err
	}
	body, err := c.do("POST", apiEndpoints["upload_request"], bytes.NewBuffer(jsonParams))
	if err != nil {
		return UploadRequest{}, err
	}

	var uploadRequest UploadRequest
	if err = json.Unmarshal(body, &uploadRequest); err != nil {
		return UploadRequest{}, err
	}
	return uploadRequest, nil
}

// Push file.
// See: https://docs.pushbullet.com/v2/pushes/
//
// Usage:
//   fileUrl, err := client.PushFile("foo.txt", "text/plain", "foo.txt")
func (c *Client) PushFile(filename, filetype, path string) (string, error) {
	req, err := c.UploadRequest(Params{
		"file_name": filename,
		"file_type": filetype,
	})
	if err != nil {
		return "", err
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("awsaccesskeyid", req.Data.AwsAccessKeyId)
	writer.WriteField("acl", req.Data.Acl)
	writer.WriteField("key", req.Data.Key)
	writer.WriteField("signature", req.Data.Signature)
	writer.WriteField("policy", req.Data.Policy)
	writer.WriteField("content-type", req.Data.ContentType)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(part, file); err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}
	uploadReq, err := http.NewRequest("POST", req.UploadUrl, body)
	uploadReq.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(uploadReq)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusNoContent {
		return "", errors.New("error uploading file")
	}
	return req.FileUrl, nil
}
