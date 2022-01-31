package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Cliente struct {
	Id   int    `json:"id"`
	Nome string `json:"nome"`
	Tipo string `json:"tipo"`
}

var Clientes []Cliente = []Cliente{
	{
		Id:   1,
		Nome: "Jose",
		Tipo: "Fisico",
	},
	{
		Id:   2,
		Nome: "CCC",
		Tipo: "Juridico",
	},
	{
		Id:   3,
		Nome: "Pedro",
		Tipo: "Especial",
	},
}

func rotaMain(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Bem vindo ao server http")
}

func listarClientes(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(rw)
	encoder.Encode(Clientes)
}
func cadastrarClientes(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	c, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal()
	}
	var novoCliente Cliente
	json.Unmarshal(c, &novoCliente)
	novoCliente.Id = len(Clientes) + 1
	Clientes = append(Clientes, novoCliente)

	encoder := json.NewEncoder(rw)
	encoder.Encode(novoCliente)
	//rw.WriteHeader(http.StatusCreated)
}
func buscarClientes(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	v := mux.Vars(r)
	id, _ := strconv.Atoi(v["id"])

	for _, C := range Clientes {
		if C.Id == id {
			json.NewEncoder(rw).Encode(C)
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

func deleteClientes(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	v := mux.Vars(r)
	id, _ := strconv.Atoi(v["id"])

	for i, C := range Clientes {
		if C.Id == id {
			Clientes = append(Clientes[0:i], Clientes[i+1:]...)
			rw.WriteHeader(http.StatusNoContent)
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

func editarClientes(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	v := mux.Vars(r)
	id, _ := strconv.Atoi(v["id"])

	c, _ := ioutil.ReadAll(r.Body)

	var Ecliente Cliente
	json.Unmarshal(c, &Ecliente)

	for i, C := range Clientes {
		if C.Id == id {
			Clientes[i] = Ecliente
			json.NewEncoder(rw).Encode(Ecliente)
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
}

func rotaConfig(rota *mux.Router) {

	rota.HandleFunc("/", rotaMain)
	rota.HandleFunc("/clientes", listarClientes).Methods("GET")
	rota.HandleFunc("/clientes/{id}", buscarClientes).Methods("GET")
	rota.HandleFunc("/clientes/", cadastrarClientes).Methods("POST")
	rota.HandleFunc("/clientes/{id}", deleteClientes).Methods("DELETE")
	rota.HandleFunc("/clientes/{id}", editarClientes).Methods("PUT")
}

func serverConfig() {
	rota := mux.NewRouter().StrictSlash(true)
	rotaConfig(rota)
	fmt.Println("Servidor esta rodando na porta 4553.")
	log.Fatal(http.ListenAndServe(":4553", rota))
}

func main() {
	serverConfig()
}
