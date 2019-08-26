package sendgrid

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
	"mvdan.cc/xurls"
)

// UnmarshalSGMail transforms sendgrid email JSON to a workable struct, determining known and unknown data
func UnmarshalSGMail(input []byte) (*SGItem, error) {
	resMap := make(map[string]interface{})

	err := json.Unmarshal(input, &resMap)
	if err != nil {
		return nil, err
	}

	var res SGMailV3
	var md mapstructure.Metadata
	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			Metadata: &md,
			Result:   &res,
		})

	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(resMap); err != nil {
		return nil, err
	}

	var html, text, url string
	for _, i := range res.Content {
		switch i.Type {
		case "text/html":
			html = i.Value
		case "text/plain":
			text = i.Value
			url = GetLink(&i.Value)
		default:
			md.Unused = append(md.Unused, fmt.Sprintf("Content.%s", i.Type))
		}
	}

	var recipients []string
	for _, v := range res.Personalizations {
		for _, to := range v.To {
			recipients = append(recipients, to.Email)
		}
	}

	mail := SGItem{
		ID:         bson.NewObjectId().Hex(),
		Recipients: recipients,
		Subject:    res.Subject,
		Full:       res,
		HTML:       html,
		Text:       text,
		URL:        url,
		CreatedAt:  time.Now(),
		Unused:     md.Unused,
	}

	return &mail, nil
}

// MarshalSGMail returns the meta items from SGItems data
func MarshalSGMail(items []SGItem) []byte {
	res, err := json.Marshal(items)
	if err != nil {
		log.Println(err)
		return nil
	}

	return res
}

// GetLink retrieves a link from a text email
func GetLink(text *string) string {
	url := xurls.Strict.FindString(*text)
	link := html.UnescapeString(url)
	return link
}
