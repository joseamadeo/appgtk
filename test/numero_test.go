package test

import (
	"testing"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// TestNewNumericoDevuelveUnEntry alineado a la derecha y expandible.
func TestNewNumericoDevuelveUnEntry(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)

	if n.Text() != "" {
		t.Errorf("un widget nuevo deberia estar vacio, obtuve %q", n.Text())
	}
	if alineacion := n.Object().Alignment(); alineacion != 1.0 {
		t.Errorf("el entry deberia estar alineado a la derecha (1.0), obtuve %v", alineacion)
	}
	if !n.Object().HExpand() {
		t.Error("el entry deberia ocupar todo el ancho disponible (hexpand)")
	}
	if n.Object().PlaceholderText() != "" {
		t.Errorf("un widget nuevo no deberia tener placeholder, obtuve %q", n.Object().PlaceholderText())
	}
}

// TestObjectEsElEntryNativoParaPoderAgregarloAUnContenedor.
func TestObjectEsElEntryNativo(t *testing.T) {
	n := nuevoNumerico(t, formatoEntero)

	var _ *gtk.Entry = n.Object()

	box := gtk.NewBox(gtk.OrientationHorizontal, 0)
	box.Append(n.Object())
	if box.FirstChild() == nil {
		t.Error("el entry devuelto por Object() no se pudo agregar a un contenedor")
	}
}

// TestSetPlaceholderText.
func TestSetPlaceholderText(t *testing.T) {
	n := nuevoNumerico(t, formatoEntero)
	n.SetPlaceholderText("Importe")

	if n.Object().PlaceholderText() != "Importe" {
		t.Errorf("placeholder = %q, se esperaba %q", n.Object().PlaceholderText(), "Importe")
	}
}

// TestSetTextFormateaComoEnExcel separando miles con punto y decimales con coma.
func TestSetTextFormateaComoEnExcel(t *testing.T) {
	casos := []struct {
		formato string
		entrada string
		quiere  string
	}{
		// formato con 2 decimales
		{formatoDecimal, "0", "0,00"},
		{formatoDecimal, "1", "1,00"},
		{formatoDecimal, "999", "999,00"},
		{formatoDecimal, "1000", "1.000,00"},
		{formatoDecimal, "1234", "1.234,00"},
		{formatoDecimal, "1234567", "1.234.567,00"},
		{formatoDecimal, "1234567890123", "1.234.567.890.123,00"},
		// separador decimal coma
		{formatoDecimal, "1234,5", "1.234,50"},
		{formatoDecimal, "1,5", "1,50"},
		{formatoDecimal, "1234,567", "1.234,57"},  // redondea
		{formatoDecimal, "1234,5678", "1.234,57"}, // redondea
		{formatoDecimal, "0,005", "0,01"},         // redondea
		{formatoDecimal, "0,004", "0,00"},         // redondea
		// el punto de la entrada se interpreta como separador de miles
		{formatoDecimal, "1.234,56", "1.234,56"},
		{formatoDecimal, "12.34.56", "123.456,00"},
		// negativos (solo mediante SetText o pegado)
		{formatoDecimal, "-1234", "-1.234,00"},
		{formatoDecimal, "-1234,25", "-1.234,25"},
		// vacio
		{formatoDecimal, "", ""},
		// formato entero
		{formatoEntero, "0", "0"},
		{formatoEntero, "123", "123"},
		{formatoEntero, "1234", "1.234"},
		{formatoEntero, "1234567", "1.234.567"},
		{formatoEntero, "1234567890", "1.234.567.890"},
		{formatoEntero, "1234,56", "1.235"}, // descarta los decimales
		{formatoEntero, "", ""},
		// formato con 3 decimales
		{formatoTripleDe, "1234,5", "1.234,500"},
		{formatoTripleDe, "1234,5678", "1.234,568"},
		{formatoTripleDe, "2000000", "2.000.000,000"},
	}

	for _, c := range casos {
		t.Run(c.formato+"/"+c.entrada, func(t *testing.T) {
			n := nuevoNumerico(t, c.formato)
			n.SetText(c.entrada)
			if n.Text() != c.quiere {
				t.Errorf("SetText(%q) con formato %q = %q, se esperaba %q",
					c.entrada, c.formato, n.Text(), c.quiere)
			}
		})
	}
}

// TestSetTextDescartaLoQueNoEsUnNumero.
func TestSetTextDescartaLoQueNoEsUnNumero(t *testing.T) {
	casos := []string{
		"abc",
		"12abc",
		"1.2.3.4x",
		" ",
		".",
		",",
		"1,2,3",
		"$ 100",
		"100%",
	}

	for _, entrada := range casos {
		t.Run(entrada, func(t *testing.T) {
			n := nuevoNumerico(t, formatoDecimal)
			n.SetText(entrada)
			if n.Text() != "" {
				t.Errorf("SetText(%q) = %q, se esperaba texto vacio", entrada, n.Text())
			}
		})
	}
}

// TestSetTextEsIdempotente: pasar el texto ya formateado vuelve a dar el mismo
// resultado, porque el punto de miles se quita antes de leer el numero.
func TestSetTextEsIdempotente(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("1234,5")
	primero := n.Text()

	n.SetText(primero)
	if n.Text() != primero {
		t.Errorf("SetText(%q) = %q, se esperaba %q (idempotente)", primero, n.Text(), primero)
	}
}

// TestSetTextInterpretaElPuntoComoSeparadorDeMiles documenta un detalle del
// formato: hay que asignar el valor crudo y nunca el texto ya formateado,
// porque "12.34.56" se interpreta como 123456.
func TestSetTextInterpretaElPuntoComoSeparadorDeMiles(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("12.34.56")

	if n.Text() != "123.456,00" {
		t.Errorf("SetText(\"12.34.56\") = %q, se esperaba %q (el punto se lee como miles)",
			n.Text(), "123.456,00")
	}
}

// TestSetTextNegativoEnFormatoEntero documenta que el separador de miles no
// contempla el signo menos: "-500" se muestra como "-.500".
func TestSetTextNegativoEnFormatoEntero(t *testing.T) {
	n := nuevoNumerico(t, formatoEntero)
	n.SetText("-500")

	if n.Text() != "-.500" {
		t.Errorf("SetText(\"-500\") con formato entero = %q, se esperaba %q",
			n.Text(), "-.500")
	}
	if n.FormatoDB() != "-500" {
		t.Errorf("FormatoDB() = %q, se esperaba %q", n.FormatoDB(), "-500")
	}
}

// TestSetTextNoAceptaEspaciosAlrededor documenta que SetText no recorta el texto.
func TestSetTextNoAceptaEspaciosAlrededor(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText(" 1234 ")

	if n.Text() != "" {
		t.Errorf("SetText(\" 1234 \") = %q, se esperaba texto vacio", n.Text())
	}
}

// TestFormatoDBDevuelveElNumeroSinMilesYConPuntoDecimal.
func TestFormatoDBDevuelveElNumeroSinMiles(t *testing.T) {
	casos := []struct {
		formato string
		entrada string
		quiere  string
	}{
		{formatoDecimal, "1234,56", "1234.56"},
		{formatoDecimal, "1.234.567,89", "1234567.89"},
		{formatoDecimal, "0", "0.00"},
		{formatoDecimal, "-1234,25", "-1234.25"},
		{formatoEntero, "1234", "1234"},
		{formatoEntero, "1.234.567", "1234567"},
		{formatoDecimal, "", ""},
	}

	for _, c := range casos {
		t.Run(c.formato+"/"+c.entrada, func(t *testing.T) {
			n := nuevoNumerico(t, c.formato)
			n.SetText(c.entrada)
			if n.FormatoDB() != c.quiere {
				t.Errorf("FormatoDB() con texto %q = %q, se esperaba %q",
					n.Text(), n.FormatoDB(), c.quiere)
			}
		})
	}
}

// TestSetFormatoCambiaLosDecimalesYReformateaElTextoActual.
func TestSetFormatoCambiaLosDecimalesYReformateaElTextoActual(t *testing.T) {
	n := nuevoNumerico(t, formatoEntero)
	n.SetText("1234567")
	if n.Text() != "1.234.567" {
		t.Fatalf("con formato entero esperaba %q, obtuve %q", "1.234.567", n.Text())
	}

	n.SetFormato(formatoDecimal)
	if n.Text() != "1.234.567,00" {
		t.Errorf("tras SetFormato(%q) esperaba %q, obtuve %q", formatoDecimal, "1.234.567,00", n.Text())
	}
	if db := n.FormatoDB(); db != "1234567.00" {
		t.Errorf("FormatoDB() = %q, se esperaba %q", db, "1234567.00")
	}

	n.SetFormato(formatoEntero)
	if n.Text() != "1.234.567" {
		t.Errorf("tras volver a formato entero esperaba %q, obtuve %q", "1.234.567", n.Text())
	}
}

// TestSetFormatoRedondeaAlPasarAFormatoEntero.
func TestSetFormatoRedondeaAlPasarAFormatoEntero(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("1234567,6")
	if n.Text() != "1.234.567,60" {
		t.Fatalf("con formato decimal esperaba %q, obtuve %q", "1.234.567,60", n.Text())
	}

	n.SetFormato(formatoEntero)

	if n.Text() != "1.234.568" {
		t.Errorf("tras SetFormato(%q) esperaba %q (redondeo), obtuve %q",
			formatoEntero, "1.234.568", n.Text())
	}
}

// TestSetFormatoConVariosDecimales.
func TestSetFormatoConVariosDecimales(t *testing.T) {
	n := nuevoNumerico(t, formatoEntero)
	n.SetFormato(formatoTripleDe)
	n.SetText("1234,5678")
	if n.Text() != "1.234,568" {
		t.Errorf("con 3 decimales esperaba %q, obtuve %q", "1.234,568", n.Text())
	}

	n.SetFormato("#.###,0000")
	n.SetText("1234,5678")
	if n.Text() != "1.234,5678" {
		t.Errorf("con 4 decimales esperaba %q, obtuve %q", "1.234,5678", n.Text())
	}
}

// TestSetFormatoCuentaLosCaracteresTrasLaComa: la cantidad de decimales es la
// cantidad de caracteres que siguen a la coma del formato.
func TestSetFormatoCuentaLosCaracteresTrasLaComa(t *testing.T) {
	n := nuevoNumerico(t, formatoEntero)

	n.SetFormato("#.###,0")
	n.SetText("1234,56")
	if n.Text() != "1.234,6" {
		t.Errorf("con el formato %q esperaba %q, obtuve %q", "#.###,0", "1.234,6", n.Text())
	}

	n.SetFormato("#.###,00")
	n.SetText("1234,56")
	if n.Text() != "1.234,56" {
		t.Errorf("con el formato %q esperaba %q, obtuve %q", "#.###,00", "1.234,56", n.Text())
	}
}

// TestSetFormatoSinComasUsaFormatoEntero.
func TestSetFormatoSinComasUsaFormatoEntero(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("1234,6")

	n.SetFormato("#.### sin decimales")
	if n.Text() != "1.235" {
		t.Errorf("un formato sin coma debe dejar solo el entero redondeado, esperaba %q, obtuve %q",
			"1.235", n.Text())
	}
}

// TestSetTextVacioNoDejaResiduo.
func TestSetTextVacioNoDejaResiduo(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("1234,5")
	n.SetText("")

	if n.Text() != "" {
		t.Errorf("SetText(\"\") = %q, se esperaba texto vacio", n.Text())
	}
	if n.FormatoDB() != "" {
		t.Errorf("FormatoDB() = %q, se esperaba texto vacio", n.FormatoDB())
	}
}

// TestSetTextEsIdempotenteAlVolverAFormatear: escribir varias veces sobre el
// mismo valor no lo altera.
func TestSetTextRepetidoNoAlteraElValor(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	for i := 0; i < 3; i++ {
		n.SetText("9876,1")
		if n.Text() != "9.876,10" {
			t.Fatalf("la %d vez esperaba %q, obtuve %q", i+1, "9.876,10", n.Text())
		}
	}
}
