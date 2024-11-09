package main

import (
	"estacionamiento/src/scenes"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Estacionamiento GOOOO")

	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(900, 400))
	myWindow.SetFixedSize(true)

	//--esena
	scene := scenes.NewScene(myWindow)
	scene.Init()

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

	//------------

	scene.AddImage(img)
	myWindow.ShowAndRun()
}
