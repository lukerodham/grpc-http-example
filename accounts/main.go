package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/lukerodham/grpc-http-example/proto-go/organisations"
	"github.com/lukerodham/grpc-http-example/proto-go/users"
	"github.com/micro/go-micro"
	"gopkg.in/gin-gonic/gin.v1"
)

var db *gorm.DB
var service micro.Service

func main() {
	// get configuration and test it.

	db = initDB(
		"mysql",
		"root@tcp(localhost:3306)/accounts?charset=utf8&parseTime=True&loc=Local",
		true)

	// setup gprc server
	service = micro.NewService(
		micro.Name("accounts"),
		micro.Version("latest"),
	)

	go setupGRPC(service)
	// setup http server
	setupHTTP()
}

func setupGRPC(s micro.Service) {

	service.Init()
	registerServices(service)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func setupHTTP() {
	router := gin.Default()
	registerHTTPServices(router)
	router.Run(":8000")
}

func registerServices(s micro.Service) {
	users.RegisterUsersHandler(s.Server(), new(UsersRPC))
}

func registerHTTPServices(e *gin.Engine) {
	new(UsersHTTP).RegisterHandlers(e, service)
}

func initDB(dbDialect string, dbDsn string, logMode bool) *gorm.DB {
	connection, err := gorm.Open(dbDialect, dbDsn)
	if err != nil {
		log.Fatal(err)
	}

	connection.LogMode(logMode)

	if err := connection.DB().Ping(); err != nil {
		log.Fatal(err)
	}

	connection.AutoMigrate(users.User{}, organisations.Organisation{})

	return connection
}
