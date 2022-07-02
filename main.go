package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	bdSql      = "mysql"
	bdPort     = "127.0.0.1:3306"
	bdUser     = "root"
	bdPassword = "root"
)

type CatsOwner struct {
	Id                                                    int
	CatName, FluffyLVL, BirthDate, Url, Owner, PhoneOwner string
}

var allCats = []CatsOwner{}

// головна сторінка
func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html", "templates/header.html",
		"templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// підключуння до бази даних
	db, err := sql.Open(bdSql, fmt.Sprintf("%s:%s@tcp(%s)/cat_counting_db", bdUser, bdPassword, bdPort))
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	log.Print("Підключення до Бази Даних!")
	defer db.Close()

	// вибірка даних з  "cat_counting_db", табличка "cats"
	res, err := db.Query("SELECT * FROM `cats`")
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}

	allCats = []CatsOwner{}
	for res.Next() {
		var c CatsOwner
		err = res.Scan(&c.Id, &c.CatName, &c.BirthDate, &c.Url, &c.FluffyLVL, &c.Owner, &c.PhoneOwner)
		if err != nil {
			// panic(err)
			fmt.Fprintln(w, err.Error())
		}
		allCats = append(allCats, c)
	}

	err = tmpl.ExecuteTemplate(w, "home", allCats)
	if err != nil {
		// panic(err)
		fmt.Fprintln(w, err.Error())
	}
}

// вивід сторінки після нажаття кнопки "Добавити нового котика"
func create(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html",
		"templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "create", nil)
}

//вивід сторінки "Альбом"
func album(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/album.html", "templates/header.html",
		"templates/footer.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "album", nil)
}

// запис інформації з форми додавання котиків
func inf_save(w http.ResponseWriter, r *http.Request) {
	// провірка: дивимося чи відповідь є POST, якщо не POST, а наприклад GET
	// то значить відповідь пуста і щось не так.
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("GET-Метод запрещен!"))
		return
	}
	// відповід з "create.html"
	catName, birthDate, fluffyLVL := r.FormValue("catName"), r.FormValue("birthDate"), r.FormValue("fluffyLVL")
	url, owner, phoneOwner := r.FormValue("URL"), r.FormValue("owner"), r.FormValue("phoneOwner")

	// підключуння до бази даних
	db, err := sql.Open(bdSql, fmt.Sprintf("%s:%s@tcp(%s)/cat_counting_db", bdUser, bdPassword, bdPort))
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	log.Print("Підключення до Бази Даних!")
	defer db.Close()

	// додавання до бази даних інформаціїї
	insert, err := db.Query(fmt.Sprintf("INSERT INTO `cats` (`cat_name`, `birth_date`, `url_photo`, "+
		"`flyffy_level`, `cat_owner`, `phone_owner`) VALUES ('%s', '%s', '%s', '%s', '%s', '%s')",
		catName, birthDate, url, fluffyLVL, owner, phoneOwner))
	if err != nil {
		fmt.Fprintln(w, err.Error())
		// log.Fatal(err)
	}
	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//
func handleRequest() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/create/", create)
	mux.HandleFunc("/inf_save/", inf_save)
	mux.HandleFunc("/album/", album)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func main() {
	handleRequest()
}
