package test

import (
	"testing"
	"time"
)

// TestSoloSeAceptanDigitosYLaComaAlTipear: el teclado solo deja pasar numeros y
// la coma decimal; letras, signo menos, punto, espacio y otros simbolos se
// borran en el momento.
func TestSoloSeAceptanDigitosYLaComaAlTipear(t *testing.T) {
	casos := []struct {
		tecla  string
		quiere string
	}{
		{"0", "0"},
		{"1", "1"},
		{"5", "5"},
		{"9", "9"},
		{"a", ""},
		{"A", ""},
		{"z", ""},
		{"$", ""},
		{"%", ""},
		{"-", ""}, // el negativo solo entra por SetText o por pegado
		{".", ""}, // el punto de miles se agrega solo al formatear
		{" ", ""},
		{"\n", ""},
		{"\t", ""},
		{",", ","}, // se admite la coma decimal
	}

	for _, c := range casos {
		t.Run("tecla_"+c.tecla, func(t *testing.T) {
			n := nuevoNumerico(t, formatoDecimal)
			tipar(n, c.tecla)
			if n.Text() != c.quiere {
				t.Errorf("tras tipear %q el texto = %q, se esperaba %q", c.tecla, n.Text(), c.quiere)
			}
		})
	}
}

// TestSoloSeAceptaUnaComaDecimal: en un formato con decimales la segunda coma
// se descarta.
func TestSoloSeAceptaUnaComaDecimal(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)

	for _, tecla := range []string{"1", ",", "2", ",", "3", ","} {
		tipar(n, tecla)
	}

	if n.Text() != "1,23" {
		t.Errorf("texto = %q, se esperaba %q (una sola coma)", n.Text(), "1,23")
	}
}

// TestElPuntoDeMilesNoSeEscribePorTeclado: el punto solo lo agrega el
// formateador, tipearlo se descarta. Se escribe letra por letra como lo haria
// el usuario: el punto nunca llega a ser parte de una insercion multiple, que
// es el camino del pegado.
func TestElPuntoDeMilesNoSeEscribePorTeclado(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	for _, tecla := range []string{"1", ".", "2"} {
		tipar(n, tecla)
	}

	if n.Text() != "12" {
		t.Errorf("tras tipear \"1.2\" el texto = %q, se esperaba %q (el punto se descarta)", n.Text(), "12")
	}
}

// TestComaSeAceptaAlTipearEnFormatoEntero documenta el comportamiento actual:
// el filtro de la coma mira el formato, asi que con formato entero la coma
// tambien se puede escribir aunque despues se descarte al formatear.
func TestComaSeAceptaAlTipearEnFormatoEntero(t *testing.T) {
	n := nuevoNumerico(t, formatoEntero)
	tipar(n, ",")

	if n.Text() != "," {
		t.Errorf("texto = %q, se esperaba %q", n.Text(), ",")
	}
	if db := n.FormatoDB(); db != "." {
		t.Errorf("FormatoDB() = %q, se esperaba %q (solo cambia los separadores, no valida)", db, ".")
	}
}

// TestPegadoDesdeElPortapapelesAceptaUnNumeroConFormatoDeExcel.
func TestPegadoDesdeElPortapapelesAceptaUnNumero(t *testing.T) {
	casos := []struct {
		portapapeles string
		quiere       string
	}{
		{"1.234,56", "1234,56"},
		{"1234,56", "1234,56"},
		{"1.234,56 ", "1234,56"},
		{"1234.56", "1234,56"},     // el punto decimal de los paises que no usan coma
		{"1,234.56", "1234,56"},    // miles con coma y decimal con punto
		{"(1.234,56)", "-1234,56"}, // los negativos contables van entre parentesis
		{"1.000.000", "1000000"},
		{"1234567", "1234567"},
		{"12 34", "1234"},  // los espacios se eliminan
		{" 1234 ", "1234"}, // se recorta
		{"$ 3.000", "3000"},
		{"3.000\n", "3000"},
		{"3.000\r\n", "3000"},
		{"-1500", "-1500"}, // por pegado si se admite el negativo
		{"12,5", "12,50"},  // los decimales se completan segun el formato
	}

	for _, c := range casos {
		t.Run(c.portapapeles, func(t *testing.T) {
			n := nuevoNumerico(t, formatoDecimal)
			pegar(n, c.portapapeles)
			if n.Text() != c.quiere {
				t.Errorf("pegar %q = %q, se esperaba %q", c.portapapeles, n.Text(), c.quiere)
			}
		})
	}
}

// TestPegadoRechazaLoQueNoEsUnNumero.
func TestPegadoRechazaLoQueNoEsUnNumero(t *testing.T) {
	casos := []string{
		"abc",
		"12abc",
		"1.2.3.4x",
		"12,34,56", // mas de una coma decimal
		"1,2,3",
		"50%",
		"$",
		"hola mundo",
		// cosas que strconv.ParseFloat acepta pero no son numeros para el entry
		"NaN",
		"Inf",
		"-Inf",
		"1e5",
		"0x10",
		"+1500",
	}

	for _, portapapeles := range casos {
		t.Run(portapapeles, func(t *testing.T) {
			n := nuevoNumerico(t, formatoDecimal)
			pegar(n, portapapeles)
			if n.Text() != "" {
				t.Errorf("pegar %q = %q, se esperaba texto vacio", portapapeles, n.Text())
			}
		})
	}
}

// TestPegadoRechazaTextoYConservaLoQueYaEstabaEscrito.
func TestPegadoRechazaTextoYConservaLoQueYaEstabaEscrito(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	tipar(n, "9999")
	antes := n.Text()

	pegar(n, "no es un numero")

	if n.Text() != antes {
		t.Errorf("tras pegar texto invalido el texto = %q, se esperaba %q", n.Text(), antes)
	}
}

// TestPegadoDeUnNumeroReemplazaElContenidoAnterior.
func TestPegadoDeUnNumeroReemplazaElContenidoAnterior(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	tipar(n, "9999")

	pegar(n, "1.234,56")

	if n.Text() != "1234,56" {
		t.Errorf("tras pegar el texto = %q, se esperaba %q (el pegado reemplaza todo)",
			n.Text(), "1234,56")
	}
}

// TestPegadoSoloConEspaciosNoSeBorra documenta que un pegado que queda vacio
// (por ejemplo "  ") se ignora y queda tal cual en el entry.
func TestPegadoSoloConEspaciosNoSeBorra(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	pegar(n, "  ")

	if n.Text() != "  " {
		t.Errorf("pegar %q = %q, se esperaba %q", "  ", n.Text(), "  ")
	}
}

// TestPegadoRealDesdeElPortapapelesDelSistema: se copia un numero al
// portapapeles y se pega como lo haria el usuario (Ctrl+V).
func TestPegadoRealDesdeElPortapapelesDelSistema(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	ponerPortapapeles(t, n, "1.234,56")
	tipar(n, "77")

	pegar(n, leerPortapapeles(t, n))

	if n.Text() != "1234,56" {
		t.Errorf("tras pegar el portapapeles el texto = %q, se esperaba %q", n.Text(), "1.234,56")
	}
}

// TestPegadoDesdeElPortapapelesIgnoraUnNumeroDelSistema.
func TestPegadoDesdeElPortapapelesIgnoraTextoNoNumerico(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	ponerPortapapeles(t, n, "hola mundo")
	tipar(n, "77")

	pegar(n, leerPortapapeles(t, n))

	if n.Text() != "77" {
		t.Errorf("tras pegar texto no numerico el texto = %q, se esperaba %q", n.Text(), "77")
	}
}

// TestAccionesDelMenuContextualSePuedenActivar: "copiar", "copiar DB" y "pegar
// especial" estan registradas en el grupo de acciones "numero" del entry.
func TestAccionesDelMenuContextualSePuedenActivar(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("1234,5")

	for _, accion := range []string{"numero.copiar-normal", "numero.copiar-db", "numero.pegar-especial"} {
		if !n.Object().ActivateAction(accion, nil) {
			t.Errorf("no se pudo activar la accion %q", accion)
		}
	}
}

// TestAccionCopiarPoneElTextoEnElPortapapeles.
func TestAccionCopiarPoneElTextoEnElPortapapeles(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("1234,5")
	esperado := n.Text()

	n.Object().ActivateAction("numero.copiar-normal", nil)

	if obtenido := leerPortapapeles(t, n); obtenido != esperado {
		t.Errorf("portapapeles = %q, se esperaba %q", obtenido, esperado)
	}
}

// TestAccionCopiarDBPoneElNumeroEnFormatoBaseDeDatos.
func TestAccionCopiarDBPoneElNumeroEnFormatoBaseDeDatos(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("1.234.567,89")
	esperado := n.FormatoDB()

	n.Object().ActivateAction("numero.copiar-db", nil)

	if obtenido := leerPortapapeles(t, n); obtenido != esperado {
		t.Errorf("portapapeles = %q, se esperaba %q", obtenido, esperado)
	}
}

// TestPegarEspecialPegaElNumeroDelPortapapeles: la accion "pegar-especial" toma
// el numero del portapapeles del sistema y reemplaza el contenido del entry.
func TestPegarEspecialPegaElNumeroDelPortapapeles(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("1234,5")
	ponerPortapapeles(t, n, "9.999,99")

	n.Object().ActivateAction("numero.pegar-especial", nil)
	esperar(2*time.Second, func() bool { return n.Text() == "9999,99" })

	if n.Text() != "9999,99" {
		t.Errorf("la accion pegar-especial dejo el texto = %q, se esperaba %q",
			n.Text(), "9999,99")
	}
}

// TestPegarEspecialIgnoraUnTextoNoNumerico: si el portapapeles no tiene un
// numero, la accion no toca lo que ya estaba en el entry.
func TestPegarEspecialIgnoraUnTextoNoNumerico(t *testing.T) {
	n := nuevoNumerico(t, formatoDecimal)
	n.SetText("1234,5")
	antes := n.Text()
	ponerPortapapeles(t, n, "hola mundo")

	n.Object().ActivateAction("numero.pegar-especial", nil)
	esperar(300*time.Millisecond, func() bool { return false })

	if n.Text() != antes {
		t.Errorf("la accion pegar-especial cambio el texto a %q, se esperaba %q",
			n.Text(), antes)
	}
}

// TestPegarEspecialRechazaLoQueNoEsUnNumeroDelPortapapeles: la accion descarta
// el portapapeles entero si no representa un numero, en vez de pegarlo.
func TestPegarEspecialRechazaLoQueNoEsUnNumeroDelPortapapeles(t *testing.T) {
	casos := []string{"abc", "12abc", "1,2,3", "50%", "hola mundo", "$", "  "}

	for _, portapapeles := range casos {
		t.Run(portapapeles, func(t *testing.T) {
			n := nuevoNumerico(t, formatoDecimal)
			n.SetText("1234,5")
			antes := n.Text()
			ponerPortapapeles(t, n, portapapeles)

			n.Object().ActivateAction("numero.pegar-especial", nil)
			esperar(300*time.Millisecond, func() bool { return false })

			if n.Text() != antes {
				t.Errorf("con el portapapeles en %q el texto = %q, se esperaba %q",
					portapapeles, n.Text(), antes)
			}
		})
	}
}
