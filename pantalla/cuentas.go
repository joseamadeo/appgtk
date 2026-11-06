package pantalla

import (
	"github.com/joseamadeo/appgtk/entradas"
	"github.com/joseamadeo/appgtk/estructuras"
	"encoding/csv"
	"os"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// 1. Definimos el objeto que representa una fila
type datosCuenta struct {
	Cuenta    string
	Imputable bool
	Nombre    string
}

type MiTabla[T any]struct{
	estructuras.TablaGenerica[T]
	UltimoIdxBusqueda int
}

var listaCuentas = []datosCuenta{}

func Cuentas(w *gtk.ApplicationWindow) {
	cajaPrincipal := gtk.NewBox(gtk.OrientationVertical, 10)
	cajaPrincipal.SetMarginStart(10)
	cajaPrincipal.SetMarginEnd(10)
	cajaPrincipal.SetMarginTop(10)
	cajaPrincipal.SetMarginBottom(10)

	lbl1 := gtk.NewLabel("Cuentas Contables")
	cajaPrincipal.Append(lbl1)

	listaCuentas, err := cargarCuentas("cuentas.txt")
	if err != nil {
		MostrarError(&w.Window, "No se pudo cargar cuentas.txt: "+err.Error())
		listaCuentas = []datosCuenta{} // Iniciar con lista vacía si hay error
	}

	colsCuentas := []estructuras.ColumnaConfig[datosCuenta]{
		{
			Titulo: "Cuenta",
			Extraer: func(dc datosCuenta) string { return dc.Cuenta },
			Comparar: func(a, b datosCuenta) int {return estructuras.CompararStrings(a.Cuenta, b.Cuenta)},
		},
		{
			Titulo: "Nombre",
			Extraer: func(dc datosCuenta) string { return dc.Nombre },
		},
		{
			Titulo: "Imputable",
			Extraer: func(dc datosCuenta) string {
				if dc.Imputable {
					return "Sí"
				}
				return "No"
			},
		},
	}
	tabla := entradas.NewTabla(listaCuentas, colsCuentas)
	tabla.Widget.SetVExpand(true)
	tabla.Widget.SetHExpand(true)

	cajaForm := gtk.NewBox(gtk.OrientationVertical, 5)
	cajaForm.SetMarginTop(10)
	cajaForm.SetMarginBottom(10)

	// Fila 1: Cuenta e imputable
	fila1 := gtk.NewBox(gtk.OrientationHorizontal, 5)
	lblCuenta := gtk.NewLabel("Cuenta:")
	entCuenta := entradas.NewCuentaContable()
	fila1.Append(lblCuenta)
	fila1.Append(entCuenta)
	lblImputable := gtk.NewLabel("Imputable:")
	chkImputable := gtk.NewCheckButtonWithLabel("¿Es Imputable?")
	fila1.Append(lblImputable)
	fila1.Append(chkImputable)
	cajaForm.Append(fila1)

	// Fila 2: Nombre
	entNombre := gtk.NewEntry()
	entNombre.SetPlaceholderText("Nombre de la cuenta")
	entNombre.SetHExpand(true)
	cajaForm.Append(entNombre)

	// Fila 3: Botones de acción
	cajaBotones := gtk.NewBox(gtk.OrientationHorizontal, 5)
	btnAgregarModificar := gtk.NewButtonWithLabel("Agregar Cuenta")
	btnAgregarModificar.AddCSSClass("suggested-action")

	btnEliminar := gtk.NewButtonWithLabel("Eliminar")
	btnEliminar.AddCSSClass("destructive-action")

	btnLimpiar := gtk.NewButtonWithLabel("Limpiar")

	btnGuardarArchivo := gtk.NewButtonWithLabel("Guardar en Archivo")
	btnGuardarArchivo.AddCSSClass("suggested-action")
	// Armamos la caja de botones
	cajaBotones.Append(btnAgregarModificar)
	cajaBotones.Append(btnEliminar)
	cajaBotones.Append(btnLimpiar)
	cajaBotones.Append(btnGuardarArchivo)
	cajaForm.Append(cajaBotones)


	// ... después de crear la tabla ...
	tabla.AlSeleccionar(func(registro datosCuenta) {
    // Cargamos los datos en los widgets del formulario
    entCuenta.SetCuenta(registro.Cuenta)
    entNombre.SetText(registro.Nombre)
    chkImputable.SetActive(registro.Imputable)
    
    // Cambiamos el texto del botón para indicar edición
    btnAgregarModificar.SetLabel("Actualizar")
	})


	btnAgregarModificar.ConnectClicked(func ()  {
		cuenta := entCuenta.GetCuenta()
		nombre := entNombre.Text()
		imputable := chkImputable.Active()
		if cuenta == "" || nombre == "" {
			MostrarError(&w.Window, "Cuenta y Nombre son obligatorios.")
			return
		}
		nuevaCuenta := datosCuenta{Cuenta: cuenta, Nombre: nombre, Imputable: imputable}
		// Verificar si la cuenta ya existe
		tabla.UltimoIdxBusqueda = -1 // Reiniciar índice de búsqueda
		existe, _ := tabla.Buscar(func(dc datosCuenta) bool {
			return dc.Cuenta == cuenta
		})
		if existe.Cuenta != "" {
			// Actualizar la cuenta existente
			existe.Nombre = nombre
			existe.Imputable = imputable
			tabla.Modificar(existe)
		} else {
			// Agregar nueva cuenta
			tabla.AgregarRegistro(nuevaCuenta)
		}
		entCuenta.SetCuenta("")
		entNombre.SetText("")
		chkImputable.SetActive(false)
		btnAgregarModificar.SetLabel("Agregar Cuenta")
	})

	btnEliminar.ConnectClicked(func ()  {
		err := tabla.EliminarSeleccionado()
		if err != nil {
			MostrarError(&w.Window, "Error al eliminar: "+err.Error())
		}
	})

	btnGuardarArchivo.ConnectClicked(func ()  {
		cuentas := tabla.Datos
		archivo, err := os.Create("cuentas.txt")
		if err != nil {
			MostrarError(&w.Window, "Error al crear cuentas.txt: "+err.Error())
			return
		}
		defer archivo.Close()
		escritura := csv.NewWriter(archivo)
		escritura.Comma = '|'
		// Escribir cabecera
		err = escritura.Write([]string{"Cuenta", "Nombre", "Imputable"})
		if err != nil {
			MostrarError(&w.Window, "Error al escribir la cabecera: "+err.Error())
			return
		}
		// Escribir datos
		for _, cuenta := range cuentas {
			imputableStr := "No"
			if cuenta.Imputable {
				imputableStr = "Sí"
			}
			err = escritura.Write([]string{cuenta.Cuenta, cuenta.Nombre, imputableStr})
			if err != nil {
				MostrarError(&w.Window, "Error al escribir la cuenta '"+cuenta.Cuenta+"': "+err.Error())
				return
			}
		}
		escritura.Flush()
		if err := escritura.Error(); err != nil {
			MostrarError(&w.Window, "Error al guardar cuentas.txt: "+err.Error())
		}
	})
	
	btnLimpiar.ConnectClicked(func ()  {
		entCuenta.SetCuenta("")
		entNombre.SetText("")
		chkImputable.SetActive(false)
		btnAgregarModificar.SetLabel("Agregar Cuenta")
	})

	cajaPrincipal.Append(tabla.Widget)
	cajaPrincipal.Append(cajaForm)
	w.SetChild(cajaPrincipal)
}

func cargarCuentas(path string) ([]datosCuenta, error) {
	archivo, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer archivo.Close()

	lectura := csv.NewReader(archivo)
	lectura.Comma = '|'

	// 1. LEER Y DESCARTAR LA CABECERA
	_, err = lectura.Read() 
	if err != nil {
		return nil, err // Retorna error si el archivo está vacío
	}

	// 2. LEER EL RESTO DE LOS DATOS
	registros, err := lectura.ReadAll()
	if err != nil {
		return nil, err
	}

	var cuentas []datosCuenta
	for _, registro := range registros {
		imputable := false
		if len(registro) >= 3 {
			if registro[2] == "Sí" {
				imputable = true
			}
			//
			cuentas = append(cuentas, datosCuenta{
				Cuenta:    registro[0],
				Nombre:    registro[1],
				Imputable: imputable,
			})
		}
	}
	return cuentas, nil
}
