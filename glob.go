package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/denisenkom/go-mssqldb"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type PostN struct {
	Title string
	Body  string
}

type Post struct {
	PostId  int32
	Title   string
	Body    string
	Created time.Time `db:"CreateDateTime"`
	Updated time.Time `db:"UpdateDateTime"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		// Serve default
	case "/posts":
		switch r.Method {
		case "GET":
			posts := getPosts()
			ps, _ := json.Marshal(posts)
			fmt.Fprint(w, string(ps))
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			checkErr(err, "Error reading request body")
			var p PostN
			err = json.Unmarshal(body, &p)
			checkErr(err, "Error unmarshalling request body")
			createPost(&p)
		default:
			fmt.Fprintf(w, "Method: %s is invalid for Route: %s", r.Method, r.URL.Path)
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func dbInit() *gorp.DbMap {
	db, err := sql.Open("mssql", "server=localhost;database=Glob;user id=sa;password=Grogan_01;log=63")
	checkErr(err, "sql.Open failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqlServerDialect{}}

	dbmap.AddTableWithName(Post{}, "Post").SetKeys(true, "PostId")

	return dbmap
}

func newPost(title, body string) Post {
	return Post{
		Title:   title,
		Body:    body,
		Created: time.Now().UTC(),
		Updated: time.Now().UTC(),
	}
}

func getPosts() (posts []Post) {
	dbmap := dbInit()
	defer dbmap.Db.Close()

	_, err := dbmap.Select(&posts, "SELECT * FROM [dbo].[Post] ORDER BY [CreateDateTime] DESC")
	checkErr(err, "Error in Sql")

	return posts
}

func createPost(p *PostN) {
	dbmap := dbInit()
	defer dbmap.Db.Close()

	post := newPost(p.Title, p.Body)

	err := dbmap.Insert(&post)
	checkErr(err, "Insert failed.")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
