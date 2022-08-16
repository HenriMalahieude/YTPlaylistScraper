package main

import (
	"log"
	"net/http"
	"os"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func main() {
	client := &http.Client{
		Transport: &transport.APIKey{Key: api_key}, //Contained within info.go which is ignored. Go make your own API Key lol
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new Youtube client: %v", err)
	}

	part := []string{"snippet"}
	call := service.PlaylistItems.List(part).PlaylistId(os.Args[1]).MaxResults(10)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error calling List Playlist Items: %v", err)
	}

	for _, item := range response.Items {
		println(item.Snippet.Title + ": " + item.Snippet.VideoOwnerChannelTitle + ".\n-Link: https://youtu.be/" + item.Snippet.ResourceId.VideoId + "\n")

	}
}
