package sendgrid

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"path"
)

const (
	// SGSend path to send sendgrid email
	SGSend = "/v3/mail/send"
	// SGAll path to retrieve all mail
	SGAll = "/v3/mail/all"
	// SGHtmlByID path to retrieve html by id
	SGHtmlByID = "/v3/mail/html/id/"
	// SGTxtByID path to retrieve text by id
	SGTxtByID = "/v3/mail/text/id/"
	// SGUrlByID path to retrieve url by id
	SGUrlByID = "/v3/mail/text/link/id/"
)

// Send process sendgrid emails and returns the request details as well as any url is present to stdout
func Send(store Sendgrid) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
		writer.WriteHeader(200)

		requestDump, err := httputil.DumpRequest(request, false)
		if err != nil {
			log.Println(err)
		}

		log.Print(string(requestDump))

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}

		if len(body) > 1 {
			item, err := UnmarshalSGMail(body)
			if err != nil {
				log.Printf("Could not process email payload %v", err)
				return
			}

			fmt.Println(fmt.Sprintf("Recipients: %v \nSubject: %s \nURL: %s \nHTML: %s \nTXT: %s \nID: %s",
				item.Recipients,
				item.Subject,
				item.URL,
				fmt.Sprintf("http://%s%s%s", request.Host, SGHtmlByID, item.ID),
				fmt.Sprintf("http://%s%s%s", request.Host, SGTxtByID, item.ID),
				item.ID,
			))

			fmt.Println()

			ctx := context.Background()
			store.Set(ctx, item)
		} else {
			fmt.Println("Body empty")
			fmt.Println()
		}
	}
}

// All returns all mail in the db
func All(store Sendgrid) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.Background()
		items := MarshalSGMail(store.All(ctx))
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
		writer.WriteHeader(200)
		writer.Write(items)
	}
}

// HTMLByID returns the html email for a given id
func HTMLByID(store Sendgrid) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.Background()
		id := path.Base(request.RequestURI)
		item := store.GetByID(ctx, &id)
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
		writer.WriteHeader(200)
		writer.Write([]byte(item.HTML))
	}
}

// TXTByID returns the text email for a given id
func TXTByID(store Sendgrid) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.Background()
		id := path.Base(request.RequestURI)
		item := store.GetByID(ctx, &id)
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
		writer.WriteHeader(200)
		writer.Write([]byte(item.Text))
	}
}

// LinkByID returns the link in an email for a given id
func LinkByID(store Sendgrid) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.Background()
		id := path.Base(request.RequestURI)
		item := store.GetByID(ctx, &id)
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
		writer.WriteHeader(200)
		writer.Write([]byte(item.URL))
	}
}
