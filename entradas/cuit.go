package entradas

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type Cuit struct {
	*gtk.Box
	entry *gtk.Entry
}

func NewCuit() *Cuit {
	box := gtk.NewBox(gtk.OrientationHorizontal, 0)
	box.AddCSSClass("linked")
	entry := gtk.NewEntry()
	entry.SetMaxLength(13)
	entry.SetPlaceholderText("XX-XXXXXXXX-X")
	box.Append(entry)
	cuit := &Cuit{
		Box:   box,
		entry: entry,
	}
	cuit.setupPasteFormatCUIT()
	return cuit
}

func (c *Cuit) GetCuit() string {
	return c.entry.Text()
}

func (c *Cuit) GetCuitSinFormato() string {
	cuit := c.entry.Text()
	return strings.ReplaceAll(cuit, "-", "")
}

func (c *Cuit) SetCuit(cIng string) {
	if !isFormattedCUIT(cIng) {
		return
	}
	c.entry.SetText(cIng)
}

func (c *Cuit) setupPasteFormatCUIT() {
	var handleID glib.SignalHandle
	handleID = c.entry.Connect("changed", func ()  {
		// 1. Bloquea el manejador para evitar recursión infinita
		c.entry.HandlerBlock(handleID)
		// 2. Obtener solo los números y limitar a 11 dígitos
		rewText := c.entry.Text()
		nums := ""
		for _, r := range rewText {
			if unicode.IsDigit(r) {
				nums += string(r)
			}
		}
		if len(nums) > 11 {
			nums = nums[:11]
		}
		// 3. formatear dinámicamente el CUIT
		formatted := ""
		for i, r := range nums {
			if i == 2 || i == 10 {
				formatted += "-"
			}
			formatted += string(r)
		}
		// 4. Actualizar el Entry
		c.entry.SetText(formatted)
		c.entry.SetPosition(len(formatted))
		// 5. Aplicar el estilo de validación
		c.EstiloEntradaCUIT()
		// 6. Desbloquear el manejador
		c.entry.HandlerUnblock(handleID)
	})
	c.EstiloEntradaCUIT()
}

func formatCUITString(input string) string {
	re := regexp.MustCompile(`\d`)
	nums := strings.Join(re.FindAllString(input, -1), "")
	if len(nums) > 8 {
		return nums[:8]
	}
	var sd strings.Builder
	for i, r := range nums {
		if i == 2 || i == 10 {
			sd.WriteRune('-')
		}
		sd.WriteRune(r)
	}
	return sd.String()
}

func isFormattedCUIT(text string) bool {
	match, _ := regexp.MatchString(`^\d{2}-\d{8}-\d{1}$`, text)
	fmt.Println(match)
	return match
}

func (dc *Cuit) EstiloEntradaCUIT() {
	// remueve todas las clases CSS de validación
	dc.entry.RemoveCSSClass("Valido")
	dc.entry.RemoveCSSClass("Invalido")
	dc.entry.RemoveCSSClass("Neutro")
	texto := dc.GetCuitSinFormato()
	if len(texto) == 0 {
		dc.entry.AddCSSClass("Neutro")
		return
	}
	if len(texto) == 11 {
		if validateCuit(texto) {
			dc.entry.AddCSSClass("Valido")
		} else {
			dc.entry.AddCSSClass("Invalido")
		}
	} else {
		dc.entry.AddCSSClass("Invalido")
	}
}

func validateCuit(cuit string) bool {
	// Un CUIT debe tener exactamente 11 dígitos numéricos
	if len(cuit) != 11 {
		return false
	}
	mult := []int{5, 4, 3, 2, 7, 6, 5, 4, 3, 2}
	suma := 0

	for i := 0; i < 10; i++ {
		val := int(cuit[i] - '0')
		suma += val * mult[i]
	}
	rev := 11 - (suma % 11)
	if rev == 11 {
		rev = 0
	} else if rev == 10 {
		// en caso 10 es especial en el CUIT, pero generalmente se resuelve
		// en la asignación del número de CUIT previo
		rev = 9
	}
	dv := int(cuit[10] - '0')

	return dv == rev
}
