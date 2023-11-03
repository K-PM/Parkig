package models

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Auto struct {
	ID         int
	rectangulo *canvas.Rectangle
	tiempo     int
	semaQ      chan bool
}

var autosSalida []*Auto

func NuevoAutoEstacionamiento() *Auto {
	rectangulo := canvas.NewRectangle(color.RGBA{R: 30, G: 30, B: 30, A: 255})
	rectangulo.SetMinSize(fyne.NewSquareSize(float32(30)))
	texto := canvas.NewText(fmt.Sprintf("%d", 0), color.RGBA{R: 0, G: 0, B: 0, A: 255})
	texto.Hide()

	auto := &Auto{
		ID:         -1,
		rectangulo: rectangulo,
		tiempo:     0,
	}

	return auto
}

func NuevoAuto(id int, sq chan bool) *Auto {
	rand.Seed(time.Now().UnixNano())
	rangoB := rand.Intn(256)
	colorRectangulo := color.RGBA{R: 0, G: 0, B: uint8(rangoB), A: 255}
	tiempo := rand.Intn(5-1) + 1
	rectangulo := canvas.NewRectangle(colorRectangulo)
	rectangulo.SetMinSize(fyne.NewSquareSize(float32(30)))

	auto := &Auto{
		ID:         id,
		rectangulo: rectangulo,
		tiempo:     tiempo,
		semaQ:      sq,
	}

	return auto
}

func (a *Auto) IniciarConteo(id int) {
	for {
		select {
		case <-a.semaQ:
			return
		default:
			if a.tiempo <= 0 {
				a.ID = id
				autosSalida = append(autosSalida, a)
				return
			}
			a.tiempo--
			time.Sleep(1 * time.Second)
		}
	}
}

func (a *Auto) ObtenerRectangulo() *canvas.Rectangle {
	return a.rectangulo
}

func (a *Auto) ReemplazarDatos(auto *Auto) {
	a.ID = auto.ID
	a.tiempo = auto.tiempo
	a.rectangulo.FillColor = auto.rectangulo.FillColor
}

func (a *Auto) ObtenerTiempo() int {
	return a.tiempo
}

func (a *Auto) ObtenerID() int {
	return a.ID
}

func ObtenerAutosEspera() []*Auto {
	return autosSalida
}

func DesencolarSalidaAutos() *Auto {
	auto := autosSalida[0]
	if !ColaSalidaAutosVacia() {
		autosSalida = autosSalida[1:]
	}
	return auto
}

func ColaSalidaAutosVacia() bool {
	return len(autosSalida) == 0
}
