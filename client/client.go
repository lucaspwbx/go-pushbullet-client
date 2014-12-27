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
	noChannelTagError = errors.New("No channel tag parameter")
	noIdenError       = errors.New("No iden parameter")
	noFileNameError   = errors.New("No file name")
	noFileTypeError   = errors.New("No file type")
)

type HttpError struct {
	Status  int
	Message string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Status: %d, Message: %s", e.Status, e.Message)
}

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

// UPDATED - 12/2014 -> need new test
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

//12/2014 - needing a review
func (c *Client) UpdateMe(params map[string]Preferences) (User, error) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return User{}, err
	}
	body, _, err := c.do("POST", apiEndpoints["me"], bytes.NewBuffer(jsonParams))
	if err != nil {
		log.Println(err)
		return User{}, err
	}

	var user User
	if err = json.Unmarshal(body, &user); err != nil {
		return User{}, err
	}
	return user, nil
}

// UPDATED - 12/2014 -> need new test
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

//UPDATED - 12/2014 - need new test
func (c *Client) Subscriptions() (Subscriptions, error) {
	body, err := c.do("GET", apiEndpoints["subscriptions"], nil)
	if err != nil {
		return Subscriptions{}, err
	}

	var subscriptions Subscriptions
	if err = json.Unmarshal(body, &subscriptions); err != nil {
		return Subscriptions{}, err
	}
	return subscriptions, nil
}

//UPDATED - 12/2014 - need new test
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

//UPDATED - 12/2014 - need new tests
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
func (c *Client) GetContacts() (Contacts, error) {
	body, err := c.do("GET", apiEndpoints["contacts"], nil)
	if err != nil {
		return Contacts{}, err
	}

	var contacts Contacts
	if err = json.Unmarshal(body, &contacts); err != nil {
		return Contacts{}, err
	}
	return contacts, nil
}

//UPDATED - 12/2014 - need new tests
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

//UPDATED - 12/2014 - need new tests
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

//UPDATED - 12/2014 - need new tests
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

//UPDATED - 12/2014 - need new tests
func (c *Client) GetDevices() (Devices, error) {
	body, err := c.do("GET", apiEndpoints["devices"], nil)
	if err != nil {
		return Devices{}, err
	}

	var devices Devices
	if err = json.Unmarshal(body, &devices); err != nil {
		return Devices{}, err
	}
	return devices, nil
}

//UPDATED - 12/2014 - need new tests
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

//UPDATED - 12/2014 - need new tests
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

//UPDATED - 12/2014 - need new tests
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

//REVIEW
func (c *Client) GetPushes() (Pushes, error) {
	//TODO add params and allow modified_after
	body, _, err := c.do("GET", apiEndpoints["pushes"], nil)
	if err != nil {
		return Pushes{}, err
	}

	var pushes Pushes
	if err = json.Unmarshal(body, &pushes); err != nil {
		return Pushes{}, err
	}
	return pushes, nil
}

//DONE for note, link, address, checklist. Missing file implementation
func (c *Client) CreatePush(params Params) (Push, error) {
	if _, ok := params["type"]; !ok {
		return Push{}, errors.New("no type")
	}
	switch params["type"] {
	case "link":
		if _, ok := params["link"]; !ok {
			return Push{}, errors.New("no link for push of type link")
		}
	case "address":
		if _, ok := params["address"]; !ok {
			return Push{}, errors.New("no address for push of type address")
		}
	case "list":
		if _, ok := params["items"]; !ok {
			return Push{}, errors.New("No items for push of type checklist")
		}
	case "file":
		if _, ok := params["file_url"]; !ok {
			return Push{}, errors.New("No url for push of type file")
		}
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Push{}, err
	}
	body, _, err := c.do("POST", apiEndpoints["pushes"], bytes.NewBuffer(jsonParams))
	if err != nil {
		return Push{}, err
	}

	var push Push
	if err = json.Unmarshal(body, &push); err != nil {
		return Push{}, err
	}
	return push, nil
}

//UPDATED - 12/2014 - need new tests
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

//UPDATED - 12/2014 - need new tests
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

//UPDATED - 12/2014 - need new tests
func (c *Client) uploadRequest(params Params) (UploadRequest, error) {
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

//UPDATED - 12/2014 - need new tests
func (c *Client) Upload(filename, filetype, path string) error {
	req, err := c.UploadRequest(Params{
		"file_name": filename,
		"file_type": filetype,
	})
	if err != nil {
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		return err
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
		return err
	}
	if _, err = io.Copy(part, file); err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	uploadReq, err := http.NewRequest("POST", req.UploadUrl, body)
	uploadReq.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(uploadReq)
	if err != nil {
		return err
	}
	//data, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(data))
	return nil
}
