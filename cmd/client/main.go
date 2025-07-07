package main

import (
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	client "slai.io/takehome/pkg/client"
)

func main() {
	log.Println("Starting client...")

	c, err := client.NewClient("./")
	if err != nil {
		log.Fatal(err)
	}

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	/* TODO: get from CLI*/
	watcher.Add("./input")

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("Watcher closed")
					return
				}

				if event.Has(fsnotify.Write) {
					log.Printf("File %s was written to", event.Name)

					content, err := os.ReadFile(event.Name)

					if err != nil {
						log.Printf("Error reading file %s: %s", event.Name, err)
						continue
					}

					success, message, err := c.Sync(event.Name, string(content))

					if err != nil {
						log.Printf("Error syncing file %s: %s", event.Name, err)
						continue
					}

					if !success {
						log.Printf("Error syncing file %s: %s", event.Name, message)
						continue
					}

					log.Printf("Synced file %s: %s", event.Name, message)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Printf("Error: %s", err)
			}
		}
	}()

	/** Starter code */

	someMessage := "hello there"
	for {

		log.Printf("Sending: '%s'", someMessage)

		value, err := c.Echo(someMessage)
		if err != nil {
			log.Fatal("Unable to send request.")
		}

		log.Printf("Received: '%s'", value)

		time.Sleep(time.Minute)
	}

}
