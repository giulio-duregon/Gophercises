package quizGame

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

const statReport = `
------------------------------------------------------------------
You got {{.CorrectAnswers}}/{{.TotalQuestions}} questions correct!
That is a percentage of {{divide .CorrectAnswers .TotalQuestions}}%
------------------------------------------------------------------
`

type report struct {
	CorrectAnswers, TotalQuestions int
}

type problem struct {
	Q, A string
}

func tempDivide(a, b int) string {
	return fmt.Sprintf("%2.2f", float64(a)*100/float64(b))
}

func setup() (*string, *bool, time.Duration) {
	quizDuration := time.Second * 30
	v := flag.Bool("verbose", false, "Prints out if you answered a question incorrectly")
	flag.DurationVar(&quizDuration, "time", time.Second*30, "The amount of seconds to complete the quiz")
	fileName := flag.String("csv", "static/problems.csv", "The path to the question/answer csv. 2 Columns, the first being question, the second answer.")
	flag.Parse()
	return fileName, v, quizDuration
}

func getCsv(fileName *string) (lines [][]string) {
	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatalf("Could not open file %s\n", *fileName)
	}

	r := csv.NewReader(file)
	lines, err = r.ReadAll()

	if err != nil || len(lines) == 0 {
		log.Fatalf("Failed to parse provided csv file: %s\n", *fileName)
	}
	return
}

func lineParse(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		problem := problem{
			// Get rid of spaces
			Q: strings.TrimSpace(line[0]),
			A: strings.TrimSpace(line[1]),
		}
		ret[i] = problem
	}
	return ret
}

func reportStats(sessionReport report, funcMap template.FuncMap) {
	// Output the results
	sessionStats := template.Must(template.New("OutputResults").Funcs(funcMap).Parse(statReport))
	err := sessionStats.Execute(os.Stdout, sessionReport)
	if err != nil {
		log.Fatalf("For whatever reason, we couldn't output the results...")
	}
}

func promptInput(answerChan chan string, i int, p *problem) {
	fmt.Printf("Problem #%d: %s = \n", i+1, p.Q)
	var answer string
	count := 0
	for count == 0 {
		n, _ := fmt.Scanf("%s\n", &answer)
		count = n
	}
	answerChan <- answer
}

func QuizProgram() {
	// Create FuncMap for template reporting
	tempFuncMap := template.FuncMap{"divide": tempDivide}

	// Parse user input
	fileName, v, duration := setup()

	// Create timer object
	timer := time.NewTimer(duration)

	// Read in lines from CSV, parse into array of problem struct
	lines := getCsv(fileName)
	problems := lineParse(lines)

	// Quiz the user on questions
	sessionReport := report{
		CorrectAnswers: 0,
		TotalQuestions: len(problems),
	}

	answerChan := make(chan string, 1)
	// Iterate over the problem set
	for i, p := range problems {
		// Run user answers async in different process
		go promptInput(answerChan, i, &p)

		// Select between user answers and timer indicating time limit reached
		select {
		case <-timer.C:
			// Time is up, stop program and output results
			println("Time limit reached")
			reportStats(sessionReport, tempFuncMap)
			return

		case answer := <-answerChan:
			// Otherwise continue with questions, increment results
			if answer == p.A {
				sessionReport.CorrectAnswers++
			} else if *v == true {
				fmt.Printf("Incorrect Answer, %s=%s, you answered %s\n", p.Q, p.A, answer)
			}
		}
	}

	// Output final statistics
	reportStats(sessionReport, tempFuncMap)
}
