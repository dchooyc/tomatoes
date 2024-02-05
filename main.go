package main

import (
	"flag"
	"fmt"
	"tomatoes/film"
)

const (
	tomatoesURL = "https://www.rottentomatoes.com/m/"
	shawshank   = "shawshank_redemption"
	defaultRoot = tomatoesURL + shawshank
)

func main() {
	root := flag.String("url", defaultRoot, "The url to begin crawling from")
	flag.Parse()

	f, err := film.GetFilm(*root)
	if err != nil {
		fmt.Println("get film failed: ", err)
	}

	fmt.Println(*f)
}
