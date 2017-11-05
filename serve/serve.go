package serve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bcollazo/pokenalysis/poke"
	"net/http"
	"strconv"
)

func Serve(list []poke.Pokemon, port int) {
	var stringPort = ":" + strconv.Itoa(port)
	fmt.Println("Serving on " + stringPort + "...")

	http.HandleFunc("/serve", handler)
	http.ListenAndServe(stringPort, nil)
}

type Job struct {
	Ids       []int  `json:ids`
	RespondTo string `json:respond_to`
}

type JobResponse struct {
	Results []poke.BestMoveSetResult `json:results`
}

func do(job Job) {
	list := poke.ReadDataFromLocal(job.Ids)
	res := poke.BestPokemons(list, 0) // doesn't matter order.

	// POST JobResponse in job.RespondTo
	jsonValue, _ := json.Marshal(JobResponse{res})
	http.Post(job.RespondTo, "application/json", bytes.NewBuffer(jsonValue))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var job Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// TODO: Go do work...
	go do(job)

	res := struct {
		Status string `json:status`
	}{"OK"}
	js, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
