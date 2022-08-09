package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// package support middleware
func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Create a new middleware chain containing the middleware specific to
	// our dynamic application routes. For now, this chain will only contain
	// the session middleware but we'll add more to it later.
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	// Add the requireAuthentication middleware to the chain.
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	// Add the requireAuthentication middleware to the chain.
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// Add the five new routes for User.
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	// Add the requireAuthentication middleware to the chain.
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}

// without package support middleware

// mux.Get("/", app.session.Enable(http.HandlerFunc(app.home)))
// mux.Get("/snippet/create", app.session.Enable(app.requireAuthentication(http.HandlerFunc(app.createSnippetForm))))

// func (app *application) routes() http.Handler {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", app.home)
// 	mux.HandleFunc("/snippet", app.showSnippet)
// 	mux.HandleFunc("/snippet/create", app.createSnippet)
//
// 	fileServer := http.FileServer(http.Dir("./ui/static/"))
// 	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
//
// 	// Wrap the existing chain with the logRequest middleware.
// 	// Wrap the existing chain with the recoverPanic middleware.
// 	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
// }
