package models

import (
	"estacionamientoGo/src/vigilante"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Definimos un mutex global para sincronizar el acceso al estacionamiento
var mu sync.Mutex

var pasar sync.Mutex

// Variable global para las posiciones del estacionamiento
var posiciones = []estacionamiento{
	{400, 300, false}, {450, 300, false}, {500, 300, false}, {550, 300, false},
	{600, 300, false}, {650, 300, false}, {700, 300, false}, {750, 300, false},
	{800, 300, false}, {850, 300, false},
	{400, 150, false}, {450, 150, false}, {500, 150, false}, {550, 150, false},
	{600, 150, false}, {650, 150, false}, {700, 150, false}, {750, 150, false},
	{800, 150, false}, {850, 150, false},
}

// Estructura de cada lugar de estacionamiento
type estacionamiento struct {
	X       int
	Y       int
	ocupado bool
}

// Ball representa un objeto con posición y estado
type Ball struct {
	posX, posY int32
	status     bool
	angulo     int32
	esperando  bool
	observers  []Observer
}

// Crear una nueva instancia de Ball
func NewBall() *Ball {
	return &Ball{posX: 0, posY: 0, status: true}
}

// Lógica principal de la Ball
func (b *Ball) Run() {

	var sigue bool = true
	var incX int32 = 10
	var rotationSpeed int32 = 5

	for sigue {
		// Buscar un lugar disponible de manera segura
		posicion := buscarLugarSeguro()
		if posicion == -1 { // Si no hay lugares disponibles
			fmt.Println("No hay lugares disponibles, creando uno nuevo...")
			b.moverAZonaDeEspera()
			time.Sleep(2 * time.Second) // Esperar un tiempo antes de volver a intentar
			continue

		}

		// Obtener la posición asignada
		mu.Lock()
		lugarSeleccionado := posiciones[posicion]
		originX := int32(lugarSeleccionado.X)
		originY := int32(lugarSeleccionado.Y)
		mu.Unlock()

		// Mover la Ball hacia el lugar asignado
		b.status = true
		b.posX = 0
		b.posY = 200
		b.angulo = 0
		for b.status {
			if b.posX != originX {
				b.posX += incX
				b.status = true
			}
			if b.posX == originX {
				b.posY = originY
				b.status = false
			}

			// Rotación de la Ball
			b.angulo = (b.angulo + rotationSpeed) % 360
			b.NotifyAll()
			time.Sleep(50 * time.Millisecond)

		}

		fmt.Printf("Ball estacionada en lugar: %+v\n", posiciones[posicion])

		fmt.Println("completadp")

		sigue = false

		time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
		mu.Lock()

		posiciones[posicion].ocupado = false
		mu.Unlock()

		fmt.Printf("Ball liberó lugar: %+v\n", posiciones[posicion])

	}

}

// Mover la Ball a la zona de espera
func (b *Ball) moverAZonaDeEspera() {
	b.esperando = true
	b.posX = 100
	b.posY = 100
	b.NotifyAll()
	fmt.Println("Ball en la zona de espera...")
}

// Buscar un lugar disponible de manera segura
func buscarLugarSeguro() int {
	mu.Lock()
	defer mu.Unlock()

	for i, lugar := range posiciones {
		if !lugar.ocupado {
			posiciones[i].ocupado = true // Marcar el lugar como ocupado
			return i
		}
	}
	return -1 // No hay lugares disponibles
}

// Agregar un nuevo lugar de manera segura
func agregarLugarSeguro() {
	mu.Lock()
	defer mu.Unlock()

	ultimo := posiciones[len(posiciones)-1]
	nuevoLugar := estacionamiento{
		X:       ultimo.X + 100, // Desplazamos en el eje X
		Y:       ultimo.Y,
		ocupado: false,
	}

	posiciones = append(posiciones, nuevoLugar)
	fmt.Println("Se agregó un nuevo lugar:", nuevoLugar)
}

// Register añade un observador a la lista
func (b *Ball) Register(observer Observer) {
	b.observers = append(b.observers, observer)
}

// Unregister elimina un observador de la lista
func (b *Ball) Unregister(observer Observer) {
	for i, o := range b.observers {
		if o == observer {
			b.observers = append(b.observers[:i], b.observers[i+1:]...)
			break
		}
	}
}

// NotifyAll notifica a todos los observadores sobre una actualización
func (b *Ball) NotifyAll() {
	for _, observer := range b.observers {
		observer.Update(Pos{X: b.posX, Y: b.posY})
	}
}

// func Destruir(b *Ball, done chan<- bool) {
func Destruir(b *Ball) {
	var desX int32 = 10
	originX := int32(100)
	b.posY = 200
	b.status = true
	for b.status {
		if b.posX != originX {
			b.posX -= desX
			b.status = true
		}
		if b.posX == originX {
			b.status = false
		}
		b.NotifyAll()
		time.Sleep(50 * time.Millisecond)

	}
	fmt.Println("DEstruir--")

}
func Estacionamiento(estadoChannel <-chan string, resultadoChannel chan<- bool, v *vigilante.Vigilante) {
	for estado := range estadoChannel {

		// Actualiza el estado del vigilante y envía el resultado
		resultado := v.ActualizarEstado(estado)

		resultadoChannel <- resultado

		fmt.Printf("Estacionamiento procesó: %s | Entrada libre: %v\n", estado, resultado)

	}
}

/*func (b *Ball) SetStatus(status bool) {
	b.status = status
}
*/
