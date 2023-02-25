package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) GetRecommenedPlaces(w http.ResponseWriter, r *http.Request) {
	place, err := app.DB.GetPlaceById("26524f97-b2c1-11ed-ae52-0a0027000007")

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
