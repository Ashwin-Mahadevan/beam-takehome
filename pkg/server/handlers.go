package server

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"slai.io/takehome/pkg/common"
)

func HandleEcho(msg []byte, client *Client) error {
	log.Println("Received ECHO request.")

	var request common.EchoRequest
	err := json.Unmarshal(msg, &request)

	if err != nil {
		log.Fatal("Invalid echo request.")
	}

	response := &common.EchoResponse{
		BaseResponse: common.BaseResponse{
			RequestId:   request.RequestId,
			RequestType: request.RequestType,
		},
		Value: request.Value,
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
	if err != nil {
		return err
	}

	return nil
}

func HandleSync(msg []byte, client *Client) error {
	log.Println("Received SYNC request.")

	var request common.SyncRequest
	err := json.Unmarshal(msg, &request)

	if err != nil {
		log.Fatal("Invalid sync request.")
	}

	log.Printf("Syncing file %s", request.Path)
	log.Printf("Content: %s", request.Content)

	os.WriteFile(fmt.Sprintf("./sync/%s", request.Path), []byte(request.Content), 0644)

	response := &common.SyncResponse{
		BaseResponse: common.BaseResponse{
			RequestId:   request.RequestId,
			RequestType: request.RequestType,
		},
		Success: true,
		Message: "Synced file.",
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
	if err != nil {
		return err
	}

	return nil
}
