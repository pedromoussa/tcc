package main

import (
	"biblioteca/safechannel"
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("== Usando SafeChannel ==")
	testSafeChannel()

	fmt.Println("\n== Usando Channels Convencionais ==")
	testConventionalChannel()
}

func testSafeChannel() {
	// Criando um SafeChannel com buffer de 2
	sc := safechannel.NewSafeChannel[string](2)

	// bufferedChannel := safechannel.NewSafeChannel[int](2)
	// unbufferedChannel := safechannel.NewSafeChannel[string](0)

	var wg sync.WaitGroup
	wg.Add(2)

 // Goroutine 1 - Enviando dados usando SafeChannel
 go func() {
	 defer wg.Done()
	 for i := 0; i < 2; i++ {
		 err := sc.Send(fmt.Sprintf("SafeChannel Mensagem %d", i+1))
		 if err != nil {
			 fmt.Println("Erro ao enviar (SafeChannel):", err)
			 return
		 }
		 fmt.Println("Enviado (SafeChannel):", i+1)
	 }
	 // Fechando o SafeChannel após envio
	 if err := sc.Close(); err != nil {
		 fmt.Println("Erro ao fechar (SafeChannel):", err)
	 }
  }()
 	// Goroutine 2 - Recebendo dados usando SafeChannel
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second) // Simulando algum processamento
		for {
			value, err := sc.Receive()
			if err != nil {
				fmt.Println("Erro ao receber (SafeChannel):", err)
				return
			}
			fmt.Println("Recebido (SafeChannel):", value)
		}
	}()

	// Aguardando as goroutines finalizarem
	wg.Wait()

	// Teste de erro: Tentar enviar para um SafeChannel já fechado
	err := sc.Send("Tentando enviar após fechar (SafeChannel)")
	if err != nil {
		fmt.Println("Erro detectado conforme esperado (SafeChannel):", err)
	}
}

func testConventionalChannel() {
	// Criando um canal convencional com buffer de 2
	ch := make(chan string, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1 - Enviando dados usando Channel Convencional
	go func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			ch <- fmt.Sprintf("Channel Mensagem %d", i+1)
			fmt.Println("Enviado (Channel):", i+1)
		}
		close(ch) // Fechando o canal após envio
	}()

	// Goroutine 2 - Recebendo dados usando Channel Convencional
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second) // Simulando algum processamento
		for value := range ch {
			fmt.Println("Recebido (Channel):", value)
		}
	}()

	// Aguardando as goroutines finalizarem
	wg.Wait()

	// Teste de erro: Tentar enviar para um canal já fechado
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Erro detectado conforme esperado (Channel):", r)
		}
	}()
	ch <- "Tentando enviar após fechar (Channel)" // Isso causará um pânico
}