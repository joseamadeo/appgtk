package pantalla

/*
#cgo pkg-config: poppler-glib
#include <poppler.h>
#include <stdlib.h>
*/
import "C"

import (
	"context"
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

var (
	pathPDF string
	actualPagina int = 0
	zoomNivel float64 = 1.0
	totalPaginas int = 0
)

func PDFVista(w *gtk.ApplicationWindow) {
	vbox := gtk.NewBox(gtk.OrientationVertical, 0)

	// Barra de Herramientas
	barraHerramienta := gtk.NewBox(gtk.OrientationHorizontal,5)
	barraHerramienta.SetMarginTop(5)
	barraHerramienta.SetMarginBottom(5)
	barraHerramienta.SetMarginStart(5)
	barraHerramienta.SetMarginEnd(5)
	vbox.Append(barraHerramienta)

	// Etiqueta de la dirección del PDF
	lblPath := gtk.NewLabel("")
	lblPath.SetHExpand(true)
	lblPath.SetXAlign(0)
	lblPath.SetEllipsize(pango.EllipsizeEnd)
	lblPath.SetMarginStart(10)
	vbox.Append(lblPath)

	// 2. Área de dibujo para mostrar el PDF
	area := gtk.NewDrawingArea()
	area.SetVExpand(true) // que ocupe el espacio disponible

	// 3. Botón con referencia explicita
	btnAnterior := gtk.NewButtonFromIconName("go-previous-symbolic")
	btnSiguiente := gtk.NewButtonFromIconName("go-next-symbolic")

	btnAnterior.SetSensitive(true)
	btnAnterior.ConnectClicked(func() {
		if actualPagina > 0 {
			actualPagina--
			area.QueueDraw() // <--- Esto es vital
			actualizarBotones(btnAnterior, btnSiguiente)
		}
	})


	btnSiguiente.ConnectClicked(func() {
		if actualPagina < totalPaginas-1 { // Validar contra el total de páginas
			actualPagina++ // Aquí podrías validar contra el total de páginas
			area.QueueDraw() // <--- Esto es vital
			actualizarBotones(btnAnterior, btnSiguiente)
		}
	})

	barraHerramienta.Append(btnAnterior)
	barraHerramienta.Append(btnSiguiente)

	btnAbrir := gtk.NewButtonWithLabel("Abrir PDF")
	btnAbrir.ConnectClicked(func() {
		dialogo := gtk.NewFileDialog()
		dialogo.SetTitle("Seleccionar PDF")
		dialogo.Open(context.Background(), &w.Window, func(res gio.AsyncResulter) {
			archivo, err := dialogo.OpenFinish(res)
			if err == nil {
				pathPDF = archivo.Path()
				actualPagina = 0 // Reiniciar a la primera página al cargar un nuevo PDF
				// Obtenemos el total de páginas real
		    totalPaginas = getPDFtotalPaginas(pathPDF) 
				lblPath.SetText(pathPDF) // Mostrar la ruta del PDF en la etiqueta
				// importante: avisar al área de dibujo que debe redibujar
				area.QueueDraw()
				actualizarBotones(btnAnterior, btnSiguiente)
			}
		})
	})
	barraHerramienta.Append(btnAbrir)

	// 3. Controles de Zoom
	btnZoomMenor := gtk.NewButtonFromIconName("zoom-in-symbolic")
	btnZoomMas := gtk.NewButtonFromIconName("zoom-out-symbolic")
	
	btnZoomMenor.ConnectClicked(func() { zoomNivel *= 1.2; area.QueueDraw() })
	btnZoomMas.ConnectClicked(func() { zoomNivel /= 1.2; area.QueueDraw() })
	
	barraHerramienta.Append(btnZoomMenor)
	barraHerramienta.Append(btnZoomMas)

	// 4. Botón Imprimir
	btnImprimir := gtk.NewButtonFromIconName("document-print-symbolic")
	btnImprimir.ConnectClicked(func() {
		if pathPDF != "" {
			imprimirPDF(w)
		}
	})
	barraHerramienta.Append(btnImprimir)


	// --- ÁREA DE DIBUJO ---
	area.SetDrawFunc(func(_ *gtk.DrawingArea, cr *cairo.Context, width, height int) {
		if pathPDF == "" {
			return
		}
		// Limpiar fondo (opcional pero recomendado)
		cr.SetSourceRGB(0.9, 0.9, 0.9) // Fondo gris claro
		cr.Paint()
		// Aplicar zoom
		cr.Save()
		cr.Scale(zoomNivel, zoomNivel)
		// Renderizar directamente usando la variable global actualizada
		success := renderPDFPagina(pathPDF, actualPagina, cr)
		cr.Restore() // Restaurar el estado original después de renderizar
		if !success {
			// Mostrar un mensaje de error en el área de dibujo
			cr.SetSourceRGB(1, 0, 0) // Rojo para el error
			cr.SelectFontFace("Sans", cairo.FontSlantNormal, cairo.FontWeightBold)
			cr.SetFontSize(14)
			cr.MoveTo(10, 30)
			cr.ShowText("Error al cargar el PDF o renderizar la página.")
			return
		}
	})
	vbox.Append(area)

	w.SetChild(vbox)
}

func actualizarBotones(btnAnterior, btnSiguiente *gtk.Button) {
    btnAnterior.SetSensitive(actualPagina > 0)
    btnSiguiente.SetSensitive(actualPagina < totalPaginas-1)
}

func getPDFtotalPaginas(s string) int {
	cPath := C.CString("file:///" + s)
	defer C.free(unsafe.Pointer(cPath))

	var err *C.GError
	doc := C.poppler_document_new_from_file(cPath, nil, &err)
	if doc == nil {
		return 0
	}
	defer C.g_object_unref(C.gpointer(doc))

	return int(C.poppler_document_get_n_pages(doc))
}

func renderPDFPagina(s string, paginaNum int, cr *cairo.Context) bool{
	cPath := C.CString("file:///" + s) // Poppler espera una URI, por eso el prefijo "file:///"
	defer C.free(unsafe.Pointer(cPath))

	var err *C.GError
	doc := C.poppler_document_new_from_file(cPath, nil, &err)
	if doc == nil {
		return false
	}
	defer C.g_object_unref(C.gpointer(doc))

	// 2. Obtener la página específica
	pagina := C.poppler_document_get_page(doc, C.int(paginaNum))
	if pagina == nil {
		return false
	}
	defer C.g_object_unref(C.gpointer(pagina))
	// Converimos el contexto de gotk4 al puntero de C que espera poppler
	nativeCtx := (*C.cairo_t)(unsafe.Pointer(cr.Native()))
	C.poppler_page_render(pagina, nativeCtx)
	return true
}

func imprimirPDF(parent *gtk.ApplicationWindow) {
	printOp := gtk.NewPrintOperation()
	printOp.SetJobName("Imprimir PDF")

	_, err := printOp.Run(gtk.PrintOperationActionPrintDialog, &parent.Window)
	if err != nil {
		// Creamos el diálogo con los 4 argumentos exactos que pide tu compilador
		dialog := gtk.NewMessageDialog(
			&parent.Window,
			gtk.DialogModal,
			gtk.MessageError, 
			gtk.ButtonsOK,
		)
		
		dialog.SetTitle("ERROR DE IMPRESIÓN")
		
		dialog.SetMarkup("Error al imprimir el documento: " + err.Error())

		// Manejo de respuesta para cerrar la ventana
		dialog.Connect("response", func() {
			dialog.Destroy()
		})

		dialog.Show()
	}
}
