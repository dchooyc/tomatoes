package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
	"tomatoes/film"
)

const (
	tomatoesURL = "https://www.rottentomatoes.com/m/"
	shawshank   = "shawshank_redemption"
	out         = "output.json"
)

type processedFilm struct {
	film *film.Film
	err  error
}

func main() {
	root := flag.String("title", shawshank, "The title to begin crawling from")
	output := flag.String("output", out, "Output location")
	maxDepth := flag.Int("depth", 2, "The depth at which to stop crawling")
	numWorkers := flag.Int("workers", 20, "The number of workers to process films")
	flag.Parse()

	file, err := os.Create(*output)
	if err != nil {
		panic(err)
	}

	queue := []string{*root}
	titleToFilm := make(map[string]*film.Film)
	findFilms(queue, titleToFilm, *maxDepth, *numWorkers)

	films := arrangeFilms(titleToFilm)

	jsonData, err := json.Marshal(films)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("writing to file failed: ", err)
	}
}

func arrangeFilms(titleToFilm map[string]*film.Film) film.Films {
	arranged := []film.Film{}

	for _, curFilm := range titleToFilm {
		if curFilm != nil {
			arranged = append(arranged, *curFilm)
		}
	}

	sort.Slice(arranged, func(i, j int) bool {
		return arranged[i].Ratings > arranged[j].Ratings
	})

	return film.Films{Films: arranged}
}

func findFilms(queue []string, titleToFilm map[string]*film.Film, maxDepth, numWorkers int) {
	for i := 1; i <= maxDepth; i++ {
		fmt.Println("depth: " + strconv.Itoa(i))
		fmt.Println("films: " + strconv.Itoa(len(queue)))
		isLast := false

		if i == maxDepth {
			isLast = true
		}

		queue = processQueue(isLast, numWorkers, queue, titleToFilm)
	}
}

func processQueue(isLast bool, numWorkers int, queue []string, titleToFilm map[string]*film.Film) []string {
	titles := make(chan string, len(queue))
	processedFilms := make(chan *processedFilm, len(queue))
	var wg sync.WaitGroup

	createWorkers(min(len(queue), numWorkers), isLast, titles, processedFilms, &wg)

	for _, title := range queue {
		wg.Add(1)
		titles <- title
	}

	close(titles)

	go func() {
		wg.Wait()
		close(processedFilms)
	}()

	collect := make(map[string]bool)

	for pFilm := range processedFilms {
		if pFilm.err != nil {
			fmt.Println(pFilm.err)
			continue
		}

		titleToFilm[pFilm.film.Title] = pFilm.film

		for _, title := range pFilm.film.SimilarFilms {
			collect[title] = true
		}
	}

	next := []string{}

	for title := range collect {
		if _, ok := titleToFilm[title]; !ok {
			next = append(next, title)
		}
	}

	return next
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func createWorkers(numWorkers int, isLast bool, titles <-chan string, processedFilms chan<- *processedFilm, wg *sync.WaitGroup) {
	for w := 0; w < numWorkers; w++ {
		go worker(w, isLast, titles, processedFilms, wg)
	}
}

func worker(workerID int, isLast bool, titles <-chan string, processedFilms chan<- *processedFilm, wg *sync.WaitGroup) {
	for title := range titles {
		url := tomatoesURL + title
		pFilm := processFilm(isLast, url)
		if pFilm.film != nil {
			fmt.Printf("Worker %d: %s\n", workerID, pFilm.film.Title)
		}
		processedFilms <- pFilm
		wg.Done()
	}
}

func processFilm(isLast bool, url string) *processedFilm {
	res := &processedFilm{}

	curFilm, err := film.GetFilm(url)
	if err != nil {
		res.err = fmt.Errorf("error getting %s: %w", url, err)
		return res
	}

	curFilm.URL = url
	res.film = curFilm
	return res
}
