package main

import (
	"html/template"
	"net/http"
	"regexp"
	"sync"
	"time"
)

type Student struct {
	FirstName string
	LastName  string
	Age       int
	Gender    string
}

type Class struct {
	Name         string
	Field        string
	Level        string
	StudentCount int
	StudentsList []Student
}

type ViewData struct {
	Count   int
	Message string
}

type UserData struct {
	LastName      string
	FirstName     string
	BirthDate     string
	Gender        string
	Age           int
	ErrorMessage  string
	IsFormSuccess bool
}

var (
	viewCount int
	class     Class
	mutex     sync.Mutex

	letterOnlyRegex = regexp.MustCompile("^[a-zA-ZÀ-ÿ\\s-]+$")
)

func init() {
	class = Class{
		Name:         "B1 Informatique",
		Field:        "Informatique",
		Level:        "Bachelor 1",
		StudentCount: 3,
		StudentsList: []Student{
			{FirstName: "Jean", LastName: "Dupont", Age: 20, Gender: "Masculin"},
			{FirstName: "Marie", LastName: "Martin", Age: 19, Gender: "Féminin"},
			{FirstName: "Pierre", LastName: "Gustav", Age: 21, Gender: "Masculin"},
		},
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/promo", promoHandler)
	http.HandleFunc("/change", changeHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/user/form", userFormHandler)
	http.HandleFunc("/user/treatment", userTreatmentHandler)
	http.HandleFunc("/user/display", userDisplayHandler)

	http.ListenAndServe(":8080", nil)
}

func promoHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	tmpl, err := template.ParseFiles("promo.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func changeHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	viewCount++
	currentCount := viewCount
	mutex.Unlock()

	var message string
	if currentCount%2 == 0 {
		message = "Le nombre de vues est pair"
	} else {
		message = "Le nombre de vues est impair"
	}

	data := ViewData{
		Count:   currentCount,
		Message: message,
	}

	tmpl, err := template.ParseFiles("change.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}

func userFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/user_form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func validateUserData(data *UserData) bool {
	if len(data.LastName) < 1 || len(data.LastName) > 32 || !letterOnlyRegex.MatchString(data.LastName) {
		return false
	}

	if len(data.FirstName) < 1 || len(data.FirstName) > 32 || !letterOnlyRegex.MatchString(data.FirstName) {
		return false
	}

	if data.Gender != "Masculin" && data.Gender != "Féminin" && data.Gender != "Autre" {
		return false
	}

	if data.BirthDate == "" {
		return false
	}

	return true
}

func userTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/user/form", http.StatusSeeOther)
		return
	}

	userData := &UserData{
		LastName:  r.FormValue("lastname"),
		FirstName: r.FormValue("firstname"),
		BirthDate: r.FormValue("birthdate"),
		Gender:    r.FormValue("gender"),
	}

	if validateUserData(userData) {
		userData.IsFormSuccess = true

		birthDate, _ := time.Parse("2006-01-02", userData.BirthDate)
		userData.Age = int(time.Since(birthDate).Hours() / 24 / 365)

		mutex.Lock()
		class.StudentsList = append(class.StudentsList, Student{
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
			Age:       userData.Age,
			Gender:    userData.Gender,
		})
		class.StudentCount++
		mutex.Unlock()

		http.Redirect(w, r, "/user/display", http.StatusSeeOther)
	} else {
		userData.ErrorMessage = "Données invalides. Veuillez vérifier vos informations."
		http.Redirect(w, r, "/user/form", http.StatusSeeOther)
	}
}

func userDisplayHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/user_display.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	lastStudent := class.StudentsList[len(class.StudentsList)-1]
	userData := UserData{
		LastName:      lastStudent.LastName,
		FirstName:     lastStudent.FirstName,
		Age:           lastStudent.Age,
		Gender:        lastStudent.Gender,
		IsFormSuccess: true,
	}

	tmpl.Execute(w, userData)
}