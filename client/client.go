package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	apiEndpoints = Endpoint{"contacts": "https://api.pushbullet.com/v2/contacts",
		"pushes":  "https://api.pushbullet.com/v2/pushes",
		"devices": "https://api.pushbullet.com/v2/devices",
		"me":      "https://api.pushbullet.com/v2/users/me",
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
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return data, resp.StatusCode, nil
}

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

func (c *Client) GetMe() (User, error) {
	//TODO add params and allow modified_after
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

func (c *Client) CreatePush(params Params) (Push, error) {
	if _, ok := params["type"]; !ok {
		return Push{}, errors.New("no type")
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

func (c *Client) DeleteDevice(params Params) (int, error) {
	id, ok := params["iden"]
	if !ok {
		return -1, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["devices"]+"/%s", id)
	_, status, err := c.do("DELETE", endpoint, nil)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return status, nil
}

func (c *Client) DeletePush(params Params) (int, error) {
	id, ok := params["iden"]
	if !ok {
		return -1, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["pushes"]+"/%s", id)
	_, status, err := c.do("DELETE", endpoint, nil)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return status, nil
}

func (c *Client) DeleteContact(params Params) (int, error) {
	id, ok := params["iden"]
	if !ok {
		return -1, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["contacts"]+"/%s", id)
	_, status, err := c.do("DELETE", endpoint, nil)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return status, nil
}
