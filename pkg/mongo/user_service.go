package mongo

import (
	root "github.com/425devon/go_rest_api/pkg"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserService struct {
	Collection *mgo.Collection
	hash       root.Hash
}

func NewUserService(session *Session, dbName string, collectionName string, hash root.Hash) *UserService {
	collection := session.GetCollection(dbName, collectionName)
	collection.EnsureIndex(userModelIndex())
	return &UserService{collection, hash}
}

func (p *UserService) CreateUser(u *root.User) error {
	user := newUserModel(u)
	hashedPassword, err := p.hash.Generate(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return p.Collection.Insert(&user)
}

func (p *UserService) GetAllUsers() ([]*root.User, error) {
	users := []userModel{}
	var rootUsers []*root.User
	err := p.Collection.Find(bson.M{}).All(&users)
	for _, user := range users {
		rootUsers = append(rootUsers, user.toRootUser())
	}
	return rootUsers, err
}

func (p *UserService) GetByUsername(username string) (*root.User, error) {
	model := userModel{}
	err := p.Collection.Find(bson.M{"username": username}).One(&model)
	return model.toRootUser(), err
}

func (p *UserService) DeleteUserByUsername(username string) error {
	err := p.Collection.Remove(bson.M{"username": username})
	return err
}

// func (p *UserService) GetAllUsers() ([]userModel, error) {
// 	users := []userModel{}
// 	// rootUsers := []*root.User{}
// 	err := p.Collection.Find(bson.M{}).All(&users)
// 	// for i, user := range users {
// 	// 	rootUsers[i] = user.toRootUser()
// 	// }
// 	return users, err
// }
