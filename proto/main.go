package main

import (
	"log"

	"github.com/micro/go-micro"
	"github.com/sipsynergy/proto-go/users"
	"golang.org/x/net/context"
)

func main() {

	service := micro.NewService(micro.Name("accounts"))

	userService := users.NewUsersClient("accounts", service.Client())
	rsp, err := userService.ListAll(context.TODO(), &users.Request{})

	log.Println(err)
	log.Println(rsp)
}
