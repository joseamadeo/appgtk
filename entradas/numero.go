package entradas

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type EntradaNumerica struct {
	*gtk.Entry
	handlerID glib.SignalHandle
	formato string // "#.###" o "#.###,00"
	decimales int
	procesando bool //Bandera para evitar bucles infinitos al formatear
	esPegado bool
}

// NewEntradaNumerica crea el widget. 
// formato: "#.###" para enteros o "#.###,00" para decimales (2 decimales)
func NewEntradaNumerica(formato string) *EntradaNumerica {
	en := &EntradaNumerica{
		Entry:   gtk.NewEntry(),
		formato: formato,
	}
	en.SetAlignment(1.0) //alineación a la derecha
	en.SetHExpand(true)
	// Calcular cuántos decimales requiere el formato (contar ceros tras la coma)
	en.actualizarConfiguracionFormato()
	// para el focus de entry
	focoCrtl := gtk.NewEventControllerFocus()
	focoCrtl.ConnectEnter(func ()  {
		en.procesando = true
		// Al entrar en foco, quitar el formato para facilitar la edición
		textoLimpio := aplicarFormatoSinMiles(en.Entry.Text())
		en.Entry.SetText(textoLimpio)
		en.procesando = false
	})
	focoCrtl.ConnectLeave(func ()  {
		en.procesando = true
		en.Entry.SetText(en.aplicarFormato(en.Entry.Text()))
		en.procesando = false
	})
	en.Entry.AddController(focoCrtl)
	// El buffer es el que realmente maneja el texto
	buffer := en.Entry.Buffer()
	buffer.Connect("inserted-text", func(pos uint, text string, length uint) {
		if en.procesando {
			return
		}
		en.procesando = true
		defer func() { en.procesando = false }()
		// Si es pegado, no validar carácter a carácter
		if length == 1 {
			// Validación del teclado
			if !(text >= "0" && text <="9") && text != "," {
				buffer.DeleteText(pos, int(length))
				return
			}
			// evitar múltiples comas
			if strings.Count(en.Entry.Text(), ",") > 1 && (strings.Count(en.formato, ",") == 1  && text == ","){
				buffer.DeleteText(pos, int(length))
				return
			}
		} else {
			// para el pegado del portapapeles
			if text == "" {
				return
			}
			textoLimpio := strings.ReplaceAll(text, " ", "")
			textoLimpio = strings.ReplaceAll(textoLimpio, ".", "")
			textoLimpio = strings.ReplaceAll(textoLimpio, ",", ".")
			if _, err := strconv.ParseFloat(textoLimpio, 64); err != nil {
				buffer.DeleteText(pos, int(length))
				return 
			}
			buffer.DeleteText(pos, int(length))
			textoLimpio = strings.ReplaceAll(textoLimpio, ".", ",")
			buffer.SetText(textoLimpio, len(textoLimpio))
		}
	})
	return en
}

func (en *EntradaNumerica) actualizarConfiguracionFormato(){
	if strings.Contains(en.formato, ",") {
		partes := strings.Split(en.formato, ",")
		en.decimales = len(partes[1])
	} else {
		en.decimales = 0
	}
}

func (en *EntradaNumerica) aplicarFormato(input string) string {
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
	temp := strings.ReplaceAll(input, ".", "")
	temp = strings.ReplaceAll(temp, ",", ".")
	return temp
}

func (en *EntradaNumerica) FormatoDB() string {
	textoLimpio := en.Entry.Text()
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

func (en *EntradaNumerica) SetFormato(formato string){
	en.formato = formato
	
	// 1. Recalcular la configuración de decimales
	en.actualizarConfiguracionFormato()
	// 2. Reaplicar el formato al texto actual
	en.procesando = true
	textoActual := en.Entry.Text()
	textoFormateado := en.aplicarFormato(textoActual)
	en.Entry.SetText(textoFormateado)
	en.procesando = false
}