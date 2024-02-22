package service

import (
	"log"
	"reflect"

	"github.com/dp3why/mrgo/backend"
	"github.com/dp3why/mrgo/constants"
	"github.com/dp3why/mrgo/model"
	"github.com/olivere/elastic/v7"
)

func CheckUser(username string, password string) (bool, error) {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewMatchQuery("username", username))
	query.Must(elastic.NewMatchQuery("password", password))
	searchResult, err := backend.ESBackend.ReadFromES(
		query,
		constants.USER_INDEX,
	)
	if err != nil {
		return false, err
	}
	var utype model.User
	for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
		u := item.(model.User)
		if u.Password == password {
			log.Default().Printf("Login as %s\n", username)
			return true, nil
		}
	}
	return false, nil
}

func AddUser(user *model.User) (bool, error) {
	query := elastic.NewTermQuery("username", user.Username)
	searchResult, err := backend.ESBackend.ReadFromES(
		query,
		constants.USER_INDEX,
	)
	if err != nil {
		return false, err
	}

	if searchResult.TotalHits() > 0 {
		return false, nil
	}

	err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
	if err != nil {
		return false, err
	}
	log.Default().Printf("User is added: %s\n", user.Username)
	return true, nil
}