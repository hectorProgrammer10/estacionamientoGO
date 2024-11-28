package scenes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type Scene struct {
	scene     fyne.Window
	container *fyne.Container
}

func NewScene(scene fyne.Window) *Scene {
	return &Scene{scene: scene, container: nil}
}

func (s *Scene) Init() {

	// Crear un rectángulo de fondo con color personalizado
	rectFondo := canvas.NewRectangle(color.White)
	rectFondo.Resize(fyne.NewSize(900, 400))
	rectFondo.Move(fyne.NewPos(0, 0))
	rect := canvas.NewImageFromFile("./assets/estacionamiento.png")
	rect.Resize(fyne.NewSize(600, 250))
	rect.Move(fyne.NewPos(300, 150))

	// Colocar el rectángulo de fondo y los widgets dentro de un contenedor
	s.container = container.NewWithoutLayout(rectFondo, rect)
	s.scene.SetContent(s.container)
}

func (s *Scene) AddWidget(widget fyne.Widget) {
	s.container.Add(widget)
	s.container.Refresh()
}

func (s *Scene) AddImage(image *canvas.Image) {
	s.container.Add(image)
	s.container.Refresh()
}
