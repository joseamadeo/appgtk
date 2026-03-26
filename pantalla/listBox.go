package pantalla

import (
	"strconv"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func ListBox(w *gtk.ApplicationWindow) {
	vboxPrincipal := gtk.NewBox(gtk.OrientationHorizontal, 5)
	
	/**************************************/
	hIzquierda := gtk.NewBox(gtk.OrientationVertical,2)
	hIzquierda.SetHExpand(true)

	lblIzquierda := gtk.NewLabel("ListBox")
	lblIzquierda.SetMarginTop(2)
	lblIzquierda.SetMarginBottom(2)
	lblIzquierda.SetMarginStart(2)
	lblIzquierda.SetMarginEnd(2)
	lblIzquierda.SetText("List Box")
	hIzquierda.Append(lblIzquierda)

	scrollList := gtk.NewScrolledWindow()
	scrollList.SetPolicy(gtk.PolicyAutomatic, gtk.PolicyAlways)
	scrollList.SetVExpand(true)
	vboxList := gtk.NewBox(gtk.OrientationVertical, 0)
	vboxList.Append(scrollList)

	listBox := gtk.NewListBox()
	for i := 1; i <= 20; i++ {
		row := gtk.NewListBoxRow()
		rowLabel := gtk.NewLabel("Elemento " + strconv.Itoa(i))
		row.SetChild(rowLabel)
		
		listBox.Append(row)
	}
	scrollList.SetChild(listBox)

	hIzquierda.Append(vboxList)



	/*********************************************/
	hDerecha := gtk.NewBox(gtk.OrientationVertical,2)
	hDerecha.SetHExpand(true)


	lblDerecha := gtk.NewLabel("ComboBox")
	lblDerecha.SetMarginTop(2)
	lblDerecha.SetMarginBottom(2)
	lblDerecha.SetMarginStart(2)
	lblDerecha.SetMarginEnd(2)
	lblDerecha.SetText("Combo Box")
	hDerecha.Append(lblDerecha)

	datos := gtk.NewStringList([]string{
    "Opción 1",
    "Opción 2",
    "Opción 3",
    "Opción 4",
    "Opción 5",
	})

	comboBox := gtk.NewDropDown(datos, nil)
	comboBox.SetMarginTop(20)
	comboBox.SetMarginBottom(20)
	comboBox.SetMarginStart(20)
	comboBox.SetMarginEnd(20)

	hDerecha.Append(comboBox)



	vboxPrincipal.Append(hIzquierda)
	vboxPrincipal.Append(hDerecha)
	w.SetChild(vboxPrincipal)

}