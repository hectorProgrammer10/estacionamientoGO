package vigilante

import (
	"fmt"
	"sync"
)

type Vigilante struct {
	mu           sync.Mutex // Mutex para proteger la variable de estado
	cond         *sync.Cond // Condición para controlar el flujo de estados
	entradaLibre bool       // Estado actual de la entrada
}

// NuevoVigilante inicializa un nuevo vigilante
func NuevoVigilante() *Vigilante {
	v := &Vigilante{entradaLibre: true}
	v.cond = sync.NewCond(&v.mu) // Crear la condición basada en el mutex
	return v
}

// ActualizarEstado actualiza el estado del vigilante
func (v *Vigilante) ActualizarEstado(estado string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Control de flujo basado en el estado
	switch estado {
	case "entrando":
		for !v.entradaLibre { // Espera hasta que la entrada esté libre
			v.cond.Wait()
		}
		v.entradaLibre = false
	case "saliendo":
		for !v.entradaLibre { // Espera hasta que la entrada esté libre (indicado por "completado")
			v.cond.Wait()
		}
		v.entradaLibre = false
	case "completado":
		v.entradaLibre = true
		v.cond.Broadcast() // Notifica a todas las gorutinas en espera
	default:
		fmt.Printf("Estado no reconocido: %s\n", estado)
	}

	return v.entradaLibre
}
