// main.go
package main

import (
	"estacionamiento/src/scenes"
	"math/rand"
	"os"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

const (
	capacidadEstacionamiento = 20
	numVehiculos             = 100
)

var (
	mutexPuerta        sync.Mutex
	lugaresDisponibles = capacidadEstacionamiento
	estacionamiento    = make([]bool, capacidadEstacionamiento) // indica si el lugar está ocupado
)

func main() {
	rand.Seed(time.Now().UnixNano())
	myApp := app.New()
	myWindow := myApp.NewWindow("Estacionamiento GOOOO")

	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(900, 400))
	myWindow.SetFixedSize(true)

	// Crear escena y cargar imagen de fondo
	scene := scenes.NewScene(myWindow)
	scene.Init()
	cargarImagenFondo(scene, "src/views/estacionamiento.png")

	// Simulación de llegada de vehículos
	for i := 0; i < numVehiculos; i++ {
		go crearVehiculo(i, scene)                                          // Lanza cada vehículo como una goroutine
		time.Sleep(time.Duration(rand.ExpFloat64() * float64(time.Second))) // Simula llegada en distribución Poisson
	}

	myWindow.ShowAndRun()
}

func cargarImagenFondo(scene *scenes.Scene, imagePath string) {
	imgFile, err := os.Open(imagePath)
	if err != nil {
		panic("No se pudo cargar la imagen: " + err.Error())
	}
	defer imgFile.Close()
	img := canvas.NewImageFromFile(imagePath)
	img.FillMode = canvas.ImageFillOriginal
	img.Resize(fyne.NewSize(600, 250))
	img.Move(fyne.NewPos(300, 150))
	scene.AddImage(img)
}

// Función que simula la creación de un vehículo
func crearVehiculo(id int, scene *scenes.Scene) {
	for {
		mutexPuerta.Lock()
		if lugaresDisponibles > 0 {
			// Hay espacio disponible, intenta entrar
			lugaresDisponibles--
			lugar := ocuparLugar()
			mutexPuerta.Unlock() // Libera el control de la puerta para que otros vehículos entren o salgan

			if lugar != -1 {
				scene.MostrarVehiculo(id, lugar)                        // Muestra el vehículo en la UI
				time.Sleep(time.Duration(rand.Intn(3)+3) * time.Second) // Tiempo de estacionamiento aleatorio entre 3 y 5 segundos

				// Salida del vehículo
				mutexPuerta.Lock()
				lugaresDisponibles++
				liberarLugar(lugar)
				scene.QuitarVehiculo(id) // Remueve el vehículo de la UI
				mutexPuerta.Unlock()
			}
			return
		} else {
			// Estacionamiento lleno, espera a que un lugar se libere
			mutexPuerta.Unlock()
			time.Sleep(500 * time.Millisecond) // Reintenta después de un breve descanso
		}
	}
}

// Busca un lugar vacío y lo marca como ocupado, devuelve el índice del lugar o -1 si no hay lugar
func ocuparLugar() int {
	for i := 0; i < len(estacionamiento); i++ {
		if !estacionamiento[i] {
			estacionamiento[i] = true
			return i
		}
	}
	return -1
}

// Libera el lugar en el estacionamiento
func liberarLugar(lugar int) {
	if lugar >= 0 && lugar < len(estacionamiento) {
		estacionamiento[lugar] = false
	}
}
