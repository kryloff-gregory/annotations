package user

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"main/model"
	"os"
)

type user struct {
	UserID         string `json:"userId"`
	UserName       string `json:"userName"`
	HashedPassword string `json:"hashedPassword"`
}

type users struct {
	Users []*user `json:"users"`
}

type Provider struct {
	userMap map[string]*model.User
}

func NewProvider(usersFileName string) (*Provider, error) {
	userMap, err := readConfig(usersFileName)
	if err != nil {
		return nil, err
	}
	return &Provider{
		userMap: userMap,
	}, nil
}

func (u *Provider) GerUserByName(username string) *model.User {
	return u.userMap[username]
}

func readConfig(usersFileName string) (map[string]*model.User, error) {
	jsonFile, err := os.Open(usersFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var users *users

	if err := json.Unmarshal(byteValue, &users); err != nil {
		return nil, err
	}
	result := make(map[string]*model.User)
	for _, usr := range users.Users {
		id, err := uuid.Parse(usr.UserID)
		if err != nil {
			return nil, err
		}
		result[usr.UserName] = &model.User{
			Name:           usr.UserName,
			ID:             id,
			HashedPassword: usr.HashedPassword,
		}
	}

	return result, err
}
