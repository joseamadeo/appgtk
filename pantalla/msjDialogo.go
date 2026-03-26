package pantalla

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func MsjDialogo(wa *gtk.ApplicationWindow, ww *gtk.Window){
	btnDialogo := gtk.NewButton()
	btnDialogo.SetLabel("Mostrar Dialogo de Mensaje")
	btnDialogo.ConnectClicked(func() {
		// create a message dialog attached to window 'ww'
		dialog := gtk.NewMessageDialog(
			ww,
			gtk.DialogModal,
			gtk.MessageInfo,
			gtk.ButtonsOKCancel,
		)
		dialog.SetTitle("Este es un diálogo de mensaje de información.")
		msg := gtk.NewLabel("Este es el contenido del diálogo de mensaje.\nPuede contener varias líneas de texto.")
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
	
	wa.SetChild(btnDialogo)
}