package pantalla

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

// MostrarError presenta un diálogo modal de error asociado a la ventana padre.
func MostrarError(parent *gtk.Window, mensaje string) {
	dialog := gtk.NewMessageDialog(
		parent,
		gtk.DialogModal,
		gtk.MessageError,
		gtk.ButtonsOK,
	)
	dialog.SetTitle("Error")
	dialog.SetObjectProperty("text", "Ocurrió un error")
	dialog.SetObjectProperty("secondary-text", mensaje)
	dialog.SetIconName("dialog-error")
	dialog.ConnectResponse(func(int) {
		dialog.Destroy()
	})
	dialog.Present()
}
