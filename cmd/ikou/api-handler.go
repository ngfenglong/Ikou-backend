package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) GetRecommenedPlaces(w http.ResponseWriter, r *http.Request) {
	place, err := app.DB.GetPlaceById("")

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	out, err := json.MarshalIndent(place, "", " ")

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
