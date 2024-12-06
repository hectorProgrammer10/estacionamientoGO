package views

import (
	"estacionamientoGo/src/models"
	"estacionamientoGo/src/scenes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
)

type Car struct {
	Car *canvas.Image
}

func NewCar() *Car {
	return &Car{Car: nil}
}

func (b *Car) AddCar(c scenes.Scene) {
	pelota := canvas.NewImageFromURI(storage.NewFileURI("./assets/autoConcurrente.png"))
	pelota.Resize(fyne.NewSize(48, 96))
	pelota.Move(fyne.NewPos(-100, 200))
	b.Car = pelota
	c.AddImage(pelota)
}

// Update define lo que el observador hará cuando reciba una notificación
func (b *Car) Update(pos models.Pos) {
	//fmt.Printf("%d : %d\n", pos.X, pos.Y)
	b.Car.Move(fyne.NewPos(float32(pos.X), float32(pos.Y)))
}

func (b *Car) RemoveCar() {
	if b.Car != nil {
		b.Car.Hide() // Ocultar la imagen antes de destruirla
		b.Car = nil  // Eliminar la referencia
	}
}
