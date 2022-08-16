package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	fmt.Printf("./tool [link/id]\n -s # or -s all\n 		 size of the playlist, say all for all items (default of 10)\n")
	fmt.Printf(" -o listname\n		 output file name (defaults to\n")
}

//TODO: Modularize so that we can create graphical app, and download operator
//TODO: Have it create the TXT files in a folder
func main() {
	if len(os.Args) <= 1 {
		outputToolUsage()
		return
	}

	var ( //Options
		max        int64  = 10
		autoMax    bool   = false
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
		if os.Args[i] == "-s" {
			if os.Args[i+1] != "all" {
				max2, err := strconv.Atoi(os.Args[i+1])
				check("Error converting to integer("+os.Args[i+1]+")\n%v", err)
				max = int64(max2)
			} else {
				autoMax = true
			}
			i++
		} else if os.Args[i] == "-o" {
			outputFile = os.Args[i+1]
			i++
		} else {
			fmt.Printf("Unrecognized flag: %v\n", os.Args[i])
			return
		}
	}

	client := &http.Client{
		Transport: &transport.APIKey{Key: api_key}, //Contained within info.go which is ignored. Go make your own API Key lol
	}

	service, err1 := youtube.New(client)
	check("Error creating new Youtube client\n%v", err1)

	playlistInfoCall := service.Playlists.List([]string{"snippet", "contentDetails"}).Id(id)

	infoResponse, err3 := playlistInfoCall.Do()
	check("Error getting Info of Playlist("+id+")\n%v", err3)

	if autoMax {
		max = infoResponse.Items[0].ContentDetails.ItemCount
	}

	playListItemsCall := service.PlaylistItems.List([]string{"snippet"}).PlaylistId(id).MaxResults(max)
	itemsResponse, err2 := playListItemsCall.Do()
	check("Error getting List of Playlist("+id+") Items\n%v", err2)

	f, err4 := os.OpenFile(outputFile+".txt", os.O_TRUNC, os.ModeAppend)
	if err4 != nil { //Most likely the file doesn't exist
		f, err4 = os.Create(outputFile + ".txt")
		check("Error creating/opening output file\n%v", err4)
	}
	defer f.Close()

	//err := f.Truncate(0)
	//check("Error truncating file\n%v", err)

	_, err := f.Seek(0, 0)
	check("Error seeking in file\n%v", err)

	f.WriteString("Playlist \"" + infoResponse.Items[0].Snippet.Title + "\": " + infoResponse.Items[0].Snippet.Description)

	for _, item := range itemsResponse.Items {
		f.WriteString("\n\n")
		f.WriteString(item.Snippet.Title + ": " + item.Snippet.VideoOwnerChannelTitle + "\n -Link: https://youtu.be/" + item.Snippet.ResourceId.VideoId)
	}

	f.Sync()
}
