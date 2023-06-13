package main

import (
	"Gophercises/src/adventure"
	"Gophercises/src/link"
	"Gophercises/src/quizGame"
	"Gophercises/src/urlShortener"
	"log"
)

const (
	QuizGame int = iota + 1
	URLShortener
	Adventure
	LinkParser
)

func exerciseSelector(program *int) func() {
	switch *program {
	case QuizGame:
		return quizGame.QuizProgram
	case URLShortener:
		return urlShortener.UrlShortProgram
	case Adventure:
		return adventure.CYOAProgram
	case LinkParser:
		return link.LinkProgram
	}
	return func() {
		log.Fatalf("Invalid argument, you entered: %d\n", *program)
	}
}

func main() {
	program := 4
	exerciseSelector(&program)()

}
