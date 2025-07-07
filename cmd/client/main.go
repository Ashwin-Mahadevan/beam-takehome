package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
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

	entries, err := os.ReadDir("./input")

	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		content, err := os.ReadFile(fmt.Sprintf("./input/%s", entry.Name()))

		if err != nil {
			log.Fatal(err)
		}

		base64Content := base64.StdEncoding.EncodeToString(content)

		success, message, err := c.Sync(entry.Name(), base64Content)

		if err != nil {
			log.Fatal(err)
		}

		if !success {
			log.Fatal(message)
		}

		log.Printf("Synced file %s: %s", entry.Name(), message)
	}


	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	watcher.Add("./input")

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("Watcher closed")
					return
				}

				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					log.Printf("File %s was written to", event.Name)

					content, err := os.ReadFile(event.Name)

					if err != nil {
						log.Printf("Error reading file %s: %s", event.Name, err)
						continue
					}

					relativePath := strings.TrimPrefix(event.Name, "input/")
					base64Content := base64.StdEncoding.EncodeToString(content)

					success, message, err := c.Sync(relativePath, base64Content)

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
