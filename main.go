package main

import (
	"estacionamientoGo/src/models"
	"estacionamientoGo/src/scenes"
	"estacionamientoGo/src/views"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	stage := myApp.NewWindow("App - Car")
	stage.CenterOnScreen()
	stage.Resize(fyne.NewSize(900, 400))
	stage.SetFixedSize(true)

	//---
	// Create scene
	scene := scenes.NewScene(stage)
	scene.Init()

	// Botón para agregar una Car que entra

	button := widget.NewButton("Iniciar", func() {
		for i := 0; i < 100; i++ {
			b1 := models.NewCar()
			car := views.NewCar()
			car.AddCar(*scene)
			b1.Register(car)

			// Goroutine principal para manejar entrada y salida
			go func() {

				b1.Run() // Manejar lógica de entrada

				models.Destruir(b1) // Animación de salida

				car.RemoveCar()
			}()

			fmt.Println("auto creado numero: ", i)
			time.Sleep(600 * time.Millisecond)
		}
	})

	button.Move(fyne.NewPos(0, 350))
	button.Resize(fyne.NewSize(100, 50))
	scene.AddWidget(button)

	stage.ShowAndRun()
}
