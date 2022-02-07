package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	conDB := "postgres://postgres:postgres@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", conDB)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	//TESTANDO CONEXÃO COM O BANCO
	errPing := db.Ping()
	if errPing != nil {
		log.Fatal(errPing.Error())
	}
	fmt.Println("Conectado no DB!")

	//CRIANDO TABELA CLIENTE
	_, errC := db.Exec(
		"CREATE TABLE IF NOT EXISTS cliente (" +
			"id SERIAL NOT NULL PRIMARY KEY," +
			"nome VARCHAR(50) NOT NULL," +
			"tipo VARCHAR(20) NOT NULL)")

	if errC != nil {
		log.Fatal(errC.Error())
	}
	fmt.Println("Tabela criada")

	//INSERIR DADOS NO BANCO
	_, errI := db.Exec(
		"INSERT INTO cliente (nome, tipo) VALUES " +
			"('João','Especial')," + "('Pedro','fisico')," + "('CCC','Juridico')")

	if errI != nil {
		log.Fatal(errI.Error())
	}
	fmt.Println("Cadastros prenchidos")

}
