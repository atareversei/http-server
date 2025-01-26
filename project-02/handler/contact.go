package handler

import (
	"fmt"
	"github.com/atareversei/network-course-projects/pkg/http"
	"github.com/atareversei/network-course-projects/project-02/entity"
	"strings"
)

func (h Handler) Contact(req http.Request, res *http.Response) {
	contacts := make([]entity.Contact, 0)
	code, _ := req.Params("code")
	if code != "" {
		contact, err := h.repo.FindByCode(code)
		if err == nil {
			contacts = append(contacts, contact)
		}
	}
	if code == "" {
		name, _ := req.Params("name")
		if name != "" {
			ctcs := h.repo.ListByName(name)
			contacts = append(contacts, ctcs...)
		}
		phone, _ := req.Params("phone")
		if phone != "" {
			ctcs := h.repo.ListByPhone(phone)
			contacts = append(contacts, ctcs...)
		}
		address, _ := req.Params("address")
		if address != "" {
			ctcs := h.repo.ListByAddress(address)
			contacts = append(contacts, ctcs...)
		}
		email, _ := req.Params("email")
		if email != "" {
			ctcs := h.repo.ListByEmail(email)
			contacts = append(contacts, ctcs...)
		}
		if name == "" && phone == "" && address == "" && email == "" {
			contacts = h.repo.ListAllContacts()
		}
	}
	var builder strings.Builder
	for _, contact := range contacts {
		tmpl := fmt.Sprintf("<li>code: %s name: %s phone: %s address: %s email: %s</li>", contact.Code, contact.Name, contact.Phone, contact.Address, contact.Email)
		builder.WriteString(tmpl)
	}
	res.SetHeader("Content-Type", "text/html")

	if len(contacts) == 0 {
		res.Write([]byte("<b>no contact found!</b>"))
	} else {
		res.Write([]byte("<ul>" + builder.String() + "</ul>"))
	}

}
