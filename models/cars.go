package models

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Car struct {
	ID         int
	rectangule *canvas.Rectangle
	text       *canvas.Text
	time       int
	semQuit    chan bool
}


var exitCars []*Car

func NewSpaceCar() *Car {

	rectangule := canvas.NewRectangle(color.RGBA{R: 30, G: 30, B: 30, A: 255})

	rectangule.SetMinSize(fyne.NewSquareSize(float32(30)))

	text := canvas.NewText(fmt.Sprintf("%d", 0), color.RGBA{R: 0, G: 0, B: 0, A: 255})
	text.Hide()

	car := &Car{
		ID:         -1,
		rectangule: rectangule,
		time:       0,
		text:       text,
	}

	return car
}

func NewCar(id int, sQ chan bool) *Car {
	rand.Seed(time.Now().UnixNano()) 
	rangB := rand.Intn(256) // Componente azul entre 0 y 255
	colorRectangle := color.RGBA{R: 0, G: 0, B: uint8(rangB), A: 255}
	
	time := rand.Intn(5-1)

	rectangule := canvas.NewRectangle(colorRectangle)
	rectangule.SetMinSize(fyne.NewSquareSize(float32(30)))

	text := canvas.NewText(fmt.Sprintf("%d", time), color.RGBA{R: 30, G: 30, B: 30, A: 255})
	text.Hide()

	car := &Car{
		ID:         id,
		rectangule: rectangule,
		time:       time,
		text:       text,
		semQuit:    sQ,
	}

	return car
}

func (c *Car) StartCount(id int) {
	for {
		select {
		case <-c.semQuit:
			return
		default:
			if c.time <= 0 {
				c.ID = id
				exitCars = append(exitCars, c)
				return
			}
			c.time--
			c.text.Text = fmt.Sprintf("%d", c.time)
			time.Sleep(1 * time.Second)
		}
	}
}

func (c *Car) GetRectangle() *canvas.Rectangle {
	return c.rectangule
}
func (c *Car) ReplaceData(car *Car) {
	c.ID = car.ID
	c.time = car.time
	c.rectangule.FillColor = car.rectangule.FillColor
	c.text.Text = car.text.Text
	c.text.Color = car.text.Color
}

func (c *Car) GetText() *canvas.Text {
	return c.text
}

func (c *Car) GetTime() int {
	return c.time
}

func (c *Car) GetID() int {
	return c.ID
}

func GetWaitCars() []*Car {
	return exitCars
}

func PopExitWaitCars() *Car {
	car := exitCars[0]
	if !WaitExitCarsIsEmpty() {
		exitCars = exitCars[1:]
	}
	return car
}

func WaitExitCarsIsEmpty() bool {
	return len(exitCars) == 0
}
