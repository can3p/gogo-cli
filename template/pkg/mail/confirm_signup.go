package mail

import (
	"fmt"
	"log"
	"os"

	"<projectrepo>/pkg/links"
	"<projectrepo>/pkg/model/core"
	"github.com/jmoiron/sqlx"
	mailjet "github.com/mailjet/mailjet-apiv3-go/v3"
	"github.com/pkg/errors"
)

func ConfirmSignup(db *sqlx.DB, user *core.User) error {
	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))

	if user.EmailConfirmSeed.String == "" {
		return errors.Errorf("cannot send confirm email for user with empty confirmation seed, user id = %s", user.ID)
	}

	link := links.AbsLink("confirm_signup", user.EmailConfirmSeed.String)
	to := user.Email

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "<projectemail>",
				Name:  "Your <projectname>",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
				},
			},
			Subject: "Welcome to <projectname>",
			TextPart: fmt.Sprintf(`
	Hi!

	Thank you for your interest in <projectname>! Please follow the link to confirm your email address

	%s`, link),
			HTMLPart: fmt.Sprintf(`
	<p>Hi!</p>

	<p>Thank you for your interest in <projectname>! Please follow the link to confirm your email address</p>

	<a href="%s">%s</a>`, link, link),
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
