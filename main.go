package main

import (
	"fmt"
	"time"
	"biblioteca/safechannel"
)

// Exemplo de uso do SafeChannel

// Pedido representa um pedido no restaurante
type Pedido struct {
	ID      int
	Cliente string
}

// Garçom recebe pedidos e os envia para a cozinha pelo SafeChannel
func Garcom(sc *safechannel.SafeChannel[Pedido]) {
	for i := 1; i <= 5; i++ {
		pedido := Pedido{ID: i, Cliente: fmt.Sprintf("Cliente %d", i)}
		err := sc.Send(pedido)
		if err != nil {
			fmt.Println("Erro ao enviar pedido:", err)
		} else {
			fmt.Printf("Garçom: Pedido %d recebido\n", pedido.ID)
		}
		time.Sleep(1 * time.Second) // Simula tempo entre pedidos
	}
	sc.Close() // Fecha o canal após todos os pedidos serem enviados
}

// Cozinha recebe pedidos e os processa
func Cozinha(sc *safechannel.SafeChannel[Pedido]) {
	for {
		pedido, err := sc.Receive()
		if err != nil {
			fmt.Println("Cozinha: Nenhum pedido restante, fechando.")
			sc.Close()
			return
		}
		fmt.Printf("Cozinha: Preparando pedido %d para %s\n", pedido.ID, pedido.Cliente)
		time.Sleep(2 * time.Second) // Simula tempo de preparo
		fmt.Printf("Cozinha: Pedido %d pronto!\n", pedido.ID)
	}
}

// Monitoramento do canal de notificações
func monitorarNotificacoes(sc *safechannel.SafeChannel[Pedido]) {
	go func() {
		for {
			notificacao, err := sc.ReadNotification()
			if err != nil {
				break
			}
			fmt.Printf("LOG: %s | Função: %s | Arquivo: %s | Linha: %d\n",
				notificacao.Message, notificacao.FuncName, notificacao.File, notificacao.Line)
		}
	}()
}

func main() {
	sc := safechannel.MakeSafechannel[Pedido](2)
	sc.EnableNotifications(5)	// Ativa o canal de notificações

	go Garcom(sc)
	go Cozinha(sc)
	monitorarNotificacoes(sc)

	time.Sleep(15 * time.Second) // Aguarda todas as goroutines finalizarem
	fmt.Println("Restaurante fechado!")
}
