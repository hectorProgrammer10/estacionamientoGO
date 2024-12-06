package models

import (
	"estacionamientoGo/src/vigilante"
	"math/rand"
	"sync"
	"time"
)

// Definimos un mutex global para sincronizar el acceso al estacionamiento
var mu sync.Mutex

var salida = true
var salidaR = true

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

// Car representa un objeto con posición y estado
type Car struct {
	posX, posY int32
	status     bool
	angulo     int32
	esperando  bool
	posicion   int
	observers  []Observer
}

// Crear una nueva instancia de Car
func NewCar() *Car {
	return &Car{posX: -100, posY: 200, status: true}
}

// Lógica principal de la Car
func (b *Car) Run() {

	var sigue bool = true
	var incX int32 = 50
	var rotationSpeed int32 = 5
	var contador int = 0
	for sigue {
		if salida {
			if salidaR {
				salida = false
				salidaR = false
				// Buscar un lugar disponible de manera segura

				posicion := buscarLugarSeguro()
				if posicion == -1 { // Si no hay lugares disponibles
					//fmt.Println("No hay lugares disponibles, creando uno nuevo...")
					b.moverAZonaDeEspera()
					time.Sleep(1 * time.Second) // Esperar un tiempo antes de volver a intentar
					continue

				}

				// Obtener la posición asignada
				mu.Lock()
				lugarSeleccionado := posiciones[posicion]
				originX := int32(lugarSeleccionado.X)
				originY := int32(lugarSeleccionado.Y)
				mu.Unlock()
				wait := true
				// Mover la Car hacia el lugar asignado
				b.status = true
				b.angulo = 0
				wait = true
				b.posX = -100
				b.posY = 200
				b.NotifyAll()
				for b.status {
					for wait {
						if b.posX != originX {
							b.posX += incX
							b.status = true
						}
						if b.posX == originX {
							b.posY = originY
							b.status = false
							wait = false
						}

						// Rotación de la Car
						b.angulo = (b.angulo + rotationSpeed) % 90
						b.NotifyAll()
						time.Sleep(50 * time.Millisecond)
					}

				}

				//fmt.Printf("Car estacionada en lugar: %+v\n", posiciones[posicion])

				//fmt.Println("completadp")

				sigue = false
				salida = true
				salidaR = true
				time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
				mu.Lock()
				salida = false
				if !salida {
					wait = false

				}
				mu.Unlock()
				b.posicion = posicion

				//fmt.Printf("Car liberó lugar: %+v\n", posiciones[posicion])
			}

		} else {
			contador++
			time.Sleep(1 * time.Second)
			//fmt.Println("esperando a que salida se desocupe salida", salida, salidaR)
		}
	}

}

// Mover la Car a la zona de espera
func (b *Car) moverAZonaDeEspera() {
	mover := true
	b.esperando = true

	x := int32(100)
	y := int32(100)
	var inX int32 = 10
	var inY int32 = 10

	for mover {
		if b.posX != x {
			b.posX += inX

		}
		if b.posY != y {
			b.posY -= inY
		}
		if b.posX == x {
			if b.posY == y {
				mover = false
				salida = true
				salidaR = true
			}
		}
		b.NotifyAll()
		time.Sleep(50 * time.Millisecond)
	}

	//fmt.Println("Car en la zona de espera...")
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

// Register añade un observador a la lista
func (b *Car) Register(observer Observer) {
	b.observers = append(b.observers, observer)
}

// Unregister elimina un observador de la lista
func (b *Car) Unregister(observer Observer) {
	for i, o := range b.observers {
		if o == observer {
			b.observers = append(b.observers[:i], b.observers[i+1:]...)
			break
		}
	}
}

// NotifyAll notifica a todos los observadores sobre una actualización
func (b *Car) NotifyAll() {
	for _, observer := range b.observers {
		observer.Update(Pos{X: b.posX, Y: b.posY})
	}
}

// func Destruir(b *Car, done chan<- bool) {
func Destruir(b *Car) {

	var desX int32 = 50
	originX := int32(-50)
	b.posY = 200
	b.status = true
	for b.status {

		if !salida {
			for b.status {
				if salidaR {
					salidaR = false
					for b.status {
						if b.posX != originX {
							b.posX -= desX
							b.status = true

						}
						if b.posX == originX {
							salida = true
							salidaR = true
							b.status = false
							posiciones[b.posicion].ocupado = false
						}
						b.NotifyAll()
						time.Sleep(50 * time.Millisecond)
					}
				} else {
					time.Sleep(1 * time.Second)
				}
			}
		} else {
			time.Sleep(1 * time.Second)
			//fmt.Println("esperando a salir : salidaR", salidaR, salida)
		}

	}

	//fmt.Println("DEstruir--")

}
func Estacionamiento(estadoChannel <-chan string, resultadoChannel chan<- bool, v *vigilante.Vigilante) {
	for estado := range estadoChannel {

		// Actualiza el estado del vigilante y envía el resultado
		resultado := v.ActualizarEstado(estado)

		resultadoChannel <- resultado

		//fmt.Printf("Estacionamiento procesó: %s | Entrada libre: %v\n", estado, resultado)

	}
}

/*func (b *Car) SetStatus(status bool) {
	b.status = status
}
*/
