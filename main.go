package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	encoder := json.NewEncoder(rw)
	encoder.Encode(Clientes)
}
func cadastrarClientes(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusCreated)

	l, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal()
	}
	var novoCliente Cliente
	json.Unmarshal(l, &novoCliente)
	novoCliente.Id = len(Clientes) + 1
	Clientes = append(Clientes, novoCliente)

	encoder := json.NewEncoder(rw)
	encoder.Encode(novoCliente)
}
func buscarClientes(rw http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(p[2])
	for _, C := range Clientes {
		if C.Id == id {
			json.NewEncoder(rw).Encode(C)
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
	return
}

func deleteClientes(rw http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(p[2])
	for i, C := range Clientes {
		if C.Id == id {
			Clientes = append(Clientes[0:i], Clientes[i+1:]...)
			rw.WriteHeader(http.StatusNoContent)
			return
		}
	}
	rw.WriteHeader(http.StatusNotFound)
	return
}

func rotasClientes(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	p := strings.Split(r.URL.Path, "/")

	switch {
	case r.Method == "GET" && len(p) == 2 || len(p) == 3 && p[2] == "":
		listarClientes(rw, r)
	case r.Method == "POST":
		cadastrarClientes(rw, r)
	case r.Method == "GET" && len(p) == 3 && p[2] != "" || len(p) == 4 && p[3] == "":
		buscarClientes(rw, r)
	case r.Method == "DELETE":
		deleteClientes(rw, r)
	}
}

func rotaConfig() {
	http.HandleFunc("/", rotaMain)
	http.HandleFunc("/clientes", rotasClientes)
	http.HandleFunc("/clientes/", rotasClientes)
}

func serverConfig() {
	rotaConfig()
	fmt.Println("Servidor esta rodando na porta 4553.")
	log.Fatal(http.ListenAndServe(":4553", nil))
}

func main() {
	serverConfig()
}
