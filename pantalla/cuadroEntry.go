package pantalla

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func CuadroEntry(w *gtk.ApplicationWindow, ww *gtk.Window) {
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
	txtEntry := gtk.NewEntry()
	grid.Attach(txtEntry, 0, 0, 1, 1) // widget, columna, fila, ancho, alto

	chk1 := gtk.NewCheckButton()
	chk1.SetLabel("Opción 1")
	grid.Attach(chk1, 0, 1, 1, 1)

	chk2 := gtk.NewCheckButton()
	chk2.SetLabel("Opción 2")
	grid.Attach(chk2, 0, 2, 1, 1)

	chk3 := gtk.NewCheckButton()
	chk3.SetLabel("Opción 3")
	grid.Attach(chk3, 0, 3, 1, 1)

	radio1 := gtk.NewCheckButton()
	radio1.SetLabel("Elección A")
	radio1.SetActive(true)
	grid.Attach(radio1, 1, 1, 1, 1)

	radio2 := gtk.NewCheckButton()
	radio2.SetLabel("Elección B")
	radio2.SetGroup(radio1)
	grid.Attach(radio2, 1, 2, 1, 1)

	radio3 := gtk.NewCheckButton()
	radio3.SetLabel("Elección C")
	radio3.SetGroup(radio1)
	grid.Attach(radio3, 1, 3, 1, 1)

	btnAceptar := gtk.NewButton()
	btnAceptar.SetLabel("Aceptar")
	grid.Attach(btnAceptar, 0, 4, 1, 1)
	btnAceptar.ConnectClicked(func() {
		seleccionados := "Seleccionaste:\n"
		//seleccionados += .Sprintf("Entrada: %s\n", txtEntry.Text())
		if chk1.Active() {
			seleccionados += "Opción 1\n"
		}
		if chk2.Active() {
			seleccionados += "Opción 2\n"
		}
		if chk3.Active() {
			seleccionados += "Opción 3\n"
		}	
		if radio1.Active() {
			seleccionados += "Elección A\n"
		}
		if radio2.Active() {
			seleccionados += "Elección B\n"
		}
		if radio3.Active() {
			seleccionados += "Elección C\n"
		}

		dialog := gtk.NewMessageDialog(
			ww,
			gtk.DialogModal,
			gtk.MessageInfo,
			gtk.ButtonsOKCancel,
		)
		dialog.SetTitle("Estos son los seleccionados.")
		msg := gtk.NewLabel(seleccionados)
		dialog.ContentArea().Append(msg)
		dialog.SetIconName("dialog-information")
		dialog.ConnectResponse(func(r int) {
			if r == -5 { // GTK_RESPONSE_OK
				//("Se presionó OK")
			} else {
				//("Se presionó Cancelar")
			}
			dialog.Destroy()
		})

		dialog.Present()
	})
}