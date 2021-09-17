package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	// String de conexão com o mysql
	ConexaoBanco = ""

	API_SECRET = []byte("")

	//Porta onde a API vai estar rodando
	Porta = 0
)

// Carregar vai inicializar as váriaveis de ambiente
func Carregar() {
	var erro error

	if erro = godotenv.Load(); erro != nil {
		log.Fatal(erro)
	}

	Porta, erro = strconv.Atoi(os.Getenv("API_PORT"))

	if erro != nil {
		Porta = 9000
	}

	ConexaoBanco = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USUARIO"),
		os.Getenv("DB_SENHA"),
		os.Getenv("DB_NOME"),
	)

	API_SECRET = []byte(os.Getenv("API_SECRET"))
}
