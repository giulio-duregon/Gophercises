package adventure

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

type storyContainer map[string]story

type Options struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
type story struct {
	Title   string    `json:"title"`
	Story   []string  `json:"story"`
	Options []Options `json:"options"`
}

func parseInput() (renderedContentPath string, storyPath, servePath, templatePath *string, port *int) {
	storyPath = flag.String("p", "static/adventure/stories.json", "Path to json file containing story information. 'Intro' key is the starting point for the story of the JSON map")
	servePath = flag.String("d", "static/adventure/", "Path to serve rendered html files and css styling")
	renderedContentPath = *servePath + "chapters/"
	templatePath = flag.String("t", "static/adventure/storyTemplate.html", "Path to json file containing story information. 'Intro' story should be the first story of json array")
	port = flag.Int("port", 8080, "The port for the file server to listen on")
	flag.Parse()
	return
}

func getStories(path *string) (bytes []byte) {
	bytes, err := os.ReadFile(*path)
	if err != nil {
		log.Fatalf("Could not read json file at path: %s\n", *path)
	}
	return

}

func parseStories(bytes []byte) (stories storyContainer) {
	err := json.Unmarshal(bytes, &stories)
	if err != nil {
		log.Fatalf("Could not unmarshal bytes from story, err: %v\n", err)
	}
	return
}

func generateStoryHTML(s storyContainer, dirPath, templatePath string) {
	// Make directory
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Couldn't make the directory: %s\n%v\n", dirPath, err)
	}
	// Create templates for stories
	adventureTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	for name, scene := range s {
		path := filepath.Join(dirPath, name+".html")
		out, err := os.Create(path)
		if err != nil {
			log.Fatalf("Could not create path: %s\n %v\n", path, err)
		}
		err = adventureTemplate.Execute(out, &scene)
		if err != nil {
			log.Fatalf("Could not execute template, %v\n", err)
		}
		err = out.Close()
		if err != nil {
			log.Fatalf("Could not close file %s, %v\n", path, err)
		}
	}
}

func CYOAProgram() {
	// Parse input
	renderedContentPath, storiesPath, servePath, templatePath, port := parseInput()
	// Retrieve stories and parse into JSON
	storyBytes := getStories(storiesPath)
	// Parse json into stories
	stories := parseStories(storyBytes)
	// Use stories in templates, create .html rendered templates
	generateStoryHTML(stories, renderedContentPath, *templatePath)

	// Create routes / handlers for stories
	adventureMux := http.NewServeMux()
	adventureMux.Handle("/", http.RedirectHandler("/adventure/chapters/intro.html", http.StatusPermanentRedirect))
	adventureMux.Handle("/adventure/", http.StripPrefix("/adventure", http.FileServer(http.Dir(*servePath))))

	// Start webserver and serve
	fmt.Printf("Listening on port: %d\n", *port)
	fmt.Printf("Link: http://localhost:%d/\n", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), adventureMux))
}
