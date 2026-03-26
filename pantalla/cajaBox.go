package pantalla

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/Tobotobo/msgbox"
)

func Caja(w *gtk.ApplicationWindow) {
	cajaVert1 := gtk.NewBox(gtk.OrientationVertical,2)
	lbl1 := gtk.NewLabel("Primer")
	cajaVert1.Append(lbl1)
	txt2 := "Segunda \n segunda linea \n tercera smdsdb. sdldb sdfbsbl dldsv"
	lbl2 := gtk.NewLabel(txt2)
	cajaVert1.Append(lbl2)

	btnMsgB1 := gtk.NewButton()
	btnMsgB1.SetLabel("Mostrar Mensaje Caja Box")
	btnMsgB1.ConnectClicked(func() {
		msgbox.Show("Mensaje Box Simple", "Caja de mensaje desde Caja Box")
	})
	cajaVert1.Append(btnMsgB1)

	btnMsgB2 := gtk.NewButton()
	btnMsgB2.SetLabel("Mostrar Mensaje Caja Box 2")
	btnMsgB2.ConnectClicked(func() {
		//bRes2 := 
		msgbox.Info().AbortRetryIgnore().Show("Información", "Este es un mensaje informativo desde Caja Box 2").IsOK()

	})
	cajaVert1.Append(btnMsgB2)

	btnMsgB3 := gtk.NewButton()
	btnMsgB3.SetLabel("Mostrar Mensaje Caja Box 3")
	btnMsgB3.ConnectClicked(func() {
		// bres3 := 
		msgbox.Err().OK().Show("Mensaje Error", "¿Desea continuar desde Caja Box 3?")
	})
	cajaVert1.Append(btnMsgB3)

	btnMsgB4 := gtk.NewButton()
	btnMsgB4.SetLabel("Mostrar Mensaje Caja Box 4")
	btnMsgB4.ConnectClicked(func() {
		// bres4 := 
		msgbox.Warn().OKCancel().Show("Mensaje Advertencia", "Esta es una advertencia desde Caja Box 4")
	})
	cajaVert1.Append(btnMsgB4)	

	btnMsgB5 := gtk.NewButton()
	btnMsgB5.SetLabel("Mostrar Mensaje Caja Box 5")
	btnMsgB5.ConnectClicked(func() {
		//bres5 := 
		msgbox.Ques().YesNo().Show("Mensaje Pregunta", "¿Desea continuar desde Caja Box 5?")
	})
	cajaVert1.Append(btnMsgB5)

	cajaVert2 := gtk.NewBox(gtk.OrientationVertical,2)
	txt3 := "Tercero \n segunda linea \n tercera smdsdb. sdldb sdfbsbl dldsv"
	lbl3 := gtk.NewLabel(txt3)
	cajaVert2.Append(lbl3)
	lbl4 := gtk.NewLabel("Cuatro")
	cajaVert2.Append(lbl4)

	//check Buttons
	check1 := gtk.NewCheckButton()
	check1.SetLabel("Opción 1")
	cajaVert2.Append(check1)

	check2 := gtk.NewCheckButton()
	check2.SetLabel("Opción 2")
	cajaVert2.Append(check2)

	cajaHori := gtk.NewBox(gtk.OrientationHorizontal, 5)
	cajaHori.Append(cajaVert1)
	cajaHori.Append(cajaVert2)

	w.SetChild(cajaHori)
}