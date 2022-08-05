package main // MicroCMS: One-file Markdown file to HTML CMS
// git clone https://github.com/tomcam/microcms
// cd microcms
// go mod init github.com/tomcam/microcms
// go mod tidy

// Example invocations
// Include the 2 css files shown
// go run main.go -styles "theme.css light-mode.css" foo.md

// Get CSS file from CDN
// go run main.go -styles "https://unpkg.com/spectre.css/dist/spectre.min.css" foo.md > foo.html
// go run main.go -styles "//writ.cmcenroe.me/1.0.4/writ.min.css" foo.md > foo.html


import (
	"bytes"
	"flag"
	"fmt"
	"github.com/yuin/goldmark"
	"io/ioutil"
	"os"
	"strings"
)

var defaultExample = `
# CMS example
hello, world.
`

var docType = `
<!DOCTYPE html>
<html lang=`


func assemble(article string, title string, language string, styles []string) string {
  var htmlFile string
  var stylesheets string
  for _, sheet := range styles {
    s := fmt.Sprintf("\t<link rel=\"stylesheet\" href=\"%s\"/>\n", sheet)
    stylesheets += s
  }
  htmlFile = docType + "\"" + language + "\">" + "\n" +
    "<head>\n" +
    "\t<meta charset=\"utf-8\">\n" +
    "\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n"  +
    "\t<title>" + title + "</title>\n" +
    stylesheets  +
    "</head>\n<body>" +
    article + 
    "</body>\n</html>"
  fmt.Println(htmlFile)
  return htmlFile
}
func main() {
	var styles string
	flag.StringVar(&styles, "styles", "", "One or more stylesheets (use quotes if more than one)")

	var templs string
	flag.StringVar(&templs, "templates", "", "One or more templates (use quotes if more than one)")

	var title string
	flag.StringVar(&title, "title", "powered by microCMS", "Contents of the HTML title tag")

	var language string
	flag.StringVar(&language, "language", "en", "HTML language designation, such as en or fr")

	flag.Parse()
	filename := flag.Arg(0)

	stylesheets := strings.Split(styles, " ")
	//templates := strings.Split(templs, ", ")
	//fmt.Printf("Filename: %v\nStylesheets: %v\nTemplates: %v\nTitle: %v", filename, stylesheets, templates, title)

	if len(os.Args) < 2 {
		// No file was provided on the command line. Use defaultExample
		if HTML, err := mdToHTML([]byte(defaultExample)); err != nil {
			//quit(err.Error(), 1)
			quit(err, 1)
		} else {
			fmt.Println(string(HTML))
			quit(err, 0)
		}
	}

	if HTML, err := mdFileToHTML(filename); err != nil {
		quit(err, 1)
	} else {
    assemble(HTML, title, language, stylesheets)
		//fmt.Println(HTML)
		quit(err, 0)
	}



}

// mdToHTML takes Markdown source as a byte slice and converts it to HTML
// using Goldmark's default settings.
func mdToHTML(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(input, &buf); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// mdFileToHTML converts a source file to an HTML string
// using Goldmark's default settings.
func mdFileToHTML(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if HTML, err := mdToHTML(bytes); err != nil {
		return "", err
	} else {
		return string(HTML), nil
	}
}

func quit(err error, exitCode int) {
	if err != nil {
		//fmt.Printf("%v ", err.Error())
		fmt.Printf("%v ", err)
	}
	os.Exit(exitCode)
}
