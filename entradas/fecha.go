package entradas

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Fecha struct {
	*gtk.Box
	entry   *gtk.Entry
	cal     *gtk.Calendar
	Min 		time.Time //Fecha mínima permitida
	Max 		time.Time //Fecha máxima permitida
	popover *gtk.Popover      // Este es el del calendario (botón)
	popMenu *gtk.PopoverMenu  // Añade este para el menú contextual (clic derecho)
	cursorPosicion int
	borrando bool // Variable para detectar si se está borrando un carácter
}

func NewFecha() *Fecha {
	box := gtk.NewBox(gtk.OrientationHorizontal, 0)
	box.AddCSSClass("linked") // Estilo visual para unir entry y button

	entry := gtk.NewEntry()
	entry.SetPlaceholderText("DD/MM/AAAA")
	entry.SetMaxLength(10) // Formato DD/MM/AAAA tiene 10 caracteres

	btn := gtk.NewButtonFromIconName("office-calendar-symbolic")
	btn.AddCSSClass("calendar-icon-button") // Clase CSS personalizada para el botón
	popover := gtk.NewPopover()
	popover.SetParent(btn)
	
	calendar := gtk.NewCalendar()
	popover.SetChild(calendar)
	btn.ConnectClicked(func() {
		popover.Popup()
	})

	// CREAR EL OBJETO dp INMEDIATAMENTE
	dp := &Fecha{
		Box:     box,
		entry:   entry,
		popover: popover,
		cal:     calendar,
		Min: 		time.Now().AddDate(-100, 0, 0), // Fecha mínima por defecto
		Max: 		time.Now().AddDate(100, 0, 0), // Fecha máxima por defecto
	}
	
	// Definí esto dentro de NewFecha para que tenga acceso a 'dp'
	procesarPegado := func() {
		cb := dp.entry.Clipboard()
		cb.ReadTextAsync(context.Background(), func(res gio.AsyncResulter) {
			text, err := cb.ReadTextFinish(res)
			if err != nil || text == "" {
				return
			}
			fechaObj, fechaStr, err := normalizarFechaEstricta(text)
			if err != nil || fechaObj.Before(dp.Min) || fechaObj.After(dp.Max) {
				glib.IdleAdd(func() {
					dp.entry.SetText("")
					dp.cursorPosicion = dp.entry.Position()
					//dp.EstiloEntradaFecha()
				})
				return
			}
			glib.IdleAdd(func() {
				dp.entry.SetText(fechaStr)
			})
			dt := glib.NewDateTimeLocal(fechaObj.Year(), int(fechaObj.Month()), fechaObj.Day(), 0, 0, 0)
			dp.cal.SelectDay(dt)
			dp.cursorPosicion = dp.entry.Position()
			//dp.EstiloEntradaFecha()
		})
	}

	grupo := gio.NewSimpleActionGroup()
	// Acción: Copiar Fecha (DD/MM/AAAA)
	copiarFecha := gio.NewSimpleAction("copiar-normal", nil)
	copiarFecha.ConnectActivate(func(_ *glib.Variant) {
		entry.Clipboard().SetText(entry.Text())
	})
	// Acción: Copiar Formato DB (YYYY-MM-DD)
	copiarDB := gio.NewSimpleAction("copiar-db", nil)
	copiarDB.ConnectActivate(func(_ *glib.Variant) {
		t, err := time.Parse("02/01/2006", entry.Text())
		if err == nil {
			entry.Clipboard().SetText(t.Format("2006-01-02"))
		}
	})
	// Acción: Pegado Especial (Tu lógica ya existente)
	pegarEspecial := gio.NewSimpleAction("pegar-especial", nil)
	pegarEspecial.ConnectActivate(func(_ *glib.Variant) {
		procesarPegado()
	})

	grupo.AddAction(copiarFecha)
	grupo.AddAction(copiarDB)
	grupo.AddAction(pegarEspecial)
	entry.InsertActionGroup("fecha", grupo)

	// 2. Construir el Menú en Español
	menuModel := gio.NewMenu()

	// Sección de Copiado
	seccionCopiar := gio.NewMenu()
	seccionCopiar.Append("Copiar Fecha (DD/MM/AAAA)", "fecha.copiar-normal")
	seccionCopiar.Append("Copiar Formato DB (AAAA/MM/DD)", "fecha.copiar-db")
	menuModel.AppendSection("", seccionCopiar)

	// Sección de Pegado
	seccionPegar := gio.NewMenu()
	seccionPegar.Append("Pegar Fecha Especial", "fecha.pegar-especial")
	menuModel.AppendSection("", seccionPegar)

	// 3. El Menú Contextual (La "Pieza Maestra")
	dp.popMenu = gtk.NewPopoverMenuFromModel(menuModel)
	dp.popMenu.SetParent(entry)
	dp.popMenu.SetHasArrow(false)

	clickGest := gtk.NewGestureClick()
	clickGest.SetButton(3)
	// Fase Capture: Intercepta antes de que el Entry muestre su menú propio
	clickGest.SetPropagationPhase(gtk.PhaseCapture) 

	clickGest.ConnectPressed(func(n int, x, y float64) {
		// Validación ANTIPANIC
		if dp == nil || dp.popMenu == nil {
			return 
		}
		clickGest.SetState(gtk.EventSequenceClaimed) // Bloquea el menú original

		dp.popMenu.Popup()
	})
	dp.entry.AddController(clickGest)
	
	ctlTeclaPress := gtk.NewEventControllerKey()
	ctlTeclaPress.SetPropagationPhase(gtk.PhaseCapture)
	ctlTeclaPress.ConnectKeyPressed(func(keyval, keycode uint, state gdk.ModifierType) (ok bool) {
		if state.Has(gdk.ControlMask) && keyval == gdk.KEY_v {
      procesarPegado()
      return false 
    }
		// 1. Permitir números (del 0 al 9 en teclado alfanumérico y numérico)
    if (keyval >= gdk.KEY_0 && keyval <= gdk.KEY_9) || (keyval >= gdk.KEY_KP_0 && keyval <= gdk.KEY_KP_9) {
			dp.cursorPosicion = dp.entry.Position()
			//dp.EstiloEntradaFecha()
      return false // Permitir: GTK procesa la inserción
    }
		if keyval == gdk.KEY_BackSpace || keyval == gdk.KEY_Delete {
			// Permitir que Backspace y Suprimir funcionen normalmente sin interferencia del código de auto-formato
			dp.borrando = true
			dp.cursorPosicion = dp.entry.Position()
			//dp.EstiloEntradaFecha() // Actualiza el estilo para reflejar el cambio
			return false // Deja que el evento se propague para que el entry lo procese
		}
		dp.borrando = false // Resetea la bandera si se presiona cualquier otra tecla que no sea Backspace o Delete
		// 3. Permitir TODAS las flechas de dirección (Izquierda, Derecha, Arriba, Abajo)
		if keyval == gdk.KEY_Left || keyval == gdk.KEY_Right || keyval == gdk.KEY_Up || keyval == gdk.KEY_Down {
			pos := dp.entry.Position()
			if keyval == gdk.KEY_Left {
				pos = max(0, pos-1)
			} else if keyval == gdk.KEY_Right {
				pos = min(len(dp.entry.Text()), pos+1)
			}
			dp.cursorPosicion = pos
			dp.entry.SetPosition(pos)
			return true
		}
		if keyval == gdk.KEY_Tab {
			dp.cursorPosicion = dp.entry.Position()
			//dp.EstiloEntradaFecha() // Actualiza el estilo para reflejar el cambio de enfoque
			return false
		}
		dp.cursorPosicion = dp.entry.Position()
		//dp.EstiloEntradaFecha() // Actualiza el estilo para reflejar cualquier otro cambio
		return true // Bloquear cualquier otra tecla
	})
	// DETECTAR CUANDO SE SUELTA LA TECLA (Clave para apagar la bandera)
	ctlTeclaPress.ConnectKeyReleased(func(keyval, keycode uint, state gdk.ModifierType) {
		if keyval == gdk.KEY_BackSpace || keyval == gdk.KEY_Delete {
			dp.borrando = false // <--- APAGAR BANDERA
		}
	})
	entry.AddController(ctlTeclaPress)

	// Al seleccionar una fecha en el calendario
	calendar.ConnectDaySelected(func() {
		fecha := calendar.Date()
		fechaFormateada := fmt.Sprintf("%02d/%02d/%04d", fecha.DayOfMonth(), fecha.Month(), fecha.Year())
		entry.SetText(fechaFormateada)
		entry.SetPosition(-1) // Mueve el cursor al final del texto
		dp.cursorPosicion = dp.entry.Position()
		//dp.EstiloEntradaFecha()
		popover.Popdown()
	})
	box.Append(entry)
	box.Append(btn)
	dp.setupAutoFormatoTeclado()

	return dp
} 

func (dp *Fecha) setupAutoFormatoTeclado() {
	var handleID glib.SignalHandle
	var textoAnterior string

	handleID = dp.entry.Connect("changed", func() {
		// 1. Obtener texto y LIMPIAR inmediatamente caracteres basura
		// Esto limpia el RLM (U+200F), saltos de línea y retornos
		rawText := dp.entry.Text()
		textoLimpio := strings.NewReplacer(
			"\u200F", "",
			"\n", "",
			"\r", "",
		).Replace(rawText)

		// El auto-formato solo se aplica al escribir (no al borrar) y con texto
		if !dp.borrando && textoLimpio != "" {
			largoActual := len(textoLimpio)
			largoAnterior := len(textoAnterior)

			esPegado := (textoLimpio != textoAnterior) && (largoActual != largoAnterior+1) && (largoActual >= largoAnterior)

			dp.entry.HandlerBlock(handleID)
			if esPegado {
				// --- SOLUCIÓN PEGAR ---
				re := regexp.MustCompile(`[^0-9]`)
				soloNumeros := re.ReplaceAllString(textoLimpio, "")
				textoFinal := formatDateString(soloNumeros)
				dp.entry.SetText(textoFinal)
				dp.entry.SetPosition(-1)
			} else if largoActual > largoAnterior {
				// --- SOLUCIÓN ESCRITURA MANUAL ---
				// Filtramos para que no entren letras ni caracteres raros
				re := regexp.MustCompile(`[^0-9/]`)
				if re.MatchString(textoLimpio) {
					textoLimpio = re.ReplaceAllString(textoLimpio, "")
					dp.entry.SetText(textoLimpio)
					dp.entry.SetPosition(-1)
				}

				// Auto-formato de barras (solo si está escribiendo)
				if (len(textoLimpio) == 2 || len(textoLimpio) == 5) && !strings.HasSuffix(textoLimpio, "/") {
					dp.entry.SetText(textoLimpio + "/")
					dp.entry.SetPosition(-1)
				}
			}
			dp.entry.HandlerUnblock(handleID)
		}
		textoAnterior = dp.entry.Text()
		dp.cursorPosicion = dp.entry.Position()
		// El estilo se aplica SIEMPRE (escribir, borrar o vaciar)
		dp.EstiloEntradaFecha()
	})
}

// Método para obtener la fecha como time.Time
func (d *Fecha) GetFecha() (time.Time, error) {
	texto := d.entry.Text()
	if texto == "" {
		return time.Time{}, fmt.Errorf("El campo de fecha está vacío")
	}
	// validar la longitud mínima
	if len(texto) != 10 {
		return time.Time{}, fmt.Errorf("La fecha debe tener el formato DD/MM/AAAA")
	}
	// Intentar parsear a objeto time.Time para validar existencia (ej. no 31/02)
	f, err := time.Parse("02/01/2006", texto)
	if err != nil {
		return time.Time{}, fmt.Errorf("la fecha ingresada no es válida: %w", err)
	}

	// Validar contra el rango Min y Max configurado en el struct
	if f.Before(d.Min) {
		return time.Time{}, fmt.Errorf("la fecha es anterior al mínimo permitido (%s)", d.Min.Format("02/01/2006"))
	}
	if f.After(d.Max) {
		return time.Time{}, fmt.Errorf("la fecha es posterior al máximo permitido (%s)", d.Max.Format("02/01/2006"))
	}
	// Si todo es correcto, devolvemos el texto tal cual está en el entry
	return f, nil
}

func (d *Fecha) GetFechaDB() (string, error) {
	// devuelve la fecha en formato AAAA-MM-DD para base de datos
	texto := d.entry.Text()
	// 1. Intentar parsear el texto actual del entry
	fecha, err := time.Parse("02/01/2006", texto)
	if err != nil {
		return "", fmt.Errorf("la fecha no es válida para base de datos: %w", err)
	}
	// 2. Validar contra el rango Min y Max configurado en el struct
	if fecha.Before(d.Min) {
		return "", fmt.Errorf("la fecha es anterior al mínimo permitido (%s)", d.Min.Format("02/01/2006"))
	}
	if fecha.After(d.Max) {
		return "", fmt.Errorf("la fecha es posterior al máximo permitido (%s)", d.Max.Format("02/01/2006"))
	}
	return fecha.Format("2006-01-02"), nil
}

func (dp *Fecha) SetFecha(fecha string) error {
	if !isFormattedFecha(fecha) {
		return fmt.Errorf("la fecha no tiene el formato correcto (DD/MM/AAAA)")
	}
	// Intentar parsear para validar que sea una fecha real
	f, err := time.Parse("02/01/2006", fecha)
	if err != nil {
		return fmt.Errorf("la fecha no es válida: %w", err)
	}
	// Validar contra el rango Min y Max configurado en el struct
	if f.Before(dp.Min) || f.After(dp.Max) {
		return fmt.Errorf("la fecha está fuera del rango permitido")
	}
	dp.cursorPosicion = dp.entry.Position()
	dp.entry.SetText(fecha)
	//dp.EstiloEntradaFecha()
	return nil
}

func isFormattedFecha(text string) bool {
	match, _ := regexp.MatchString(`^\d{2}/\d{2}/\d{4}$`, text)
	return match
}

func formatDateString(input string) string {
	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`\d`)
	nums := strings.Join(re.FindAllString(input, -1), "")
	// Si tenemos 6 dígitos (DDMMYY), lo convertimos a 8 (DDMM20YY)
	if len(nums) > 8 {
		nums = nums[:8]
	}
	var sd strings.Builder
	for i, r := range nums {
		if i == 2 || i == 4 {
			sd.WriteRune('/')
		}
		sd.WriteRune(r)
	}
	return sd.String()
}

func (dp *Fecha) EstiloEntradaFecha() {
	text := dp.entry.Text()
	// remueve todas las clases CSS de validación
	dp.entry.RemoveCSSClass("Valido")
	dp.entry.RemoveCSSClass("Invalido")
	dp.entry.RemoveCSSClass("Neutro")
	dp.entry.SetTooltipText("")

	if text == "" {
		dp.entry.AddCSSClass("Neutro")
		dp.entry.SetTooltipText("Ingrese una fecha DD/MM/AAAA")
		dp.entry.SetPosition(dp.cursorPosicion)
		return
	}
	f, err := time.Parse("02/01/2006", text)
	if err == nil && len(text) == 10 {
		// Aquí es donde se valida el rango que definiste
		if f.Before(dp.Min) {
			dp.entry.AddCSSClass("Invalido")
			dp.entry.SetTooltipText("Fecha minima: " + dp.Min.Format("02/01/2006"))
		} else if f.After(dp.Max) {
			dp.entry.AddCSSClass("Invalido")
			dp.entry.SetTooltipText("La fecha máxima: " + dp.Max.Format("02/01/2006"))
		} else {
			dp.entry.AddCSSClass("Valido")
			dp.entry.SetTooltipText("Fecha correcta")
		}
	} else {
		dp.entry.AddCSSClass("Invalido")
		dp.entry.SetTooltipText("Fecha inválida")
	}
	dp.entry.SetPosition(dp.cursorPosicion)
}

func (dp *Fecha) EsValido() bool{
	text := dp.entry.Text()
	if text == "" {
		return false
	}
	f, err := time.Parse("02/01/2006", text)
	if err == nil && len(text) == 10 {
		// Aquí es donde se valida el rango que definiste
		if f.Before(dp.Min) {
			return false
		} else if f.After(dp.Max) {
			return  false
		} else {
			return true
		}
	} else {
		return false
	}
}

func normalizarFechaEstricta(entrada string) (time.Time, string, error) {
	// Limpiar espacios y normalizar separadores (por si usan - o .)
	entrada = strings.ReplaceAll(entrada, "-", "/")
	entrada = strings.ReplaceAll(entrada, ".", "/")
	entrada = strings.TrimSpace(entrada)

	formatos := []string{
		"02/01/2006", // DD/MM/YYYY
		"2006/01/02", // YYYY/MM/DD
	}

	for _, f := range formatos {
		t, err := time.Parse(f, entrada)
		if err == nil {
			// REGLA 1: Validar que el año tenga 4 dígitos si el formato lo pide
			// Esto evita que "2" sea un año válido en "2006"
			if (f == "02/01/2006" || f == "2006/01/02") && t.Year() < 1000 {
				continue
			}

			// REGLA 2: Validar reconstrucción
			// time.Parse "arregla" fechas como 31/04 a 01/05. 
			// Si al formatear de vuelta no coincide con la entrada, la fecha es inválida.
			if t.Format(f) != entrada {
				continue
			}
			return t, t.Format("02/01/2006"), nil
		}
	}
	return time.Time{}, "", fmt.Errorf("formato no reconocido")
}
