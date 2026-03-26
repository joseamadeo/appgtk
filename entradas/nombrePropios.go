package entradas

import (
	"strings"
	"unicode"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type NombrePropio struct {
	*gtk.Box
	entry *gtk.Entry
	mayusculas bool
	MaxLength int
	// Guardamos el ID del handler para poder bloquearlo/desbloquearlo globalmente
	changeHandlerID glib.SignalHandle 
}

func NewNombrePropio(placeholder string, soloMayuscula bool, largomaximo int) *NombrePropio {
	box := gtk.NewBox(gtk.OrientationHorizontal, 0)
	box.AddCSSClass("linked")

	entry := gtk.NewEntry()
	entry.SetPlaceholderText(placeholder)
	entry.SetHExpand(true)
	//Aplicamos el limite de caracteres a widget de GTK
	if largomaximo > 0 {
		entry.SetMaxLength(largomaximo)
	}
	box.Append(entry)
	n := &NombrePropio{
		Box:        box,
		entry:      entry,
		mayusculas: soloMayuscula,
		MaxLength:  largomaximo,
	}
	n.setupLogic()
	return n
}

func (n *NombrePropio) SoloMayusculas(soloMayuscula bool) {
	n.mayusculas = soloMayuscula
	// Llama a la nueva función auxiliar para re-procesar el texto existente
	n.processCurrentText() 
}

// Función auxiliar que contiene la lógica central de procesamiento
func (n *NombrePropio) processCurrentText() {
	// 1. Obtener el texto actual
	textoOriginal := n.entry.Text()
  var textoProcesado string
  // 2. Filtra para que solo acepte letras y espacios
  filtrado := ""
  for _, r := range textoOriginal {
    if unicode.IsLetter(r) || unicode.IsSpace(r) {
      filtrado += string(r)
    }
  }
  // 3. Si Mayusculas es true, convertir a mayúsculas
  if n.mayusculas {
  	textoProcesado = strings.ToUpper(filtrado)
  } else {
    // Formato Nombre Propio: "juan perez" -> "Juan Perez"
		// Usamos strings.ToLower primero para asegurar que el resto sea minúscula
		palabras := strings.Fields(strings.ToLower(filtrado))
		for i, p := range palabras {
			runas := []rune(p)
			if len(runas) > 0 {
				runas[0] = unicode.ToUpper(runas[0])
				palabras[i] = string(runas)
			}
		}
		textoProcesado = strings.Join(palabras, " ")
		
		// Preservar espacio final si el usuario está escribiendo uno
		if strings.HasSuffix(filtrado, " ") {
			textoProcesado += " "
		}
  }
  // 4. Actualizar solo si hubo cambios
  if textoOriginal != textoProcesado {
  	// Bloquear temporalmente el handler para evitar recursión infinita
    n.entry.HandlerBlock(n.changeHandlerID)
		pos := n.entry.Position()
    n.entry.SetText(textoProcesado)
    n.entry.SetPosition(pos)
    // Desbloquear el handler
    n.entry.HandlerUnblock(n.changeHandlerID)
  }
  // 5. Estilo de validación
  n.validarEstilo()
}

func (n *NombrePropio) setupLogic() {
	n.changeHandlerID = n.entry.Connect("changed", func () {
		glib.IdleAdd(func() bool {
			n.processCurrentText()
			return false
		})
	})
}

func (n *NombrePropio) validarEstilo() {
	n.entry.RemoveCSSClass("Valido")
	n.entry.RemoveCSSClass("Invalido")
	n.entry.RemoveCSSClass("Neutro")
	texto := strings.TrimSpace(n.entry.Text())
	if len(texto) == 0 {
		n.entry.AddCSSClass("Neutro")
	} else if len(texto) >= 3 {
		n.entry.AddCSSClass("Valido")
	} else {
		n.entry.AddCSSClass("Invalido")
	}
}

// Métodos para obtener y establecer el texto
func (n *NombrePropio) GetTexto() string {
	return strings.TrimSpace(n.entry.Text())
}

func (n *NombrePropio) SetTexto(texto string) {
	n.entry.SetText(texto)
}