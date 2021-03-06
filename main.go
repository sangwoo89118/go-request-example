package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// initialize data
type person struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Address *address `json:"address,omitempty"`
}
type address struct {
	City  string `json:"city"`
	State string `json:"state"`
}

var people []person

// ****************** //
// ****GET METHOD**** //
// ****************** //
func homePage(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	welcomeMessage := `
		<div style="display:grid;height:75vh">
			<h1 style="text-align:center; margin:auto"> Welcome to SQL ANGELES </h1>
		</div>`

	fmt.Fprintf(w, welcomeMessage)
}

func getPeople(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	json.NewEncoder(w).Encode(people)
}

func getPerson(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	personID := ps.ByName("id")
	for _, item := range people {
		if item.ID == personID {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	fmt.Fprintf(w, "<h1>No DATA</h1>")
}

// ******************* //
// ****POST METHOD**** //
// ******************* //
func createPerson(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	newID := ps.ByName("id")
	var newPerson person
	json.NewDecoder(r.Body).Decode(&newPerson)
	newPerson.ID = string(newID)
	people = append(people, newPerson)
	json.NewEncoder(w).Encode(people)
}

// ********************* //
// ****DELETE METHOD**** //
// ********************* //
func deletePerson(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	personID := ps.ByName("id")
	found := false
	for index, item := range people {
		if item.ID == personID {
			people = append(people[:index], people[index+1:]...)
			found = true
		}
	}

	if !found {
		fmt.Fprintf(w, "No person with ID of "+personID)
		return
	}
	json.NewEncoder(w).Encode(people)
}

// ****************** //
// ****PUT METHOD**** //
// ****************** //
func updatePerson(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	personID := ps.ByName("id")
	found := false

	for index, item := range people {
		if item.ID == personID {
			updatedPerson := &people[index]
			updatedPerson.ID = personID
			json.NewDecoder(r.Body).Decode(&updatedPerson)
			json.NewEncoder(w).Encode(updatedPerson)
			found = true
		}
	}

	if !found {
		fmt.Fprintf(w, "No person with ID of "+personID)
	}

}

// main func
func main() {
	// initialize router
	router := httprouter.New()

	// add data
	people = append(people, person{ID: "1", Name: "Sangwoo", Age: 29, Address: &address{City: "Los Angeles", State: "CA"}})
	people = append(people, person{ID: "2", Name: "Paul", Age: 28, Address: &address{City: "Irvine", State: "CA"}})

	// routers
	router.GET("/", homePage)
	router.GET("/people", getPeople)
	router.GET("/people/:id", getPerson)
	router.POST("/people/:id", createPerson)
	router.DELETE("/people/:id", deletePerson)
	router.PUT("/people/:id", updatePerson)

	log.Println("server is running on port localhost:8000")

	// listens port 8000 and add router
	http.ListenAndServe(":8000", router)
}
