package admin

import (
	"fmt"
	"log"
	"os"

	"<projectrepo>/pkg/model/core"
	mailjet "github.com/mailjet/mailjet-apiv3-go/v3"
)

var NotifyAddress string = "dpetroff@gmail.com"

func NotifyNewUser(user *core.User) {
	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "<projectemail>",
				Name:  "Your <projectname>",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: NotifyAddress,
				},
			},
			Subject: "New User on <projectname>",
			TextPart: fmt.Sprintf(`
	Hi!

	New user alert:

	* ID: %s
	* Email: %s`, user.ID, user.Email),
			HTMLPart: fmt.Sprintf(`
	<p>Hi!</p>

	<p>New user alert:</p>

	<ul>
		<li>ID: %s</li>
		<li>Email: %s</li>
	</ul>`, user.ID, user.Email),
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)

	if err != nil {
		log.Fatal(err)
	}
}

func NotifySignupConfirmed(user *core.User) {
	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "<projectemail>",
				Name:  "Your <projectname>",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: NotifyAddress,
				},
			},
			Subject: "New User confirmed email on <projectname>",
			TextPart: fmt.Sprintf(`
	Hi!

	New conrirmed user alert:

	* ID: %s
	* Email: %s`, user.ID, user.Email),
			HTMLPart: fmt.Sprintf(`
	<p>Hi!</p>

	<p>New conrirmed user alert:</p>

	<ul>
		<li>ID: %s</li>
		<li>Email: %s</li>
	</ul>`, user.ID, user.Email),
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)

	fmt.Println(res)

	if err != nil {
		log.Fatal(err)
	}
}

func NotifyThrowAwayEmailSignupAttempt(email string) {
	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "<projectemail>",
				Name:  "Your <projectname>",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: NotifyAddress,
				},
			},
			Subject: "An attempt to use a throwaway email domain on <projectname>",
			TextPart: fmt.Sprintf(`
	Hi!

	A user has just tried to use a throwaway email: %s`, email),
			HTMLPart: fmt.Sprintf(`
	<p>Hi!</p>

	<p>A user has just tried to use a throwaway email: %s</p>
	`, email),
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)

	if err != nil {
		log.Fatal(err)
	}
}
