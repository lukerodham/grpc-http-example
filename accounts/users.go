package main

import (
	"errors"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/google/jsonapi"
	micro "github.com/micro/go-micro"
	"github.com/sipsynergy/go-sipsynergy/utils"
	users "github.com/sipsynergy/proto-go/users"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gin-gonic/gin.v1"
)

/********************************************
 * gRPC Logic
 ********************************************/

// UsersRPC struct handles container the gRPC functions neatly.
type UsersRPC struct{}

// ListAll is a gRPC function to list all users.
func (u *UsersRPC) ListAll(ctx context.Context, req *users.Request, rsp *users.ListResponse) error {

	foundUsers, err := new(Users).FindAll()

	if err != nil {
		return err
	}

	new := make([]*users.User, len(foundUsers))
	for i, v := range foundUsers {
		new[i] = &v
	}

	rsp.Users = new

	return nil
}

// Show is the gRPC function to fetch one user.
func (u *UsersRPC) Show(ctx context.Context, req *users.ShowRequest, rsp *users.SingleResponse) error {

	foundUser, err := new(Users).FindOne(req.GetID())

	if err != nil {
		return err
	}

	rsp.User = foundUser

	return nil
}

/********************************************
 * HTTP Logic
 ********************************************/

// UsersHTTP is the struct we to set HTTP routing logic to.
type UsersHTTP struct {
	Service micro.Service
}

// RegisterHandlers will register routes for HTTP access.
func (u *UsersHTTP) RegisterHandlers(e *gin.Engine, service micro.Service) {
	u.Service = service

	e.GET("/users", u.ListAll)
	e.POST("/users", u.Create)
}

// ListAll here is the http method for listing all users.
func (u *UsersHTTP) ListAll(c *gin.Context) {
	rsp, err := u.getGRPCService().ListAll(c, &users.Request{})

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if len(rsp.Users) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Writer.WriteHeader(200)
	c.Header("Content-Type", "application/json")

	if err := jsonapi.MarshalManyPayload(c.Writer, rsp.Users); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	return
}

// Show is the HTTP method to handle fetching one user.
func (u *UsersHTTP) Show(c *gin.Context) {
	rsp, err := u.getGRPCService().Show(c, &users.ShowRequest{ID: c.Param("id")})

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Writer.WriteHeader(200)
	c.Header("Content-Type", "application/json")

	if err := jsonapi.MarshalOnePayload(c.Writer, rsp.User); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	return
}

// Create handles the creation of the entity.
func (u *UsersHTTP) Create(c *gin.Context) {
	user := new(users.User)

	if err := jsonapi.UnmarshalPayload(c.Request.Body, user); err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	new(Users).Create(user)

	c.Writer.WriteHeader(200)
	c.Header("Content-Type", "application/json")

	if err := jsonapi.MarshalOnePayload(c.Writer, user); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	return
}

func (u *UsersHTTP) getGRPCService() users.UsersClient {
	return users.NewUsersClient("accounts", u.Service.Client())
}

/********************************************
 * Storage Logic
 ********************************************/

// Users is the struct we use to set logic functions to.
type Users struct{}

// FindAll returns a list of users from storage.
func (u *Users) FindAll() ([]users.User, error) {
	var (
		users []users.User
		err   error
	)

	r := db.Preload("Organisation").Find(&users)

	if r.RecordNotFound() {
		err = errors.New("no users found")
	}

	return users, err
}

// FindOne returns a single user from storage.
func (u *Users) FindOne(ID string) (*users.User, error) {
	var (
		user users.User
		err  error
	)

	r := db.Preload("Organisation").First(&user, "id = ?", ID)

	if r.RecordNotFound() {
		err = errors.New("user not found")
	}

	return &user, err
}

// Create handles create the storage row for the user.
func (u *Users) Create(user *users.User) *users.User {
	user = u.beforeCreate(user)

	db.Create(user)

	return user
}

// beforeCreate handles setting the custom human id.
func (u *Users) beforeCreate(user *users.User) *users.User {
	identifier := utils.GenerateHumanID("USR")

	if db.Select("1").First(new(users.User), "id = ?", identifier).RecordNotFound() == false {
		return u.beforeCreate(user)
	}

	// reset password.
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 13)
	if err != nil {
		panic(err)
	}

	newPassword := string(passwordHash)
	log.Println(newPassword)
	user.Password = newPassword
	user.ID = identifier

	return user
}
