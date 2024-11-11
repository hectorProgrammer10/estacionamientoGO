package scenes

import (
	"fmt"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type Scene struct {
	scene     fyne.Window
	container *fyne.Container
	vehiculos map[int]*canvas.Image
}

func NewScene(scene fyne.Window) *Scene {
	return &Scene{scene: scene, container: nil, vehiculos: make(map[int]*canvas.Image)}
}

func (s *Scene) Init() {
	// Crear fondo blanco
	rect := canvas.NewRectangle(color.White)
	rect.Resize(fyne.NewSize(900, 400))
	rect.Move(fyne.NewPos(0, 0))
	s.container = container.NewWithoutLayout(rect)
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

var posiciones = []fyne.Position{
	// Lugares 1-10 (en la fila superior)
	{X: 100, Y: 250}, {X: 150, Y: 250}, {X: 200, Y: 250}, {X: 250, Y: 250},
	{X: 300, Y: 250}, {X: 350, Y: 250}, {X: 400, Y: 250}, {X: 450, Y: 250},
	{X: 500, Y: 250}, {X: 550, Y: 250},
	// Lugares 11-20 (en la segunda fila)
	{X: 100, Y: 300}, {X: 150, Y: 300}, {X: 200, Y: 300}, {X: 250, Y: 300},
	{X: 300, Y: 300}, {X: 350, Y: 300}, {X: 400, Y: 300}, {X: 450, Y: 300},
	{X: 500, Y: 300},
}

func (s *Scene) MostrarVehiculo(id, lugar int) {
	imgPath := "assets/autoConcurrente.png"
	imgFile, err := os.Open(imgPath)
	if err != nil {
		fmt.Println("Error cargando la imagen del vehículo:", err)
		return
	}
	defer imgFile.Close()

	vehiculo := canvas.NewImageFromFile(imgPath)
	vehiculo.Resize(fyne.NewSize(48, 48))
	vehiculo.Move(posiciones[lugar])
	s.container.Add(vehiculo)
	s.vehiculos[id] = vehiculo
	s.container.Refresh()
}

func (s *Scene) QuitarVehiculo(id int) {
	if vehiculo, exists := s.vehiculos[id]; exists {
		s.container.Remove(vehiculo)
		delete(s.vehiculos, id)
		s.container.Refresh()
	}
}
