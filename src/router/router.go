package router

import (
	"api-redesocial/src/router/rotas"

	"github.com/gorilla/mux"
)

//Vai retornar um router com as rotas configuradas
func Gerar() *mux.Router {
	r := mux.NewRouter()
	return rotas.Configurar(r)
}
