package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type BookReview struct {
	Title     string
	Author    string
	URL       string
	Status    string
	Rating    float64
	Review    string
	Date_Read string
	Genre     string
}

var TEMPLATE string = `Title: Books
Date: 2020-02-02
Slug: books
Authors: hvnsweeting
Summary: Book notes

See [Recfile]({attach}books.rec)

[GNU Recutils](https://www.gnu.org/software/recutils/manual/recutils.html#Top)

Discuss on [HackerNews](https://news.ycombinator.com/item?id=22153665)
{{ range $r := . }}
## {{ $r.Rating}}⭐ {{ $r.Title }} by {{ $r.Author }}
{{ end }}
`

func main() {
	f, err := os.Open("content/pages/books.rec")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	var review BookReview
	var rs []BookReview
	for {
		if scanner.Scan() == false {
			break
		}
		t := scanner.Text()
		if strings.HasPrefix(t, "Title:") {
			title := strings.TrimPrefix(t, "Title:")
			review.Title = strings.Trim(title, " ")
		} else if strings.HasPrefix(t, "Author:") {
			author := strings.TrimPrefix(t, "Author:")
			review.Author = strings.Trim(author, " ")
		} else if strings.HasPrefix(t, "Rating:") {
			rate := strings.TrimPrefix(t, "Rating:")
			rate = strings.Trim(rate, " ")
			r, err := strconv.ParseFloat(rate, 64)
			if err != nil {
				log.Fatal(err)
			}
			review.Rating = r

			rs = append(rs, review)
			// TODO how to create NEW review when got new book
			// this one is prone to use old book value if the new book missing the field.
		}

		// fmt.Printf("%s\n", review)
	}
	// fmt.Printf("%v\n", rs)

	t := template.New("namename")
	_, err = t.Parse(TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(os.Stdout, rs)
	if err != nil {
		log.Fatal(err)
	}
}
