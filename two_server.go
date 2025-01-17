package main

import (
	"flag"
	"github.com/atareversei/network-course-projects/pkg/cli"
	"github.com/atareversei/network-course-projects/pkg/http"
	"github.com/atareversei/network-course-projects/pkg/http/server"
	"os"
	"strings"
)

type Contact struct {
	code    string
	name    string
	phone   string
	address string
	email   string
}

func loadData(contacts *[]Contact) {
	data, err := os.ReadFile("./two_data.txt")
	if err != nil {
		cli.Error("server could not find the data file", err)
		os.Exit(1)
	}

	for _, record := range strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n") {
		contactInfo := strings.Split(record, ",")
		var contact Contact
		contact.code = contactInfo[0]
		contact.name = contactInfo[1]
		contact.phone = contactInfo[2]
		contact.address = contactInfo[3]
		contact.email = contactInfo[4]

		*contacts = append(*contacts, contact)
	}
}

func main() {
	portFlag := flag.Int("port", 8080, "Port to serve")
	flag.Parse()
	port := *portFlag

	contacts := make([]Contact, 0)
	loadData(&contacts)

	s := server.New(port)
	s.Get("/contacts", contactHandler)
	s.All("/", notFoundHandler)
	s.Start()
}

func contactHandler(req http.Request, res *http.Response) {
	res.SetHeader("Content-Type", "text/html")
	res.Write([]byte("<h1>Hello, world!</h1>"))
}

func notFoundHandler(req http.Request, res *http.Response) {

}
