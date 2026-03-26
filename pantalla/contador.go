package pantalla

import (
	"strconv"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func Contador(w *gtk.ApplicationWindow, ww *gtk.Window) {
	var valor int
	vbox := gtk.NewBox(gtk.OrientationHorizontal, 0)
	w.SetChild(vbox)

	label := gtk.NewLabel("0")
	label.SetMarginTop(20)
	label.SetMarginBottom(20)
	label.SetMarginStart(20)
	label.SetMarginEnd(20)
	label.SetText("Contador")
	vbox.Append(label)

	txtContador := gtk.NewEntry()
	txtContador.SetMarginTop(20)
	txtContador.SetMarginBottom(20)
	txtContador.SetMarginStart(20)
	txtContador.SetMarginEnd(20)
	txtContador.SetText("0")
	txtContador.SetEditable(false)
	vbox.Append(txtContador)

	btnContador := gtk.NewButton()
	btnContador.SetLabel("Contar")
	btnContador.SetMarginTop(20)
	btnContador.SetMarginBottom(20)
	btnContador.SetMarginStart(20)
	btnContador.SetMarginEnd(20)
	btnContador.ConnectClicked(func() {
		valor++
		txtContador.SetText(strconv.Itoa(valor))
	})
	vbox.Append(btnContador)
}