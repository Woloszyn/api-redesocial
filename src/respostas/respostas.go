package respostas

import (
	"encoding/json"
	"log"
	"net/http"
)

// Json retorna uma resposta em json para a requisição
func Json(w http.ResponseWriter, statusCode int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if dados != nil {
		if erro := json.NewEncoder(w).Encode(dados); erro != nil {
			log.Fatal(erro)
		}
	}
}

func Erro(w http.ResponseWriter, statusCode int, erro error) {
	Json(w, statusCode, struct {
		Erro string `json:"erro"`
	}{
		Erro: erro.Error(),
	})
}
