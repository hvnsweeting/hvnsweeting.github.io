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
	Free      bool
	Status    string
	Rating    int64
	Review    string
	Date_Read string
	Genre     string
}

var TEMPLATE string = `Title: Books review
Date: 2020-02-02
Slug: books
Authors: hvnsweeting
Summary: Books reviews

See [Recfile]({attach}books.rec)

[GNU Recutils](https://www.gnu.org/software/recutils/manual/recutils.html#Top)

Discuss on [HackerNews](https://news.ycombinator.com/item?id=22153665)

The language I used for book title is the language I read the book.
Tên tác phẩm dùng ngôn ngữ nào thì tôi đọc cuốn sách bằng ngôn ngữ đó.

{{ range $r := . }}
##{{ $r.Rating}}⭐ [{{ $r.Title }} by {{ $r.Author }}{{ if $r.Free }} [Free]{{end}}]({{ $r.URL }})
*{{ $r.Status }}* *{{ $r.Date_Read }}*

{{ $r.Review }}
{{ end }}`

func main() {
	f, err := os.Open("content/pages/books.rec")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	review := BookReview{}
	var rs []BookReview
	for scanner.Scan() {
		t := scanner.Text()

		if strings.HasPrefix(t, "Title:") {
			title := strings.TrimPrefix(t, "Title:")
			review.Title = strings.Trim(title, " ")
		} else if strings.HasPrefix(t, "Author:") {
			author := strings.TrimPrefix(t, "Author:")
			review.Author = strings.Trim(author, " ")
		} else if strings.HasPrefix(t, "Date_Read:") {
			review.Date_Read = strings.Trim(strings.TrimPrefix(t, "Date_Read:"), " ")
		} else if strings.HasPrefix(t, "Status:") {
			review.Status = strings.Trim(strings.TrimPrefix(t, "Status:"), " ")
		} else if strings.HasPrefix(t, "URL:") {
			review.URL = strings.Trim(strings.TrimPrefix(t, "URL:"), " ")
			if strings.Contains(review.URL, "standardebook") {
				review.Free = true
			}
		} else if strings.HasPrefix(t, "Review:") {
			review.Review = strings.Trim(strings.TrimPrefix(t, "Review:"), " ")
		} else if strings.HasPrefix(t, "Rating:") {
			rate := strings.TrimPrefix(t, "Rating:")
			rate = strings.Trim(rate, " ")
			r, err := strconv.ParseInt(rate, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			review.Rating = r

		} else if strings.Trim(t, "") == "" {
			if review.Title != "" {
				rs = append(rs, review)
			}
			review = BookReview{}
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
