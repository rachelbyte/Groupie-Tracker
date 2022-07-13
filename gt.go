package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// Tracker slice of Artists as there are many Artists in the API.
type Tracker struct {
	Artists []Artist
}

// Artist struct containing relevant information about each artist such as image date and location.
type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Relations    Relation `json:"relations"`
}

/*The Relation struct links the dates and locations together
it is easier to unmarshal than dates and locations separately*/
type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

/* Relation starts with an index which includes the whole map
has to be a slice  of Relations with the Id index and datesLocations*/

type RelStrct struct {
	Index []Relation `json:"index"`
}

// Naming the variables I am using
var (
	res      Tracker
	t        *template.Template
	tpl      *template.Template
	relStrct RelStrct
)

// Main this handles the mainpage, linking the CSS in the static folder and also listening on port 3000
func main() {
	http.HandleFunc("/", Art)
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":3000", nil)
}

// This function handles getting the data from the API and unmarshalling it into a JSON struct
func Art(w http.ResponseWriter, r *http.Request) {
	// Getting the API data
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	// Using IO to read the file
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// Unmarshalling the data []byte into the struct artists
	json.Unmarshal(body, &res.Artists)

	// Parsing the html error file
	tpl, err = template.ParseFiles("err.html")
	if err != nil {
		log.Fatal(err)
	}
	// Parsing the main or artists file
	t, err := template.ParseFiles("artists.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.Execute(w, http.StatusInternalServerError)
		return

	}
	// If Mainpage is not / then 404 not found
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		tpl.Execute(w, http.StatusNotFound)
		return
	}
	// If not method get could be Post then return an error
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		tpl.Execute(w, http.StatusMethodNotAllowed)
		return
	}

	Rel()
	// Making sure the datesLocations is the same in artists and in Relations
	for i := range relStrct.Index {
		res.Artists[i].Relations.DatesLocations = relStrct.Index[i].DatesLocations
	}

	t.Execute(w, res.Artists)
}

// Same things as Artist unmarshalling the Relations struct

func Rel() {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &relStrct)
}
