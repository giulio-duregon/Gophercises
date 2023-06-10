package main

import (
	"Gophercises/src/quizGame"
	"Gophercises/src/urlShortener"
	"flag"
)

const (
	QuizGame int = iota + 1
	URLShortener
)

func exerciseSelector() {

}

func main() {
	program := flag.Int("p", 1, "Selects the exercise that you want to run, 1-based indexing")
	flag.Parse()
	switch *program {
	case QuizGame:
		quizGame.QuizProgram()
	case URLShortener:
		urlShortener.UrlShortProgram()
	}

}
