package util

import (
	"crypto/tls"
	"net/smtp"
	"log"
	"fmt"
	config "../config"
)

var auth smtp.Auth
var to []string
var tconfig *tls.Config

//init  mail config
func init () {

}

func SendMail (msg []byte, DesignerEmail string) {
	
	auth = smtp.PlainAuth("", config.Gf.Username, config.Gf.Password, config.Gf.SmtpHost)
	to = config.Gf.CoptTo
	to = append(to, DesignerEmail)

	err := smtp.SendMail(config.Gf.SmtpHost+":"+config.Gf.Port, auth, config.Gf.From, to, msg)
	if err != nil {
		fmt.Println(err, to)
	}
}

func SendMailOneByOne (msg []byte, DesignerEmail string) {
	
	auth = smtp.PlainAuth("", config.Gf.Username, config.Gf.Password, config.Gf.SmtpHost)
	to = config.Gf.CoptTo
	to = append(to, DesignerEmail)

	//为了找出 Invalid RCPT TO address provided
	for _, tmpto := range to {
		var tmptoo []string
		tmptoo = append(tmptoo, tmpto)
		err := smtp.SendMail(config.Gf.SmtpHost+":"+config.Gf.Port, auth, config.Gf.From, tmptoo, msg)
		if err != nil {
			fmt.Println(err, tmptoo)
		} else {
			fmt.Println("Send to ", tmptoo, "success!!!")
		}
	}

}

func Test2() {
	msg := []byte("To: marui@cmcm.com\r\n" +
		"Subject: 邮件测试\r\n" +
		"\r\n" +
		"测试要什么body，啊？.\r\n")

	// Connect to the remote SMTP server.
	c, err := smtp.Dial("email-smtp.us-west-2.amazonaws.com:25")
	if err != nil {
		log.Fatal("Dial:", err)
	}

	err = c.StartTLS(tconfig)
	if err != nil {
		log.Fatal("StartTLS:", err)
	}

	err = c.Auth(auth)
	if err != nil {
		log.Fatal("Auth:", err)
	}

	// Set the sender and recipient first
	if err := c.Rcpt("wujunjian@cmcm.com"); err != nil {
		log.Fatal("Rcpt:", err)
	}

	if err := c.Mail("themedesign@cmcm.com"); err != nil {
		log.Fatal("Mail:", err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal("Data:", err)
	}
	_, err = fmt.Fprintf(wc, string(msg))
	if err != nil {
		log.Fatal("Send:", err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal("Close:", err)
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal("Quit:", err)
	}
}

