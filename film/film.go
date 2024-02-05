package film

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const (
	FilmPosterIndicator = "Watch trailer for "
)

type Film struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	PosterUrl string `json:"poster_url"`
	Rating    string `json:"rating"`
	Year      string `json:"year"`
	Genre     string `json:"genre"`
	Runtime   string `json:"runtime"`
	Score     string `json:"score"`
	Ratings   string `json:"ratings"`
}

func GetFilm(url string) (*Film, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get failed: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html failed: %w", err)
	}

	film := &Film{}

	extractFilmInfo(doc, film)

	return film, nil
}

func extractFilmInfo(n *html.Node, curFilm *Film) {
	if n.Type == html.ElementNode && n.Data == "h1" {
		extractTitle(n, curFilm)
	}

	if n.Type == html.ElementNode && n.Data == "p" {
		extractDetails(n, curFilm)
	}

	if n.Type == html.ElementNode && n.Data == "rt-img" {
		extractPosterUrl(n, curFilm)
	}

	if n.Type == html.ElementNode && n.Data == "score-board-deprecated" {
		extractRating(n, curFilm)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractFilmInfo(c, curFilm)
	}
}

func extractTitle(n *html.Node, curFilm *Film) {
	for _, attr := range n.Attr {
		if attr.Key == "class" && attr.Val == "title" {
			tNode := n.FirstChild
			if tNode != nil && tNode.Type == html.TextNode {
				curFilm.Title = tNode.Data
				break
			}
		}
	}
}

func extractPosterUrl(n *html.Node, curFilm *Film) {
	correctImg := false
	src := ""
	for _, attr := range n.Attr {
		if attr.Key == "alt" && strings.HasPrefix(attr.Val, FilmPosterIndicator) {
			correctImg = true
		}

		if attr.Key == "src" {
			src = attr.Val
		}
	}

	if correctImg {
		curFilm.PosterUrl = src
	}
}

func extractRating(n *html.Node, curFilm *Film) {
	for _, attr := range n.Attr {
		if attr.Key == "rating" {
			curFilm.Rating = attr.Val
		}
	}
}

func extractDetails(n *html.Node, curFilm *Film) {
	for _, attr := range n.Attr {
		if attr.Key == "data-qa" && attr.Val == "score-panel-subtitle" {
			textNode := n.FirstChild
			if textNode != nil && textNode.Type == html.TextNode {
				parts := strings.Split(textNode.Data, ", ")
				curFilm.Year = parts[0]
				curFilm.Genre = parts[1]
				curFilm.Runtime = parts[2]
				break
			}
		}
	}
}
