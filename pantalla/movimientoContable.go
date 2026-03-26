package pantalla

import (
	"appgtk/entradas"
	"appgtk/estructuras"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type movimiento struct {
	NumeroCuenta 	string
	NombreCuenta 	string
	Fecha     		time.Time
	Importe   		string
	Tipo					string // "I" para Ingreso, "E" para Egreso
}

type tablaMovimientos[T any]struct{
	estructuras.TablaGenerica[T]
	UltimoIdxBusqueda int
}

var listaMovimientos = []movimiento{}

var mapaCuentas = make(map[uint]estructuras.CuentaInfo)

func MovimientoContable(w *gtk.ApplicationWindow) {
	cajaPrincipal := gtk.NewBox(gtk.OrientationVertical, 10)
	cajaPrincipal.SetMarginStart(10)
	cajaPrincipal.SetMarginEnd(10)
	cajaPrincipal.SetMarginTop(10)
	cajaPrincipal.SetMarginBottom(10)

	lbl1 := gtk.NewLabel("Cuentas Contables")
	cajaPrincipal.Append(lbl1)
	// --- 1. WIDGET TABLA (ColumnView) ---
	// Aquí iría la lógica de tu ColumnView para mostrar registros existentes
	colsCuentas := []estructuras.ColumnaConfig[movimiento]{
    {Titulo: "N° Cuenta", Extraer: func(m movimiento) string { return m.NumeroCuenta }},
    {Titulo: "Nombre Cuenta", Extraer: func(m movimiento) string { return m.NombreCuenta }},
    {Titulo: "Fecha", Extraer: func(m movimiento) string { return m.Fecha.Format("02/01/2006") }}, // o m.Fecha.Format(...) si es time.Time
    {Titulo: "Importe", Extraer: func(m movimiento) string { return m.Importe }},
    {Titulo: "Tipo", Extraer: func(m movimiento) string { return m.Tipo }},
	}
	lista := []movimiento{} // Aquí cargarías tus movimientos desde un archivo o base de datos
	tabla := entradas.NewTabla(lista, colsCuentas)
	inicializarArchivoYTabla(tabla)
	// --- 2. FORMULARIO ---
	formGrid := gtk.NewGrid()
	formGrid.SetColumnSpacing(20)
	formGrid.SetRowSpacing(15)
	// Fila 1 : Cuenta Contable, Importe
	// Columna 0: Cuenta
	vboxCuenta := gtk.NewBox(gtk.OrientationVertical, 5)
	lblCuenta := gtk.NewLabel("Cuenta Contable:")
	cmbCuenta := gtk.NewDropDown(nil, nil)
	actualizarComboCuentas(cmbCuenta)
	vboxCuenta.Append(lblCuenta)
	vboxCuenta.Append(cmbCuenta)
	formGrid.Attach(vboxCuenta, 0, 0, 1, 1)
	// columna 1 Importe
	vboxImporte := gtk.NewBox(gtk.OrientationVertical, 5)
	lblImporte := gtk.NewLabel("Importe:")
	entryImporte := entradas.NewEntradaNumerica("#.###,00")
	entryImporte.SetPlaceholderText("0,00")
	vboxImporte.Append(lblImporte)
	vboxImporte.Append(entryImporte)
	formGrid.Attach(vboxImporte, 1, 0, 1, 1)

	// Fila 2: Fecha | Tipo (ingreso/Egreso)
	// Columna 0: Fecha
	vboxFecha := gtk.NewBox(gtk.OrientationVertical, 5)
	lblFecha := gtk.NewLabel("Fecha (dd/mm/yyyy):")
	entryFecha := entradas.NewFecha()
	entryFecha.Min, _ = time.Parse("02/01/2006", "01/01/2026") // Fecha mínima permitida
	entryFecha.Max = time.Now() // Fecha máxima permitida es hoy
	vboxFecha.Append(lblFecha)
	vboxFecha.Append(entryFecha)
	formGrid.Attach(vboxFecha, 0, 1, 1, 1)

	// Columna 1: Swicth para Tipo
	vboxTipo := gtk.NewBox(gtk.OrientationVertical, 5)
	lblTipoHeader := gtk.NewLabel("Tipo de Movimiento")
	hboxSwitch := gtk.NewBox(gtk.OrientationHorizontal, 10)
	lblEstado := gtk.NewLabel("Egreso (E)")
	switchTipo := gtk.NewSwitch()
	switchTipo.SetHExpand(false)
	switchTipo.SetHAlign(gtk.AlignStart)
	switchTipo.SetVAlign(gtk.AlignCenter)
	switchTipo.ConnectStateSet(func(state bool) (ok bool) {
		if state {
			lblEstado.SetText("Ingreso (I)")
		} else {
			lblEstado.SetText("Egreso (E)")
		}
		return false
	})
	hboxSwitch.Append(switchTipo)
	hboxSwitch.Append(lblEstado)

	vboxTipo.Append(lblTipoHeader)
	vboxTipo.Append(hboxSwitch)
	formGrid.Attach(vboxTipo, 1, 1, 1, 1)

	// Fila 3: Botón Agregar
	btnAgregar := gtk.NewButtonWithLabel("Agregar Registro")
	btnAgregar.AddCSSClass("suggested-action")
	btnAgregar.ConnectClicked(func ()  {
		manejarAgregar(tabla, cmbCuenta, entryFecha, entryImporte, switchTipo)
	})
	formGrid.Attach(btnAgregar, 0, 2, 2, 1)

	// Ensamblado Final
	cajaPrincipal.Append(tabla.Widget)
	cajaPrincipal.Append(gtk.NewSeparator(gtk.OrientationHorizontal))
	cajaPrincipal.Append(formGrid)
		
	w.SetChild(cajaPrincipal)
}

// Función para procesar cuentas.txt: Filtrar por 'Sí' y ordenar por Fecha
func actualizarComboCuentas(dd *gtk.DropDown) {
	file, err := os.Open("./cuentas.txt")
	if err != nil {	return}
	defer file.Close()

	model := gtk.NewStringList(nil)
	mapaCuentas = make(map[uint]estructuras.CuentaInfo)
	var contador uint = 0

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {} 	// Saltar cabecera si existe
	// Leer línea por línea
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "|")
		if len(parts) < 3 { continue}
		// Filtrado por columna Imputable 'Sí'
		if strings.TrimSpace(parts[2]) == "Sí" {
			id := strings.TrimSpace(parts[0]) // "1-00-00-00"
			nombre := strings.TrimSpace(parts[1])
			// Agregamos solo el nombre al aspecto visual 
			model.Append(nombre)
			// Guardamos la estructura en nuestro mapa GO
			mapaCuentas[contador] = estructuras.CuentaInfo{
				NumeroCuenta: id,
				NombreCuenta: nombre,
			}
			contador++
		}
	}
	dd.SetModel(model)
}

func manejarAgregar(tabla *entradas.MiTabla[movimiento], cmbCuenta *gtk.DropDown, entryFecha *entradas.Fecha, entryImporte *entradas.EntradaNumerica, swTipo *gtk.Switch) {
 	idx := cmbCuenta.Selected()
  if idx == gtk.InvalidListPosition { 
    return // O manejar error si no hay nada seleccionado
  }
	info, existe := mapaCuentas[idx]
	if !existe {
		fmt.Println("Índice seleccionado no encontrado en mapaCuentas:", idx)
		return
	}

	fechaStr, _ := entryFecha.GetFecha()
	fechaDB, _ := entryFecha.GetFechaDB()
	
	letraTipo := "E"
	if swTipo.Active() {
		letraTipo = "I"
	}

	nuevoMovimiento := movimiento{
		NumeroCuenta: info.NumeroCuenta,
		NombreCuenta: info.NombreCuenta,
		Fecha:       fechaStr,
		Importe:     entryImporte.Text(),
		Tipo:        letraTipo,
	}

	// --- GUARDADO FÍSICO (APPEND) ---
	f, err := os.OpenFile("movimientos.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		// Formateamos la línea igual que la lectura
		linea := fmt.Sprintf("%s|%s|%s|%s|%s\n", 
			nuevoMovimiento.NumeroCuenta, 
			nuevoMovimiento.NombreCuenta, 
			fechaDB, 
			nuevoMovimiento.Importe, 
			nuevoMovimiento.Tipo)
		
		f.WriteString(linea)
		f.Close()
	}

	tabla.AgregarRegistro(nuevoMovimiento)

	entryImporte.SetText("")
}

func cargarMovimientosIniciales(tabla *entradas.MiTabla[movimiento]) {
	file, err := os.Open("movimientos.txt")
	if err != nil {
		return // Si no existe, simplemente no carga nada
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		partes := strings.Split(scanner.Text(), "|")
		if len(partes) < 5 { continue }

		// Parsear la fecha guardada (ajustá el formato si es necesario)
		f, _ := time.Parse("02/01/2006", partes[2])

		m := movimiento{
			NumeroCuenta: partes[0],
			NombreCuenta: partes[1],
			Fecha:        f,
			Importe:      partes[3],
			Tipo:         partes[4],
		}
		tabla.AgregarRegistro(m)
	}
}

func inicializarArchivoYTabla(tabla *entradas.MiTabla[movimiento]) {
	nombreArchivo := "movimientos.txt"
	cabecera := "Cuenta|Nombre|Fecha|Importe|Tipo\n"

	// 1. Intentar abrir o crear el archivo
	f, err := os.OpenFile(nombreArchivo, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	// 2. Si el archivo está vacío, escribir la cabecera
	info, _ := f.Stat()
	if info.Size() == 0 {
		f.WriteString(cabecera)
	}

	// 3. Leer el contenido para la tabla
	scanner := bufio.NewScanner(f)
	if scanner.Scan() {} // Saltar la línea de cabecera

	for scanner.Scan() {
		linea := scanner.Text()
		if linea == "" { continue }
		
		p := strings.Split(linea, "|")
		if len(p) < 5 { continue }

		// Convertir fecha del archivo a time.Time
		fec, err := time.Parse("2006-01-02", p[2])
		if err != nil {
			fmt.Println("Error al parsear fecha:", err)
			continue
		}

		m := movimiento{
			NumeroCuenta: p[0],
			NombreCuenta: p[1],
			Fecha:        fec,
			Importe:      p[3],
			Tipo:         p[4],
		}
		tabla.AgregarRegistro(m)
	}
}
