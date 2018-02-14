package helpers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

//MailNewRegistration mails QR code to registered participants
func MailNewRegistration(email, aid string) error {

	if err := GenerateQR(aid); err != nil {
		return err
	}

	from := mail.NewEmail("Abacus 2018", "mail@abacus.org.in")
	subject := "Registration Confirmaton"
	to := mail.NewEmail("Participant", email)
	m := mail.NewV3MailInit(from, subject, to)
	m.Personalizations[0].SetSubstitution("-abacusid-", aid)
	m.SetTemplateID("1cd142f6-eabe-4960-95ce-cc123869d8fb")
	aqr := mail.NewAttachment()
	dat, err := ioutil.ReadFile(fmt.Sprintf("qr/%s.png", aid))
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(dat))
	aqr.SetContent(encoded)
	aqr.SetType("image/png")
	aqr.SetFilename(fmt.Sprintf("%s.png", aid))
	aqr.SetDisposition("attachment")
	aqr.SetContentID("QR code")
	m.AddAttachment(aqr)
	apikey := os.Getenv("MAIL_API_KEY")
	request := sendgrid.GetRequest(apikey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err = sendgrid.API(request)
	if err != nil {
		return err
	}
	//fmt.Println(response.StatusCode)
	// fmt.Println(response.Body)
	// fmt.Println(response.Headers)
	return nil
}

//MailRegistrationLink mails registration link to new participants
func MailRegistrationLink(email string, id int) error {
	from := mail.NewEmail("Abacus 2018", "mail@abacus.org.in")
	subject := "Registration Link"
	to := mail.NewEmail("Participant", email)
	m := mail.NewV3MailInit(from, subject, to)
	m.Personalizations[0].SetSubstitution("-profilelink-", strconv.Itoa(id))
	m.SetTemplateID("610c3fbf-de23-4883-83a7-077fdc3cdd28")
	apikey := os.Getenv("MAIL_API_KEY")
	request := sendgrid.GetRequest(apikey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	if err != nil {
		return err
	}
	//fmt.Println(response.StatusCode)
	// fmt.Println(response.Body)
	// fmt.Println(response.Headers)
	return nil
}
