package main

import (
	"net/http"
	"html/template"
	"log"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"io/ioutil"
	"strconv"
	"os/exec"
)

type Data struct {
	UpdateId string
	Message struct {
		Text string
	}
}
type Embed struct {
	Embed template.HTML
}

func main() {
	num := lastPostNumber()
	embed := fetch(num)
	if len(embed) > 0 || num == "267" {
		writeToHtml(embed)
		writeLastPostNumber(strconv.Itoa(plusOne(num)))
		//pushToGitHub()
	}
}

func pushToGitHub() {
	cmdName := "git"
	commitArgs := []string{"commit", "--all"}
	pushArgs := []string{"push"}
	if out, err := exec.Command(cmdName, commitArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git commit command: ", out)
		os.Exit(1)
	}
	if out, err := exec.Command(cmdName, pushArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git push command: ", out)
		os.Exit(1)
	}
}
func lastPostNumber() string {
	b, err := ioutil.ReadFile("num.txt") // just pass the file name
	if err != nil {
		fmt.Println(err)
	}

	return string(b)
}

func plusOne(postNumber string) int{
	i, err := strconv.Atoi(postNumber)
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	i++
	return i
}

func writeLastPostNumber(number string) {
	f, err := os.Create("num.txt")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	f.WriteString(number)
}

func writeToHtml(embedString template.HTML) {
	t, err := template.ParseFiles("template.html")
	if err != nil {
		log.Print(err)
		return
	}

	f, err := os.Create("telegram.html")
	if err != nil {
		log.Println("create file: ", err)
		return
	}
	e := Embed{
		Embed: embedString,
	}

	err = t.Execute(f, e)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
	f.Close()
}

func fetch(lastPostNumber string) template.HTML {
	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf("https://t.me/dkollection/%d/", plusOne(lastPostNumber)))
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var embed string
	// Find the embed items
	doc.Find("#embed_code_field").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		embed = s.Text()
		if embed == "" {
			fmt.Println("No embed")
			return
		}
	})
	return template.HTML(embed)
}
