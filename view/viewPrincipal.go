package view

import (
	"fyne.io/fyne/v2"
	"Estacionamiento/utils"
)
type VistaPrincipal struct {
	Window fyne.Window
}

func NuevaVistaPrincipal(Window fyne.Window) *VistaPrincipal {
	vistaPrincipal := &VistaPrincipal{
		Window: Window,
	}
	vistaPrincipal.IniciarSimulacionEstacionamiento()
	return vistaPrincipal
}

func (v *VistaPrincipal) IniciarSimulacionEstacionamiento() {
	utils.NuevaVistaEstacionamiento(v.Window)
}
