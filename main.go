package main

import (
	"os"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Estacionamiento GOOOO")

	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(900, 400))
	myWindow.SetFixedSize(true)

	//fondo blanco
	rect := canvas.NewRectangle(color.White)
	rect.Resize(fyne.NewSize(900, 400))
	//estacionamiento
	imagePath := "src/views/estacionamiento.png"
	imgFile, err := os.Open(imagePath)
	if err != nil {
		panic("No se pudo cargar la imagen: " + err.Error())
	}
	defer imgFile.Close()
	img := canvas.NewImageFromFile(imagePath)
	img.FillMode = canvas.ImageFillOriginal
	img.Resize(fyne.NewSize(600, 250))

	img.Move(fyne.NewPos(300, 150))

	background := container.NewWithoutLayout(img)
	//------------

	// Si necesitas añadir widgets sobre el fondo, usa un contenedor adicional
	content := container.NewMax(
		rect, background, // Fondo
		// Aquí puedes agregar otros widgets o contenedores como `canvas.Text` o `widget.Button`
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
