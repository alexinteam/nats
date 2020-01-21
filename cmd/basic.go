package cmd

import (
	"flag"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"runtime"
)

var (
	publishQueue   string
	subscribeQueue string
	messageParam   string
)

var pubCmd = &cobra.Command{
	Use:   "pub",
	Short: "Start pub nats-pub",
	Long:  "Start publish nats-pub",
	Run: func(cmd *cobra.Command, args []string) {
		var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
		var authUser = flag.String("u", "nats", "The nats server authentication user for clients")
		var authPassword = flag.String("p", "", "The nats server authentication password for clients")

		log.SetFlags(0)
		flag.Parse()
		args = flag.Args()
		nc, err := nats.Connect(*urls, nats.UserInfo(*authUser, *authPassword))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connected to NATS server: " + *urls)
		if len(args) < 2 {
			usage_prod()
		}
		msg := []byte(messageParam)
		nc.Publish(publishQueue, msg)
		nc.Flush()
		if err := nc.LastError(); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Published [%s] : '%s'\n", publishQueue, msg)
		}
	},
}

var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "Start pub nats-pub",
	Long:  "Start subscribe nats-pub",
	Run: func(cmd *cobra.Command, args []string) {
		var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
		var authUser = flag.String("u", "nats", "The nats server authentication user for clients")
		var authPassword = flag.String("p", "", "The nats server authentication password for clients")
		//var command = flag.String("c", "", "Whether to produce or consume a message")
		log.SetFlags(0)
		flag.Parse()
		args = flag.Args()
		nc, err := nats.Connect(*urls, nats.UserInfo(*authUser, *authPassword))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connected to NATS server: " + *urls)
		if len(args) < 1 {
			usage_con()
		}
		nc.Subscribe(subscribeQueue, func(msg *nats.Msg) {
			log.Printf("Received message '%s\n", string(msg.Data)+"'")
		})
		nc.Flush()
		if err := nc.LastError(); err != nil {
			log.Fatal(err)
		}
		log.Printf("Listening on [%s]\n", subscribeQueue)
		runtime.Goexit()
	},
}

func init() {
	runtime.GOMAXPROCS(1)
	rootCmd.AddCommand(subCmd)
	rootCmd.AddCommand(pubCmd)

	pubCmd.PersistentFlags().StringVar(&publishQueue, "publishQueue", "", "publishQueue=test")
	_ = viper.BindPFlag("publishQueue", subCmd.PersistentFlags().Lookup("publishQueue"))

	pubCmd.PersistentFlags().StringVar(&messageParam, "message", "", "message=testMessage")
	_ = viper.BindPFlag("message", subCmd.PersistentFlags().Lookup("message"))

	subCmd.PersistentFlags().StringVar(&subscribeQueue, "subscribeQueue", "", "subscribeQueue=test")
	_ = viper.BindPFlag("subscribeQueue", subCmd.PersistentFlags().Lookup("subscribeQueue"))
}

func usage_prod() {
	log.Fatalf("Usage: nats-pub [-s server (%s)] [-u user (%s)] [-p password (%s)] -c produce <subject> <msg> \n", nats.DefaultURL, "nats", "S3Cr3TP@5w0rD")
}

func usage_con() {
	log.Fatalf("Usage: nats-pub [-s server (%s)] [-u user (%s)] [-p password (%s)] -c consume <subject> \n", nats.DefaultURL, "nats", "S3Cr3TP@5w0rD")
}
