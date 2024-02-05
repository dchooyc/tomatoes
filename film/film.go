package film

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const (
	FilmPosterIndicator = "Watch trailer for "
	MovieIndicator      = "/m/"
)

type Films struct {
	Films []Film `json:"films"`
}

type Film struct {
	Title         string   `json:"title"`
	URL           string   `json:"url"`
	PosterUrl     string   `json:"poster_url"`
	MediaType     string   `json:"media_type"`
	Rating        string   `json:"rating"`
	Year          string   `json:"year"`
	Genre         string   `json:"genre"`
	Runtime       string   `json:"runtime"`
	AudienceScore int      `json:"audience_score"`
	TomatoScore   int      `json:"tomato_score"`
	Ratings       int      `json:"ratings"`
	SimilarFilms  []string `json:"similar_films"`
}

func GetFilm(url string) (*Film, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	resp, err := http.DefaultClient.Do(req)
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
	if n.Type == html.ElementNode {
		if n.Data == "h1" {
			extractTitle(n, curFilm)
		}

		if n.Data == "a" {
			extractSimilar(n, curFilm)
			extractRatings(n, curFilm)
		}

		if n.Data == "p" {
			extractYearGenreRuntime(n, curFilm)
		}

		if n.Data == "rt-img" {
			extractPosterUrl(n, curFilm)
		}

		if n.Data == "score-board-deprecated" {
			extractScoreRatingMedia(n, curFilm)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractFilmInfo(c, curFilm)
	}
}

func extractSimilar(n *html.Node, curFilm *Film) {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			url := attr.Val

			if strings.HasPrefix(url, MovieIndicator) {
				parts := strings.Split(url[3:], "/")
				if len(parts) == 1 {
					curFilm.SimilarFilms = append(curFilm.SimilarFilms, parts[0])
				}
			}

			break
		}
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

func extractScoreRatingMedia(n *html.Node, curFilm *Film) {
	for _, attr := range n.Attr {
		if attr.Key == "rating" {
			curFilm.Rating = attr.Val
		}

		if attr.Key == "audiencescore" {
			val, _ := strconv.Atoi(attr.Val)
			curFilm.AudienceScore = val
		}

		if attr.Key == "mediatype" {
			curFilm.MediaType = attr.Val
		}

		if attr.Key == "tomatometerscore" {
			val, _ := strconv.Atoi(attr.Val)
			curFilm.TomatoScore = val
		}
	}
}

func extractYearGenreRuntime(n *html.Node, curFilm *Film) {
	for _, attr := range n.Attr {
		if attr.Key == "data-qa" && attr.Val == "score-panel-subtitle" {
			textNode := n.FirstChild
			if textNode != nil && textNode.Type == html.TextNode {
				parts := strings.Split(textNode.Data, ", ")
				for i := 0; i < len(parts); i++ {
					if i == 0 {
						curFilm.Year = parts[i]
					}

					if i == 1 {
						curFilm.Genre = parts[i]
					}

					if i == 2 {
						curFilm.Runtime = parts[i]
					}
				}
				break
			}
		}
	}
}

func extractRatings(n *html.Node, curFilm *Film) {
	for _, attr := range n.Attr {
		if attr.Key == "data-qa" && attr.Val == "audience-rating-count" {
			textNode := n.FirstChild
			if textNode != nil && textNode.Type == html.TextNode {
				parts := strings.Fields(textNode.Data)
				strVal := strings.Split(parts[0], "+")[0]
				val := strings.Join(strings.Split(strVal, ","), "")
				ratings, _ := strconv.Atoi(val)
				curFilm.Ratings = ratings
				break
			}
		}
	}
}
