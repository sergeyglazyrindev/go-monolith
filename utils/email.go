package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"math/big"
	"net/smtp"
	"reflect"
	"strings"
	"time"
)

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

type SentEmail struct {
	From    string
	To      []string
	Subject string
	CC      []string
	Message []byte
}

type SentEmailsDuringTestsType struct {
	SentEmails []SentEmail
}

func (se *SentEmailsDuringTestsType) AddSentEmail(from string, to []string, subject string, CC []string, message []byte) {
	se.SentEmails = append(se.SentEmails, SentEmail{
		From:    from,
		To:      to,
		Subject: subject,
		CC:      CC,
		Message: message,
	})
}

func (se *SentEmailsDuringTestsType) ClearTestEmails() {
	se.SentEmails = []SentEmail{}
}

func (se *SentEmailsDuringTestsType) IsAnyMatchedEmailSent(expectedEmail *SentEmail) bool {
	for i := range se.SentEmails {
		storedSentEmail := se.SentEmails[i]
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

func (se *SentEmailsDuringTestsType) IsAnyEmailSentWithStringInBodyOrSubject(expectedEmail *SentEmail) bool {
	for i := range se.SentEmails {
		storedSentEmail := se.SentEmails[i]
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

var SentEmailsDuringTests = SentEmailsDuringTestsType{
	SentEmails: make([]SentEmail, 0),
}

// @todo rework
// SendEmail sends email using system configured variables
func SendEmail(from string, to []string, cc []string, bcc []string, subject string, body string) error {
	if !core.CurrentConfig.InTests && (core.CurrentConfig.D.GoMonolith.EmailUsername == "" || core.CurrentConfig.D.GoMonolith.EmailPassword == "" || core.CurrentConfig.D.GoMonolith.EmailSMTPServer == "" || core.CurrentConfig.D.GoMonolith.EmailSMTPServerPort == 0) {
		errMsg := "Email not sent because email global variables are not set"
		core.Trail(core.CRITICAL, errMsg)
		return fmt.Errorf(errMsg)
	}

	// Get the domain name of sender
	domain := strings.Split(from, "@")
	if !core.CurrentConfig.InTests && len(domain) < 2 {
		return nil
	}
	domain[0] = strings.TrimSpace(domain[0])
	domain[0] = strings.TrimSuffix(domain[0], ">")

	// Construct the email
	MIME := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := "From: " + from + "\r\n"
	msg += "To: " + strings.Join(to, ",") + "\r\n"
	if len(cc) > 0 {
		msg += "CC: " + strings.Join(cc, ",") + "\r\n"
	}
	msg += "Date: " + time.Now().UTC().Format(time.RFC1123Z) + "\r\n"
	msg += "Message-ID: " + fmt.Sprintf("<%s-%s-%s-%s-%s@%s>", GenerateBase32(8), GenerateBase32(4), GenerateBase32(4), GenerateBase32(4), GenerateBase32(12), domain[0]) + "\r\n"
	msg += "Subject: " + subject + "\r\n"
	msg += MIME + "\r\n"
	msg += strings.Replace(body, "\n", "<br/>", -1)
	msg += "\r\n"
	// Append CC and BCC
	to = append(to, cc...)
	to = append(to, bcc...)

	if !core.CurrentConfig.InTests {
		go func() {
			err := smtp.SendMail(fmt.Sprintf("%s:%d", core.CurrentConfig.D.GoMonolith.EmailSMTPServer, core.CurrentConfig.D.GoMonolith.EmailSMTPServerPort),
				smtp.PlainAuth("", core.CurrentConfig.D.GoMonolith.EmailUsername, core.CurrentConfig.D.GoMonolith.EmailPassword, core.CurrentConfig.D.GoMonolith.EmailSMTPServer),
				core.CurrentConfig.D.GoMonolith.EmailFrom, to, []byte(msg))

			if err != nil {
				core.Trail(core.CRITICAL, "Email was not sent. %s", err)
			}
		}()
	} else {
		SentEmailsDuringTests.AddSentEmail(from, to, subject, cc, []byte(body))
	}
	return nil
}
