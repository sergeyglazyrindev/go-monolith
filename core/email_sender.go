package core

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"reflect"
	"strings"
	"time"
)

type IEmailSender interface {
	SetFrom(from string)
	SetSubject(subject string)
	SetBody(body string)
	AddRecipient(to string)
	AddCC(cc string)
	AddBCC(bcc string)
	SendEmail() error
}

type SentEmail struct {
	From    string
	To      []string
	Subject string
	CC      []string
	Message []byte
}

//
//var SentEmailsDuringTests = SentEmailsDuringTestsType{
//	SentEmails: make([]SentEmail, 0),
//}

// GenerateBase64 generates a base64 string of length length
func GenerateBase64(length int) string {
	base := new(big.Int)
	base.SetString("64", 10)

	base64 := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	tempKey := ""
	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, base)
		tempKey += string(base64[int(index.Int64())])
	}
	return tempKey
}

// GenerateBase32 generates a base64 string of length length
func GenerateBase32(length int) string {
	base := new(big.Int)
	base.SetString("32", 10)

	base32 := "234567abcdefghijklmnopqrstuvwxyz"
	tempKey := ""
	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, base)
		tempKey += string(base32[int(index.Int64())])
	}
	return tempKey
}

type SMTPEmailSender struct {
	from string
	to []string
	cc []string
	bcc []string
	subject string
	body string
}

func (ses *SMTPEmailSender) SetFrom(from string) {
	ses.from = from
}

func (ses *SMTPEmailSender) SetSubject(subject string) {
	ses.subject = subject
}

func (ses *SMTPEmailSender) SetBody(body string) {
	ses.body = body
}

func (ses *SMTPEmailSender) AddRecipient(to string) {
	if ses.to == nil {
		ses.to = make([]string, 0)
	}
	ses.to = append(ses.to, to)
}

func (ses *SMTPEmailSender) AddCC(cc string) {
	if ses.cc == nil {
		ses.cc = make([]string, 0)
	}
	ses.cc = append(ses.cc, cc)
}

func (ses *SMTPEmailSender) AddBCC(bcc string) {
	if ses.bcc == nil {
		ses.bcc = make([]string, 0)
	}
	ses.bcc = append(ses.bcc, bcc)
}

func (ses *SMTPEmailSender) SendEmail() error{
	// Get the domain name of sender
	domain := strings.Split(ses.from, "@")
	domain[0] = strings.TrimSpace(domain[0])
	domain[0] = strings.TrimSuffix(domain[0], ">")

	// Construct the email
	MIME := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := "From: " + ses.from + "\r\n"
	msg += "To: " + strings.Join(ses.to, ",") + "\r\n"
	if ses.cc != nil && len(ses.cc) > 0 {
		msg += "CC: " + strings.Join(ses.cc, ",") + "\r\n"
	}
	msg += "Date: " + time.Now().UTC().Format(time.RFC1123Z) + "\r\n"
	msg += "Message-ID: " + fmt.Sprintf("<%s-%s-%s-%s-%s@%s>", GenerateBase32(8), GenerateBase32(4), GenerateBase32(4), GenerateBase32(4), GenerateBase32(12), domain[0]) + "\r\n"
	msg += "Subject: " + ses.subject + "\r\n"
	msg += MIME + "\r\n"
	msg += strings.Replace(ses.body, "\n", "<br/>", -1)
	msg += "\r\n"
	// Append CC and BCC
	to := make([]string, 0)
	to = append(to, ses.to...)
	to = append(to, ses.cc...)
	to = append(to, ses.bcc...)

	err := smtp.SendMail(fmt.Sprintf("%s:%d", CurrentConfig.D.GoMonolith.EmailSMTPServer, CurrentConfig.D.GoMonolith.EmailSMTPServerPort),
		smtp.PlainAuth("", CurrentConfig.D.GoMonolith.EmailUsername, CurrentConfig.D.GoMonolith.EmailPassword, CurrentConfig.D.GoMonolith.EmailSMTPServer),
		CurrentConfig.D.GoMonolith.EmailFrom, to, []byte(msg))

	if err != nil {
		Trail(CRITICAL, "Email was not sent. %s", err)
		return err
	}
	return nil
}

type EmailSenderForTests struct {
	from string
	to []string
	cc []string
	bcc []string
	subject string
	body string
	SentEmails []SentEmail
}

func (ses *EmailSenderForTests) SetFrom(from string) {
	ses.from = from
}

func (ses *EmailSenderForTests) SetSubject(subject string) {
	ses.subject = subject
}

func (ses *EmailSenderForTests) SetBody(body string) {
	ses.body = body
}

func (ses *EmailSenderForTests) AddRecipient(to string) {
	if ses.to == nil {
		ses.to = make([]string, 0)
	}
	ses.to = append(ses.to, to)
}

func (ses *EmailSenderForTests) AddCC(cc string) {
	if ses.cc == nil {
		ses.cc = make([]string, 0)
	}
	ses.cc = append(ses.cc, cc)
}

func (ses *EmailSenderForTests) AddBCC(bcc string) {
	if ses.bcc == nil {
		ses.bcc = make([]string, 0)
	}
	ses.bcc = append(ses.bcc, bcc)
}

func (ses *EmailSenderForTests) SendEmail() error{
	if ses.SentEmails == nil {
		ses.SentEmails = make([]SentEmail, 0)
	}
	ses.SentEmails = append(ses.SentEmails, SentEmail{
		From:    ses.from,
		To:      ses.to,
		Subject: ses.subject,
		CC:      ses.cc,
		Message: []byte(ses.body),
	})
	return nil
}

func (ses *EmailSenderForTests) ClearTestEmails() {
	ses.SentEmails = make([]SentEmail, 0)
}

func (ses *EmailSenderForTests) IsAnyMatchedEmailSent(expectedEmail *SentEmail) bool {
	for i := range ses.SentEmails {
		storedSentEmail := ses.SentEmails[i]
		match := false
		if len(expectedEmail.Subject) > 0 {
			match = expectedEmail.Subject == storedSentEmail.Subject
		}
		if len(expectedEmail.To) > 0 {
			match = reflect.DeepEqual(expectedEmail.To, storedSentEmail.To)
		}
		if len(expectedEmail.From) > 0 {
			match = expectedEmail.From == storedSentEmail.From
		}
		if len(expectedEmail.CC) > 0 {
			match = reflect.DeepEqual(expectedEmail.CC, storedSentEmail.CC)
		}
		if match {
			return true
		}
	}
	return false
}

func (ses *EmailSenderForTests) IsAnyEmailSentWithStringInBodyOrSubject(expectedEmail *SentEmail) bool {
	for i := range ses.SentEmails {
		storedSentEmail := ses.SentEmails[i]
		match := false
		if len(expectedEmail.Subject) > 0 {
			match = strings.Contains(storedSentEmail.Subject, expectedEmail.Subject)
		}
		if len(expectedEmail.Message) > 0 {
			match = bytes.Contains(storedSentEmail.Message, expectedEmail.Message)
		}
		if match {
			return true
		}
	}
	return false
}

type EmailSenderFactory struct {
	MakeEmailSender func() IEmailSender
}

var ProjectEmailSenderFactory *EmailSenderFactory
var TestEmailSender *EmailSenderForTests
func init() {
	ProjectEmailSenderFactory = &EmailSenderFactory{}
	ProjectEmailSenderFactory.MakeEmailSender = func() IEmailSender {
		return &SMTPEmailSender{to: make([]string, 0), cc: make([]string, 0), bcc: make([]string, 0)}
	}
	TestEmailSender = &EmailSenderForTests{}
}
