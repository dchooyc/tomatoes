package main

import (
	"fmt"
	"tomatoes/film"
)

func main() {
	target := "https://www.rottentomatoes.com/m/shawshank_redemption"
	f, err := film.GetFilm(target)
	if err != nil {
		fmt.Println("get film failed: ", err)
	}
	fmt.Println(*f)
}
