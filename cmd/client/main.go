package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
	client "slai.io/takehome/pkg/client"
)

func main() {
	log.Println("Starting client...")

	c, err := client.NewClient("./")
	if err != nil {
		log.Fatal(err)
	}

	value, err :=c.Echo("Hello, World!")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received: '%s'", value)

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	/* TODO: get from CLI*/
	watcher.Add("./input")

	for {
	    select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("Watcher closed")
				return
			}

			if event.Has(fsnotify.Write) {
				log.Printf("File %s was written to", event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Printf("Error: %s", err)
		}
    }

	/** Starter code */

	// someMessage := "hello there"
	// for {

	// 	log.Printf("Sending: '%s'", someMessage)

	// 	value, err := c.Echo(someMessage)
	// 	if err != nil {
	// 		log.Fatal("Unable to send request.")
	// 	}

	// 	log.Printf("Received: '%s'", value)

	// 	time.Sleep(time.Second)
	// }

}
