package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/coreos/go-log/log"
	"github.com/go-martini/martini"
	omxplayer "github.com/longnguyen11288/go-omxplayer"
)

var dataDir string

type File struct {
	Filename string `json:"filename"`
}

type Files []string

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

func FilesHandler(w http.ResponseWriter) {
	var files Files
	osFiles, _ := ioutil.ReadDir(dataDir)
	for _, f := range osFiles {
		files = append(files, f.Name())
	}
	output, err := json.Marshal(files)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, output)
}

func main() {
	flag.StringVar(&dataDir, "data-dir", ".", "Data directory for videos")
	flag.Parse()

	os.Chdir(dataDir)
	player := omxplayer.New()
	m := martini.Classic()
	m.Map(&player)
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Get("/files", FilesHandler)
	m.Post("/playfile", PlayFileHandler)
	m.Post("/stopfile", StopFileHandler)
	m.Run()
}
