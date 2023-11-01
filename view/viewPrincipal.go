package view

import (
	"fyne.io/fyne/v2"
	"Estacionamiento/utils"
)

type MainView struct {
	window fyne.Window
}

func NewMainView(window fyne.Window) *MainView {
	MainView := &MainView{
		window: window,
	}
	MainView.StartParkingSimulation()
	return MainView
}


func (m *MainView) StartParkingSimulation() {
	utils.NewParkingView(m.window)
}
