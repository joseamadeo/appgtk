package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/joseamadeo/appgtk/entradas"
)

// Formatos soportados por el widget.
const (
	formatoEntero   = "#.###"
	formatoDecimal  = "#.###,00"
	formatoTripleDe = "#.###,000"
)

// TestMain inicializa GTK una sola vez para toda la suite.
func TestMain(m *testing.M) {
	os.Setenv("GSK_RENDERER", "cairo")
	os.Setenv("GTK_THEME", "Adwaita:light")
	gtk.Init()
	os.Exit(m.Run())
}

// nuevoNumerico crea el widget con el formato pedido y verifica que sea usable.
func nuevoNumerico(t *testing.T, formato string) *entradas.Numerico {
	t.Helper()
	n := entradas.NewNumerico(formato)
	if n == nil || n.Object() == nil {
		t.Fatalf("NewNumerico(%q) no devolvio un entry valido", formato)
	}
	return n
}

// tipar simula que el usuario escribe con el teclado al final del texto actual.
func tipar(n *entradas.Numerico, teclas string) {
	n.Object().Buffer().InsertText(uint(len(n.Text())), teclas, len(teclas))
}

// pegar simula el pegado desde el portapapeles: GtkEntry inserta el texto del
// portapapeles en el buffer y ahi es donde el widget valida que sea un numero.
func pegar(n *entradas.Numerico, portapapeles string) {
	n.Object().Buffer().InsertText(uint(len(n.Text())), portapapeles, len(portapapeles))
}

// ponerPortapapeles deja un valor en el portapapeles del sistema.
func ponerPortapapeles(t *testing.T, n *entradas.Numerico, texto string) {
	t.Helper()
	cb := portapapelesDe(t, n)
	cb.SetText(texto)
	esperar(2*time.Second, func() bool {
		actual, err := leer(cb)
		return err == nil && actual == texto
	})
}

// leerPortapapeles lee el texto del portapapeles del sistema.
func leerPortapapeles(t *testing.T, n *entradas.Numerico) string {
	t.Helper()
	texto, err := leer(portapapelesDe(t, n))
	if err != nil {
		t.Fatalf("no se pudo leer el portapapeles: %v", err)
	}
	return texto
}

func portapapelesDe(t *testing.T, n *entradas.Numerico) *gdk.Clipboard {
	t.Helper()
	cb := n.Object().Clipboard()
	if cb == nil {
		t.Skip("el portapapeles no esta disponible en este entorno")
	}
	return cb
}

func leer(cb *gdk.Clipboard) (string, error) {
	var (
		texto   string
		err     error
		termino bool
	)
	cb.ReadTextAsync(context.Background(), func(res gio.AsyncResulter) {
		texto, err = cb.ReadTextFinish(res)
		termino = true
	})
	esperar(2*time.Second, func() bool { return termino })
	return texto, err
}

// esperar bombea el bucle principal de GLib hasta que se cumpla la condicion
// o se agote el tiempo de espera.
func esperar(timeout time.Duration, cond func() bool) bool {
	ctx := glib.MainContextDefault()
	limite := time.Now().Add(timeout)
	for time.Now().Before(limite) {
		if cond() {
			return true
		}
		if ctx.Pending() {
			ctx.Iteration(false)
			continue
		}
		time.Sleep(2 * time.Millisecond)
	}
	return cond()
}
