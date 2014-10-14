package main

import (
	"fmt"
	"os"

	"github.com/lucasweiblen/pushbulletclient/client"
)

func main() {
	///cli := client.NewClient("token")
	//resp, code, _ := cli.GetContacts()
	//fmt.Println(resp)
	//fmt.Println(code)
	//resp, code, _ = cli.CreateContact(client.Params{"name": "joao", "email": "foo@joao.com"})
	//fmt.Println(resp)
	//fmt.Println(code)
	//resp, code, _ = cli.UpdateContact(client.Params{"iden": "ujDigsFMxWesjAjjvzVOLY", "name": "blindpigs"})
	//fmt.Println(resp)
	//fmt.Println(code)
	//resp, code, _ = cli.DeleteContact(client.Params{"iden": "ujDigsFMxWesjz0F6E05sa"})
	//fmt.Println(resp)
	//fmt.Println(code)
	cli := client.NewClient("swbpcaTIjyV5eAYZnjfL2GZqFiiqrBHH")
	//contact, _ := cli.CreateContact(client.Params{"name": "joao", "email": "foo@joao.com"})
	//fmt.Println(contact)
	//updated, _ := cli.UpdateContact(client.Params{"iden": "ujvSxVpCjh6sjAiVsKnSTs", "name": "foo_updated"})
	//fmt.Println(updated)
	//device, _ := cli.CreateDevice(client.Params{"nickname": "teste", "type": "stream"})
	//fmt.Println(device)
	//updated, _ := cli.UpdateDevice(client.Params{"iden": "ujvSxVpCjh6sjz2gRnO9Aq", "nickname": "deviceupdated"})
	//bla, _ := cli.DeleteDevice(client.Params{"iden": "ujvSxVpCjh6sjAsoeMFET6"})
	//fmt.Println(updated)
	//pushes, _ := cli.GetPushes()
	//fmt.Println(pushes.Pushes)
	//b, _ := cli.CreatePush(client.Params{"type": "link", "title": "baz", "body": "foo", "url": "blaa"})
	//fmt.Println(b)
	//r, _ := cli.UpdatePush(client.Params{"iden": "ujvSxVpCjh6sjAcyrQ9Cmq", "title": "modified"})
	//fmt.Println(r)
	//o, _ := cli.CreatePush(client.Params{"type": "address", "name": "ok", "address": "bla"})
	//fmt.Println(o)
	//c, _ := cli.CreatePush(client.Params{"type": "list", "title": "titulo", "items": "bla"})
	//subs, _ := cli.Subscriptions()
	//for _, v := range subs.Subscriptions {
	//fmt.Println(v.Iden)
	//}
	//teste, _ := cli.Subscribe(client.Params{"channel_tag": "bla"})
	//fmt.Println(teste)
	//teste, _ := cli.GetChannel(client.Params{"tag": "bla"})
	//fmt.Println(teste)
	//teste, _ := cli.Unsubscribe(client.Params{"iden": "ujvSxVpCjh6sjAqnjXGNrw"})
	//fmt.Println(teste)
	//c, _ := cli.CreatePush(client.Params{"type": "list", "title": "titulo", "items": []string{"foo", "bar"}})
	//fmt.Println(c)
	//req, _ := cli.UploadRequest(client.Params{"file_name": "image.png", "file_type": "image/png"})
	//fmt.Println("FileURL: ", req.UploadUrl)
	//fmt.Println("Data access key: ", req.Data.AwsAccessKeyId)
	//fmt.Println("Acl: ", req.Data.Acl)
	//fmt.Println("Kye: ", req.Data.Key)
	//fmt.Println("Signature: ", req.Data.Signature)
	//fmt.Println("Policyt: ", req.Data.Policy)
	//fmt.Println("Content-Type: ", req.Data.ContentType)

	//up, err := cli.Upload()
	//if err != nil {
	//fmt.Println(err)
	//return
	//}
	//fmt.Println(up)
	fmt.Println(cli)
	//fmt.Println(os.Getwd())
	path, _ := os.Getwd()
	path += "/teste.txt"
	//fmt.Println(path)
	up, err := cli.Upload(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(up)
}
