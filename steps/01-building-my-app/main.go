package main

import "github.com/sirupsen/logrus"

func main() {
	var log = logrus.New()
	log.Level = logrus.DebugLevel

	log.WithFields(logrus.Fields{
		"app": "picosec",
	}).Info("01-building-my-app")
	log.Info("Olá mundo!")

	// Problems
	log.Warning("A imagem não é otimizada para Golang")
	log.Warning("A imagem possui muitos binários")
	log.Warning("A imagem é desnecessariamente grande")
	log.Warning("O usuário padrão é o 'root'")
	log.Warning("O usuário possui mais privilégios do que o necessário")
	log.Warning("Existem muitas camadas (layers)")
}
