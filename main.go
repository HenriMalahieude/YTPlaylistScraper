package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func check(msg string, er error) {
	if er != nil {
		log.Fatalf(msg, er)
	}
}

func outputToolUsage() {
	fmt.Printf("Playlist Tool\n")
	fmt.Printf("./tool [link/id]\n")
	fmt.Printf(" -o listname\n		 output file name (defaults to\n")
}

//TODO: Modularize so that we can create graphical app, and download operator
//TODO: Have it create the TXT files in a folder
func main() {
	//----------------------------------------------Arguments
	if len(os.Args) <= 1 {
		outputToolUsage()
		return
	}

	var ( //Options
		//max        int64  = 10
		//autoMax    bool   = false
		id         string = ""
		outputFile string = "playlist"
	)

	//os.Args[0] == executable name anyways
	id = os.Args[1]

	//Parse the ID
	if strings.Index(id, "=") != -1 {
		id = strings.Split(id, "=")[1]
	}

	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "-o" {
			outputFile = os.Args[i+1]
			i++
		} else {
			fmt.Printf("Unrecognized flag: %v\n", os.Args[i])
			return
		}
	}

	outputFile = strings.Split(outputFile, ".")[0]

	//----------------------------------------------Main Code
	//Open the file first
	f, err4 := os.OpenFile(outputFile+".txt", os.O_TRUNC, os.ModeAppend)
	if err4 != nil { //Most likely the file doesn't exist
		f, err4 = os.Create(outputFile + ".txt")
		check("Error creating/opening output file\n%v", err4)
	}
	defer f.Close() //EZ

	_, err := f.Seek(0, 0)
	check("Error seeking in file\n%v", err)

	//Open networking
	client := &http.Client{
		Transport: &transport.APIKey{Key: api_key}, //Contained within info.go which is gitignored. Go make your own API Key lol
	}

	service, err := youtube.New(client)
	check("Error creating new Youtube client\n%v", err)

	//Get Playlist Info
	infoResponse, err := service.Playlists.List([]string{"snippet", "contentDetails"}).Id(id).Do()
	check("Error getting Info of Playlist("+id+")\n%v", err)

	f.WriteString("Playlist \"" + infoResponse.Items[0].Snippet.Title + "\"\nDescription:" + infoResponse.Items[0].Snippet.Description)

	//Get Playlist Items, Paging each of them
	playListItemsCall := service.PlaylistItems.List([]string{"snippet"}).PlaylistId(id).MaxResults(50)
	err = playListItemsCall.Pages(context.Background(), func(resp *youtube.PlaylistItemListResponse) error {
		for _, item := range resp.Items {
			f.WriteString("\n\n")
			f.WriteString(fmt.Sprint(item.Snippet.Position) + ". " + item.Snippet.Title + ": " + item.Snippet.VideoOwnerChannelTitle + "\n -Link: https://youtu.be/" + item.Snippet.ResourceId.VideoId)
		}

		return nil
	})
}
