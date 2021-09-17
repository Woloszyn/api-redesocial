package main

import (
	"api-redesocial/src/config"
	"api-redesocial/src/router"
	"fmt"
	"log"
	"net/http"
)

// func init() {
// 	chave := make([]byte, 64)
// 	if _, erro := rand.Read(chave); erro != nil {
// 		log.Fatal(erro)
// 	}
// 	stringBase64 := base64.StdEncoding.EncodeToString(chave)
// 	fmt.Print(stringBase64)
// }

func main() {
	config.Carregar()
	fmt.Println("Escutando na porta ", config.Porta)
	r := router.Gerar()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))
}
