package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"runtime"
	"time"

	stan "github.com/nats-io/stan.go"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
)

var streamPubCmd = &cobra.Command{
	Use:   "streamPub",
	Short: "Start pub nats",
	Long:  "Start publish nats",
	Run: func(cmd *cobra.Command, args []string) {
		publisher, err := nats.NewStreamingPublisher(
			nats.StreamingPublisherConfig{
				ClusterID: "test-cluster",
				ClientID:  "example-publisher",
				StanOptions: []stan.Option{
					stan.NatsURL("nats://127.0.0.1:4222"),
				},
				Marshaler: nats.GobMarshaler{},
			},
			watermill.NewStdLogger(false, false),
		)
		if err != nil {
			panic(err)
		}

		publishMessages(publisher)
	},
}

var streamSubCmd = &cobra.Command{
	Use:   "streamSub",
	Short: "Start streamSub nats",
	Long:  "Start subscribe nats",
	Run: func(cmd *cobra.Command, args []string) {
		subscriber, err := nats.NewStreamingSubscriber(
			nats.StreamingSubscriberConfig{
				ClusterID:        "test-cluster",
				ClientID:         "example-subscriber",
				QueueGroup:       "example",
				DurableName:      "my-durable",
				SubscribersCount: 4, // how many goroutines should consume messages
				CloseTimeout:     time.Minute,
				AckWaitTimeout:   time.Second * 30,
				StanOptions: []stan.Option{
					stan.NatsURL("nats://127.0.0.1:4222"),
				},
				Unmarshaler: nats.GobMarshaler{},
			},
			watermill.NewStdLogger(false, false),
		)
		if err != nil {
			panic(err)
		}

		messages, err := subscriber.Subscribe(context.Background(), "group1.topic")
		if err != nil {
			panic(err)
		}

		go process(messages)
		runtime.Goexit()
	},
}

func init() {
	runtime.GOMAXPROCS(1)
	rootCmd.AddCommand(streamPubCmd)
	rootCmd.AddCommand(streamSubCmd)

	//pubCmd.PersistentFlags().StringVar(&publishQueue, "publishQueue", "", "publishQueue=test")
	//_ = viper.BindPFlag("publishQueue", subCmd.PersistentFlags().Lookup("publishQueue"))
	//
	//pubCmd.PersistentFlags().StringVar(&message, "message", "", "message=testMessage")
	//_ = viper.BindPFlag("message", subCmd.PersistentFlags().Lookup("message"))
	//
	//subCmd.PersistentFlags().StringVar(&subscribeQueue, "subscribeQueue", "", "subscribeQueue=test")
	//_ = viper.BindPFlag("subscribeQueue", subCmd.PersistentFlags().Lookup("subscribeQueue"))
}

func publishMessages(publisher message.Publisher) {
	for {
		msg := message.NewMessage(watermill.NewUUID(), []byte("Hello, world!"))

		if err := publisher.Publish("group1.topic", msg); err != nil {
			panic(err)
		}
		fmt.Println("Message", "Hello, world!", "Published")

		time.Sleep(time.Second * 5)
	}
}

func process(messages <-chan *message.Message) {
	for msg := range messages {
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))
		//fmt.Println("received message: %s, payload: %s", msg.UUID, string(msg.Payload))

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
