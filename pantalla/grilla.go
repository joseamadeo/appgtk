package pantalla

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func Grilla(w *gtk.ApplicationWindow){
	grid := gtk.NewGrid()
	grid.SetRowSpacing(10) //espacio entre filas
	grid.SetColumnSpacing(10) // espacio entre columnas
	grid.SetMarginTop(20)
	grid.SetMarginBottom(20)
	grid.SetMarginStart(20)
	grid.SetMarginEnd(20)

	// Establece el Grid como el único hijo de la ventana.
	w.SetChild(grid)
	// Crea una etiqueta y la adjunta en la celda (0, 0).
	// La etiqueta ocupa una sola celda.
	label := gtk.NewLabel("Etiqueta principal:")
	grid.Attach(label, 0, 0, 1, 1) // widget, columna, fila, ancho, alto

	// Crea un Entry y lo adjunta en la celda (1, 0).
	// Ocupa una celda de ancho y una de alto.
	entry := gtk.NewEntry()
	grid.Attach(entry, 1, 0, 1, 1)

	// Crea un botón que abarca dos columnas.
	// Se adjunta en la celda (0, 1) y tiene un 'ancho' de 2.
	longButton := gtk.NewButtonWithLabel("Botón que abarca 2 columnas")
	longButton.ConnectClicked(func() {
		// ("¡Botón largo clicado!")
	})
	grid.Attach(longButton, 0, 1, 2, 1)

	// Crea un botón estándar.
	// Se adjunta en la celda (0, 2).
	button1 := gtk.NewButtonWithLabel("Botón 1")
	button1.ConnectClicked(func() {
		// ("¡Botón 1 clicado!")
	})
	grid.Attach(button1, 0, 2, 1, 1)

	// Crea otro botón estándar.
	// Se adjunta en la celda (1, 2).
	button2 := gtk.NewButtonWithLabel("Botón 2")
	button2.ConnectClicked(func() {
		// ("¡Botón 2 clicado!")
	})
	grid.Attach(button2, 1, 2, 1, 1)
}