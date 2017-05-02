package main

import (
	"log"

	"github.com/lukerodham/grpc-http-example/proto-go/users"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

func main() {

	service := micro.NewService(micro.Name("accounts"))

	userService := users.NewUsersClient("accounts", service.Client())
	rsp, err := userService.ListAll(context.TODO(), &users.Request{})

	log.Println(err)
	log.Println(rsp)
}
