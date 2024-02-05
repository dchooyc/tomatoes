package film

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
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
