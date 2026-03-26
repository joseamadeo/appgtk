package pantalla

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/Tobotobo/msgbox"
)

func EstiloCSS(w *gtk.ApplicationWindow) {
	vbox := gtk.NewBox(gtk.OrientationVertical, 0)
	w.SetChild(vbox)

	linea1 := gtk.NewBox(gtk.OrientationHorizontal, 0)
	vbox.Append(linea1)
	linea2 := gtk.NewBox(gtk.OrientationHorizontal, 0)
	vbox.Append(linea2)
	linea3 := gtk.NewBox(gtk.OrientationHorizontal, 0)
	vbox.Append(linea3)

  btnPrimary := gtk.NewButtonWithLabel("Guardar (btn-primary)")
  btnPrimary.AddCSSClass("btn-primary")
	btnPrimary.ConnectClicked(func() {
		msgbox.Show("Hice clic en Guardar")
	})
  linea1.Append(btnPrimary)
	
	btnSecundario := gtk.NewButtonWithLabel("Cancelar (btn-secondary)")
	btnSecundario.AddCSSClass("btn-secondary")
	btnSecundario.ConnectClicked(func() {	
		msgbox.Show("Hice clic en Cancelar")
	})
	linea1.Append(btnSecundario)
	
	btnExito := gtk.NewButtonWithLabel("Exito (btn-exito)")
	btnExito.AddCSSClass("btn-success")
	btnExito.ConnectClicked(func() {
		msgbox.Show("Hice clic en Exito")
	})
	linea1.Append(btnExito)

	btnPeligro := gtk.NewButtonWithLabel("Peligro (btn-danger)")
	btnPeligro.AddCSSClass("btn-danger")
	btnPeligro.ConnectClicked(func() {
		msgbox.Show("Hice clic en Peligro")
	})
	linea2.Append(btnPeligro)

	btnAlerta := gtk.NewButtonWithLabel("Alerta (btn-warning)")
	btnAlerta.AddCSSClass("btn-warning")
	btnAlerta.ConnectClicked(func() {
		msgbox.Show("Hice clic en Alerta")
	})
	linea2.Append(btnAlerta)

	btnInfo := gtk.NewButtonWithLabel("Info (btn-info)")
	btnInfo.AddCSSClass("btn-info")
	btnInfo.ConnectClicked(func() {
		msgbox.Show("Hice clic en Info")
	})
	linea2.Append(btnInfo)

	btnLiviano := gtk.NewButtonWithLabel("Liviano (btn-light)")
	btnLiviano.AddCSSClass("btn-light")
	btnLiviano.ConnectClicked(func() {
		msgbox.Show("Hice clic en Liviano")
	})
	linea3.Append(btnLiviano)

	btnOscuro := gtk.NewButtonWithLabel("Oscuro (btn-dark)")
	btnOscuro.AddCSSClass("btn-dark")
	btnOscuro.ConnectClicked(func() {
		msgbox.Show("Hice clic en Oscuro")
	})
	linea3.Append(btnOscuro)
}
