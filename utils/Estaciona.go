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

var semRenderNewCarWait chan bool
var semQuit chan bool

type ParkingView struct {
	window               fyne.Window
	waitRectangleStation [models.MaxWait]*canvas.Rectangle
}

var parking *models.Parking

func NewParkingView(window fyne.Window) *ParkingView {
	parkingView := &ParkingView{window: window}
	semQuit = make(chan bool)
	semRenderNewCarWait = make(chan bool)
	parking = models.NewParking(semRenderNewCarWait, semQuit)
	parkingView.MakeScene()
	parkingView.StartSimulation()

	return parkingView
}

func (p *ParkingView) MakeScene() {
	bgCanvas := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	bgCanvas.Resize(fyne.NewSize(300,200))

	containerParkingView := container.New(layout.NewVBoxLayout())
	containerParkingOut := container.New(layout.NewHBoxLayout())

	containerParkingOut.Add(p.MakeWaitStation())
	containerParkingOut.Add(p.MakeExitStation())
	containerParkingView.Add(p.MakeEnterAndExitStation())
	containerParkingView.Add(p.MakeParking())

	
	p.window.SetContent(containerParkingView)	
	p.window.Resize(fyne.NewSize(300, 200))
}

func (p *ParkingView) MakeParking() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parking.MakeParking()
	parkingArray := parking.GetParking()
	for i := 0; i < len(parkingArray); i++ {
		
		parkingContainer.Add(container.NewCenter(parkingArray[i].GetRectangle()))
	}
	return container.NewCenter(parkingContainer)
}

func (p *ParkingView) MakeWaitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	for i := len(p.waitRectangleStation) - 1; i >= 0; i-- {
		car := models.NewSpaceCar().GetRectangle()
		p.waitRectangleStation[i] = car
		p.waitRectangleStation[i].Hide()
		parkingContainer.Add(p.waitRectangleStation[i])
	}
	return parkingContainer
}

func (p *ParkingView) MakeExitStation() *fyne.Container {
	out := parking.MakeOutStation()
	return container.NewCenter(out.GetRectangle())
}

func (p *ParkingView) MakeEnterAndExitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parkingContainer.Add(layout.NewSpacer())
	entrace := parking.MakeEntraceStation()
	parkingContainer.Add(entrace.GetRectangle())
	parkingContainer.Add(layout.NewSpacer())

	exit := parking.MakeExitStation()
	
	parkingContainer.Add(exit.GetRectangle())
	parkingContainer.Add(layout.NewSpacer())
	return container.NewCenter(parkingContainer)
}

func (p *ParkingView) RenderNewCarWaitStation() {
	for {
		select {
		case <-semQuit:
			return
		case <-semRenderNewCarWait:
			waitCars := parking.GetWaitCars()
			for i := len(waitCars) - 1; i >= 0; i-- {
				
					p.waitRectangleStation[i].Show()
					p.waitRectangleStation[i].FillColor = waitCars[i].GetRectangle().FillColor
				
			}
			p.window.Content().Refresh()
		}
	}
}

func (p *ParkingView) StartSimulation() {
	go parking.GenerateCars()
	go parking.OutCarToExit()
	go parking.CheckParking()
	go p.RenderUpdate()
	go p.RenderNewCarWaitStation()

}


func (p *ParkingView) RenderUpdate() {
	for {
		select {
		case <-semQuit:
			return
		default:
			p.window.Content().Refresh()
			time.Sleep(1 * time.Second)
		}
	}
}

