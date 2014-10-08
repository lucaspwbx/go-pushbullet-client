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
	Iden            string `json:"iden"`
	Email           string `json:"email"`
	EmailNormalized string `json:"email_normalized"`
	Name            string `json:"name"`
	ImageUrl        string `json:"image_url"`
	Preferences     struct {
		Onboarding struct {
			App       bool
			Friends   bool
			Extension bool
		} `json:"onboarding"`
		Social bool `json:"social"`
	} `json:"preferences"`
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

type Note struct {
	title string
	body  string
	kind  string
}

type Link struct {
	title string
	body  string
	url   string
	kind  string
}

type Address struct {
	name    string
	address string
	kind    string
}

type List struct {
	title string
	items []string
	kind  string
}

type File struct {
	fileName string
	fileType string
	fileUrl  string
	body     string
	kind     string
}

type Client struct {
	token      string
	HttpClient *http.Client
}

type Params map[string]interface{}

type Endpoint map[string]string

func NewNote(title, body string) *Note {
	return &Note{title: title, body: body, kind: "note"}
}

func NewLink(title, body, url string) *Link {
	return &Link{title: title, body: body, url: url, kind: "link"}
}

func NewAddress(name, address string) *Address {
	return &Address{name: name, address: address, kind: "address"}
}

func NewList(title string, items ...string) *List {
	return &List{title: title, items: items, kind: "list"}
}

func NewFile(fname, ftype, furl, body string) *File {
	return &File{fileName: fname, fileType: ftype, fileUrl: furl, body: body, kind: "file"}
}

func NewClient(token string) *Client {
	httpClient := &http.Client{}
	return &Client{token: token, HttpClient: httpClient}
}
