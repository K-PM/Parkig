package models

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

var mutexEntradaSalida sync.Mutex

const (
	lambda         = 2.0
	MaxEspera    int = 10
	MaxEstacionamiento int = 20
)

type Estacionamiento struct {
	colaEspera            []*Auto
	estacionamiento             [MaxEstacionamiento]*Auto
	entrada             *Auto
	salida                *Auto
	salidaEstacionamiento                 *Auto
	semaEspera             chan bool
	semRenderNuevoAutoEspera chan bool
}

func NuevoEstacionamiento(senaENCW chan bool, sq chan bool) *Estacionamiento {
	estacionamiento := &Estacionamiento{
		semRenderNuevoAutoEspera: senaENCW,
		semaEspera:             sq,
	}
	return estacionamiento
}

func (e *Estacionamiento) CrearEstacionamiento() {
	for i := range e.estacionamiento {
		auto := NuevoAutoEstacionamiento()
		e.estacionamiento[i] = auto
	}
}

func (e *Estacionamiento) CrearSalidaEstacionamiento() *Auto {
	e.salidaEstacionamiento = NuevoAutoEstacionamiento()
	return e.salidaEstacionamiento
}

func (e *Estacionamiento) CrearSalida() *Auto {
	e.salida = NuevoAutoEstacionamiento()
	return e.salida
}

func (e *Estacionamiento) CrearEntrada() *Auto {
	e.entrada = NuevoAutoEstacionamiento()
	return e.entrada
}

//PARA Poisson
func (e *Estacionamiento) GenerarAutos() {
	cantidadAutos := 0
	for {
		for cantidadAutos < 100{
			tiempoInterarribo := -math.Log(1-rand.Float64()) / lambda
			time.Sleep(time.Duration(tiempoInterarribo * float64(time.Second)))

			if len(e.colaEspera) < MaxEspera {
				auto := NuevoAuto(cantidadAutos, e.semaEspera)
				e.colaEspera = append(e.colaEspera, auto)
				e.semRenderNuevoAutoEspera <- true
				cantidadAutos++
			}
		}
	}
}

func (e *Estacionamiento) VerificarEstacionamiento() {
	for {
		select {
		case <-e.semaEspera:
			return
		default:
			indice := e.BuscarEspacio()
			if indice != -1 && !e.ColaEsperaVacia() {
				mutexEntradaSalida.Lock()
				e.MoverAEntrada()
				e.MoverAEstacionamiento(indice)
				mutexEntradaSalida.Unlock()
			}
		}
	}
}

func (e *Estacionamiento) MoverAEntrada() {
	auto := e.DesencolarAutosEspera()
	e.entrada.ReemplazarDatos(auto)
	time.Sleep(1 * time.Second)
}

func (e *Estacionamiento) MoverAEstacionamiento(indice int) {
	e.estacionamiento[indice].ReemplazarDatos(e.entrada)

	e.entrada.ReemplazarDatos(NuevoAutoEstacionamiento())
	go e.estacionamiento[indice].IniciarConteo(indice)
	time.Sleep(1 * time.Second)
}

func (e *Estacionamiento) SalidaAutoASalida() {
	for {
		select {
		case <-e.semaEspera:
			return
		default:
			if !ColaSalidaAutosVacia() {
				mutexEntradaSalida.Lock()
				auto := DesencolarSalidaAutos()

				e.MoverASalida(auto.ID)
				e.MoverASalidaEstacionamiento()
				mutexEntradaSalida.Unlock()

				time.Sleep(1 * time.Second)
				e.salidaEstacionamiento.ReemplazarDatos(NuevoAutoEstacionamiento())
			}
		}
	}
}

func (e *Estacionamiento) MoverASalida(indice int) {
	e.salida.ReemplazarDatos(e.estacionamiento[indice])
	e.estacionamiento[indice].ReemplazarDatos(NuevoAutoEstacionamiento())
	time.Sleep(1 * time.Second)
}

func (e *Estacionamiento) MoverASalidaEstacionamiento() {
	e.salidaEstacionamiento.ReemplazarDatos(e.salida)
	e.salida.ReemplazarDatos(NuevoAutoEstacionamiento())
	time.Sleep(1 * time.Second)
}

func (e *Estacionamiento) BuscarEspacio() int {
	for s := range e.estacionamiento {
		if e.estacionamiento[s].ObtenerID() == -1 {
			return s
		}
	}
	return -1
}

func (e *Estacionamiento) DesencolarAutosEspera() *Auto {
	auto := e.colaEspera[0]
	if !e.ColaEsperaVacia() {
		e.colaEspera = e.colaEspera[1:]
	}
	return auto
}

func (e *Estacionamiento) ColaEsperaVacia() bool {
	return len(e.colaEspera) == 0
}

func (e *Estacionamiento) ObtenerAutosEspera() []*Auto {
	return e.colaEspera
}

func (e *Estacionamiento) ObtenerAutoEntrada() *Auto {
	return e.entrada
}

func (e *Estacionamiento) ObtenerAutoSalida() *Auto {
	return e.salida
}

func (e *Estacionamiento) ObtenerEstacionamiento() [MaxEstacionamiento]*Auto {
	return e.estacionamiento
}

func (e *Estacionamiento) LimpiarEstacionamiento() {
	for i := range e.estacionamiento {
		e.estacionamiento[i] = nil
	}
}
