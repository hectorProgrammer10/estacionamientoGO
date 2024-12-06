package main

import (
	"estacionamientoGo/src/models"
	"estacionamientoGo/src/scenes"
	"estacionamientoGo/src/views"
	"estacionamientoGo/src/vigilante"
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

// Variables globales sincronizadas
var mu sync.Mutex
var ballSaliendo bool = false
var ballEntrando bool = false
var cond = sync.NewCond(&mu)

func main() {
	myApp := app.New()
	stage := myApp.NewWindow("App - Ball")
	stage.CenterOnScreen()
	stage.Resize(fyne.NewSize(900, 400))
	stage.SetFixedSize(true)

	//---
	// Create scene
	scene := scenes.NewScene(stage)
	scene.Init()

	// Botón para agregar una ball que entra

	button := widget.NewButton("Iniciar", func() {
		for i := 0; i < 100; i++ {
			v := vigilante.NuevoVigilante()
			estadoChannel := make(chan string)
			resultadoChannel := make(chan bool)
			b1 := models.NewBall()
			ball := views.NewBall()
			ball.AddBall(*scene)
			b1.Register(ball)
			go models.Estacionamiento(estadoChannel, resultadoChannel, v)
			// Goroutine principal para manejar entrada y salida
			go func() {
				mu.Lock()
				// Esperar a que no haya balls saliendo
				for ballSaliendo {
					cond.Wait()
				}
				ballEntrando = true
				mu.Unlock()

				// Entrar al estacionamiento
				estadoChannel <- "entrando"
				b1.Run() // Manejar lógica de entrada
				estadoChannel <- "completado"

				mu.Lock()
				ballEntrando = false
				cond.Broadcast() // Notificar a otras goroutines
				mu.Unlock()

				// Simular salida
				//fmt.Println("Ball saliendo...")
				mu.Lock()
				ballSaliendo = true
				mu.Unlock()

				estadoChannel <- "saliendo"
				models.Destruir(b1) // Animación de salida
				estadoChannel <- "completado"

				mu.Lock()
				ballSaliendo = false
				cond.Broadcast() // Notificar a otras goroutines
				mu.Unlock()

				close(estadoChannel)
				ball.RemoveBall()
			}()

			// Manejo del resultado
			go func() {
				for resultado := range resultadoChannel {
					if resultado {
						//fmt.Println("Acceso denegado.")
					} else {

						//fmt.Println("Acceso permitido.")
					}
				}
			}()
			fmt.Println("auto creado numero: ", i)
			time.Sleep(500 * time.Microsecond)
		}
	})

	button.Move(fyne.NewPos(0, 350))
	button.Resize(fyne.NewSize(100, 50))
	scene.AddWidget(button)

	stage.ShowAndRun()
}
