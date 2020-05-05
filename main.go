package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

// isExtraneousInfo filters string input based on a given string slice
// TODO: return error if error
func isExtraneousInfo(text string, list []string) bool {
	for _, el := range list {
		// fmt.Println(strings.Contains(text, el))
		// fmt.Printf("%s %s", text, el)
		if strings.Contains(text, el) {
			return true
		}
	}
	return false
}

// writeSliceToFile iterates through a []string slice
// and creates a file at the provided filePath
// TODO: return error if error
func writeSliceToFile(filePath string, slice []string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for _, el := range slice {
		fmt.Fprintln(f, el)
	}

}

func main() {
	c := colly.NewCollector()

	// question bank slice
	qb := []string{}
	// answer bank slice
	answers := []string{}

	// Find and visit all links
	c.OnHTML(".post", func(e *colly.HTMLElement) {
		extraneous := []string{
			"Exam-Sample",
			"Quickly",
			"Set",
			"Get your Absolutely",
			"See author's",
			"An expert on R&D",
			"Correct Answer of the above",
			"Correct Answers",
			"Access the Full Database of all Questions",
		}
		// response := e.ChildTexts("p")
		// fmt.Println(e.Text)
		e.ForEach("p", func(count int, e *colly.HTMLElement) {
			text := e.Text
			if !isExtraneousInfo(text, extraneous) {
				qb = append(qb, text)
			}
			// fmt.Println(e.Text)
		})

	})

	c.OnHTML("table", func(e *colly.HTMLElement) {
		extraneous := []string{
			"Question No.",
			"Correct Answer",
		}

		e.ForEach("tr", func(count int, e *colly.HTMLElement) {
			text := e.Text
			if !isExtraneousInfo(text, extraneous) {
				// trimmedText := strings.TrimSpace(text)
				// trimmedText = strings.ReplaceAll(trimmedText, " ", "")
				answers = append(answers, text)
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// max question number is 1060
	for min, max, increment := 61, 70, 10; max <= 1060; {
		url := fmt.Sprintf("https://www.softwaretestinggenius.com/istqb-certification-exam-sample-papers-q-%d-to-%d/", min, max)
		c.Visit(url)
		min = min + increment
		max = max + increment
	}

	writeSliceToFile("vault/qb.txt", qb)
	writeSliceToFile("vault/answers.txt", answers)
}
