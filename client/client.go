package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func (c *Client) GetDevices() (map[string]interface{}, int, error) {
	req, err := http.NewRequest("GET", apiEndpoints["devices"], nil)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}

func (c *Client) GetContacts() (map[string]interface{}, int, error) {
	req, err := http.NewRequest("GET", apiEndpoints["contacts"], nil)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}

func (c *Client) GetPushes() (map[string]interface{}, int, error) {
	//TODO add params and allow modified_after
	req, err := http.NewRequest("GET", apiEndpoints["pushes"], nil)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}

func (c *Client) CreatePush(params Params) (map[string]interface{}, int, error) {
	if _, ok := params["type"]; !ok {
		return nil, -1, err
	}
	jsonified_params, err := json.Marshal(params)
	if err != nil {
		return nil, -1, err
	}
	req, err := http.NewRequest("POST", apiEndpoints["pushes"], bytes.NewBuffer(jsonified_params))
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}

func (c *Client) CreateDevice(params Params) (map[string]interface{}, int, error) {
	jsonified_params, err := json.Marshal(params)
	if err != nil {
		return nil, -1, err
	}
	req, err := http.NewRequest("POST", apiEndpoints["devices"], bytes.NewBuffer(jsonified_params))
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}

func (c *Client) CreateContact(params Params) (map[string]interface{}, int, error) {
	jsonified_params, err := json.Marshal(params)
	if err != nil {
		return nil, -1, err
	}
	req, err := http.NewRequest("POST", apiEndpoints["contacts"], bytes.NewBuffer(jsonified_params))
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}

func (c *Client) UpdateDevice(params Params) (map[string]interface{}, int, error) {
	id, ok := params["iden"]
	if !ok {
		return nil, -1, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["devices"]+"/%s", id)
	jsonified_params, err := json.Marshal(params)
	if err != nil {
		return nil, -1, err
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonified_params))
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}

func (c *Client) UpdatePush(params Params) (map[string]interface{}, int, error) {
	id, ok := params["iden"]
	if !ok {
		return nil, -1, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["pushes"]+"/%s", id)
	jsonified_params, err := json.Marshal(params)
	if err != nil {
		return nil, -1, err
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonified_params))
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}

func (c *Client) UpdateContact(params Params) (map[string]interface{}, int, error) {
	id, ok := params["iden"]
	if !ok {
		return nil, -1, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["contacts"]+"/%s", id)
	jsonified_params, err := json.Marshal(params)
	if err != nil {
		return nil, -1, err
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonified_params))
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}

func (c *Client) DeleteDevice(params Params) (map[string]interface{}, int, error) {
	id, ok := params["iden"]
	if !ok {
		return nil, -1, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["devices"]+"/%s", id)
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}

func (c *Client) DeletePush(params Params) (map[string]interface{}, int, error) {
	id, ok := params["iden"]
	if !ok {
		return nil, -1, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["pushes"]+"/%s", id)
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}

func (c *Client) DeleteContact(params Params) (map[string]interface{}, int, error) {
	id, ok := params["iden"]
	if !ok {
		return nil, -1, errors.New("No id")
	}
	delete(params, "iden")
	endpoint := fmt.Sprintf(apiEndpoints["contacts"]+"/%s", id)
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	req.SetBasicAuth(c.token, "")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, -1, err
	}
	return result, resp.StatusCode, nil
}
