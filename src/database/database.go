package database

import (
	"api-redesocial/src/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Driver
)

// Conectar abre a conexão com o bd e retorna
func Conectar() (*sql.DB, error) {
	db, erro := sql.Open("mysql", config.ConexaoBanco)

	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil

}
