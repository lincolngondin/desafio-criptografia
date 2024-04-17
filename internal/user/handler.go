package user

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type userService interface {
	CreateUser(userDocument string, creditCardToken string, value int64) error
	GetUserByID(id int64) (*User, error)
	DeleteUserById(id int64) error
	UpdateUserById(id int64, user *User) (*User, error)
}

type handler struct {
	svc userService
}

func NewHandler(sv userService) *handler {
	return &handler{
		svc: sv,
	}
}

type RequestData struct {
	UserDocument    string `json:"user_document"`
	CreditCardToken string `json:"credit_card_token"`
	Value           int64  `json:"value"`
}

func (reqdata *RequestData) isValid() bool {
	return reqdata.UserDocument != "" && reqdata.CreditCardToken != ""
}

type ResponseError struct {
	Error string `json:"erro"`
}

func NewResponseErr(e string) *ResponseError {
	return &ResponseError{e}
}

func (hnd *handler) CreateUserHandler(response http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(response)
	decoder := json.NewDecoder(request.Body)

	postData := &RequestData{}
	decoderErr := decoder.Decode(postData)

	if decoderErr != nil {
		response.WriteHeader(http.StatusBadRequest)
		encoder.Encode(NewResponseErr(decoderErr.Error()))
		return
	}

	if !postData.isValid() {
		response.WriteHeader(http.StatusBadRequest)
		encoder.Encode(NewResponseErr("Os campos user_document ou credit_card_token não podem estar vazios!"))
		return
	}

	err := hnd.svc.CreateUser(postData.UserDocument, postData.CreditCardToken, postData.Value)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		encoder.Encode(NewResponseErr(err.Error()))
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func getId(idStr string) (int64, error) {
	id, idConvErr := strconv.ParseInt(idStr, 10, 64)
	if idConvErr != nil {
		return 0, idConvErr
	}
	return id, nil
}

func (hnd *handler) GetUserHandler(response http.ResponseWriter, request *http.Request) {
	id, getIdErr := getId(request.PathValue("id"))
	if getIdErr != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	user, err := hnd.svc.GetUserByID(id)
	if err == ErrUserNotFound {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(response)
	encoder.Encode(user)
}

func (hnd *handler) DeleteUserHandler(response http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(response)
	id, idErr := getId(request.PathValue("id"))
	if idErr != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	err := hnd.svc.DeleteUserById(id)
	if err == ErrUserNotFound {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(err.Error())
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (hnd *handler) UpdateUserHandler(response http.ResponseWriter, request *http.Request) {
	id, getIdErr := getId(request.PathValue("id"))
	if getIdErr != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	encoder := json.NewEncoder(response)
	decoder := json.NewDecoder(request.Body)
	newUser := &RequestData{}
	decoderErr := decoder.Decode(newUser)
	if decoderErr != nil {
		response.WriteHeader(http.StatusBadRequest)
		encoder.Encode(ResponseError{"Entre o valor de todos os campos!"})
		return
	}

	if !newUser.isValid() {
		response.WriteHeader(http.StatusBadRequest)
		encoder.Encode(NewResponseErr("Os campos user_document ou credit_card_token não podem estar vazios!"))
		return
	}

	user := NewUser(newUser.UserDocument, newUser.CreditCardToken, newUser.Value)
	user.Id = id

	_, err := hnd.svc.UpdateUserById(id, user)
	if err == ErrUserNotFound {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ResponseError{err.Error()})
		return
	}

	response.WriteHeader(http.StatusOK)
	encoder.Encode(user)
}
