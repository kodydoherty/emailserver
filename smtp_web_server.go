package main

import (
	"net/http"
	"net/smtp"

	"text/template"
	"bytes"
	"fmt"
)


const emailAddress = "xxxx@gmail.com"
const password = "******"
const emailServer = "smtp.gmail.com"
const serverPort = 587


const subject = "Someone has sent you a message" 
var auth = smtp.PlainAuth(
	"",
	emailAddress,
	password,
	emailServer,
)

var emailTemplate *template.Template

const emailBody = ` 
From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}


Someone has sent you a message using your contact form:

Name: {{.yourname}}
Email: {{.email}}
Website: {{.website}}

Message:
{{.message}}

`

func setUpTemplate() {
	emailTemplate = template.New("emailTemplate")
	template.Must( emailTemplate.Parse(emailBody) )
}

func main() {

	setUpTemplate()
	
	http.Handle( "/", http.FileServer( http.Dir("site") ) )
	http.HandleFunc("/send", sendEmailHandler )
	http.ListenAndServe(":8000", nil)
}




func sendEmailHandler(w http.ResponseWriter, req *http.Request) {


	req.ParseForm()
	params := req.Form

	data := map[string]string{}
	data["yourname"] = params.Get("yourname")
	data["email"] = params.Get("email")
	data["website"] = params.Get("website")
	data["message"] = params.Get("message")

	data["Subject"] = subject
	data["To"] = emailAddress
	data["From"] = emailAddress

	var body bytes.Buffer
	emailTemplate.Execute(&body, data) 

    
    err := smtp.SendMail(
            fmt.Sprintf("%s:%d", emailServer, serverPort),
            auth,
            emailAddress,
            []string{emailAddress},
            body.Bytes(),
    )
    if err != nil {
            println("Error sending email")
    }
    http.Redirect(w, req, "/", http.StatusFound)

}