package main

import (
	"net/http"
)

func (app *application) exampleHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}
	if !user.Activated {
		app.inactiveAccountResponse(w, r)
		return
	}
	// The rest of the handler logic goes here...
}
