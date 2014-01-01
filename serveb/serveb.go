// Example serveb serves bundled data over HTTP.
//
// Usage is:
//
//    serveb <laddr>
//
// Where "<laddr>" is the TCP local network address to listen for HTTP
// connections to. Example:
//
//    serveb :8080
//
package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
)

var index_tmpl = `
<html>
<head><title>Bundle contents</title></head>
<body>

<h1>Bundle contents:</h1>

{{range .}}<div><a href="{{.Name}}">{{.Name}}</a></div>
{{end}}
</body>
`

var tmpl *template.Template

func Usage(cmd string) {
	fmt.Fprintf(os.Stderr, "Usage is: %s <local addr>\n", cmd)
}

func show_index(w http.ResponseWriter) {
	dir := _bundleIdx.Dir("")
	tmpl.Execute(w, dir)
}

func serve_entry(w http.ResponseWriter, name string) {
	e := _bundleIdx.Entry(name)
	if e == nil {
		http.Error(w,
			"404 entry not found: "+name,
			http.StatusNotFound)
		return
	}
	rb, err := e.Open(0)
	if err != nil {
		http.Error(w,
			"Internal Error: "+err.Error(),
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, rb)
	_ = rb.Close()
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		show_index(w)
	} else {
		serve_entry(w, r.URL.Path[1:])
	}
}

func main() {
	var err error

	if len(os.Args) != 2 {
		Usage(path.Base(os.Args[0]))
		os.Exit(1)
	}
	tmpl = template.Must(template.New("index").Parse(index_tmpl))
	http.HandleFunc("/", handler)
	err = http.ListenAndServe(os.Args[1], nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
