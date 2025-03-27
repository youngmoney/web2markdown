package main

import (
	"flag"
	"fmt"
	readability "github.com/go-shiori/go-readability"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func fetch(uri string, agent string) (string, error) {
	if agent == "curl-browser" {
		return ExecuteCommand(`curl -L -A 'Mozilla/5.0 (Macintosh; Intel Mac OS X 14_7_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3 Safari/605.1.15' -s "$1"`, []string{uri}, "")
	}
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}
	if agent != "" {
		request.Header.Set("User-Agent", agent)
	}
	request.Header.Set("Accept", "*/*")
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	all, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(all), err
}

func command(uri string, agent string, minContent int) {
	content, err := fetch(uri, agent)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error fetching:", err)
		os.Exit(1)
	}

	parsedURL, err := url.Parse(uri)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing uri:", err)
		os.Exit(1)
	}

	article, err := readability.FromReader(strings.NewReader(content), parsedURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing:", err)
		os.Exit(1)
	}
	var body = article.Content
	if body == "" {
		body = content
	}
	converted, err := ExecuteCommand(`pandoc --from "html" --extract-media=/dev/null --to commonmark-raw_attribute-raw_html --lua-filter <(echo "function Image(i) return pandoc.Span({}) end")`, []string{}, body)
	ExitIfNonZero(err)
	if len(converted) > minContent {
		fmt.Println("#", article.Title)
		if article.Byline != "" {
			fmt.Println()
			fmt.Println("By", article.Byline)
		}
		if article.PublishedTime != nil {
			fmt.Println()
			fmt.Println("Published", article.PublishedTime.Format(time.DateOnly))

		}
		fmt.Println()
		fmt.Print(converted)
	} else {
		fmt.Fprint(os.Stderr, "not enough content:\n", converted)
		os.Exit(1)
	}
}

func main() {
	agent := flag.String("user-agent", "", "user agent")
	minContent := flag.Int("min-content", 200, "the minimum content length after conversion to output")
	flag.Parse()
	uri := flag.Arg(0)
	if uri == "" {
		fmt.Fprintln(os.Stderr, "web2markdown <URI>")
		os.Exit(1)
	}

	command(uri, *agent, *minContent)
}
