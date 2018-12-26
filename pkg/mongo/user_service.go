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

func (p *UserService) CreateUser(u *root.User) (_id string, err error) {
	user := newUserModel(u)
	user.Id = bson.NewObjectId()
	var uid = user.Id.Hex()
	hashedPassword, err := p.hash.Generate(user.Password)
	if err != nil {
		return uid, err
	}
	user.Password = hashedPassword
	return uid, p.Collection.Insert(&user)
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

func (p *UserService) GetUserById(id string) (*root.User, error) {
	model := userModel{}
	err := p.Collection.FindId(bson.ObjectIdHex(id)).One(&model)
	return model.toRootUser(), err
}

func (p *UserService) UpdateUser(u *root.User) error {
	err := p.Collection.UpdateId(bson.ObjectIdHex(u.Id), &u)
	return err
}

func (p *UserService) DeleteUserById(id string) error {
	oid := bson.ObjectIdHex(id)
	err := p.Collection.RemoveId(oid)
	return err
}
