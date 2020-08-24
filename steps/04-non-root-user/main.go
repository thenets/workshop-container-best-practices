package main

import "github.com/sirupsen/logrus"

func main() {
	var log = logrus.New()
	log.Level = logrus.DebugLevel

	log.WithFields(logrus.Fields{
		"app": "picosec",
	}).Info("03-choosing-a-better-image")
	log.Info("Olá mundo!")

	// Problems
	log.Debug("A imagem é otimizada para Golang")
	log.Debug("A imagem possui poucos binários")
	log.Debug("A imagem é pequena")
	log.Debug("O usuário padrão não é o 'root'")
	log.Warning("O usuário possui mais privilégios do que o necessário")
	log.Warning("Existem muitas camadas (layers)")
}
