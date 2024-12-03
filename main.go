package main

import (
	"estacionamientoGo/src/models"
	"estacionamientoGo/src/scenes"
	"estacionamientoGo/src/views"
	"estacionamientoGo/src/vigilante"
	"fmt"
	_ "fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	var salida sync.Mutex
	// Ejecutamos el estacionamiento como un hilo (gorutina)
	//Estacionamiento(estadoChannel, resultadoChannel, v)
	myApp := app.New()
	stage := myApp.NewWindow("App - Ball")
	stage.CenterOnScreen()
	stage.Resize(fyne.NewSize(900, 400))
	stage.SetFixedSize(true)

	//---
	// Create scene
	scene := scenes.NewScene(stage)
	scene.Init()

	//Add a new widget

	button := widget.NewButton("Entrar", func() {

		// Inicializamos el vigilante
		v := vigilante.NuevoVigilante()

		// Canales para comunicar los estados y recibir resultados
		estadoChannel := make(chan string)
		resultadoChannel := make(chan bool)
		// Creamos el objeto observado
		b1 := models.NewBall()
		// Add Balon (Observador)
		ball := views.NewBall()
		ball.AddBall(*scene)
		//Registramos a ball como observador de la goroutine b1
		b1.Register(ball)
		go models.Estacionamiento(estadoChannel, resultadoChannel, v)
		go func() {
			fmt.Println("Run--")
			fmt.Println("Entrando--")
			estadoChannel <- "entrando"
			salida.Lock()
			b1.Run()

			estadoChannel <- "completado"
			fmt.Println("saliendo--")
			salida.Lock()
			estadoChannel <- "saliendo"

			models.Destruir(b1) //animación salir
			estadoChannel <- "completado"
			close(estadoChannel)

			ball.RemoveBall()

		}()
		go func() {
			for resultado := range resultadoChannel {
				if !resultado {
					fmt.Println("Acceso permitido.")
					salida.Unlock()

				} else {
					fmt.Println("Notificación: Acceso denegado.")

				}
			}
		}()

	})
	button.Move(fyne.NewPos(0, 350))
	button.Resize(fyne.NewSize(100, 50))
	scene.AddWidget(button)

	stage.ShowAndRun()

}
