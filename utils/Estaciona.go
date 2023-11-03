package utils

import (
	"Estacionamiento/models"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var semRenderNuevoAutoEspera chan bool
var semSalir chan bool
var estacionamiento *models.Estacionamiento

type VistaEstacionamiento struct {
	ventana                fyne.Window
	rectangulosEstacionEspera [models.MaxEspera]*canvas.Rectangle
	containerVistaEstacionamiento *fyne.Container
}

func NuevaVistaEstacionamiento(ventana fyne.Window) *VistaEstacionamiento {
	vistaEstacionamiento := &VistaEstacionamiento{ventana: ventana}
	vistaEstacionamiento.inicializar()
	return vistaEstacionamiento
}

func (v *VistaEstacionamiento) inicializar() {
	v.inicializarCanales()
	v.inicializarEstacionamiento()
	v.inicializarInterfazUsuario()
	v.iniciarSimulacion()
}

func (v *VistaEstacionamiento) inicializarCanales() {
	semSalir = make(chan bool)
	semRenderNuevoAutoEspera = make(chan bool)
}

func (v *VistaEstacionamiento) inicializarEstacionamiento() {
	estacionamiento = models.NuevoEstacionamiento(semRenderNuevoAutoEspera, semSalir)
	estacionamiento.CrearEstacionamiento()
}

func (v *VistaEstacionamiento) inicializarInterfazUsuario() {
	v.crearEscena()
	v.ventana.Resize(fyne.NewSize(300, 200))
}

func (v *VistaEstacionamiento) iniciarSimulacion() {
	go estacionamiento.GenerarAutos()
	go estacionamiento.SalidaAutoASalida()
	go estacionamiento.VerificarEstacionamiento()
	go v.actualizarRenderizado()
	go v.actualizarNuevoAutoEstacionEspera()
}

func (v *VistaEstacionamiento) crearEscena() {
	canvasFondo := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	canvasFondo.Resize(fyne.NewSize(300, 200))

	v.containerVistaEstacionamiento = container.New(layout.NewVBoxLayout()) 
	containerSalidaEstacionamiento := container.New(layout.NewHBoxLayout())

	containerSalidaEstacionamiento.Add(v.crearEstacionEspera())
	containerSalidaEstacionamiento.Add(v.crearEstacionSalida())
	v.containerVistaEstacionamiento.Add(v.crearEstacionEntradaYSalida())
	v.containerVistaEstacionamiento.Add(v.crearEstacionamiento())

	v.ventana.SetContent(v.containerVistaEstacionamiento)	
	v.ventana.Resize(fyne.NewSize(300, 200))
}

func (v *VistaEstacionamiento) actualizarRenderizado() {
	for {
		select {
		case <-semSalir:
			return
		default:
			v.ventana.Content().Refresh()
			time.Sleep(1 * time.Second)
		}
	}
}

func (v *VistaEstacionamiento) actualizarNuevoAutoEstacionEspera() {
	for {
		select {
		case <-semSalir:
			return
		case <-semRenderNuevoAutoEspera:
			autosEspera := estacionamiento.ObtenerAutosEspera()
			for i := len(autosEspera) - 1; i >= 0; i-- {
				v.rectangulosEstacionEspera[i].Show()
				v.rectangulosEstacionEspera[i].FillColor = autosEspera[i].ObtenerRectangulo().FillColor
			}
			v.ventana.Content().Refresh()
		}
	}
}

func (v *VistaEstacionamiento) crearEstacionamiento() *fyne.Container {
	containerEstacionamiento := container.New(layout.NewGridLayout(5))
	estacionamiento.CrearEstacionamiento()
	estacionamientoArray := estacionamiento.ObtenerEstacionamiento()
	for i := 0; i < len(estacionamientoArray); i++ {
		containerEstacionamiento.Add(container.NewCenter(estacionamientoArray[i].ObtenerRectangulo()))
	}
	return container.NewCenter(containerEstacionamiento)
}

func (v *VistaEstacionamiento) crearEstacionEspera() *fyne.Container {
	containerEstacionamiento := container.New(layout.NewGridLayout(5))
	for i := len(v.rectangulosEstacionEspera) - 1; i >= 0; i-- {
		auto := models.NuevoAutoEstacionamiento().ObtenerRectangulo()
		v.rectangulosEstacionEspera[i] = auto
		v.rectangulosEstacionEspera[i].Hide()
		containerEstacionamiento.Add(v.rectangulosEstacionEspera[i])
	}
	return containerEstacionamiento
}

func (v *VistaEstacionamiento) crearEstacionSalida() *fyne.Container {
	salida := estacionamiento.CrearSalidaEstacionamiento()
	return container.NewCenter(salida.ObtenerRectangulo())
}

func (v *VistaEstacionamiento) crearEstacionEntradaYSalida() *fyne.Container {
	containerEstacionamiento := container.New(layout.NewGridLayout(5))
	containerEstacionamiento.Add(layout.NewSpacer())
	entrada := estacionamiento.CrearEntrada()
	containerEstacionamiento.Add(entrada.ObtenerRectangulo())
	containerEstacionamiento.Add(layout.NewSpacer())

	salida := estacionamiento.CrearSalida()
	
	containerEstacionamiento.Add(salida.ObtenerRectangulo())
	containerEstacionamiento.Add(layout.NewSpacer())
	return container.NewCenter(containerEstacionamiento)
}
