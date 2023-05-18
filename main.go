package main

import (
	"NotificationServer/services"
)

func main() {
	client := services.GetClient()
	services.GetMessage(10, client)
	services.GetDeadLetterMessage(client)

	// to := "leebang31698@gmail.com"
	// subject := "Hello from ticketWatcher"
	// body := "This is the body of the email."

	// err := services.SendEmailResolver(to, subject, body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
