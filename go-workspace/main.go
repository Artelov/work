package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "pooja"
	password = "pooja"
	dbname   = "hpe"
)

var tpl *template.Template
var db *sql.DB
var err error

func init() {
	tpl = template.Must(template.ParseGlob("Templates/*"))
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	http.Handle("/public/CSS/", http.StripPrefix("/public/CSS", http.FileServer(http.Dir("./public/CSS"))))
	http.Handle("/public/img/", http.StripPrefix("/public/img", http.FileServer(http.Dir("./public/img"))))
	http.HandleFunc("/", dashboard)
	http.HandleFunc("/index", index)
	http.HandleFunc("/addAllergy", addAllergy)
	http.HandleFunc("/addHospital", addHospital)
	http.HandleFunc("/addMedication", addMedication)
	http.HandleFunc("/addSymptom", addSymptom)
	http.HandleFunc("/addImmu", addImmu)
	http.HandleFunc("/addImmu/process", addImmuprocess)
	http.HandleFunc("/contactus", contactus)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/signup/process", signupprocess)
	http.HandleFunc("/login/process", loginprocess)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/tutorial", tutorial)
	http.ListenAndServe(":200", nil)
}

type Users struct {
	Username string
	Password string
	Email    string
}

func logout(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	fmt.Println(err)
}
func index(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	fmt.Println(err)
}
func signupprocess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	bk := Users{}
	bk.Username = r.FormValue("username")
	bk.Password = r.FormValue("password")
	bk.Email = r.FormValue("email")

	_, err = db.Exec("INSERT INTO users(username, password,email) VALUES($1,$2,$3)", bk.Username, bk.Password, bk.Email)
	if err != nil {
		http.Error(w, "Server error, unable to create your account.", 500)
		return
	}
	http.Redirect(w, r, "/dashboard", 301)
	w.Write([]byte("User created!"))
	return

}

type Immu struct {
	Vaccine    string
	Protection string
	Date       string
	Note       string
	file_name  string
	file_data  []byte
}

func addImmuprocess(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		bk := Immu{}
		bk.Vaccine = r.FormValue("vaccine")
		bk.Protection = r.FormValue("protection")
		bk.Date = r.FormValue("date")
		bk.Note = r.FormValue("note")
		bk.file_name = r.FormValue("file")
		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)
		var size int64 = handler.Size
		bytes := make([]byte, size)
		buffer := bufio.NewReader(file)
		_, err = buffer.Read(bytes)
		fmt.Println("printing bytes")
		fmt.Println(bytes)

		_, err = db.Exec("INSERT INTO Immu(Vaccine, Protection,Date,Note,file_name,file_data) VALUES($1,$2,$3,$4,$5,$6)", bk.Vaccine, bk.Protection, bk.Date, bk.Note, bk.file_name, bytes)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Server error, unable to create your account.", 500)
			return
		}

	}
	http.Redirect(w, r, "/", 301)
	w.Write([]byte("Done uploaded to your activity"))
}
func loginprocess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=$1 AND password=$2", username, password).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(w, r, "/login", 301)
		return
	}

	http.Redirect(w, r, "/dashboard", 301)

}

func login(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	fmt.Println(err)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "dashboard.gohtml", nil)
	fmt.Println(err)
}
func addAllergy(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "addAllergy.gohtml", nil)
	fmt.Println(err)
}

func addHospital(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "addHospital.gohtml", nil)
	fmt.Println(err)
}
func addImmu(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "addImmu.gohtml", nil)
	fmt.Println(err)
}

func addMedication(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "addMedication.gohtml", nil)
	fmt.Println(err)
}
func addSymptom(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "addSymptom.gohtml", nil)
	fmt.Println(err)
}

func signup(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "signup.gohtml", nil)
	fmt.Println(err)

}

func tutorial(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "tutorial.gohtml", nil)
	fmt.Println(err)
}
func contactus(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "contactus.gohtml", nil)
	fmt.Println(err)
}
