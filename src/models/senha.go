package models

//Senha representa o formato da requisição para alteração da senha
type Senha struct {
	Nova  string `json:"nova"`
	Atual string `json:"atual"`
}
