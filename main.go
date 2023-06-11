package main

import (
	"Gophercises/src/adventure"
	"Gophercises/src/quizGame"
	"Gophercises/src/urlShortener"
	"log"
)

const (
	QuizGame int = iota + 1
	URLShortener
	Adventure
)

func exerciseSelector(program *int) func() {
	switch *program {
	case QuizGame:
		return quizGame.QuizProgram
	case URLShortener:
		return urlShortener.UrlShortProgram
	case Adventure:
		return adventure.CYOAProgram
	}
	return func() {
		log.Fatalf("Invalid argument, you entered: %d\n", *program)
	}
}

func main() {
	program := 3
	exerciseSelector(&program)()

}
