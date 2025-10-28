package repository

import (
	"errors"
	"os"
	"strings"

	"github.com/atareversei/http-server/examples/simple/entity"
	"github.com/atareversei/http-server/internal/cli"
)

var RecordNotFound = errors.New("record not found")

type Repo struct {
	contacts []entity.Contact
}

func New() Repo {
	contacts := make([]entity.Contact, 0)
	loadData(&contacts)
	return Repo{contacts: contacts}
}

func loadData(contacts *[]entity.Contact) {
	data, err := os.ReadFile("./repository/data.txt")
	if err != nil {
		cli.Error("server could not find the data file", err)
		os.Exit(1)
	}
	for _, record := range strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n") {
		contactInfo := strings.Split(record, ",")
		var contact entity.Contact
		contact.Code = contactInfo[0]
		contact.Name = contactInfo[1]
		contact.Phone = contactInfo[2]
		contact.Address = contactInfo[3]
		contact.Email = contactInfo[4]
		*contacts = append(*contacts, contact)
	}
}

func (r *Repo) FindByCode(code string) (entity.Contact, error) {
	for _, contact := range r.contacts {
		if contact.Code == code {
			return contact, nil
		}
	}
	return entity.Contact{}, RecordNotFound
}

func (r *Repo) ListByName(name string) []entity.Contact {
	ctcs := make([]entity.Contact, 0)
	for _, contact := range r.contacts {
		if strings.Contains(strings.ToLower(contact.Name), name) {
			ctcs = append(ctcs, contact)
		}
	}
	return ctcs
}

func (r *Repo) ListByPhone(phone string) []entity.Contact {
	ctcs := make([]entity.Contact, 0)
	for _, contact := range r.contacts {
		if strings.Contains(contact.Phone, phone) {
			ctcs = append(ctcs, contact)
		}
	}
	return ctcs
}

func (r *Repo) ListByAddress(address string) []entity.Contact {
	ctcs := make([]entity.Contact, 0)
	for _, contact := range r.contacts {
		if strings.Contains(strings.ToLower(contact.Address), address) {
			ctcs = append(ctcs, contact)
		}
	}
	return ctcs
}

func (r *Repo) ListByEmail(email string) []entity.Contact {
	ctcs := make([]entity.Contact, 0)
	for _, contact := range r.contacts {
		if strings.Contains(strings.ToLower(contact.Email), email) {
			ctcs = append(ctcs, contact)
		}
	}
	return ctcs
}

func (r *Repo) ListAllContacts() []entity.Contact {
	return r.contacts
}
