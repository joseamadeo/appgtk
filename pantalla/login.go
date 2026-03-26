package pantalla

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

func ShowLogin(app *gtk.Application, onSuccess func()) {
	win := gtk.NewApplicationWindow(app)
	win.SetTitle("Login")
	win.SetDefaultSize(300, 200)
	box := gtk.NewBox(gtk.OrientationVertical, 10)
	box.SetMarginTop(20)
	box.SetMarginBottom(20)
	box.SetMarginStart(20)
	box.SetMarginEnd(20)

	user := gtk.NewEntry()
	user.SetPlaceholderText("Usuario")
	user.SetText("admin") // Para pruebas, puedes eliminar esta línea en producción

	pass := gtk.NewPasswordEntry()
	pass.SetText("1234") // Para pruebas, puedes eliminar esta línea en producción

	btn := gtk.NewButtonWithLabel("Login")
	btn.ConnectClicked(func() {
		// Simulación de validación
		if user.Text() == "admin" && pass.Text() == "1234" {
			win.Close() // Cerrar la ventana de login
			onSuccess()  // Llamar a la función de éxito
		} else {
			user.SetText("")
			pass.SetText("")
			// create a message dialog attached to window 'ww'
			dialog := gtk.NewMessageDialog(
				&win.Window,
				gtk.DialogModal,
				gtk.MessageError,
				gtk.ButtonsOK,
			)
			dialog.SetIconName("dialog-error")
			dialog.SetObjectProperty("text", "Error de acceso")
			dialog.SetObjectProperty("secondary-text", "Usuario o contraseña incorrectos.")
			dialog.ConnectResponse(func(responseId int) {
				dialog.Destroy()
			})
			dialog.Present()
		}
	})

	box.Append(user)
	box.Append(pass)
	box.Append(btn)
	win.SetChild(box)

	win.Present()
}