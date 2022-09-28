package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	connection.DatabaseConnect()
	route := mux.NewRouter()

	route.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")

	route.HandleFunc("/project", myProject).Methods("GET")
	route.HandleFunc("/project/{id}", myProjectDetail).Methods("GET")
	route.HandleFunc("/form-project", myProjectForm).Methods("GET")
	route.HandleFunc("/add-project", myProjectData).Methods("POST")
	route.HandleFunc("/form-edit-project/{id}", myProjectFormEditProject).Methods("GET")
	route.HandleFunc("/edit-project/{id}", myProjectEdited).Methods("POST")
	route.HandleFunc("/delete-project/{id}", myProjectDelete).Methods("GET")

	route.HandleFunc("/contact", contact).Methods(("GET"))

	fmt.Println("Server running at localhost port 8000")
	http.ListenAndServe("localhost:8000", route)
}

type StructInputDataForm struct {
	Id              int
	ProjectName     string
	StartDate       time.Time
	EndDate         time.Time
	StartDateFormat string
	EndDateFormat   string
	Description     string
	Techno          []string
	Duration        string
}

func home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("views/index.html")
	if err != nil {
		panic(err)
	}
	template.Execute(w, nil)
}

func myProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/myProject.html")

	var result []StructInputDataForm
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, projectname, startdate, enddate, description, technology FROM db_myprojects")
	for data.Next() {
		var each = StructInputDataForm{}
		err := data.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Description, &each.Techno)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		each.Duration = ""

		hour := 1
		day := hour * 24
		week := hour * 24 * 7
		month := hour * 24 * 30
		year := hour * 24 * 365
		differHour := each.EndDate.Sub(each.StartDate).Hours()
		var differHours int = int(differHour)
		days := differHours / day
		weeks := differHours / week
		months := differHours / month
		years := differHours / year
		if differHours < week {
			each.Duration = strconv.Itoa(int(days)) + " Days"
		} else if differHours < month {
			each.Duration = strconv.Itoa(int(weeks)) + " Weeks"
		} else if differHours < year {
			each.Duration = strconv.Itoa(int(months)) + " Months"
		} else if differHours > year {
			each.Duration = strconv.Itoa(int(years)) + " Years"
		}

		result = append(result, each)
	}

	response := map[string]interface{}{
		"Projects": result,
	}

	if err == nil {
		tmpl.Execute(w, response)
	} else {
		w.Write([]byte("Message: "))
		w.Write([]byte(err.Error()))
	}
}

func myProjectDetail(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/myProjectDetail.html")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	ProjectDetail := StructInputDataForm{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, projectname, startdate, enddate, description, technology FROM db_myprojects WHERE id=$1", id).Scan(
		&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Description, &ProjectDetail.Techno)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}
	ProjectDetail.StartDateFormat = ProjectDetail.StartDate.Format("2006-01-02")
	ProjectDetail.EndDateFormat = ProjectDetail.EndDate.Format("2006-01-02")
	ProjectDetail.Duration = ""

	hour := 1
	day := hour * 24
	week := hour * 24 * 7
	month := hour * 24 * 30
	year := hour * 24 * 365
	differHour := ProjectDetail.EndDate.Sub(ProjectDetail.StartDate).Hours()
	var differHours int = int(differHour)
	days := differHours / day
	weeks := differHours / week
	months := differHours / month
	years := differHours / year
	if differHours < week {
		ProjectDetail.Duration = strconv.Itoa(int(days)) + " Days"
	} else if differHours < month {
		ProjectDetail.Duration = strconv.Itoa(int(weeks)) + " Weeks"
	} else if differHours < year {
		ProjectDetail.Duration = strconv.Itoa(int(months)) + " Months"
	} else if differHours > year {
		ProjectDetail.Duration = strconv.Itoa(int(years)) + " Years"
	}

	response := map[string]interface{}{
		"Project": ProjectDetail,
	}

	if err == nil {
		tmpl.Execute(w, response)
	} else {
		panic(err)
	}
}

func myProjectForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/myProjectForm.html")
	if err == nil {
		tmpl.Execute(w, nil)
	} else {
		panic(err)
	}
}

func myProjectData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	// projectName := r.PostForm.Get("projectName")
	// startDate := r.PostForm.Get("startDate")
	// endDate := r.PostForm.Get("endDate")
	// description := r.PostForm.Get("description")
	var projectName string
	var startDate string
	var endDate string
	var description string
	var techno []string
	fmt.Println(r.Form)
	for i, values := range r.Form {
		fmt.Printf("type of values is %T\n", values)
		fmt.Println(values)
		fmt.Println(i)
		for _, value := range values {
			if i == "projectName" {
				projectName = value
			}
			if i == "startDate" {
				startDate = value
			}
			if i == "endDate" {
				endDate = value
			}
			if i == "description" {
				description = value
			}
			if i == "techno" {
				techno = append(techno, value)
				fmt.Printf("type of value is %T\n", value)
			}
		}
	}
	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO db_myprojects(projectname, startdate, enddate, description, technology) VALUES ($1, $2, $3, $4, $5)", projectName, startDate, endDate, description, techno)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}

func myProjectFormEditProject(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/myProjectFormEditProject.html")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	ProjectEdit := StructInputDataForm{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, projectname, startdate, enddate, description, technology FROM db_myprojects WHERE id=$1", id).Scan(
		&ProjectEdit.Id, &ProjectEdit.ProjectName, &ProjectEdit.StartDate, &ProjectEdit.EndDate, &ProjectEdit.Description, &ProjectEdit.Techno)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}
	ProjectEdit.StartDateFormat = ProjectEdit.StartDate.Format("2006-01-02")
	ProjectEdit.EndDateFormat = ProjectEdit.EndDate.Format("2006-01-02")

	response := map[string]interface{}{
		"Project": ProjectEdit,
	}

	if err == nil {
		tmpl.Execute(w, response)
	} else {
		panic(err)
	}
}

func myProjectEdited(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	var projectName string
	var startDate string
	var endDate string
	var description string
	var techno []string
	fmt.Println(r.Form)
	for i, values := range r.Form {
		for _, value := range values {
			if i == "projectName" {
				projectName = value
			}
			if i == "startDate" {
				startDate = value
			}
			if i == "endDate" {
				endDate = value
			}
			if i == "description" {
				description = value
			}
			if i == "techno" {
				techno = append(techno, value)
			}
		}
	}
	_, err = connection.Conn.Exec(context.Background(), "UPDATE db_myprojects SET projectname=$1, startdate=$2, enddate=$3, description=$4, technology=$5 WHERE id=$6", projectName, startDate, endDate, description, techno, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}

func myProjectDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM db_myprojects WHERE id=$1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}
	http.Redirect(w, r, "/project", http.StatusMovedPermanently)
}

func contact(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/contact.html")
	if err == nil {
		tmpl.Execute(w, nil)
	} else {
		panic(err)
	}
}
