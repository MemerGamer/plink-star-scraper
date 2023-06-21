package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Define a struct to represent the JSON structure
type SearchResult struct {
	Results []struct {
		Name        string   `json:"name"`
		URL         string   `json:"url"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
		Score       float64  `json:"score"`
		Stars       int      `json:"stars"`
	} `json:"results"`
	Total int `json:"total"`
}

func parseStars(stars string) (int, error) {
	// Convert stars from string representation (e.g., "1.2k") to integer
	stars = strings.ToLower(stars)
	if strings.HasSuffix(stars, "k") {
		// Remove "k" suffix
		stars = strings.TrimSuffix(stars, "k")

		// Parse the float value
		floatVal, err := strconv.ParseFloat(stars, 64)
		if err != nil {
			return 0, err
		}

		// Multiply by 1000 to get the actual count
		starsInt := int(floatVal * 1000)
		return starsInt, nil
	}

	// setting up -1 value for stars if the stars is empty aka not found
	if stars == "" {
		stars = "-1"
	}

	// Parse the integer value directly
	starsInt, err := strconv.Atoi(stars)
	if err != nil {
		return 0, err
	}

	return starsInt, nil
}

func fetchSelector(url string) (int, error) {
	// Make an HTTP GET request to the GitHub URL
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	var body strings.Builder
	_, err = io.Copy(&body, resp.Body)
	if err != nil {
		return 0, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body.String()))
	if err != nil {
		return 0, err
	}

	selector := doc.Find("#repo-stars-counter-star").Text()
	stars, err := parseStars(selector)
	if err != nil {
		return 0, err
	}

	return stars, nil
}

var (
	inputFilename  = "search-4-telescope.json"
	outputFilename = "search-4-telescope-updated.json"
)

func main() {
	// Open the input JSON file
	inputFile, err := os.Open(inputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	// Read the input JSON data from the file
	var inputBody strings.Builder
	_, err = io.Copy(&inputBody, inputFile)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON data into the SearchResult struct
	var searchResult SearchResult
	err = json.Unmarshal([]byte(inputBody.String()), &searchResult)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the results and fetch the desired selector from each URL
	for i := range searchResult.Results {
		result := &searchResult.Results[i]

		// If the URL does not start with "https://github.com",
		// try to fetch "https://github.com/" + "name"
		if !strings.HasPrefix(result.URL, "https://github.com") {
			baseURL, err := url.Parse("https://github.com/")
			if err != nil {
				log.Fatal(err)
			}
			newURL, err := url.Parse(result.Name)
			if err != nil {
				log.Fatal(err)
			}
			result.URL = baseURL.ResolveReference(newURL).String()
		}

		// Fetch the selector and update the star count
		stars, err := fetchSelector(result.URL)
		if err != nil {
			log.Printf("Error fetching stars for %s: %v", result.URL, err)
			result.Stars = -1
			continue
		}
		result.Stars = stars

	}

	// Create the output JSON file
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Marshal the updated SearchResult struct into JSON
	outputData, err := json.MarshalIndent(searchResult, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Write the JSON data to the output file
	_, err = outputFile.Write(outputData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated star counts written to %s\n", outputFilename)
}
