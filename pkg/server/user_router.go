package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	root "github.com/425devon/go_rest_api/pkg"

	"github.com/gorilla/mux"
)

type userRouter struct {
	userService root.UserService
}

func NewUserRouter(u root.UserService, router *mux.Router) *mux.Router {
	userRouter := userRouter{u}

	router.HandleFunc("/", userRouter.getAllUsersHandler).Methods("GET")
	router.HandleFunc("/", userRouter.createUserHandler).Methods("PUT")
	router.HandleFunc("/{_id}", userRouter.updateUserHandler).Methods("PUT")
	router.HandleFunc("/{_id}", userRouter.getUserHandler).Methods("GET")
	router.HandleFunc("/{_id}", userRouter.deleteUserHandler).Methods("DELETE")
	return router
}

func (ur *userRouter) createUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user, err := decodeUser(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	_, err = ur.userService.CreateUser(&user)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Json(w, http.StatusCreated, err)
}

func (ur *userRouter) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	users, err := ur.userService.GetAllUsers()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
	}
	Json(w, http.StatusOK, users)
}

func (ur *userRouter) getUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	log.Println(vars)
	id := vars["_id"]

	user, err := ur.userService.GetUserById(id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	Json(w, http.StatusOK, user)
}

func (ur *userRouter) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	//defer r.Body.Close()
	vars := mux.Vars(r)
	log.Println(vars)
	user, err := decodeUser(r)
	if err != nil {
		Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = ur.userService.UpdateUser(&user)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	Json(w, http.StatusOK, map[string]string{"result": "success"})
}

func (ur *userRouter) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	log.Println(vars)
	id := vars["_id"]

	err := ur.userService.DeleteUserById(id)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
	}
	Json(w, http.StatusOK, nil)
}

func decodeUser(r *http.Request) (root.User, error) {
	var u root.User
	if r.Body == nil {
		return u, errors.New("No Request Body")
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	return u, err
}
