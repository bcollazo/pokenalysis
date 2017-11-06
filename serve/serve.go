package serve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bcollazo/pokenalysis/poke"
	"io"
	"net/http"
)

type JobRequest struct {
	Ids       []int  `json:"ids"`
	RespondTo string `json:"respond_to"`
}

type JobResponse struct {
	Results []poke.BestMoveSetResult `json:"results"`
}

var resultsChannel chan poke.BestMoveSetResult
var srv *http.Server

func StartMaster(ids []int, host string, port string, machines []string) {
	fmt.Println("Collecting results on " + port + "...")
	srv := &http.Server{Addr: port}
	http.HandleFunc("/", masterHandler)
	go func() { srv.ListenAndServe() }()

	go sendWork(ids, host+port, machines)

	results := []poke.BestMoveSetResult{}
	resultsChannel = make(chan poke.BestMoveSetResult, len(ids))
	for _, _ = range ids {
		r := <-resultsChannel
		results = append(results, r)
		poke.PrintBattlePokemon(r.PokemonName, r.MoveSet)
	}
	srv.Shutdown(nil)
}

func masterHandler(w http.ResponseWriter, r *http.Request) {
	var result JobResponse
	json.NewDecoder(r.Body).Decode(&result)
	defer r.Body.Close()

	// Keep collecting results.
	for _, r := range result.Results {
		resultsChannel <- r
	}

	io.WriteString(w, http.StatusText(200))
}

func sendWork(ids []int, respondTo string, machines []string) {
	subRanges := poke.DivideWork(ids, len(machines))

	for i, m := range machines {
		job := JobRequest{subRanges[i], respondTo}
		jsonValue, _ := json.Marshal(job)

		fmt.Println("Sending work to " + m)
		http.Post("http://"+m, "application/json", bytes.NewBuffer(jsonValue))
	}
	fmt.Println("Done sending work...")
}

// Assumes has all data locally.
func StartWorker(port string) {
	fmt.Println("Serving on " + port + "...")

	http.HandleFunc("/", serveHandler)
	http.ListenAndServe(port, nil)
}

func serveHandler(w http.ResponseWriter, r *http.Request) {
	var job JobRequest
	json.NewDecoder(r.Body).Decode(&job)
	defer r.Body.Close()

	go do(job)
	io.WriteString(w, http.StatusText(200))
}

func do(job JobRequest) {
	list := poke.ReadDataFromLocal(job.Ids)
	results := poke.BestPokemons(list, 0) // doesn't matter order.

	// POST JobResponse in job.RespondTo.  Ignore any errors.
	response := JobResponse{results}
	jsonValue, _ := json.Marshal(response)
	fmt.Println("Sending results to " + job.RespondTo)
	http.Post("http://"+job.RespondTo, "application/json", bytes.NewBuffer(jsonValue))
}
