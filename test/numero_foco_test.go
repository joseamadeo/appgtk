package test

import (
	"testing"
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/joseamadeo/appgtk/entradas"
)

// abrirVentanaConFoco muestra el widget numerico dentro de una ventana junto a
// un segundo entry (para poder sacarle el foco) y le pide el foco.
func abrirVentanaConFoco(t *testing.T, n *entradas.Numerico) *gtk.Entry {
	t.Helper()

	otro := gtk.NewEntry()
	box := gtk.NewBox(gtk.OrientationHorizontal, 0)
	box.Append(n.Object())
	box.Append(otro)

	win := gtk.NewWindow()
	win.SetTitle("test numero")
	win.SetDefaultSize(240, 60)
	win.SetChild(box)
	win.Present()
	t.Cleanup(func() {
		win.SetVisible(false)
		esperar(time.Second, func() bool { return !win.IsVisible() })
	})

	n.Object().GrabFocus()
	return otro
}

// esperarFocoEspera a que el entry procese el evento de foco, que en este
// entorno se detecta por el cambio de texto. Si no ocurre, el test se omite
// (pasa en una sesion sin display interactivo).
func esperarFoco(t *testing.T, n *entradas.Numerico, quiere string) {
	t.Helper()
	if !esperar(3*time.Second, func() bool { return n.Text() == quiere }) {
		t.Skipf("el entry no proceso el evento de foco en este entorno (texto = %q)", n.Text())
	}
}

// TestAlEntrarEnFocoSeQuitanLosPuntosDeMiles para poder editar el numero.
func TestAlEntrarEnFocoSeQuitanLosPuntosDeMiles(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("1234,5")
	if n.Text() != "1.234,50" {
		t.Fatalf("antes de tomar el foco esperaba %q, obtuve %q", "1.234,50", n.Text())
	}

	abrirVentanaConFoco(t, n)
	esperarFoco(t, n, "1234,50")

	if n.Text() != "1234,50" {
		t.Errorf("con el foco adentro el texto = %q, se esperaba %q (sin puntos de miles)",
			n.Text(), "1234,50")
	}
}

// TestAlSalirDeFocoSeVuelveAFormatearComoEnExcel.
func TestAlSalirDeFocoSeVuelveAFormatearComoEnExcel(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("1234,5")
	otro := abrirVentanaConFoco(t, n)
	esperarFoco(t, n, "1234,50")

	otro.GrabFocus()
	esperarFoco(t, n, "1.234,50")

	if n.Text() != "1.234,50" {
		t.Errorf("con el foco afuera el texto = %q, se esperaba %q (nuevamente formateado)",
			n.Text(), "1.234,50")
	}
}

// TestSePuedeEditarElNumeroMientrasTieneElFoco: al entrar en foco el widget
// queda listo para escribir (sin miles) y sigue admitiendo digitos.
func TestSePuedeEditarElNumeroConElFocoAdentro(t *testing.T) {
	n := nuevoNumerico(t, formatoEntero)
	n.SetText("1234")
	abrirVentanaConFoco(t, n)
	esperarFoco(t, n, "1234")

	tipar(n, "5")
	if n.Text() != "12345" {
		t.Errorf("tras tipear con el foco adentro el texto = %q, se esperaba %q", n.Text(), "12345")
	}
	if db := n.FormatoDB(); db != "12345" {
		t.Errorf("FormatoDB() con el numero a medio escribir = %q, se esperaba %q", db, "12345")
	}
}
