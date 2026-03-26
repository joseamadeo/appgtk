package entradas

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Fecha struct {
	*gtk.Box
	entry   *gtk.Entry
	popover *gtk.Popover
	cal     *gtk.Calendar
	Min 		time.Time //Fecha mínima permitida
	Max 		time.Time //Fecha máxima permitida
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
	// Configuración inicial del calendario
	dp := &Fecha{
		Box:     box,
		entry:   entry,
		popover: popover,
		cal:     calendar,
		Min: 		time.Now().AddDate(-100, 0, 0), // Fecha mínima por defecto
		Max: 		time.Now().AddDate(100, 0, 0), // Fecha máxima por defecto
	}

	btn.ConnectClicked(func() {
		popover.Popup()
	})

	// Al seleccionar una fecha en el calendario
	calendar.ConnectDaySelected(func() {
		fecha := calendar.Date()
		fechaFormateada := fmt.Sprintf("%02d/%02d/%04d", fecha.DayOfMonth(), fecha.Month(), fecha.Year())
		entry.SetText(fechaFormateada)
		entry.SetPosition(-1) // Mueve el cursor al final del texto
		dp.EstiloEntradaFecha()
		popover.Popdown()
	})
	box.Append(entry)
	box.Append(btn)
	dp.setupAutoFormatoTeclado()

	return dp
}

func (dp *Fecha) setupAutoFormatoTeclado() {
	var handleID glib.SignalHandle
	var largoPrevio int // Variable para rastrear si borramos o escribimos

	handleID = dp.entry.Connect("changed", func() {
		texto := dp.entry.Text()
		largoActual := len(texto)
		dp.entry.HandlerBlock(handleID)
		// 1. Si es un pegado desde el portapapeles, formateamos el bloque completo
		if largoActual > largoPrevio + 1 {
			texto = formatDateString(texto)
			dp.entry.SetText(texto)
			dp.entry.SetPosition(len(texto))
		} else if largoActual > largoPrevio {
			// 2. Si es escritura normal (Teclado)
			if (largoActual == 2 || largoActual == 5) {
				if !strings.HasSuffix(texto, "/") {
					dp.entry.SetText(texto + "/")
					dp.entry.SetPosition(-1)
				}
			}
		}
		largoPrevio = len(dp.entry.Text())
		dp.entry.HandlerUnblock(handleID)
		// 3. Validar Estilo (Min/Max)
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

func (d *Fecha) SetFecha(fecha string) error {
	if !isFormattedFecha(fecha) {
		return fmt.Errorf("la fecha no tiene el formato correcto (DD/MM/AAAA)")
	}
	// Intentar parsear para validar que sea una fecha real
	f, err := time.Parse("02/01/2006", fecha)
	if err != nil {
		return fmt.Errorf("la fecha no es válida: %w", err)
	}
	// Validar contra el rango Min y Max configurado en el struct
	if f.Before(d.Min) || f.After(d.Max) {
		return fmt.Errorf("la fecha está fuera del rango permitido")
	}
	d.entry.SetText(fecha)
	d.EstiloEntradaFecha()
	return nil
}

func isFormattedFecha(text string) bool {
	match, _ := regexp.MatchString(`^\d{2}/\d{2}/\d{4}$`, text)
	return match
}

func formatDateString(input string) string {
	re := regexp.MustCompile(`\d`)
	nums := strings.Join(re.FindAllString(input, -1), "")
	if len(nums) > 8 {
		return nums[:8]
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
			// Actualizar el calendario interno
			dp.cal.SetYear(f.Year())
			dp.cal.SetMonth(int(f.Month()) - 1) // Los meses en gtk.Calendar son 0-indexed
			dp.cal.SetDay(f.Day())
		}
	} else {
		dp.entry.AddCSSClass("Invalido")
		dp.entry.SetTooltipText("Fecha inválida")
	}
}

func (dp *Fecha) setupPasteFormatFecha() {
	var handletID glib.SignalHandle
	handletID = dp.entry.Connect("changed", func ()  {
		text := dp.entry.Text()
		if text == "" || len(text) == 10 {
			dp.EstiloEntradaFecha()
			return
		}
		formatted := formatDateString(text)
		dp.entry.HandlerBlock(handletID)
		
		glib.IdleAdd(func() {
			dp.entry.SetText(formatted)
			dp.entry.SetPosition(len(dp.entry.Text()))
		})
		dp.EstiloEntradaFecha()
	})
}
