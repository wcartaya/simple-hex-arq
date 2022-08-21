package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	js "users-ms/serializers/json"
	ms "users-ms/serializers/msgpack"
	"users-ms/domain"
)

type UserHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	userService domain.UserService
}

func NewHandler(userService domain.UserService) UserHandler {
	return &handler{userService: userService}
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) serializer(contentType string) domain.UserSerializer {
	if contentType == "application/x-msgpack" {
		return &ms.User{}
	}
	return &js.User{}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "user:id")
	contentType := r.Header.Get("Content-Type")
	user, err := h.userService.Find(userId)
	if err != nil {
		if errors.Cause(err) == domain.ErrUserNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	output, err := h.serializer(contentType).Encode(user)
	if err != nil{
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, user)
	//w.WriteHeader(200)
	//w.Write(output)
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	user, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = h.userService.Store(user)
	if err != nil {
		if errors.Cause(err) == domain.ErrUserInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseBody, err := h.serializer(contentType).Encode(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(w, contentType, responseBody, http.StatusCreated)
}