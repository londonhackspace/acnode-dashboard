package main

import (
	"github.com/londonhackspace/acnode-dashboard/auth"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
)

type LoginTemplateArgs struct {
	BaseTemplateArgs
	Error string
	Next  string
}

var loginTemplate *template.Template = nil

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if loginTemplate == nil {
		loginTemplate = getTemplate("login.gohtml")
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Err(err).Msg("")
			w.WriteHeader(500)
		}

		user := r.FormValue("username")
		password := r.FormValue("password")
		next := r.FormValue("next")

		if len(user) == 0 || len(password) == 0 {
			args := LoginTemplateArgs{
				BaseTemplateArgs: GetBaseTemplateArgs(),
				Error:            "Please specify a valid username and password!",
				Next:             next,
			}
			loginTemplate.ExecuteTemplate(w, "login.gohtml", args)
		} else {
			if auth.AuthenticateUser(w, user, password) {
				http.Redirect(w, r, next, 302)
			} else {
				args := LoginTemplateArgs{
					BaseTemplateArgs: GetBaseTemplateArgs(),
					Error:            "Invalid Credentials",
					Next:             next,
				}
				loginTemplate.ExecuteTemplate(w, "login.gohtml", args)
			}
		}
	} else {
		nextArgs, ok := r.URL.Query()["next"]
		var next string = "/"
		if ok {
			next = nextArgs[0]
		}

		// if the user is already logged in, just redirect
		if ok, _ := auth.CheckAuthUser(w, r); ok {
			http.Redirect(w, r, next, 302)
			return
		}

		args := LoginTemplateArgs{
			BaseTemplateArgs: GetBaseTemplateArgs(),
			Error:            "",
			Next:             next,
		}

		loginTemplate.ExecuteTemplate(w, "login.gohtml", args)
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	auth.Logout(w, r)

	nextArgs, ok := r.URL.Query()["next"]
	var next string = "/"
	if ok {
		next = nextArgs[0]
	}

	http.Redirect(w, r, next, 302)
}
