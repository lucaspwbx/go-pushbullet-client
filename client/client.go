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
)

func (c *Client) do(method, endpoint string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	fmt.Println(req)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	return data, resp.StatusCode, nil
}

//DONE
func (c *Client) GetMe() (User, error) {
	body, _, err := c.do("GET", apiEndpoints["me"], nil)
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

//TODO - test manually and with unit test - not working
func (c *Client) UpdateMe(params Params) (User, error) {
	pref, ok := params["preferences"]
	if !ok {
		return User{}, errors.New("No preferences")
	}
	//delete(params, "preferences")

	//jsonParams, err := json.Marshal(params)
	jsonParams, err := json.Marshal(pref)
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

//DONE - need unit test
func (c *Client) Subscribe(params Params) (Subscription, error) {
	if _, ok := params["channel_tag"]; !ok {
		return Subscription{}, errors.New("no channel tag parameter")
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Subscription{}, err
	}
	body, status, err := c.do("POST", apiEndpoints["subscriptions"], bytes.NewBuffer(jsonParams))
	fmt.Println(status)
	if err != nil {
		log.Println(err)
		return Subscription{}, err
	}

	var subscription Subscription
	if err = json.Unmarshal(body, &subscription); err != nil {
		return Subscription{}, err
	}
	return subscription, nil
}

//DONE - need unit test
func (c *Client) Subscriptions() (Subscriptions, error) {
	body, _, err := c.do("GET", apiEndpoints["subscriptions"], nil)
	if err != nil {
		log.Println(err)
		return Subscriptions{}, err
	}

	var subscriptions Subscriptions
	if err = json.Unmarshal(body, &subscriptions); err != nil {
		return Subscriptions{}, err
	}
	return subscriptions, nil
}

//DONE - need unit tests
func (c *Client) GetChannel(params Params) (Channel, error) {
	tag, ok := params["tag"]
	if !ok {
		return Channel{}, errors.New("No tag")
	}
	endpoint := fmt.Sprintf(apiEndpoints["channels"]+"?tag=%s", tag)
	body, _, err := c.do("GET", endpoint, nil)
	if err != nil {
		log.Println(err)
		return Channel{}, err
	}

	var channel Channel
	if err = json.Unmarshal(body, &channel); err != nil {
		return Channel{}, err
	}
	return channel, nil
}

//DONE - need unit test
func (c *Client) Unsubscribe(params Params) (int, error) {
	id, ok := params["iden"]
	if !ok {
		return -1, errors.New("No id")
	}
	endpoint := fmt.Sprintf(apiEndpoints["subscriptions"]+"/%s", id)
	_, status, err := c.do("DELETE", endpoint, nil)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return status, nil
}

//DONE - Review case of active/non active contacts
func (c *Client) GetContacts() (Contacts, error) {
	body, _, err := c.do("GET", apiEndpoints["contacts"], nil)
	if err != nil {
		log.Println(err)
		return Contacts{}, err
	}

	var contacts Contacts
	if err = json.Unmarshal(body, &contacts); err != nil {
		return Contacts{}, err
	}
	return contacts, nil
}

//DONE
func (c *Client) CreateContact(params Params) (Contact, error) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Contact{}, err
	}
	body, _, err := c.do("POST", apiEndpoints["contacts"], bytes.NewBuffer(jsonParams))
	if err != nil {
		log.Println(err)
		return Contact{}, err
	}

	var contact Contact
	if err = json.Unmarshal(body, &contact); err != nil {
		return Contact{}, err
	}
	return contact, nil
}

//DONE
func (c *Client) UpdateContact(params Params) (Contact, error) {
	id, ok := params["iden"]
	if !ok {
		return Contact{}, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["contacts"]+"/%s", id)

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Contact{}, err
	}
	body, _, err := c.do("POST", endpoint, bytes.NewBuffer(jsonParams))
	if err != nil {
		log.Println(err)
		return Contact{}, err
	}

	var contact Contact
	if err = json.Unmarshal(body, &contact); err != nil {
		return Contact{}, err
	}
	return contact, nil
}

//DONE
func (c *Client) DeleteContact(params Params) (int, error) {
	id, ok := params["iden"]
	if !ok {
		return -1, errors.New("No id")
	}
	endpoint := fmt.Sprintf(apiEndpoints["contacts"]+"/%s", id)
	_, status, err := c.do("DELETE", endpoint, nil)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return status, nil
}

//DONE - need to review edge case
func (c *Client) GetDevices() (Devices, error) {
	body, _, err := c.do("GET", apiEndpoints["devices"], nil)
	if err != nil {
		log.Println(err)
		return Devices{}, err
	}

	var devices Devices
	if err = json.Unmarshal(body, &devices); err != nil {
		return Devices{}, err
	}
	return devices, nil
}

//DONE - need to review parameters that should pass or not
func (c *Client) CreateDevice(params Params) (Device, error) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Device{}, err
	}
	body, _, err := c.do("POST", apiEndpoints["devices"], bytes.NewBuffer(jsonParams))
	if err != nil {
		log.Println(err)
		return Device{}, err
	}

	var device Device
	if err = json.Unmarshal(body, &device); err != nil {
		return Device{}, err
	}
	return device, nil
}

//DONE
func (c *Client) UpdateDevice(params Params) (Device, error) {
	id, ok := params["iden"]
	if !ok {
		return Device{}, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["devices"]+"/%s", id)

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Device{}, err
	}
	body, _, err := c.do("POST", endpoint, bytes.NewBuffer(jsonParams))
	if err != nil {
		log.Println(err)
		return Device{}, err
	}

	var device Device
	if err = json.Unmarshal(body, &device); err != nil {
		return Device{}, err
	}
	return device, nil
}

//DONE
func (c *Client) DeleteDevice(params Params) (int, error) {
	id, ok := params["iden"]
	if !ok {
		return -1, errors.New("No id")
	}
	endpoint := fmt.Sprintf(apiEndpoints["devices"]+"/%s", id)
	_, status, err := c.do("DELETE", endpoint, nil)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return status, nil
}

//REVIEW
func (c *Client) GetPushes() (Pushes, error) {
	//TODO add params and allow modified_after
	body, _, err := c.do("GET", apiEndpoints["pushes"], nil)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return Push{}, err
	}

	var push Push
	if err = json.Unmarshal(body, &push); err != nil {
		return Push{}, err
	}
	return push, nil
}

func (c *Client) UpdatePush(params Params) (Push, error) {
	id, ok := params["iden"]
	if !ok {
		return Push{}, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["pushes"]+"/%s", id)

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return Push{}, err
	}
	body, _, err := c.do("POST", endpoint, bytes.NewBuffer(jsonParams))
	if err != nil {
		log.Println(err)
		return Push{}, err
	}

	var push Push
	if err = json.Unmarshal(body, &push); err != nil {
		return Push{}, err
	}
	return push, nil
}

//DONE
func (c *Client) DeletePush(params Params) (int, error) {
	id, ok := params["iden"]
	if !ok {
		return -1, errors.New("No id")
	}
	endpoint := fmt.Sprintf(apiEndpoints["pushes"]+"/%s", id)
	_, status, err := c.do("DELETE", endpoint, nil)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return status, nil
}

//DONE - need unit test
func (c *Client) UploadRequest(params Params) (UploadRequest, error) {
	if _, ok := params["file_name"]; !ok {
		return UploadRequest{}, errors.New("No file name")
	}
	if _, ok := params["file_type"]; !ok {
		return UploadRequest{}, errors.New("no file type")
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return UploadRequest{}, err
	}
	body, _, err := c.do("POST", apiEndpoints["upload_request"], bytes.NewBuffer(jsonParams))
	if err != nil {
		log.Println(err)
		return UploadRequest{}, err
	}

	var uploadRequest UploadRequest
	if err = json.Unmarshal(body, &uploadRequest); err != nil {
		return UploadRequest{}, err
	}
	return uploadRequest, nil
}

func (c *Client) Upload(path string) (int, error) {
	req, err := c.UploadRequest(Params{
		"file_name": "teste.txt",
		"file_type": "text/plain",
	})
	if err != nil {
		return -1, errors.New("Bad upload request")
	}
	file, err := os.Open(path)
	if err != nil {
		return -1, errors.New("no file")
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
		return -1, errors.New("Error creating form file")
	}
	if _, err = io.Copy(part, file); err != nil {
		return -1, errors.New("Error copying file")
	}
	err = writer.Close()
	if err != nil {
		return -1, errors.New("Error closing file")
	}
	upload_req, err := http.NewRequest("POST", req.UploadUrl, body)
	upload_req.Header.Set("Content-Type", writer.FormDataContentType())
	fmt.Println(upload_req)
	if err != nil {
		return -1, errors.New("Error creating POST REQUEST")
	}
	client := &http.Client{}
	resp, err := client.Do(upload_req)
	if err != nil {
		return -1, errors.New("Error doing request")
	}
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))

	return resp.StatusCode, nil
}
