package main

import (
	"flag"
	"fmt"
	"github.com/nats-io/go-nats"
	"log"
	"runtime"
)

func usage_prod() {
	log.Fatalf("Usage: nats-pub [-s server (%s)] [-u user (%s)] [-p password (%s)] -c produce <subject> <msg> \n", nats.DefaultURL, "nats", "S3Cr3TP@5w0rD")
}

func usage_con() {
	log.Fatalf("Usage: nats-pub [-s server (%s)] [-u user (%s)] [-p password (%s)] -c consume <subject> \n", nats.DefaultURL, "nats", "S3Cr3TP@5w0rD")
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var authUser = flag.String("u", "nats", "The nats server authentication user for clients")
	var authPassword = flag.String("p", "", "The nats server authentication password for clients")
	var command = flag.String("c", "", "Whether to produce or consume a message")
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()
	if *command == "" {
		log.Fatalf("Error: Indicate the command using '-command' flag")
	}
	if *command != "produce" && *command != "consume" {
		log.Fatalf("Error: Supported commands are: consume & produce")
	}
	nc, err := nats.Connect(*urls, nats.UserInfo(*authUser, *authPassword))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to NATS server: " + *urls)
	if *command == "produce" {
		if len(args) < 2 {
			usage_prod()
		}
		subj, msg := args[0], []byte(args[1])
		nc.Publish(subj, msg)
		nc.Flush()
		if err := nc.LastError(); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Published [%s] : '%s'\n", subj, msg)
		}
	}
	if *command == "consume" {
		if len(args) < 1 {
			usage_con()
		}
		subj := args[0]
		nc.Subscribe(subj, func(msg *nats.Msg) {
			log.Printf("Received message '%s\n", string(msg.Data)+"'")
		})
		nc.Flush()
		if err := nc.LastError(); err != nil {
			log.Fatal(err)
		}
		log.Printf("Listening on [%s]\n", subj)
		runtime.Goexit()
	}
}
