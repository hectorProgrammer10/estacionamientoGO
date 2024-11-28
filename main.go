package main

import (
	"estacionamientoGo/src/models"
	"estacionamientoGo/src/scenes"
	"estacionamientoGo/src/views"
	_ "fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

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

	//Add a new widget

	button := widget.NewButton("Click", func() {

		// Creamos el objeto observado
		b1 := models.NewBall()
		// Add Balon (Observador)
		ball := views.NewBall()
		ball.AddBall(*scene)
		//Registramos a ball como observador de la goroutine b1
		b1.Register(ball)
		go func() {
			b1.Run()
			//fmt.Println("se acabo Run")
			models.Destruir(b1)
			ball.RemoveBall()
		}()

	})
	button.Move(fyne.NewPos(0, 350))
	button.Resize(fyne.NewSize(100, 50))
	scene.AddWidget(button)

	stage.ShowAndRun()
}
