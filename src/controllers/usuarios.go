package controllers

import (
	"api-redesocial/src/auth"
	"api-redesocial/src/database"
	"api-redesocial/src/models"
	"api-redesocial/src/repositorios"
	"api-redesocial/src/respostas"
	"api-redesocial/src/seguranca"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//CriarUsuario insere um usuário no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := ioutil.ReadAll(r.Body)

	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario

	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("cadastro"); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario.ID, erro = repositorio.Criar(usuario)

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.Json(w, http.StatusCreated, usuario)
}

//BuscarUsuarios retornará todos os usuários do banco de dados
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, erro := database.Conectar()

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	usuarios, erro := repositorio.Buscar(nomeOuNick)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.Json(w, http.StatusOK, usuarios)
}

//BuscarUsuario retornará um usuário do banco de dados
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)

	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	usuario, erro := repositorio.FindUsuarioById(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.Json(w, http.StatusOK, usuario)
}

//AlterarUsuario irá atualizar as informações do usuário
func AlterarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)

	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioToken, erro := auth.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar um usuário que não seja o logado"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario

	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	if erro = repositorio.Alterar(usuarioID, usuario); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.Json(w, http.StatusNoContent, nil)
}

//DeletarUsuario irá remover o usuário do banco de dados
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)

	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioToken, erro := auth.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioToken {
		respostas.Erro(w, http.StatusForbidden, errors.New("não é possível atualizar um usuário que não seja o logado"))
		return
	}

	db, erro := database.Conectar()

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	if erro = repositorio.Delete(usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	respostas.Json(w, http.StatusNoContent, nil)
}

func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	usuarioToken, erro := auth.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioToken == usuarioID {
		respostas.Erro(w, http.StatusBadRequest, errors.New("um usuário não pode seguir ele mesmo"))
		return
	}

	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro := repositorio.SeguirUsuario(usuarioToken, usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.Json(w, http.StatusNoContent, nil)
}

func DeixarDeSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	usuarioToken, erro := auth.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioToken == usuarioID {
		respostas.Erro(w, http.StatusBadRequest, errors.New("um usuário não pode deixar de seguir ele mesmo"))
		return
	}

	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro := repositorio.DeixarDeSeguir(usuarioToken, usuarioID); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.Json(w, http.StatusNoContent, nil)
}

func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, erro := database.Conectar()

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	usuarios, erro := repositorio.BuscarSeguidores(usuarioID, nomeOuNick)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.Json(w, http.StatusOK, usuarios)
}

func BuscarPessoasQueSigo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, erro := database.Conectar()

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	usuarios, erro := repositorio.BuscarPessoasQueSigo(usuarioID, nomeOuNick)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.Json(w, http.StatusOK, usuarios)
}

func AlterarSenha(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	usuarioToken, erro := auth.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioToken != usuarioID {
		respostas.Erro(w, http.StatusForbidden, errors.New("você não pode mudar a senha de outro usuário"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	var senha models.Senha
	if erro = json.Unmarshal(corpoRequisicao, &senha); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()

	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	senhaSalva, erro := repositorio.BuscarSenha(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = seguranca.VerificarSenha(senha.Atual, senhaSalva); erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, errors.New("a senha atual não condiz com o banco"))
		return
	}

	novaSenha, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = repositorio.SalvarNovaSenha(usuarioID, string(novaSenha)); erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	respostas.Json(w, http.StatusNoContent, nil)
}
