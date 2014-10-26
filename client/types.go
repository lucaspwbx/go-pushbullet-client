package client

import "net/http"

type Subscription struct {
	Iden    string  `json:"iden"`
	Active  bool    `json:"active"`
	Channel Channel `json:"channel"`
}

type Subscriptions struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type Channel struct {
	Iden        string `json:"iden"`
	Tag         string `json:"tag"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
}

type Device struct {
	Iden         string `json:"iden"`
	PushToken    string `json:"push_token"`
	AppVersion   int    `json:"app_version"`
	FingerPrint  string `json:"fingerprint"`
	Active       bool   `json:"active"`
	Nickname     string `json:"nickname"`
	Manufacturer string `json:"manufacturer"`
	Type         string `json:"type"`
	Model        string `json:"model"`
	Pushable     bool   `json:"pushable"`
}

type Devices struct {
	Devices []Device `json:"devices"`
}

type Contact struct {
	Iden            string `json:"iden"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	EmailNormalized string `json:"email_normalized"`
	Active          bool   `json:"active"`
}

type Contacts struct {
	Contacts []Contact `json:"contacts"`
}

type Push struct {
	Iden                    string `json:"iden"`
	Type                    string `json:"type"`
	Title                   string `json:"title"`
	Body                    string `json:"body"`
	Url                     string `json:"url"`
	Active                  bool   `json:"active"`
	Dismissed               bool   `json:"dismissed"`
	SenderIden              string `json:"sender_iden"`
	SenderEmail             string `json:"sender_email"`
	SenderEmailNormalized   string `json:"sender_email_normalized"`
	ReceiverIden            string `json:"receiver_iden"`
	ReceiverEmail           string `json:"receiver_email"`
	ReceiverEmailNormalized string `json:"receiver_email_normalized"`
}

type Pushes struct {
	Pushes []Push `json:"pushes"`
}

type User struct {
	Iden            string      `json:"iden"`
	Email           string      `json:"email"`
	EmailNormalized string      `json:"email_normalized"`
	Name            string      `json:"name"`
	ImageUrl        string      `json:"image_url"`
	Preferences     Preferences `json:"preferences"`
}

type Preferences struct {
	Onboarding struct {
		App       bool `json:"app,omitempty"`
		Friends   bool `json:"friends,omitempty"`
		Extension bool `json:"extension,omitempty"`
	} `json:"onboarding"`
	Social bool `json:"social,omitempty"`
}

type UploadRequest struct {
	FileType  string `json:"file_type"`
	FileName  string `json:"file_name"`
	FileUrl   string `json:"file_url"`
	UploadUrl string `json:"upload_url"`
	Data      struct {
		AwsAccessKeyId string `json:"awsaccesskeyid"`
		Acl            string `json:"acl"`
		Key            string `json:"key"`
		Signature      string `json:"signature"`
		Policy         string `json:"policy"`
		ContentType    string `json:"content-type"`
	} `json:"data"`
}

type Client struct {
	token      string
	HttpClient *http.Client
}

type Params map[string]interface{}

type Endpoint map[string]string

func NewClient(token string) *Client {
	httpClient := &http.Client{}
	return &Client{token: token, HttpClient: httpClient}
}
