package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Cliente struct {
	Id   int    `json:"id"`
	Nome string `json:"nome"`
	Tipo string `json:"tipo"`
}

var db *sql.DB

func rotaMain(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Bem vindo ao server http")
}

func listarClientes(rw http.ResponseWriter, r *http.Request) {
	reg, err := db.Query("SELECT id, nome, tipo FROM cliente")

	if err != nil {
		log.Println("Listar Clientes: " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var Clientes []Cliente
	for reg.Next() {
		var C Cliente
		errscan := reg.Scan(&C.Id, &C.Nome, &C.Tipo)
		if errscan != nil {
			log.Println("Listar Clientes: " + errscan.Error())
			continue
		}

		Clientes = append(Clientes, C)
	}

	encoder := json.NewEncoder(rw)
	encoder.Encode(Clientes)
}

func cadastrarClientes(rw http.ResponseWriter, r *http.Request) {

	c, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var novoCliente Cliente
	json.Unmarshal(c, &novoCliente)

	//verificando se existe informação para inserir
	if len(novoCliente.Nome) == 0 || len(novoCliente.Tipo) == 0 {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(rw).Encode("Nome e Tipo Obrigatorios!")
		return
	}

	result, inserterro := db.Exec("INSERT INTO cliente (nome, tipo) VALUES ($1, $2)", novoCliente.Nome, novoCliente.Tipo)

	//pegando o ultimo id gerado pelo banco
	idgerado, lasterr := result.LastInsertId()

	//verificando erro de inserção e do lasinsertid
	if inserterro != nil || lasterr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	novoCliente.Id = int(idgerado)

	//respondendo para cliente satus e informação inserida
	rw.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(rw)
	encoder.Encode(novoCliente)
}
func buscarClientes(rw http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id, _ := strconv.Atoi(v["id"])

	var C Cliente
	reg := db.QueryRow("SELECT id, nome, tipo FROM cliente WHERE id = $1", id)
	errscan := reg.Scan(&C.Id, &C.Nome, &C.Tipo)

	if errscan != nil {
		log.Println("Buscas Cliente: " + errscan.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(rw).Encode(C)

}

func deleteClientes(rw http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id, _ := strconv.Atoi(v["id"])

	reg := db.QueryRow("SELECT id FROM cliente WHERE id = $1", id)
	var c Cliente
	erroscan := reg.Scan(&c.Id)

	if erroscan != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	_, erroExec := db.Exec("DELETE FROM cliente WHERE id = $1", id)

	if erroExec != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)

}

func editarClientes(rw http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id, _ := strconv.Atoi(v["id"])

	e, _ := ioutil.ReadAll(r.Body)

	var eCliente Cliente
	json.Unmarshal(e, &eCliente)

	reg := db.QueryRow("SELECT id, nome, tipo FROM cliente WHERE id = $1", id)
	var C Cliente
	erroscan := reg.Scan(&C.Id, &C.Nome, &C.Tipo)

	if erroscan != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	_, erroExec := db.Exec("UPDATE cliente SET nome = $1, tipo = $2 WHERE id = $3", eCliente.Nome, eCliente.Tipo, id)

	if erroExec != nil {
		log.Println("modificando livro: " + erroExec.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(eCliente)

}

func rotaConfig(rota *mux.Router) {

	rota.HandleFunc("/", rotaMain)
	rota.HandleFunc("/clientes", listarClientes).Methods("GET")
	rota.HandleFunc("/clientes/{id}", buscarClientes).Methods("GET")
	rota.HandleFunc("/clientes/", cadastrarClientes).Methods("POST")
	rota.HandleFunc("/clientes/{id}", deleteClientes).Methods("DELETE")
	rota.HandleFunc("/clientes/{id}", editarClientes).Methods("PUT")
}

func jsonMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func serverConfig() {
	rota := mux.NewRouter().StrictSlash(true)
	rota.Use(jsonMW)
	rotaConfig(rota)
	fmt.Println("Servidor esta rodando na porta 4553.")
	log.Fatal(http.ListenAndServe(":4553", rota))
}

func dbConnect() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")

	if err != nil {
		log.Fatal(err.Error())
	}

	errPing := db.Ping()
	if errPing != nil {
		log.Fatal(errPing.Error())
	}
}

func main() {
	dbConnect()
	serverConfig()
}
