package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/joho/godotenv"
)

type NotificationResponse struct {
	Email    string `json:"email"`
	FilmName string `json:"filmName"`
	UserName string `json:"userName"`
}

func GetClient() *azservicebus.Client {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	} else {
		client, err := azservicebus.NewClientFromConnectionString(os.Getenv("SERVICEBUS_CONNECTION_STRING"), nil)
		if err != nil {
			panic(err)
		}
		return client
	}
	return nil
}

func GetMessage(count int, client *azservicebus.Client) error {
	receiver, err := client.NewReceiverForQueue("ticketwatchernotificationqueue", nil) //Change myqueue to env var
	if err != nil {
		return errors.New("error creating receiver")
	}
	defer receiver.Close(context.TODO())
	for {
		messages, err := receiver.ReceiveMessages(context.TODO(), count, nil)
		if err != nil {
			return errors.New("error receiving messages")
		}

		for _, message := range messages {
			body := message.Body
			fmt.Printf("%s\n", string(body))

			var notificationInfo NotificationResponse
			err = json.Unmarshal(body, &notificationInfo)
			if err != nil {
				return errors.New("error unmarshalling json")
			}

			subject := "The film " + notificationInfo.FilmName + " is now available!"
			emailBody := "Hello " + notificationInfo.UserName + ",\n\nThe film " + notificationInfo.FilmName + " is now available!\n\nRegards,\nTicketWatcher"
			err := SendEmailResolver(notificationInfo.Email, subject, emailBody)
			if err != nil {
				return errors.New("error sending email")
			}

			err = receiver.CompleteMessage(context.TODO(), message, nil)
			if err != nil {
				return errors.New("error completing message")
			}
		}
	}

}

func DeadLetterMessage(client *azservicebus.Client) {
	deadLetterOptions := &azservicebus.DeadLetterOptions{
		ErrorDescription: to.Ptr("exampleErrorDescription"),
		Reason:           to.Ptr("exampleReason"),
	}

	receiver, err := client.NewReceiverForQueue(os.Getenv("AZURE_SERVICEBUS_QUEUE_NAME"), nil)
	if err != nil {
		panic(err)
	}
	defer receiver.Close(context.TODO())

	messages, err := receiver.ReceiveMessages(context.TODO(), 1, nil)
	if err != nil {
		panic(err)
	}

	if len(messages) == 1 {
		err := receiver.DeadLetterMessage(context.TODO(), messages[0], deadLetterOptions)
		if err != nil {
			panic(err)
		}
	}
}

func GetDeadLetterMessage(client *azservicebus.Client) {
	receiver, err := client.NewReceiverForQueue(
		os.Getenv("AZURE_SERVICEBUS_QUEUE_NAME"),
		&azservicebus.ReceiverOptions{
			SubQueue: azservicebus.SubQueueDeadLetter,
		},
	)
	if err != nil {
		panic(err)
	}
	defer receiver.Close(context.TODO())

	messages, err := receiver.ReceiveMessages(context.TODO(), 1, nil)
	if err != nil {
		panic(err)
	}

	for _, message := range messages {
		fmt.Printf("DeadLetter Reason: %s\nDeadLetter Description: %s\n", *message.DeadLetterReason, *message.DeadLetterErrorDescription) //change to struct an unmarshal into it
		err := receiver.CompleteMessage(context.TODO(), message, nil)
		if err != nil {
			panic(err)
		}
	}
}
