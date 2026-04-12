package pantalla

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/joseamadeo/appgtk/entradas"
)

func ValidarWidget(wa *gtk.ApplicationWindow, ww *gtk.Window) {
	vbox := gtk.NewBox(gtk.OrientationVertical, 0)
	lblTitulo := gtk.NewLabel("Validar Widget")
	vbox.Append(lblTitulo)

	hboxCuentaContable :=gtk.NewBox(gtk.OrientationHorizontal, 10)
	lblCC := gtk.NewLabel("Cuenta contable")
	txtCC := entradas.NewCuentaContable()
	hboxCuentaContable.Append(lblCC)
	hboxCuentaContable.Append(txtCC)
	vbox.Append(hboxCuentaContable)

	hboxCuit := gtk.NewBox(gtk.OrientationHorizontal, 10)
	lblCuit := gtk.NewLabel(("CUTI"))
	txtCuit := entradas.NewCuit()
	hboxCuit.Append(lblCuit)
	hboxCuit.Append(txtCuit)
	vbox.Append(hboxCuit)

	hboxFecha :=gtk.NewBox(gtk.OrientationHorizontal, 10)
	lblFecha := gtk.NewLabel("Fecha")
	txtFecha := entradas.NewFecha()
	hboxFecha.Append(lblFecha)
	hboxFecha.Append(txtFecha)
	vbox.Append(hboxFecha)

	hboxBotones := gtk.NewBox(gtk.OrientationHorizontal,10)
	hboxBotones.SetVAlign(gtk.AlignBaselineCenter)
	btnAceptar := gtk.NewButton()
	btnAceptar.SetLabel("Validar Campos")
	hboxBotones.Append(btnAceptar)
	vbox.Append(hboxBotones)

	btnAceptar.ConnectClicked(func ()  {
				// create a message dialog attached to window 'ww'
		dialog := gtk.NewMessageDialog(
			ww,
			gtk.DialogModal,
			gtk.MessageInfo,
			gtk.ButtonsOKCancel,
		)
		msj := "Estos son errores"
		// Verificar las validaciones
		if txtCC.EsValido() == false {
			msj += "\n - El Cuenta Contable esta incorrecto"
		}
		if txtCuit.EsValido() == false {
			msj += "\n - El CUIT esta incorrecto"
		}
		if txtFecha.EsValido() == false {
			msj += "\n - La fecha es incorrecta"
		}
		if msj == "Estos son errores" {
			msj = "NO hay errores"
		}

		dialog.SetTitle("Estos Validar.")
		msg := gtk.NewLabel(msj)
		dialog.ContentArea().Append(msg)
		dialog.SetIconName("dialog-information")

		dialog.ConnectResponse(func(r int) {
			if r == -5 { // GTK_RESPONSE_OK

			} else {

			}
			dialog.Destroy()
		})
		dialog.Present()
	})

	wa.SetChild(vbox)

}
