package entradas

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Numerico struct {
	entry      *gtk.Entry
	handlerID  glib.SignalHandle
	formato    string // "#.###" o "#.###,00"
	decimales  int
	procesando bool //Bandera para evitar bucles infinitos al formatear
	esPegado   bool
	popMenu    *gtk.PopoverMenu // Añade este para el menú contextual (clic derecho)
}

// NewNumerico crea el widget.
// formato: "#.###" para enteros o "#.###,00" para decimales (2 decimales)
func NewNumerico(formato string) *Numerico {
	en := &Numerico{
		entry:   gtk.NewEntry(),
		formato: formato,
	}
	en.entry.SetAlignment(1.0) //alineación a la derecha
	en.entry.SetHExpand(true)
	// Calcular cuántos decimales requiere el formato (contar ceros tras la coma)
	en.actualizarConfiguracionFormato()
	// para el focus de entry
	focoCrtl := gtk.NewEventControllerFocus()
	focoCrtl.ConnectEnter(func() {
		en.procesando = true
		// Al entrar en foco, quitar el formato para facilitar la edición
		textoLimpio := aplicarFormatoSinMiles(en.entry.Text())
		en.entry.SetText(textoLimpio)
		en.procesando = false
	})
	focoCrtl.ConnectLeave(func() {
		en.procesando = true
		en.entry.SetText(en.aplicarFormato(en.entry.Text()))
		en.procesando = false
	})
	en.entry.AddController(focoCrtl)
	// El buffer es el que realmente maneja el texto
	buffer := en.entry.Buffer()
	buffer.Connect("inserted-text", func(pos uint, text string, length uint) {
		if en.procesando {
			return
		}
		en.procesando = true
		defer func() { en.procesando = false }()
		// Si es pegado, no validar carácter a carácter
		if length == 1 {
			// Validación del teclado
			if !(text >= "0" && text <= "9") && text != "," {
				buffer.DeleteText(pos, int(length))
				return
			}
			// evitar múltiples comas (el texto ya incluye la que se acaba de
			// teclear, asi que solo se descarta si hay mas de una)
			if text == "," && strings.Count(en.entry.Text(), ",") > 1 {
				buffer.DeleteText(pos, int(length))
				return
			}
		} else {
			// para el pegado del portapapeles
			limpio, esNumero := en.normalizarPortapapeles(text)
			if !esNumero {
				buffer.DeleteText(pos, int(length))
				return
			}
			if limpio == "" {
				return
			}
			buffer.DeleteText(pos, int(length))
			buffer.SetText(limpio, len(limpio))
		}
	})
	// Definí esto dentro de NewNumerico para que tenga acceso a la instancia "en" y su entry
	// procesarPegado lee el portapapeles y, si contiene un número, reemplaza el
	// contenido del entry con ese valor. Devuelve true cuando se hace cargo del
	// pegado, para que el GtkEntry no inserte el texto crudo otra vez.
	procesarPegado := func() bool {
		cb := en.entry.Clipboard()
		if cb == nil {
			return false
		}
		cb.ReadTextAsync(context.Background(), func(res gio.AsyncResulter) {
			text, err := cb.ReadTextFinish(res)
			if err != nil {
				return
			}
			limpio, esNumero := en.normalizarPortapapeles(text)
			// solo se pega si lo que hay en el portapapeles es un número
			if !esNumero || limpio == "" {
				return
			}
			en.procesando = true
			en.entry.SetText(limpio)
			en.procesando = false
		})
		return true
	}
	// Menú contextual
	grupo := gio.NewSimpleActionGroup()
	// Acción: Copiar Numero
	copiar := gio.NewSimpleAction("copiar-normal", nil)
	copiar.Connect("activate", func() {
		en.entry.Clipboard().SetText(en.entry.Text())
	})
	// Acción: Copiar Formato DB (YYYY-MM-DD)
	copiarDB := gio.NewSimpleAction("copiar-db", nil)
	copiarDB.Connect("activate", func() {
		t := strings.ReplaceAll(en.entry.Text(), ".", "")
		t = strings.ReplaceAll(t, ",", ".")
		en.entry.Clipboard().SetText(t)
	})
	// Acción: Pegado Especial
	pegarEspecial := gio.NewSimpleAction("pegar-especial", nil)
	pegarEspecial.Connect("activate", func() {
		procesarPegado()
	})

	grupo.AddAction(copiar)
	grupo.AddAction(copiarDB)
	grupo.AddAction(pegarEspecial)
	en.entry.InsertActionGroup("numero", grupo)

	// 2. Construir el Menú en Español
	menuModel := gio.NewMenu()

	// Sección de Copiado
	seccionCopiar := gio.NewMenu()
	seccionCopiar.Append("Copiar", "numero.copiar-normal")
	seccionCopiar.Append("Copiar Formato DB", "numero.copiar-db")
	menuModel.AppendSection("", seccionCopiar)

	// Sección de Pegado
	seccionPegar := gio.NewMenu()
	seccionPegar.Append("Pegar Especial", "numero.pegar-especial")
	menuModel.AppendSection("", seccionPegar)

	// 3. El Menú Contextual (La "Pieza Maestra")
	en.popMenu = gtk.NewPopoverMenuFromModel(menuModel)
	en.popMenu.SetParent(en.entry)
	en.popMenu.SetHasArrow(false)

	clickGest := gtk.NewGestureClick()
	clickGest.SetButton(3)
	// Fase Capture: Intercepta antes de que el Entry muestre su menú propio
	clickGest.SetPropagationPhase(gtk.PhaseCapture)

	clickGest.ConnectPressed(func(n int, x, y float64) {
		// Validación ANTIPANIC
		if en == nil || en.popMenu == nil {
			return
		}
		clickGest.SetState(gtk.EventSequenceClaimed) // Bloquea el menú original

		en.popMenu.Popup()
	})
	en.entry.AddController(clickGest)

	ctlPegadoEntry := gtk.NewEventControllerKey()
	ctlPegadoEntry.SetPropagationPhase(gtk.PhaseCapture)
	ctlPegadoEntry.ConnectKeyPressed(func(keyval, keycode uint, state gdk.ModifierType) (ok bool) {
		if state.Has(gdk.ControlMask) && keyval == gdk.KEY_v {
			// devolver true consume el atajo para que GtkEntry no pegue
			// el texto crudo del portapapeles por su cuenta
			return procesarPegado()
		}
		// Con Shift esa tecla es "Supr" (avance), no el separador decimal
		if state.Has(gdk.ShiftMask) {
			return false
		}
		// el punto del teclado numérico se escribe como coma decimal
		if !en.procesando && (keyval == gdk.KEY_KP_Decimal || keyval == gdk.KEY_KP_Separator) {
			return en.insertarComa()
		}
		return false
	})
	en.entry.AddController(ctlPegadoEntry)
	return en
}

// insertarComa coloca la coma decimal en la posición del cursor (o en lugar de
// lo que esté seleccionado) y deja el cursor detrás de ella. Se usa cuando se
// presiona el punto del teclado numérico, para poder teclear los decimales sin
// tener que ir a la coma. Devuelve true porque siempre consume la tecla.
func (en *Numerico) insertarComa() bool {
	// solo se admite una coma por número: la tecla se consume y no se escribe
	if strings.Contains(en.entry.Text(), ",") {
		return true
	}
	texto := en.entry.Text()
	pos := en.entry.Position()
	// si hay texto seleccionado, la coma lo reemplaza
	if ini, fin, ok := en.entry.SelectionBounds(); ok {
		pos, texto = ini, texto[:ini]+texto[fin:]
	}
	en.procesando = true
	en.entry.SetText(texto[:pos] + "," + texto[pos:])
	en.entry.SetPosition(pos + 1)
	en.procesando = false
	return true
}

// Object devuelve el widget GTK interno (necesario para Append o Attach)
func (en *Numerico) Object() *gtk.Entry {
	return en.entry
}

// Text devuelve el texto actual del entry
func (en *Numerico) Text() string {
	return en.entry.Text()
}

// SetText asigna texto al entry pasando por el formateador
func (en *Numerico) SetText(texto string) {
	en.procesando = true
	en.entry.SetText(en.aplicarFormato(texto))
	en.procesando = false
}

// SetPlaceholderText expone la propiedad nativa del Entry
func (en *Numerico) SetPlaceholderText(texto string) {
	en.entry.SetPlaceholderText(texto)
}

func (en *Numerico) actualizarConfiguracionFormato() {
	if strings.Contains(en.formato, ",") {
		partes := strings.Split(en.formato, ",")
		en.decimales = len(partes[1])
	} else {
		en.decimales = 0
	}
}

func (en *Numerico) aplicarFormato(input string) string {
	if input == "" {
		return ""
	}
	// Normalizar para procesar (punto para decimal interno)
	temp := strings.ReplaceAll(input, ".", "")
	temp = strings.ReplaceAll(temp, ",", ".")
	valFloat, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		return ""
	}
	var textoFormateado string
	if en.decimales > 0 {
		// formatea decimales
		formaroPrintf := fmt.Sprintf("%%.%df", en.decimales)
		textoFormateado = fmt.Sprintf(formaroPrintf, valFloat)
		// Separar miles y decimales
		partes := strings.Split(textoFormateado, ".")
		enteroFormateado := formatMilesString(partes[0])
		textoFormateado = enteroFormateado + "," + partes[1]
	} else {
		// formatear solo enteros
		enteroStr := fmt.Sprintf("%.0f", valFloat)
		textoFormateado = formatMilesString(enteroStr)
	}
	return textoFormateado
}

func aplicarFormatoSinMiles(input string) string {
	// Normalizar para procesar (punto para decimal interno)
	return strings.ReplaceAll(input, ".", "")
}

// normalizarPortapapeles convierte lo que haya en el portapapeles en el texto
// plano que espera el entry: sin puntos de miles, con la coma como separador
// decimal y con tantos decimales como pide el formato. Acepta el formato
// argentino ("1.234,56"), el de los paises que usan el punto ("1,234.56") y el
// negativo entre parentesis ("(1.234,56)").
//
// El segundo valor dice si lo que venia en el portapapeles es un numero: si es
// falso hay que borrar el texto pegado. Si es true pero el numero viene vacio,
// el portapapeles solo traia decoracion y se deja el entry como estaba.
func (en *Numerico) normalizarPortapapeles(bruto string) (string, bool) {
	// se saca la decoracion: saltos de linea, espacios y signos de moneda
	t := strings.NewReplacer("\n", "", "\r", "", "\t", "", " ", "",
		"$", "", "€", "", "USD", "").Replace(bruto)
	t = strings.TrimSpace(t)
	// los negativos contables van entre parentesis
	negativo := strings.HasPrefix(t, "(") && strings.HasSuffix(t, ")")
	t = strings.Trim(t, "()")
	// el signo menos solo se acepta adelante
	if strings.HasPrefix(t, "-") {
		negativo, t = true, t[1:]
	} else if strings.Contains(t, "-") {
		return "", false
	}
	if t == "" {
		return "", true // solo era decoracion
	}

	// si el numero entero esta escrito con grupos de miles ("1.234.567") no
	// hay parte decimal: el ultimo punto tambien es de miles
	entero, decimal := t, ""
	if !milesValidos(t) {
		i := strings.LastIndexAny(t, ".,")
		if i < 0 || !soloDigitos(t[i+1:]) || !milesValidos(t[:i]) {
			return "", false
		}
		entero, decimal = t[:i], t[i+1:]
	}
	entero = strings.NewReplacer(".", "", ",", "").Replace(entero)
	if negativo {
		entero = "-" + entero
	}
	// los miles y el redondeo los agrega aplicarFormato al salir del campo
	if en.decimales == 0 || decimal == "" {
		return entero, true
	}
	for len(decimal) < en.decimales {
		decimal += "0"
	}
	return entero + "," + decimal, true
}

// milesValidos informa si el entero esta escrito con grupos de tres digitos
// separados siempre por el mismo signo: valen "1234", "1.234" y "1.234.567",
// pero no "12.34", "1.2.3" ni "1,234.56".
func milesValidos(texto string) bool {
	if strings.Count(texto, ".") > 0 && strings.Count(texto, ",") > 0 {
		return false
	}
	grupos := strings.FieldsFunc(texto, func(r rune) bool {
		return r == '.' || r == ','
	})
	if len(grupos) == 0 || !soloDigitos(grupos[0]) {
		return false
	}
	if len(grupos) == 1 {
		return true
	}
	if len(grupos[0]) > 3 {
		return false
	}
	for _, g := range grupos[1:] {
		if len(g) != 3 || !soloDigitos(g) {
			return false
		}
	}
	return true
}

func soloDigitos(texto string) bool {
	for _, c := range texto {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (en *Numerico) FormatoDB() string {
	textoLimpio := en.entry.Text()
	textoLimpio = strings.ReplaceAll(textoLimpio, ".", "")
	textoLimpio = strings.ReplaceAll(textoLimpio, ",", ".")
	return textoLimpio
}

// Función auxiliar para insertar puntos de miles
func formatMilesString(entero string) string {
	var resultado []string
	n := len(entero)

	for i, v := range entero {
		if i > 0 && (n-i)%3 == 0 {
			resultado = append(resultado, ".")
		}
		resultado = append(resultado, string(v))
	}
	return strings.Join(resultado, "")
}

func (en *Numerico) SetFormato(formato string) {
	en.formato = formato

	// 1. Recalcular la configuración de decimales
	en.actualizarConfiguracionFormato()
	// 2. Reaplicar el formato al texto actual
	en.procesando = true
	textoActual := en.entry.Text()
	textoFormateado := en.aplicarFormato(textoActual)
	en.entry.SetText(textoFormateado)
	en.procesando = false
}
