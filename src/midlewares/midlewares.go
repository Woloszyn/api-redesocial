package midlewares

import (
	"api-redesocial/src/auth"
	"api-redesocial/src/respostas"
	"log"
	"net/http"
)

// Logger  escreve informações da rota no terminal
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %s %s %s", r.Method, r.RequestURI, r.Host)
		next(w, r)
	}
}

// Autenticar verifica se o usuário fazendo a requisição está autenticado
func Autenticar(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := auth.ValidarToken(r); erro != nil {
			respostas.Erro(w, http.StatusUnauthorized, erro)
			return
		}
		next(w, r)
	}
}
