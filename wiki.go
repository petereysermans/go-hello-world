package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil{
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

/*
func getTitle(w http.ResponseWriter, r *http.Request) (string, error ) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	
	if m = nil {
		http.NotFound(w, r)
		return "", error.New("Invalid Page Title")
	}

	return m[2], nil
}
*/

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, m[2])
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page){
	err := templates.ExecuteTemplate(w, tmpl + ".html", p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/*
	err = t.Execute(w, p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}*/
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)

	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
	/*
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, p)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	*/
}

func editHandler(w http.ResponseWriter, r *http.Request, title string){
	p, err := loadPage(title)
	
	if err != nil {
		p = &Page{Title: title }
	}

	renderTemplate(w, "edit", p)

	/*
	t, _ = template.ParseFiles("edit.html")
	t.Execute(w, p)

	fmt.Fprintf(w, 
		"<h1>Editing %s</h1" + 
		"<form action\"save/%s\" method=\"POST\">" +
		"<textarea name=\"body\">%s</textarea><br>" +
		"<input type=\"submit\" value=\"Save\">" + 
		"</form>",
		p.Title, p.Title, p.Body)
	*/
}

func main() {
	/*
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))
	*/

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}
