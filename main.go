package main

import (
	"fmt"
	"log"

	"github.com/lucasweiblen/pushbulletclient/client"
)

func main() {
	//200
	cli := client.NewClient("swbpcaTIjyV5eAYZnjfL2GZqFiiqrBHH")
	user, err := cli.GetMe()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(user)

	//400
	cli2 := client.NewClient("foo")
	_, err = cli2.GetMe()
	if err != nil {
		fmt.Println(err)
	}

	// returns error -> no channel tag parameter
	_, err = cli.Subscribe(client.Params{})
	if err != nil {
		fmt.Println(err)
	}

	//subscribe to channel tag
	subscription, err := cli.Subscribe(client.Params{"channel_tag": "jblow"})
	if err != nil {
		//	log.Fatalln(err)
		fmt.Println(err)
	}
	fmt.Println(subscription)

	//getting subscriptions
	subs, err := cli.Subscriptions()
	if err != nil {
		//	log.Fatalln(err)
		fmt.Println(err)
	}
	fmt.Println(subs)

	ch, err := cli.GetChannel(client.Params{"tag": "jblow"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)

	//removing subscription
	//err = cli.Unsubscribe2(client.Params{"iden": "ujvSxVpCjh6sjAgWOzmngO"})
	//if err != nil {
	//fmt.Println(err)
	//return
	//}
	//fmt.Println("subscription removed with sucess")

	fmt.Println("----------")
	contacts, err := cli.GetContacts()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(contacts)

	//contact, err := cli.CreateContact(client.Params{"name": "joao", "email": "joao@foo.com"})
	//if err != nil {
	//fmt.Println(err)
	//return
	//}
	//fmt.Println(contact)

	//err = cli.DeleteContact(client.Params{"iden": "ujvSxVpCjh6sjArHrh8WLA"})
	//if err != nil {
	//fmt.Println(err)
	//return
	//}
	//fmt.Println("Removed")

	//_, err = cli.CreateDevice(client.Params{"nickname": "foobar"})
	//if err != nil {
	//fmt.Println(err)
	//}
	devices, _ := cli.GetDevices()
	fmt.Println(devices)

	newReq, err := cli.UploadRequest(client.Params{"file_name": "teste", "file_type": "text"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(newReq)
}
