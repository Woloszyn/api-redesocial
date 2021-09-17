package models

import (
	"api-redesocial/src/seguranca"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

//Usuário representa um usuário utilizando a rede social
type Usuario struct {
	ID       uint64    `json:"id,omitempty"`
	Nome     string    `json:"nome,omitempty"`
	Nick     string    `json:"nick,omitempty"`
	Email    string    `json:"email,omitempty"`
	Senha    string    `json:"senha,omitempty"`
	CriadoEm time.Time `json:"criadoEm,omitempty"`
}

// Preparar vai chamar os métodos de validar e preparar
func (usuario *Usuario) Preparar(etapa string) error {
	if erro := usuario.validar(etapa); erro != nil {
		return erro
	}
	if erro := usuario.formatar(etapa); erro != nil {
		return erro
	}
	return nil
}

func (usuario *Usuario) validar(etapa string) error {
	if usuario.Nome == "" {
		return errors.New("o nome é obrigatório e não pode estar em branco")
	}
	if usuario.Email == "" {
		return errors.New("o e-mail é obrigatório e não pode estar em branco")
	}

	if erro := checkmail.ValidateFormat(usuario.Email); erro != nil {
		return errors.New("o e-mail está em um formato inválido")
	}
	if usuario.Nick == "" {
		return errors.New("o nick é obrigatório e não pode estar em branco")
	}
	if usuario.Senha == "" && etapa == "cadastro" {
		return errors.New("o senha é obrigatória e não pode estar em branco")
	}
	return nil
}

func (usuario *Usuario) formatar(etapa string) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Email = strings.TrimSpace(usuario.Email)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	if etapa == "cadastro" {
		senhaComHash, erro := seguranca.Hash(usuario.Senha)
		if erro != nil {
			return erro
		}

		usuario.Senha = string(senhaComHash)
	}
	return nil
}
