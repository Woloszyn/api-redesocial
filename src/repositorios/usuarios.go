package repositorios

import (
	"api-redesocial/src/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Usuarios representa um repositório de usuários
type Usuarios struct {
	db *sql.DB
}

// NovoRepositorioDeUsuarios cria um repositório de usuários
func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

//Criar insere um usuário no bd
func (repositorio Usuarios) Criar(usuario models.Usuario) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO usuarios(nome, nick, email, senha) VALUES(?, ?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}

	defer statement.Close()

	resultado, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()

	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

// Buscar traz todos os usuários que atendem a um nome ou nick
func (repositorio Usuarios) Buscar(nickOuNome string) ([]models.Usuario, error) {
	nickOuNome = fmt.Sprintf("%%%s%%", nickOuNome)
	linhas, erro := repositorio.db.Query("SELECT id, nome, nick, email, criadoEm FROM usuarios WHERE nome LIKE ? OR nick LIKE ?", nickOuNome, nickOuNome)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var usuarios []models.Usuario

	for linhas.Next() {
		var usuario models.Usuario
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (repositorio Usuarios) FindUsuarioById(usuarioID uint64) (models.Usuario, error) {
	linhas, erro := repositorio.db.Query("SELECT id, nome, nick, email, criadoEm FROM usuarios WHERE id = ?", usuarioID)
	if erro != nil {
		return models.Usuario{}, erro
	}
	defer linhas.Close()

	var usuario models.Usuario
	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nick,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return models.Usuario{}, erro
		}
	}
	return usuario, nil
}

func (repositorio Usuarios) Alterar(usuarioID uint64, usuario models.Usuario) error {
	statement, erro := repositorio.db.Prepare("UPDATE usuarios SET nome = ?, email = ?, nick = ? WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(usuario.Nome, usuario.Email, usuario.Nick, usuarioID); erro != nil {
		return erro
	}
	return nil
}

func (repositorio Usuarios) Delete(usuarioID uint64) error {
	statement, erro := repositorio.db.Prepare("DELETE FROM usuarios WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(usuarioID); erro != nil {
		return erro
	}
	return nil
}

// BuscarPorEmail busca um usuário por email e retorna o seu id e senha com hash
func (repositorio Usuarios) BuscarPorEmail(email string) (models.Usuario, error) {
	linha, erro := repositorio.db.Query("SELECT id, senha from usuarios where email = ?", email)
	if erro != nil {
		return models.Usuario{}, erro
	}
	defer linha.Close()
	var usuario models.Usuario
	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Senha); erro != nil {
			return models.Usuario{}, erro
		}
	}

	return usuario, nil
}

func (repositorio Usuarios) DeixarDeSeguir(seguidorID uint64, usuarioID uint64) error {
	ok, erro := repositorio.verificaJaSegueUsuario(seguidorID, usuarioID)
	if erro != nil {
		return erro
	}

	if ok {
		statement, erro := repositorio.db.Prepare("DELETE FROM seguidores WHERE usuario_id = ? AND seguidor_id = ?")
		if erro != nil {
			return erro
		}
		defer statement.Close()
		if _, erro = statement.Exec(usuarioID, seguidorID); erro != nil {
			return erro
		}
		return nil
	}
	return errors.New("você não segue este usuário")
}

func (repositorio Usuarios) SeguirUsuario(seguidorID uint64, usuarioID uint64) error {
	ok, erro := repositorio.verificaJaSegueUsuario(seguidorID, usuarioID)
	if erro != nil {
		return erro
	}

	if !ok {
		statement, erro := repositorio.db.Prepare("INSERT INTO seguidores (usuario_id, seguidor_id) VALUES (?,?)")
		if erro != nil {
			return erro
		}
		defer statement.Close()
		if _, erro = statement.Exec(usuarioID, seguidorID); erro != nil {
			return erro
		}
		return nil
	}
	return errors.New("você já segue o usuário informado")
}

//verificaJaSegueUsuario faz uma busca no banco de dados e verifica se o usuário já é seguido
func (repositorio Usuarios) verificaJaSegueUsuario(seguidorID uint64, usuarioID uint64) (bool, error) {
	linha, erro := repositorio.db.Query("SELECT COUNT(1) as segue from seguidores where usuario_id = ? AND seguidor_id = ?", usuarioID, seguidorID)
	if erro != nil {
		return false, erro
	}
	defer linha.Close()
	total := 0
	if linha.Next() {
		if erro = linha.Scan(&total); erro != nil {
			return false, erro
		}
	}
	return total > 0, nil
}

func (repositorio Usuarios) BuscarSeguidores(usuarioID uint64, nomeOuNick string) ([]models.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", strings.ToLower(nomeOuNick))
	linhas, erro := repositorio.db.Query(`
			SELECT usuarios.id, usuarios.nome, usuarios.nick, usuarios.criadoEm FROM seguidores
			INNER JOIN usuarios ON (seguidores.seguidor_id = usuarios.id)
			WHERE seguidores.usuario_id = ?
			AND (LOWER(usuarios.nome) LIKE ? OR LOWER(usuarios.nick) LIKE ?)
		`, usuarioID, nomeOuNick, nomeOuNick)
	if erro != nil {
		return []models.Usuario{}, erro
	}

	var seguidores []models.Usuario
	for linhas.Next() {
		var seguidor models.Usuario
		if erro := linhas.Scan(&seguidor.ID, &seguidor.Nome, &seguidor.Nick, &seguidor.CriadoEm); erro != nil {
			return []models.Usuario{}, erro
		}
		seguidores = append(seguidores, seguidor)
	}

	return seguidores, nil
}

func (repositorio Usuarios) BuscarPessoasQueSigo(usuarioID uint64, nomeOuNick string) ([]models.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", strings.ToLower(nomeOuNick))
	linhas, erro := repositorio.db.Query(`
			SELECT usuarios.id, usuarios.nome, usuarios.nick, usuarios.criadoEm FROM seguidores
			INNER JOIN usuarios ON (seguidores.usuario_id = usuarios.id)
			WHERE seguidores.seguidor_id = ?
			AND (LOWER(usuarios.nome) LIKE ? OR LOWER(usuarios.nick) LIKE ?)
		`, usuarioID, nomeOuNick, nomeOuNick)
	if erro != nil {
		return []models.Usuario{}, erro
	}

	var seguidores []models.Usuario
	for linhas.Next() {
		var seguidor models.Usuario
		if erro := linhas.Scan(&seguidor.ID, &seguidor.Nome, &seguidor.Nick, &seguidor.CriadoEm); erro != nil {
			return []models.Usuario{}, erro
		}
		seguidores = append(seguidores, seguidor)
	}

	return seguidores, nil
}

//BuscarSenha trás a senha deste usuario por id
func (repositorio Usuarios) BuscarSenha(usuario_id uint64) (string, error) {
	linha, erro := repositorio.db.Query("SELECT senha FROM usuarios WHERE id = ?", usuario_id)
	if erro != nil {
		return "", erro
	}
	defer linha.Close()
	var senha = ""
	if linha.Next() {
		if erro = linha.Scan(&senha); erro != nil {
			return "", erro
		}
	}
	return senha, nil
}

func (repositorio Usuarios) SalvarNovaSenha(usuarioID uint64, senha string) error {
	statement, erro := repositorio.db.Prepare("UPDATE usuarios SET senha = ? WHERE id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(senha, usuarioID); erro != nil {
		return erro
	}
	return nil
}
