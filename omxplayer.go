package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-martini/martini"
	omxplayer "github.com/longnguyen11288/go-omxplayer"
)

type File struct {
	Filename string `json:"filename"`
}

func PlayFileHandler(player *omxplayer.OmxPlayer,
	w http.ResponseWriter, r *http.Request) {
	var file File
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &file)
	err = player.PlayFile(file.Filename)
	if err != nil {
		fmt.Fprint(w, `{ "error": "`+err.Error()+`" }`)
		return
	}
	fmt.Fprint(w, `{ "success": "true" }`)
}

func StopFileHandler(player *omxplayer.OmxPlayer, w http.ResponseWriter) {
	err := player.StopFile()
	if err != nil {
		fmt.Fprint(w, `{ "error": "`+err.Error()+`" }`)
		return
	}
	fmt.Fprint(w, `{ "success": "true" }`)
}

func main() {
	var dataDir string
	flag.StringVar(&dataDir, "data-dir", ".", "Data directory for videos")
	flag.Parse()

	os.Chdir(dataDir)
	player := omxplayer.New()
	m := martini.Classic()
	m.Map(&player)
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Post("/playfile", PlayFileHandler)
	m.Post("/stopfile", StopFileHandler)
	m.Run()
}
