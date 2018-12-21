package mongo_test

import (
	"log"
	"testing"

	"github.com/425devon/go_rest_api/pkg"
	"github.com/425devon/go_rest_api/pkg/mock"
	"github.com/425devon/go_rest_api/pkg/mongo"
)

const (
	mongoURL           = "localhost:27017"
	dbName             = "test_db"
	userCollectionName = "user"
)

func Test_UserService(t *testing.T) {
	t.Run("CreateUser", createUser_should_insert_user_into_mongo)
	t.Run("DeleteUserByUsername", DeleteUserByUsername_should_remove_user)
}

func createUser_should_insert_user_into_mongo(t *testing.T) {
	//Arrange
	session, err := mongo.NewSession(mongoURL)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()
	mockHash := mock.Hash{}
	userService := mongo.NewUserService(session.Copy(), dbName, userCollectionName, &mockHash)

	user := root.User{
		Username: "integration_test_user",
		Password: "integration_test_password",
	}

	//Act
	err = userService.CreateUser(&user)

	//Assert
	if err != nil {
		t.Errorf("Unable to create user: %s", err)
	}
	var results []root.User
	session.GetCollection(dbName, userCollectionName).Find(nil).All(&results)

	count := len(results)
	if count != 1 {
		t.Errorf("Incorrect number of results. Expected `1`, got: `%v`", count)
	}
	if results[0].Username != user.Username {
		t.Errorf("Incorrect username expected `%s`, got: `%s`", user.Username, results[0].Username)
	}
}

func DeleteUserByUsername_should_remove_user(t *testing.T) {
	//Arrange
	session, err := mongo.NewSession(mongoURL)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()
	mockHash := mock.Hash{}
	userService := mongo.NewUserService(session.Copy(), dbName, userCollectionName, &mockHash)
	user := root.User{
		Username: "integration_test_user",
		Password: "integration_test_password",
	}

	//Act
	userService.CreateUser(&user)
	err = userService.DeleteUserByUsername(user.Username)

	//Assert
	if err != nil {
		t.Errorf("Unanble to delete user: %s", err)
	}

}
