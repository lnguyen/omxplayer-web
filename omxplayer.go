package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	player.PlayFile(file.Filename)
	fmt.Fprint(w, `{ "success": "true" }`)
}

func main() {
	player := omxplayer.New()
	m := martini.Classic()
	m.Map(player)
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Post("/playfile", PlayFileHandler)
	m.Run()
}
