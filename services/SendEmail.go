package services

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmailResolver(to, subject, body string) error {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	} else {
		from := mail.NewEmail("ticketWatcher", "playgroundevelope@gmail.com")
		subject := subject
		to := mail.NewEmail("User", to)
		htmlContent := fmt.Sprintf("<div>%s</div>", body)
		message := mail.NewSingleEmail(from, subject, to, body, htmlContent)
		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		response, err := client.Send(message)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
			fmt.Println("Email sent successfully!")
		}
	}

	return nil

}
