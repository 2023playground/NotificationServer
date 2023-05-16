package main

import (
	"NotificationServer/services"
	"log"
)

func main() {
	to := "leebang31698@gmail.com"
	subject := "Hello from ticketWatcher"
	body := "This is the body of the email."

	err := services.SendEmailResolver(to, subject, body)
	if err != nil {
		log.Fatal(err)
	}
}
