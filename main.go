package main

//Para compilar y copiar la carpeta
//  go build -o appgtk.exe -ldflags="-H=windowsgui" main.go

import (
	"appgtk/pantalla"
	_ "embed"
	"fmt"
	"os"
	"strconv"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed style.css
var cssData string

func main() {
	// 1. Configuración de entorno (Debe ir antes de NewApplication)
	dir, _ := os.Getwd()
	os.Setenv("XDG_DATA_DIRS", dir + "\\share")
	os.Setenv("GSK_RENDERER", "cairo") 
	os.Setenv("GTK_THEME", "Adwaita:light")

	// os.Setenv("GTK_DEBUG", "interactive") // Habilita el modo interactivo para depuración

	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.builder", gio.ApplicationFlagsNone)
	app.ConnectStartup(loadCss())
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		fmt.Println("Error al iniciar la aplicación:", strconv.Itoa(code))
		os.Exit(code)
	}
}

func loadCss() func() { 
	return func() {
		cssProvider := gtk.NewCSSProvider()
		cssProvider.LoadFromString(cssData) 
		gtk.StyleContextAddProviderForDisplay(
			gdk.DisplayGetDefault(),
			cssProvider,
			gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
		)
	}
}

func activate(app *gtk.Application) {
	// --- NUEVO: Configuración de Iconos ---
	// Usamos el display actual para obtener los settings y el tema
	display := gdk.DisplayGetDefault()
	if display != nil {
		// Configuramos el nombre del tema
		settings := gtk.SettingsGetForDisplay(display)
		settings.SetObjectProperty("gtk-icon-theme-name", "Adwaita")

		// Forzamos a GTK a buscar en tu carpeta local 'share/icons'
		iconTheme := gtk.IconThemeGetForDisplay(display)
		iconTheme.AddSearchPath("share\\icons") 
		// fmt.Println("Rutas de iconos:", iconTheme.SearchPath()) // Agrega este print para depurar
	}else{
		fmt.Println("**** No se pudo obtener el display para configurar los iconos.")
	}

	// Mostrar la ventana de login
	pantalla.ShowLogin(app ,func() {
		// 1 - Crear la ventana principal
		win := gtk.NewApplicationWindow(app)

		win.SetTitle("Ejemplo de la librería GTK")
		win.SetDefaultSize(800, 700)

		// 2 - Acciones del menú
		accionGrilla := gio.NewSimpleAction("grilla", nil)
		accionGrilla.ConnectActivate(func(_ *glib.Variant) {
			pantalla.Grilla(win)
		})
		win.AddAction(accionGrilla)

		accionPanel := gio.NewSimpleAction("panel", nil)
		accionPanel.ConnectActivate(func(_ *glib.Variant) {
			pantalla.Panel(win)
		})
		win.AddAction(accionPanel)

		accionCaja := gio.NewSimpleAction("caja", nil)
		accionCaja.ConnectActivate(func(_ *glib.Variant) {
			pantalla.Caja(win)
		})
		win.AddAction(accionCaja)

		accionCerrar := gio.NewSimpleAction("cerrar", nil)
		accionCerrar.ConnectActivate(func(_ *glib.Variant) {
			win.Close()
		})
		win.AddAction(accionCerrar)

		// 3 - Crear el modelo de Menú
		mnuArchivo := gio.NewMenu()
		mnuArchivo.Append("Grilla - Grid", "win.grilla")
		mnuArchivo.Append("Panel - Paned", "win.panel")
		mnuArchivo.Append("Caja - Box", "win.caja")
		//mnuArchivo.AppendSeparator()
		mnuArchivo.Append("Cerrar", "win.cerrar")

		// Segunda opcion de crear el menú
		w1 := gtk.NewWindow()
		accionMsjDialogo := gio.NewSimpleAction("msjDialogo", nil)
		accionMsjDialogo.ConnectActivate(func(_ *glib.Variant) {
			pantalla.MsjDialogo(win, w1)
		})
		win.AddAction(accionMsjDialogo)

		accionTextoEntry := gio.NewSimpleAction("textoEntry", nil)
		accionTextoEntry.ConnectActivate(func(_ *glib.Variant) {
			pantalla.CuadroEntry(win, w1)
		})
		win.AddAction(accionTextoEntry)

		mnuWidget := gio.NewMenu()
		mnuWidget.Append("Mensaje Dialogo", "win.msjDialogo")
		mnuWidget.Append("Cuadro de Texto Entry", "win.textoEntry")

		//tercera opcion de crear el menú
		accionContador := gio.NewSimpleAction("contador", nil)
		accionContador.ConnectActivate(func(_ *glib.Variant) {
			pantalla.Contador(win, w1)
		})
		win.AddAction(accionContador)

		accionListBox := gio.NewSimpleAction("listBox", nil)
		accionListBox.ConnectActivate(func(_ *glib.Variant) {
			pantalla.ListBox(win)
		})
		win.AddAction(accionListBox)

		accionEstiloCSS := gio.NewSimpleAction("estiloCSS", nil)
		accionEstiloCSS.ConnectActivate(func(_ *glib.Variant) {
			pantalla.EstiloCSS(win)
		})
		win.AddAction(accionEstiloCSS)

		accionTabs := gio.NewSimpleAction("tabs", nil)
		accionTabs.ConnectActivate(func(_ *glib.Variant) {
			pantalla.Tabs(win, w1)
		})
		win.AddAction(accionTabs)

		accionPDF := gio.NewSimpleAction("pdf", nil)
		accionPDF.ConnectActivate(func(_ *glib.Variant) {
			pantalla.PDFVista(win)
		})
		win.AddAction(accionPDF)

		mnuEjemplos := gio.NewMenu()
		mnuEjemplos.Append("Contador", "win.contador")
		mnuEjemplos.Append("ListBox", "win.listBox")
		mnuEjemplos.Append("EstiloCSS", "win.estiloCSS")
		mnuEjemplos.Append("Tabs", "win.tabs")
		mnuEjemplos.Append("PDF", "win.pdf")

		accionCuentas := gio.NewSimpleAction("cuentas", nil)
		accionCuentas.ConnectActivate(func(_ *glib.Variant) {
			pantalla.Cuentas(win)
		})
		win.AddAction(accionCuentas)

		accionMovimientoContable := gio.NewSimpleAction("movimientoContable", nil)
		accionMovimientoContable.ConnectActivate(func(_ *glib.Variant) {
			pantalla.MovimientoContable(win)
		})
		win.AddAction(accionMovimientoContable)

		accionInformePDF := gio.NewSimpleAction("informePDF", nil)
		accionInformePDF.ConnectActivate(func(_ *glib.Variant) {
			pantalla.InformePDF(win)
		})
		win.AddAction(accionInformePDF)

		mnuContable := gio.NewMenu()
		mnuContable.Append("Cuentas", "win.cuentas")
		mnuContable.Append("Movimiento Contable", "win.movimientoContable")
		mnuContable.Append("Informe PDF", "win.informePDF")
		
		// 4 - Crear la barra de menú principal
		menubar := gio.NewMenu()
		menubar.AppendSubmenu("Archivo", mnuArchivo)
		menubar.AppendSubmenu("Widget", mnuWidget)
		menubar.AppendSubmenu("Ejemplos", mnuEjemplos)
		menubar.AppendSubmenu("Contable", mnuContable)

		// 5 - Establece la barra de menú en la ventana
		app.SetMenubar(menubar)
		win.SetShowMenubar(true)

		lbl := gtk.NewLabel("Esta es un label como hijo")
		win.SetChild(lbl)

		// 6 - Presenta la ventana
		win.Present()
	})
}

func Salir() {
	// Falta completar
}
