package main

import (
  "flag"
  "html/template"
  "io"
  "log"
  "net/http"
  "sync"
)

var (
  addr = flag.String("addr", "127.0.0.1:8081", "addr to bind to")
)

const index = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
    <form method="POST" action="/">
      <input type="text" name="author" />
      <textarea name="message"></textarea>
      <input type="submit" />
    </form>

		{{range .Posts}}
      <p>
        <div><strong>{{ .Author }}</strong></div>
        <div>{{ .Message }}</div>
      </p>
    {{else}}
      <div><strong>no rows</strong></div>
    {{end}}
	</body>
</html>`

type Index struct{
  Title string
  Posts []Post
}

type Post struct{
  Author string
  Message string
}

type Server struct {
  Pages map[string]func(io.Writer) error

  mu sync.Mutex
  Posts []Post
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  log.Printf("got request %s", req.URL)
  if req.Method == "POST" {
    if err := req.ParseForm(); err != nil {
      res.WriteHeader(400)
      return
    }
    post := Post{
      Author: req.PostFormValue("author"),
      Message: req.PostFormValue("message"),
    }
    s.mu.Lock()
    s.Posts = append(s.Posts, post)
    s.mu.Unlock()
  }

  fn, ok := s.Pages[req.URL.Path]
  if !ok {
    res.WriteHeader(404)
    return
  }
  if err := fn(res); err != nil {
    res.WriteHeader(500)
  }
  //res.Write([]byte("hi there!"))
}

func main() {
  tmpl, err := template.New("index").Parse(index)
  if err != nil {
    log.Fatal(err)
  }

  s := Server{
    Pages: map[string]func(io.Writer) error {
      "/": func(w io.Writer) error {
        return tmpl.Execute(w, Index{
          Title: "My Cool Guetbook",
        })
      },
    },
  }

  log.Fatal(http.ListenAndServe(*addr, &s))
}
