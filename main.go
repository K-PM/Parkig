package main

import (
    "Estacionamiento/view"
    "fyne.io/fyne/v2/app"
)

func main() {
    myApp := app.New()
    mainWindow := myApp.NewWindow("Estacionamiento")
    view.NuevaVistaPrincipal(mainWindow)
    mainWindow.ShowAndRun()
}
