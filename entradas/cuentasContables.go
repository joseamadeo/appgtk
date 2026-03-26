package entradas

import (
	"strings"
	"unicode"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type CuentaContable struct {
	*gtk.Box
	entry *gtk.Entry
}

func NewCuentaContable() *CuentaContable {
	box := gtk.NewBox(gtk.OrientationHorizontal, 0)
	box.AddCSSClass("linked")
	
	entry := gtk.NewEntry()
	// Formato 0-00-00-00 tiene 10 caracteres totales
	entry.SetMaxLength(10) 
	entry.SetPlaceholderText("0-00-00-00")
	
	box.Append(entry)
	
	cc := &CuentaContable{
		Box:   box,
		entry: entry,
	}
	
	cc.setupFormatCuenta()
	return cc
}

func (cc *CuentaContable) GetCuenta() string {
	return cc.entry.Text()
}

func (cc *CuentaContable) GetCuentaSinFormato() string {
	return strings.ReplaceAll(cc.entry.Text(), "-", "")
}

func (cc *CuentaContable) SetCuenta(cuenta string) {
	// Solo números, máximo 7 dígitos
	nums := ""
	for _, r := range cuenta {
		if unicode.IsDigit(r) {
			nums += string(r)
		}
	}
	if len(nums) > 7 {
		nums = nums[:7]
	}
	// Formatear dinámicamente: 0-00-00-00
	formatted := ""
	for i, r := range nums {
		if i == 1 || i == 3 || i == 5 {
			formatted += "-"
		}
		formatted += string(r)
	}
	cc.entry.SetText(formatted)
	cc.EstiloEntrada()
}

func (cc *CuentaContable) setupFormatCuenta() {
	var handleID glib.SignalHandle
	handleID = cc.entry.Connect("changed", func() {
		// 1. Bloquea para evitar recursión
		cc.entry.HandlerBlock(handleID)

		// 2. Filtrar solo números y limitar a 7 dígitos
		rawText := cc.entry.Text()
		nums := ""
		for _, r := range rawText {
			if unicode.IsDigit(r) {
				nums += string(r)
			}
		}
		if len(nums) > 7 {
			nums = nums[:7]
		}

		// 3. Formatear dinámicamente: 0-00-00-00
		// Guiones en posiciones de índice: 1, 3, 5 (después de 1er, 3er y 5to número)
		formatted := ""
		for i, r := range nums {
			if i == 1 || i == 3 || i == 5 {
				formatted += "-"
			}
			formatted += string(r)
		}

		// 4. Actualizar Entry y posición del cursor
		cc.entry.SetText(formatted)
		cc.entry.SetPosition(int(len(formatted)))

		// 5. Estilo visual
		cc.EstiloEntrada()

		// 6. Desbloquear
		cc.entry.HandlerUnblock(handleID)
	})
}

func (cc *CuentaContable) EstiloEntrada() {
	cc.entry.RemoveCSSClass("Valido")
	cc.entry.RemoveCSSClass("Invalido")
	cc.entry.RemoveCSSClass("Neutro")

	texto := cc.GetCuentaSinFormato()
	
	if len(texto) == 0 {
		cc.entry.AddCSSClass("Neutro")
		return
	}

	// Validación básica: una cuenta completa debe tener 7 dígitos
	if len(texto) == 7 {
		cc.entry.AddCSSClass("Valido")
	} else {
		cc.entry.AddCSSClass("Invalido")
	}
}