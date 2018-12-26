package mongo_test

import (
	"log"
	"testing"

	root "github.com/425devon/go_rest_api/pkg"
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
	t.Run("GetUserById", getUserById_should_return_user_by_Id)
	t.Run("GetAllUsers", get_all_users_should_return_all_users)
	t.Run("UpdateUser", updateUser_should_update_user)
	t.Run("DeleteUserById", deleteUserById_should_remove_user)
}

func createUser_should_insert_user_into_mongo(t *testing.T) {
	//Arrange
	session := newSession()
	userService := newUserService(session)
	defer dropAndCloseDB(session)

	user := root.User{
		Username: "integration_test_user",
		Password: "integration_test_password",
	}

	//Act
	_, err := userService.CreateUser(&user)

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

func get_all_users_should_return_all_users(t *testing.T) {
	//Arrange
	session := newSession()
	userService := newUserService(session)
	defer dropAndCloseDB(session)

	users := []root.User{
		{
			Username: "integration_test_user",
			Password: "integration_test_password",
		},
		{
			Username: "user2",
			Password: "integration_test_password2",
		},
		{
			Username: "user3",
			Password: "integration_test_password3",
		},
	}

	//Act
	for _, user := range users {
		userService.CreateUser(&user)
	}
	retrievedUsers, err := userService.GetAllUsers()
	if err != nil {
		t.Errorf("Unable to retrieve all users: %s", err)
	}

	//Assert
	if len(retrievedUsers) != 3 {
		t.Errorf("Expecteed to retrieve `3` users. Got: `%d`", len(retrievedUsers))
	}
}

func getUserById_should_return_user_by_Id(t *testing.T) {
	//Arrange
	session := newSession()
	userService := newUserService(session)
	defer dropAndCloseDB(session)

	user := root.User{
		Username: "integration_test_user",
		Password: "integration_test_password",
	}

	//Act
	uid, err := userService.CreateUser(&user)
	recievedUser, err := userService.GetUserById(uid)

	//Assert
	if err != nil {
		t.Errorf("Error retrieving user by id: %s", err)
	}
	if recievedUser.Id != uid {
		t.Errorf("Expected user Id to match. Wanted: `%s` Got: `%s`", uid, recievedUser.Id)
	}
}

func updateUser_should_update_user(t *testing.T) {
	//Arrange
	session := newSession()
	userService := newUserService(session)
	defer dropAndCloseDB(session)

	user := root.User{
		Username: "Devon_Test",
		Password: "ChangeMe",
	}

	//Act
	uid, err := userService.CreateUser(&user)
	recievedUser, err := userService.GetUserById(uid)
	recievedUser.Password = "MuchBetter"
	err = userService.UpdateUser(recievedUser)
	updatedUser, err := userService.GetUserById(uid)
	if err != nil {
		t.Error(err)
	}

	//Assert
	if user.Password == updatedUser.Password {
		t.Error("Passwords shoud not match!")
	}
	if updatedUser.Password != "MuchBetter" {
		t.Errorf("Incorrect Password, expected: `MuchBetter` got: `%s`", updatedUser.Password)
	}

}

func deleteUserById_should_remove_user(t *testing.T) {
	//Arrange
	session := newSession()
	userService := newUserService(session)
	defer dropAndCloseDB(session)

	user := root.User{
		Username: "integration_test_user",
		Password: "integration_test_password",
	}

	//Act
	uid, _ := userService.CreateUser(&user)
	err := userService.DeleteUserById(uid)

	//Assert
	if err != nil {
		t.Errorf("Unanble to delete user: %s", err)
	}

}

func newSession() *mongo.Session {
	session, err := mongo.NewSession(mongoURL)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	return session
}

func newUserService(session *mongo.Session) *mongo.UserService {
	mockHash := mock.Hash{}
	userService := mongo.NewUserService(session.Copy(), dbName, userCollectionName, &mockHash)
	return userService
}

func dropAndCloseDB(session *mongo.Session) {
	session.DropDatabase(dbName)
	session.Close()
}
