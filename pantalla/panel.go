package pantalla

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func Panel(w *gtk.ApplicationWindow){
	vbox:= gtk.NewBox(gtk.OrientationVertical, 0)
	w.SetChild(vbox)

	paned1 := gtk.NewPaned(gtk.OrientationVertical)
	paned1.SetMarginTop(4)
	paned1.SetMarginBottom(4)
	paned1.SetMarginStart(4)
	paned1.SetMarginEnd(4)
	paned1.SetVExpand(true)
	vbox.Append(paned1)
	lbl1 := gtk.NewLabel("Arriba Panel 1")
	paned1.SetStartChild(lbl1)
	lbl2 := gtk.NewLabel("Abajo Panel 1")
	paned1.SetEndChild(lbl2)
	
	paned2 := gtk.NewPaned(gtk.OrientationHorizontal)
	paned2.SetMarginTop(4)
	paned2.SetMarginBottom(4)
	paned2.SetMarginStart(4)
	paned2.SetMarginEnd(4)
	vbox.Append(paned2)
	lbl3 := gtk.NewLabel("Izquierda ")
	paned2.SetStartChild(lbl3)
	lbl4 := gtk.NewLabel("Derecha")
	paned2.SetEndChild(lbl4)
}