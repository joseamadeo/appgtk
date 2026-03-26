package estructuras

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type ColumnaConfig[T any] struct {
	Titulo  string
	Extraer func(T) string
	Comparar func(a, b T) int 
}

type TablaGenerica[T any] struct {
	Widget *gtk.ScrolledWindow
	Store  *gtk.StringList
	Datos  []T
	Seleccion *gtk.SingleSelection
}
