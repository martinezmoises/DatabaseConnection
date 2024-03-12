package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("/home/moises/Downloads/Test1/HomePage.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) createUserForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("/home/moises/Downloads/Test1/AddUsers.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/user", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Retrieve form values
	username := r.PostForm.Get("Username")
	firstName := r.PostForm.Get("FirstName")
	lastName := r.PostForm.Get("LastName")
	nickname := r.PostForm.Get("Nickname")
	sex := r.PostForm.Get("Sex")

	// Validate form fields
	errors := make(map[string]string)

	if strings.TrimSpace(username) == "" {
		errors["Username"] = "Username cannot be left blank"
	} else if utf8.RuneCountInString(username) > 50 {
		errors["Username"] = "Username is too long (maximum is 50 characters)"
	}

	if strings.TrimSpace(firstName) == "" {
		errors["FirstName"] = "First Name cannot be left blank"
	} else if utf8.RuneCountInString(firstName) > 50 {
		errors["FirstName"] = "First Name is too long (maximum is 50 characters)"
	}

	if strings.TrimSpace(lastName) == "" {
		errors["LastName"] = "Last Name cannot be left blank"
	} else if utf8.RuneCountInString(lastName) > 50 {
		errors["LastName"] = "Last Name is too long (maximum is 50 characters)"
	}

	if strings.TrimSpace(nickname) == "" {
		errors["Nickname"] = "Nickname cannot be left blank"
	} else if utf8.RuneCountInString(nickname) > 50 {
		errors["Nickname"] = "Nickname is too long (maximum is 50 characters)"
	}

	if strings.TrimSpace(sex) == "" {
		errors["Sex"] = "Sex cannot be left blank"
	}

	if len(errors) > 0 {
		// Handle validation errors (e.g., display error messages)
		fmt.Fprint(w, errors)
		return
	}

	// Insert data into the database
	s := `
        INSERT INTO users (Username, FirstName, LastName, Nickname, Sex)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err = app.db.Exec(s, username, firstName, lastName, nickname, sex)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Redirect to the home page after successful insertion
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) displayListings(w http.ResponseWriter, r *http.Request) {
	rows, err := app.db.Query("SELECT username, firstname, lastname, nickname, sex FROM users")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type User struct {
		Username  string
		FirstName string
		LastName  string
		NickName  string
		Sex       string
	}

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Username, &user.FirstName, &user.LastName, &user.NickName, &user.Sex)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		tmpl, err := template.New("no-data").Parse(noDataTemplate)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	tmpl, err := template.ParseFiles("listings.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, users)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

var noDataTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>No Data Found</title>
</head>
<body>
    <header>
        <h1>No Data Found</h1>
    </header>

    <main>
        <p>Nothing here to see.</p>
    </main>
</body>
</html>
`
