package entradas

import (
	"fmt"
	"strconv"
	"appgtk/estructuras"
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/core/glib"      // Para los objetos
	glibv2 "github.com/diamondburned/gotk4/pkg/glib/v2" // Para el tipo de función
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type MiTabla[T any]struct{
	estructuras.TablaGenerica[T]
	UltimoIdxBusqueda int
} 

func NewTabla[T any](lista []T, columnas []estructuras.ColumnaConfig[T]) *MiTabla[T] {
	t := &MiTabla[T]{
		TablaGenerica: estructuras.TablaGenerica[T]{
			Datos: lista,
		},
		UltimoIdxBusqueda: -1,
	}
  store := gtk.NewStringList(nil)
  for i := 0; i < len(lista); i++ {
    store.Append(strconv.Itoa(i))
  }
  modelo := gtk.NewSortListModel(store, nil)
  seleccion := gtk.NewSingleSelection(modelo)
  tabla := gtk.NewColumnView(seleccion)

	for _, conf := range columnas {
		c := conf // Captura de variable para el closure
		factory := gtk.NewSignalListItemFactory()
		factory.ConnectSetup(func(obj *glib.Object) {
			celda := obj.Cast().(*gtk.ColumnViewCell)
			celda.SetChild(gtk.NewLabel(""))
		})
		factory.ConnectBind(func(obj *glib.Object) {
			celda := obj.Cast().(*gtk.ColumnViewCell)
			item := celda.Item()
			if item == nil { return }

			idxStr := item.Cast().(*gtk.StringObject).String()
			idx, _ := strconv.Atoi(idxStr)

			if idx < 0 || idx >= len(t.Datos) {
				if label, ok := celda.Child().(*gtk.Label); ok {
					label.SetText("Error: Indice")
				}
				return
			}
			texto := c.Extraer(t.Datos[idx])
			if label, ok := celda.Child().(*gtk.Label); ok {
				label.SetText(texto)
			}
		})

		col := gtk.NewColumnViewColumn(c.Titulo, &factory.ListItemFactory)
		if c.Comparar != nil {
			sorter := gtk.NewCustomSorter(glibv2.CompareDataFunc(func(aPtr, bPtr unsafe.Pointer) int {
        idxA, _ := strconv.Atoi(glib.Take(aPtr).Cast().(*gtk.StringObject).String())
        idxB, _ := strconv.Atoi(glib.Take(bPtr).Cast().(*gtk.StringObject).String())
	      // Accedemos a t.Datos para que el Sorter no de Panic al ordenar nuevos registros
        return c.Comparar(t.Datos[idxA], t.Datos[idxB])
    	}))
			col.SetSorter(&sorter.Sorter)
		}
		tabla.AppendColumn(col)
	}
  // Configuración de scroll y modelo
	if s := tabla.Sorter(); s != nil {
		modelo.SetSorter(s)
	}
	scr := gtk.NewScrolledWindow()
	scr.SetChild(tabla)
	scr.SetVExpand(true)
	scr.SetHExpand(true)

	t.Widget = scr
	t.Store = store
	t.Datos = lista
	t.Seleccion = seleccion

	return t
}

// AlSeleccionar permite registrar una función que se ejecutará cada vez que
// el usuario cambie la fila seleccionada.
func (t *MiTabla[T]) AlSeleccionar(callback func(registro T)) {
	t.Seleccion.ConnectSelectionChanged(func(pos uint, nItems uint) {
		// Obtenemos el registro usando el método que ya tenés
		registro, ok := t.GetSeleccionado()
		if ok {
			callback(registro)
		}
	})
}

func (t *MiTabla[T]) GetSeleccionado() (T, bool) {
	var vacio T
	pos := t.Seleccion.Selected()
	if pos == gtk.InvalidListPosition {
		return vacio, false
	}
	obj := t.Seleccion.Item(pos)
	if obj == nil { return vacio, false }
	idxStr := obj.Cast().(*gtk.StringObject).String()
	idx, _ := strconv.Atoi(idxStr)
	// VALIDACIÓN DE SEGURIDAD
	if idx < 0 || idx >= len(t.Datos) {
		return vacio, false
	}
	return t.Datos[idx], true
}

func (t *MiTabla[T]) AgregarRegistro(nuevo T) {
	t.Datos = append(t.Datos, nuevo)
	t.Store.Append(strconv.Itoa(len(t.Datos) - 1))
}

func (t *MiTabla[T]) Modificar(nuevo T) error {
	pos := t.Seleccion.Selected()
	if pos == gtk.InvalidListPosition {
		return fmt.Errorf("no hay nada seleccionado")
	}

	obj := t.Seleccion.Item(pos)
	if obj == nil {
		return fmt.Errorf("error al recuperar objeto")
	}
	
	idxStr := obj.Cast().(*gtk.StringObject).String()
	idx, _ := strconv.Atoi(idxStr)

	if idx < 0 || idx >= len(t.Datos) {
		return fmt.Errorf("índice fuera de rango")
	}
	t.Datos[idx] = nuevo
	t.Store.Splice(uint(idx), 1, []string{strconv.Itoa(idx)})
	
	return nil
}

func (t *MiTabla[T]) EliminarSeleccionado() error {
	pos := t.Seleccion.Selected()
	if pos == gtk.InvalidListPosition {
		return fmt.Errorf("no hay ningún registro seleccionado para eliminar")
	}

	obj := t.Seleccion.Item(pos)
	idxStr := obj.Cast().(*gtk.StringObject).String()
	idx, _ := strconv.Atoi(idxStr)

	if idx < 0 || idx >= len(t.Datos) {
		return fmt.Errorf("índice fuera de rango")
	}
	t.Datos = append(t.Datos[:idx], t.Datos[idx+1:]...)
	t.Store.Splice(0, t.Store.NItems(), nil) // Limpia todo el Store
	for i := 0; i < len(t.Datos); i++ {
		t.Store.Append(strconv.Itoa(i)) // Re-mapea índices 0, 1, 2...
	}

	return nil
}

func (t *MiTabla[T]) Buscar(criterio func(T) bool) (T, bool) {
	var vacio T
	total := len(t.Datos)
	if total == 0 {
		return vacio, false
	}

	// Calculamos desde dónde empezar (el siguiente al último o 0)
	inicio := t.UltimoIdxBusqueda + 1
	if inicio >= total {
		inicio = 0
	}

	// Buscamos desde el puntero hasta el final del array
	for i := inicio; i < total; i++ {
		if criterio(t.Datos[i]) {
			t.UltimoIdxBusqueda = i
			t.sincronizarSeleccionVisual(i)
			return t.Datos[i], true
		}
	}

	// Si llegamos al final y no hubo coincidencia, reseteamos el puntero
	t.UltimoIdxBusqueda = -1
	return vacio, false
}

// Función auxiliar para marcar la fila en la pantalla
func (t *MiTabla[T]) sincronizarSeleccionVisual(indiceReal int) {
	modelo := t.Seleccion.Model()
	n := modelo.NItems()
	idxStr := strconv.Itoa(indiceReal)

	for i := uint(0); i < n; i++ {
		obj := modelo.Item(i)
		if obj == nil { continue }
		if obj.Cast().(*gtk.StringObject).String() == idxStr {
			t.Seleccion.SetSelected(i)
			break
		}
	}
}
