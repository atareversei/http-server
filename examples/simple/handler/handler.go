package handler

import "github.com/atareversei/http-server/examples/simple/entity"

type Repository interface {
	FindByCode(code string) (entity.Contact, error)
	ListByName(name string) []entity.Contact
	ListByPhone(phone string) []entity.Contact
	ListByAddress(address string) []entity.Contact
	ListByEmail(email string) []entity.Contact
	ListAllContacts() []entity.Contact
}

type Handler struct {
	repo Repository
}

func New(repo Repository) Handler {
	return Handler{repo: repo}
}