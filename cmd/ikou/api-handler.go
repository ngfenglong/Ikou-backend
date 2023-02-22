package main

import "net/http"

func (app *application) GetRecommenedPlaces(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Printf("It is printing")
}
