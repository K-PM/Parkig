package models

import (
	"sync"
	"time"
)

var (
	mutexExitEnter sync.Mutex
)

const (
	lambda         = 2.0
	MaxWait    int = 10
	MaxParking int = 20
)

type Parking struct {
	semGenerateCar chan bool
    carID          int
	waitCars            []*Car
	parking             [MaxParking]*Car
	entrace             *Car
	exit                *Car
	out                 *Car
	semQuit             chan bool
}

func NewParking(sENCW chan bool, sQ chan bool) *Parking {
	parking := &Parking{
		semGenerateCar: make(chan bool),
        carID:          0,
        waitCars:       make([]*Car, 0),
	}
	return parking
}

func (p *Parking) MakeParking() {
	for i := range p.parking {
		car := NewSpaceCar()
		p.parking[i] = car
	}
}

func (p *Parking) MakeOutStation() *Car {
	p.out = NewSpaceCar()
	return p.out
}

func (p *Parking) MakeExitStation() *Car {
	p.exit = NewSpaceCar()
	return p.exit
}

func (p *Parking) MakeEntraceStation() *Car {
	p.entrace = NewSpaceCar()
	return p.entrace
}

func (p *Parking) GenerateCars() {
	
	for i := 0; i < 20; i++ {
        select {
        case <-p.semGenerateCar:
            return
        default:
            car := NewCar(p.carID, p.semQuit)
            go car.StartCount(p.carID)
            p.carID++
            p.waitCars = append(p.waitCars, car)
            time.Sleep(time.Second) // Espera un segundo antes de generar el siguiente vehículo
        }
    }
}


func (p *Parking) CheckParking() {
	for {
		select {
		case <-p.semQuit:
			return
		default:
			index := p.SearchSpace()
			if index != -1 && !p.WaitCarsIsEmpty() {
				mutexExitEnter.Lock()
				p.MoveToEntrace()
				p.MoveToPark(index)
				mutexExitEnter.Unlock()
			}
		}
	}
}

func (p *Parking) MoveToEntrace() {
	car := p.PopWaitCars()
	p.entrace.ReplaceData(car)
	time.Sleep(1 * time.Second)
}

func (p *Parking) MoveToPark(index int) {
	p.parking[index].ReplaceData(p.entrace)
	p.parking[index].text.Show()

	p.entrace.ReplaceData(NewSpaceCar())
	go p.parking[index].StartCount(index)
	time.Sleep(1 * time.Second)

}

func (p *Parking) OutCarToExit() {
	for {
		select {
		case <-p.semQuit:
			return
		default:
			if !WaitExitCarsIsEmpty() {
				mutexExitEnter.Lock()
				car := PopExitWaitCars()

				p.MoveToExit(car.ID)
				p.MoveToOut()
				mutexExitEnter.Unlock()

				time.Sleep(1 * time.Second)
				p.out.ReplaceData(NewSpaceCar())
			}
		}
	}
}

func (p *Parking) MoveToExit(index int) {
	p.exit.ReplaceData(p.parking[index])
	p.parking[index].text.Hide()
	p.parking[index].ReplaceData(NewSpaceCar())
	time.Sleep(1 * time.Second)
}

func (p *Parking) MoveToOut() {
	p.out.ReplaceData(p.exit)
	p.exit.ReplaceData(NewSpaceCar())
	time.Sleep(1 * time.Second)
}

func (p *Parking) SearchSpace() int {
	for s := range p.parking {
		if p.parking[s].GetID() == -1 {
			return s
		}
	}
	return -1
}

func (p *Parking) PopWaitCars() *Car {
	car := p.waitCars[0]
	if !p.WaitCarsIsEmpty() {
		p.waitCars = p.waitCars[1:]
	}
	return car
}

func (p *Parking) WaitCarsIsEmpty() bool {
	return len(p.waitCars) == 0
}

func (p *Parking) GetWaitCars() []*Car {
	return p.waitCars
}

func (p *Parking) GetEntraceCar() *Car {
	return p.entrace
}

func (p *Parking) GetExitCar() *Car {
	return p.exit
}

func (p *Parking) GetParking() [MaxParking]*Car {
	return p.parking
}

func (p *Parking) ClearParking() {
	for i := range p.parking {
		p.parking[i] = nil
	}
}
