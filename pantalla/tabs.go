package pantalla

import (
	"appgtk/entradas"
	"fmt"
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func Tabs(w *gtk.ApplicationWindow, ww *gtk.Window) {
	notebook := gtk.NewNotebook()
	// --- Tab 1 ---
	page1, label1 := crearTab1()
	notebook.AppendPage(page1, label1)
	// --- Tab 2 ---
	page2, label2 := crearTab2()
	notebook.AppendPage(page2, label2)
	// --- Tab 3 ---
	page3, label3 := crearTab3()
	notebook.AppendPage(page3, label3)
	// --- Tab 4 ---
	page4, label4 := crearTab4()
	notebook.AppendPage(page4, label4)
	// -- Tab 5 ---
	page5, label5 := crearTab5()
	notebook.AppendPage(page5, label5)
	// -- Tab 6 ---
	page6, label6 := crearTab6()
	notebook.AppendPage(page6, label6)

	w.SetChild(notebook)	
}

func crearTab1() (contenido gtk.Widgetter, label gtk.Widgetter) {
	labelWidget := gtk.NewLabel("Fechas Calendario")
	principalBox := gtk.NewBox(gtk.OrientationVertical, 10)
	principalBox.SetMarginTop(10)
	principalBox.SetMarginBottom(10)
	principalBox.SetMarginStart(10)
	principalBox.SetMarginEnd(10)

	calendario := gtk.NewCalendar()
	principalBox.Append(calendario)
	lblFecha := gtk.NewLabel("Selecciona una fecha: Ninguna")
	principalBox.Append(lblFecha)
	btnSelectDate := gtk.NewButtonWithLabel("Obtener fecha seleccionada")
	btnSelectDate.ConnectClicked(func() {
		dateTime := calendario.Date() // GtkCalendar.Date() devuelve un *gio.DateTime
		// Convertir GDateTime a un objeto time.Time de Go
		// El *gio.DateTime de gotk4 tiene un método GoTime() para esto
		año := dateTime.Year()
		mes := time.Month(dateTime.Month())
		día := dateTime.DayOfMonth()
		goTime := time.Date(año, time.Month(mes), día, 0, 0, 0, 0, time.Local)	
		// Formatear la fecha como una cadena legible (ejemplo de formato de Go: YYYY-MM-DD)
		//formatoFecha := goTime.Format("2 January 2006")
		//formatoFecha := goTime.Format("02/01/2006")
		formatoFecha := goTime.Format("20060102")
		lblFecha.SetText(fmt.Sprintf("Fecha seleccionada: %s", formatoFecha))
	})
	principalBox.Append(btnSelectDate)

	return principalBox, labelWidget
}

func crearTab2() (contenido gtk.Widgetter, label gtk.Widgetter) {
	labelWidget := gtk.NewLabel("Grid - Tabla")
	contenidoWidget := gtk.NewGrid()
	// Establecer espaciado entre filas y columnas (opcional)
	contenidoWidget.SetRowSpacing(6)   
	contenidoWidget.SetColumnSpacing(6)

	// Crear algunos botones para añadir a la cuadrícula
	button1 := gtk.NewButtonWithLabel("Botón 1")
	button2 := gtk.NewButtonWithLabel("Botón 2")
	button3 := gtk.NewButtonWithLabel("Botón 3 (ocupa 2 filas)")
	button4 := gtk.NewButtonWithLabel("Botón 4 (ocupa 2 columnas)")

	// Adjuntar los widgets a la cuadrícula usando Attach(widget, columna, fila, ancho, alto)
	// Botón 1 en (0, 0), ocupa 1 columna, 1 fila
	contenidoWidget.Attach(button1, 0, 0, 1, 1)

	// Botón 2 en (1, 0), ocupa 1 columna, 1 fila
	contenidoWidget.Attach(button2, 1, 0, 1, 1)

	// Botón 3 en (2, 0), ocupa 1 columna, 2 filas
	contenidoWidget.Attach(button3, 2, 0, 1, 2)

	// Botón 4 en (0, 1), ocupa 2 columnas, 1 fila (debajo del Botón 1 y Botón 2)
	contenidoWidget.Attach(button4, 0, 1, 2, 1)

	lblFecha := gtk.NewLabel("Ingrese una fecha (DD/MM/AAAA):")
	contenidoWidget.Attach(lblFecha, 0, 2, 2, 1)

	fechaEntry := entradas.NewFecha()
	contenidoWidget.Attach(fechaEntry, 2, 2, 1, 1)

	lblResultado := gtk.NewLabel("Fecha ingresada: Ninguna")
	btnTab2ObtenerFecha:= gtk.NewButtonWithLabel("Obtener fecha ingresada")
	btnTab2ObtenerFecha.ConnectClicked(func() {
		entrada, err := fechaEntry.GetFecha()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		lblResultado.SetText(fmt.Sprintf("Fecha ingresada: %s", entrada))
	})
	contenidoWidget.Attach(btnTab2ObtenerFecha, 0, 3, 1, 1)
	contenidoWidget.Attach(lblResultado, 1, 3, 2, 1)

	lblSetFecha := gtk.NewLabel("Establecer fecha a (DD/MM/AAAA):")
	contenidoWidget.Attach(lblSetFecha, 0, 4, 1, 1)
	fechaSetEntry := gtk.NewEntry()
	fechaSetEntry.SetPlaceholderText("DD/MM/AAAA")
	fechaSetEntry.SetMaxLength(10)
	contenidoWidget.Attach(fechaSetEntry, 1, 4, 1, 1)
	btnTab2SetFecha:= gtk.NewButtonWithLabel("Establecer fecha")
	btnTab2SetFecha.ConnectClicked(func() {
		entrada := fechaSetEntry.Text()
		fechaEntry.SetFecha(entrada)
	})
	contenidoWidget.Attach(btnTab2SetFecha, 2, 4, 1, 1)

	lblFechaDB := gtk.NewLabel("Fecha para BD (AAAA-MM-DD):")
	contenidoWidget.Attach(lblFechaDB, 0, 5, 2, 1)
	btnObtenerFechaDB:= gtk.NewButtonWithLabel("Obtener fecha para BD")
	btnObtenerFechaDB.ConnectClicked(func() {
		entrada, err := fechaEntry.GetFechaDB()
		if err != nil {
			lblFechaDB.SetText(fmt.Sprintf("Error al obtener fecha para BD: %s", err))
			return
		}
		lblFechaDB.SetText(fmt.Sprintf("Fecha para BD (AAAA-MM-DD): %s", entrada))
	})
	contenidoWidget.Attach(btnObtenerFechaDB, 2, 5, 1, 1)

	return contenidoWidget, labelWidget
}

func crearTab3() (contenido gtk.Widgetter, label gtk.Widgetter) {
	labelWidget := gtk.NewLabel("Widgets Propios")
	contenidoWidget := gtk.NewGrid()
	// Establecer espaciado entre filas y columnas (opcional)
	contenidoWidget.SetRowSpacing(6)   
	contenidoWidget.SetColumnSpacing(6)
	lblCuit := gtk.NewLabel("Ingrese una CUIT (XX-XXXXXXXX-X):")
	contenidoWidget.Attach(lblCuit, 0, 2, 2, 1)

	cuitEntry := entradas.NewCuit()
	contenidoWidget.Attach(cuitEntry, 2, 2, 1, 1)
	lblResultado := gtk.NewLabel("CUIT ingresada: Ninguna")
	btnTab2ObtenerFecha:= gtk.NewButtonWithLabel("Obtener CUIT ingresada")
	btnTab2ObtenerFecha.ConnectClicked(func() {
		entrada := cuitEntry.GetCuit()
		lblResultado.SetText(fmt.Sprintf("CUIT ingresada: %s", entrada))
	})
	contenidoWidget.Attach(btnTab2ObtenerFecha, 0, 3, 1, 1)
	contenidoWidget.Attach(lblResultado, 1, 3, 2, 1)

	lblSetCUIT := gtk.NewLabel("Establecer CUIT a (XX-XXXXXXXX-X):")
	contenidoWidget.Attach(lblSetCUIT, 0, 4, 1, 1)
	cuitSetEntry := gtk.NewEntry()
	cuitSetEntry.SetPlaceholderText("XX-XXXXXXXX-X")
	cuitSetEntry.SetMaxLength(13)
	contenidoWidget.Attach(cuitSetEntry, 1, 4, 1, 1)
	btnTab2SetCUIT:= gtk.NewButtonWithLabel("Establecer CUIT")
	btnTab2SetCUIT.ConnectClicked(func() {
		entrada := cuitSetEntry.Text()
		cuitEntry.SetCuit(entrada)
	})
	contenidoWidget.Attach(btnTab2SetCUIT, 2, 4, 1, 1)

	lblCUITSinFormato := gtk.NewLabel("CUIT sin formato (XXXXXXXXXXX):")
	contenidoWidget.Attach(lblCUITSinFormato, 0, 5, 2, 1)
	btnObtenerCUITSinFormato:= gtk.NewButtonWithLabel("Obtener CUIT sin formato")
	btnObtenerCUITSinFormato.ConnectClicked(func() {
		entrada := cuitEntry.GetCuitSinFormato()
		lblCUITSinFormato.SetText(fmt.Sprintf("CUIT sin formato (XXXXXXXXXXX): %s", entrada))
	})
	contenidoWidget.Attach(btnObtenerCUITSinFormato, 2, 5, 1, 1)
	return contenidoWidget, labelWidget
}

func crearTab4() (contenido gtk.Widgetter, label gtk.Widgetter) {
	labelWidget := gtk.NewLabel("Nombre Propio")
	contenidoWidget := gtk.NewGrid()
	// Establecer espaciado entre filas y columnas (opcional)
	contenidoWidget.SetRowSpacing(6)   
	contenidoWidget.SetColumnSpacing(6)
	// Función auxiliar para un Attach más limpio
	attach := func(c gtk.Widgetter, l, t, w, h int) {
			contenidoWidget.Attach(c, l, t, w, h)
	}
	
	lblNombre := gtk.NewLabel("Ingrese un Nombre Propio:")
	attach(lblNombre, 0, 0, 1, 1)
	nombreEntry := entradas.NewNombrePropio("Ingrese nombre...", true, 50)
	attach(nombreEntry, 1, 0, 2, 1)
	
	lblResultado := gtk.NewLabel("Nombre ingresado: Ninguno")
	btnObtenerNombre:= gtk.NewButtonWithLabel("Obtener Nombre ingresado")
	btnObtenerNombre.ConnectClicked(func() {
		entrada := nombreEntry.GetTexto()
		lblResultado.SetText(fmt.Sprintf("Nombre ingresado: %s", entrada))
	})
	attach(btnObtenerNombre, 0, 3, 1, 1)
	attach(lblResultado, 1, 3, 2, 1)

	txtIngNombre := gtk.NewEntry()
	txtIngNombre.SetPlaceholderText("Ingrese un nombre...")
	txtIngNombre.SetMaxLength(50)
	btnGetNombre:= gtk.NewButtonWithLabel("Mostrar nombre ingresado widget")
	btnGetNombre.ConnectClicked(func() {
		nombreEntry.SetTexto(txtIngNombre.Text())
	})
	attach(txtIngNombre, 0, 4, 1, 1)
	attach(btnGetNombre, 1, 4, 1, 1)

	lblMayusculas := gtk.NewLabel("mayúsculas automáticamente.")
	attach(lblMayusculas, 0, 5, 1, 1)
	switchMayusculas := gtk.NewSwitch()
	switchMayusculas.SetActive(true)
	attach(switchMayusculas, 1, 5, 1, 1)
	switchMayusculas.Connect("notify::active",func(s *gtk.Switch) {
		state := s.Active()
		nombreEntry.SoloMayusculas(state)
	})

	return contenidoWidget, labelWidget
}

func crearTab5() (contenido gtk.Widgetter, label gtk.Widgetter) {
	labelWidget := gtk.NewLabel("Entrada Numérica")
	contenidoWidget := gtk.NewGrid()
	// Establecer espaciado entre filas y columnas (opcional)
	contenidoWidget.SetRowSpacing(6)   
	contenidoWidget.SetColumnSpacing(6)
	lblNumero := gtk.NewLabel("Ingrese un número (formato #.###,00)")
	contenidoWidget.Attach(lblNumero, 0, 2, 2, 1)
	numeroEntry := entradas.NewEntradaNumerica("#.###,00")
	contenidoWidget.Attach(numeroEntry, 2, 2, 1, 1)
	lblResultado := gtk.NewLabel("Número ingresado: Ninguno")
	btnObtenerNumero:= gtk.NewButtonWithLabel("Obtener número ingresado")
	btnObtenerNumero.ConnectClicked(func() {
		entrada := numeroEntry.Text()
		lblResultado.SetText(fmt.Sprintf("Número ingresado: %s", entrada))
	})
	contenidoWidget.Attach(btnObtenerNumero, 0, 3, 1, 1)
	contenidoWidget.Attach(lblResultado, 1, 3, 2, 1)

	// formato de número para base de datos
	lblNumeroDB := gtk.NewLabel("Número para BD (formato sin miles, punto decimal):")
	btnObtenerNumeroDB:= gtk.NewButtonWithLabel("Obtener número para BD")
	btnObtenerNumeroDB.ConnectClicked(func() {
		entrada := numeroEntry.FormatoDB()
		lblNumeroDB.SetText(fmt.Sprintf("Número para BD (formato sin miles, punto decimal): %s", entrada))
	})
	contenidoWidget.Attach(lblNumeroDB, 2, 4, 2, 1)
	contenidoWidget.Attach(btnObtenerNumeroDB, 0, 4, 1, 1)

	// Cambio de formato
	txtFormato := gtk.NewEntry()
	txtFormato.SetPlaceholderText("Ingrese formato (#.### o #.###,00)")
	txtFormato.SetMaxLength(10)
	btnCambiarFormato:= gtk.NewButtonWithLabel("Cambiar formato")
	btnCambiarFormato.ConnectClicked(func() {
		lblNumero.SetText("Ingrese un número (formato " + txtFormato.Text() + ")")
		numeroEntry.SetFormato(txtFormato.Text())
	})
	contenidoWidget.Attach(btnCambiarFormato, 0, 5, 1, 1)
	contenidoWidget.Attach(txtFormato, 2, 5, 1, 1)

	return contenidoWidget, labelWidget
}

// Tab 6: Cuenta Contable con formato dinámico
func crearTab6() (contenido gtk.Widgetter, label gtk.Widgetter) {
	labelWidget := gtk.NewLabel("Cuenta Contable")
	contenidoWidget := gtk.NewGrid()
	// Establecer espaciado entre filas y columnas (opcional)
	contenidoWidget.SetRowSpacing(6)
	contenidoWidget.SetColumnSpacing(6)
	lblCuenta := gtk.NewLabel("Ingrese una cuenta contable (formato 0-00-00-00):")
	contenidoWidget.Attach(lblCuenta, 0, 2, 2, 1)
	cuentaEntry := entradas.NewCuentaContable()
	contenidoWidget.Attach(cuentaEntry.Box, 2, 2, 1, 1)
	lblResultado := gtk.NewLabel("Cuenta ingresada: Ninguna")
	btnObtenerCuenta := gtk.NewButtonWithLabel("Obtener cuenta ingresada")
	btnObtenerCuenta.ConnectClicked(func() {
		entrada := cuentaEntry.GetCuenta()
		lblResultado.SetText(fmt.Sprintf("Cuenta ingresada: %s", entrada))
	})
	contenidoWidget.Attach(btnObtenerCuenta, 0, 3, 1, 1)
	contenidoWidget.Attach(lblResultado, 1, 3, 2, 1)
	lblCuentaSinFormato := gtk.NewLabel("Cuenta sin formato (XXXXXXXXXX):")
	contenidoWidget.Attach(lblCuentaSinFormato, 0, 4, 2, 1)
	btnObtenerCuentaSinFormato := gtk.NewButtonWithLabel("Obtener cuenta sin formato")
	btnObtenerCuentaSinFormato.ConnectClicked(func() {
		entrada := cuentaEntry.GetCuentaSinFormato()
		lblCuentaSinFormato.SetText(fmt.Sprintf("Cuenta sin formato (XXXXXXXXXX): %s", entrada))
	})
	contenidoWidget.Attach(btnObtenerCuentaSinFormato, 2, 4, 1, 1)
	return contenidoWidget, labelWidget
}