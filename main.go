package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/cdipaolo/sentiment"
)

func readReviews(fp string) ([]string, error) {
	file, err := os.Open(fp)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var reviews []string
	for scanner.Scan() {
		reviews = append(reviews, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return reviews, nil
}

func processReview(review string, results chan<- string) {
	cleanedReview := preprocessText(review)
	sentimentScore := analyzeSentiment(cleanedReview)
	results <- fmt.Sprintf("Review: %s, Sentiment score: %d", review, sentimentScore)
}

func preprocessText(text string) string {
	cleanedText := strings.ToLower(text)
	cleanedText = strings.ReplaceAll(cleanedText, ".", "")
	cleanedText = strings.ReplaceAll(cleanedText, ",", "")
	// more processing steps as needed
	return cleanedText
}

func analyzeSentiment(review string) uint8 {
	model, err := sentiment.Restore()
	if err != nil {
		fmt.Println(err)
	}
	analysis := model.SentimentAnalysis(review, sentiment.English)
	return analysis.Score
}

func main() {
	reviews, err := readReviews("./reviews.txt")
	if err != nil {
		log.Fatal(err)
	}

	results := make(chan string)
	var wg sync.WaitGroup

	for _, review := range reviews {
		wg.Add(1)
		go func(review string) {
			defer wg.Done()
			processReview(review, results)
		}(review)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
