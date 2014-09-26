package main

import (
	"fmt"

	"github.com/lucasweiblen/pushbulletclient/client"
)

func main() {
	cli := client.NewClient("token")
	resp, code, _ := cli.GetContacts()
	fmt.Println(resp)
	fmt.Println(code)
	resp, code, _ = cli.CreateContact(client.Params{"name": "joao", "email": "foo@joao.com"})
	fmt.Println(resp)
	fmt.Println(code)
	resp, code, _ = cli.UpdateContact(client.Params{"iden": "ujDigsFMxWesjAjjvzVOLY", "name": "blindpigs"})
	fmt.Println(resp)
	fmt.Println(code)
	//resp, code, _ = cli.DeleteContact(client.Params{"iden": "ujDigsFMxWesjz0F6E05sa"})
	//fmt.Println(resp)
	//fmt.Println(code)
}
