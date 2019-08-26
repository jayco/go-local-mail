package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/jayco/go-local-email/internal/store"

	"github.com/jayco/go-local-email/internal/sendgrid"
)

// rootHandler lets you know that we are up
func rootHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
	writer.WriteHeader(200)
	requestDump, err := httputil.DumpRequest(request, false)

	if err != nil {
		log.Println(err)
	}

	log.Println(string(requestDump))
	writer.Write([]byte("go-local-mail loves you.\n"))
}

// Serve starts a echo email server
func Serve(port *string, db *store.Client) {
	log.Printf("Starting server, listening on port %s \n", *port)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc(sendgrid.SGSend, sendgrid.Send(sendgrid.NewSendGrid(db)))
	http.HandleFunc(sendgrid.SGAll, sendgrid.All(sendgrid.NewSendGrid(db)))
	http.HandleFunc(sendgrid.SGHtmlByID, sendgrid.HTMLByID(sendgrid.NewSendGrid(db)))
	http.HandleFunc(sendgrid.SGTxtByID, sendgrid.TXTByID(sendgrid.NewSendGrid(db)))
	http.HandleFunc(sendgrid.SGUrlByID, sendgrid.LinkByID(sendgrid.NewSendGrid(db)))
	http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
}
