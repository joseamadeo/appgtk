package pantalla

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/jung-kurt/gofpdf"
)

// Definimos una estructura para los totales
type Saldo struct {
	Ingreso float64
	Egreso  float64
}

func InformePDF(w *gtk.ApplicationWindow) {
	dir, err := os.Getwd()
	nombreArchivo := "informe_movimientos.pdf"
	
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	GenerarEstructuraPDF(pdf)
	err = pdf.OutputFileAndClose(nombreArchivo)
	cajaVert := gtk.NewBox(gtk.OrientationVertical, 10)
	cajaVert.SetMarginTop(20)

	if err != nil {
		fmt.Printf("Error al generar el PDF: %v\n", err)
		return
	} else {
		lbl := gtk.NewLabel(fmt.Sprintf("PDF generado exitosamente: %s", nombreArchivo))
		btnAbrir := gtk.NewButtonWithLabel("Abrir PDF")
		btnAbrir.ConnectClicked(func() {
			// Aquí puedes agregar la lógica para abrir el PDF generado
			fmt.Printf("Abriendo el PDF: %s\n", nombreArchivo)
			abrirArchivo(nombreArchivo)
		})

		area := gtk.NewDrawingArea()
		area.SetVExpand(true)
		area.SetHExpand(true)

		// --- ÁREA DE DIBUJO ---
		area.SetDrawFunc(func(_ *gtk.DrawingArea, cr *cairo.Context, width, height int) {
			if nombreArchivo == "" {
				return
			}
			// Limpiar fondo (opcional pero recomendado)
			cr.SetSourceRGB(0.9, 0.9, 0.9) // Fondo gris claro
			cr.Paint()
			// Aplicar zoom
			cr.Save()
			cr.Scale(zoomNivel, zoomNivel)
			// Renderizar directamente usando la variable global actualizada
			success := renderPDFPagina(filepath.Join(dir, nombreArchivo), actualPagina, cr)
			cr.Restore() // Restaurar el estado original después de renderizar
			fmt.Printf("Renderizando PDF: %s, Página: %d\n", filepath.Join(dir, nombreArchivo), actualPagina)
			fmt.Printf("Renderizado exitoso: %v\n", success)
			if !success {
				// Mostrar un mensaje de error en el área de dibujo
				cr.SetSourceRGB(1, 0, 0) // Rojo para el error
				cr.SelectFontFace("Sans", cairo.FontSlantNormal, cairo.FontWeightBold)
				cr.SetFontSize(14)
				cr.MoveTo(10, 30)
				cr.ShowText("Error al cargar el PDF o renderizar la página.")
				return
			}
		})
		
		cajaVert.Append(lbl)
		cajaVert.Append(btnAbrir)
		cajaVert.Append(area)
	}
	w.SetChild(cajaVert)
}

// Función auxiliar para abrir el archivo según el Sistema Operativo
func abrirArchivo(archivo string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", archivo}
	case "darwin": // macOS
		cmd = "open"
		args = []string{archivo}
	default: // Linux y otros
		cmd = "xdg-open"
		args = []string{archivo}
	}
	exec.Command(cmd, args...).Start()
}

func GenerarEstructuraPDF(pdf *gofpdf.Fpdf) {
	// Estructura interna para los totales
	type Totales struct {
		Nombre  string
		Egreso  float64
		Ingreso float64
	}

	reporte := make(map[string]*Totales)
	var listaCuentas []string

	// --- PASO A: Leer Cuentas.txt (Base del reporte) ---
	fCuentas, _ := os.Open("Cuentas.txt")
	scannerC := bufio.NewScanner(fCuentas)
	scannerC.Scan() // Saltar cabecera
	for scannerC.Scan() {
		p := strings.Split(scannerC.Text(), "|")
		if len(p) < 2 { continue }
		cuenta := strings.TrimSpace(p[0])
		nombre := strings.TrimSpace(p[1])
		
		reporte[cuenta] = &Totales{Nombre: nombre, Egreso: 0, Ingreso: 0}
		listaCuentas = append(listaCuentas, cuenta)
	}
	fCuentas.Close()

	// --- PASO B: Ordenar Cuentas ---
	sort.Strings(listaCuentas)

	// --- PASO C: Sumar movimientos.txt ---
	fMovs, _ := os.Open("movimientos.txt")
	scannerM := bufio.NewScanner(fMovs)
	scannerM.Scan() // Saltar cabecera
	for scannerM.Scan() {
		p := strings.Split(scannerM.Text(), "|")
		if len(p) < 5 { continue }
		
		cuenta := strings.TrimSpace(p[0])
		tipo := strings.TrimSpace(p[4]) // "I" o "E"
		
		// Limpiar formato de número (ej: 343.445,45 -> 343445.45)
		sVal := strings.ReplaceAll(p[3], ".", "")
		sVal = strings.ReplaceAll(sVal, ",", ".")
		valor, _ := strconv.ParseFloat(sVal, 64)

		if datos, existe := reporte[cuenta]; existe {
				if tipo == "I" {
						datos.Ingreso += valor
				} else if tipo == "E" {
						datos.Egreso += valor
				}
		}
	}
	fMovs.Close()

	// 1. Primero, ordenamos la lista de cuentas de mayor a menor longitud (01-01-01-01 -> 01-00-00-00)
	// Esto asegura que las "hijas" sumen a las "madres" antes de que las madres sumen a las "abuelas"
	sort.Slice(listaCuentas, func(i, j int) bool {
			return len(listaCuentas[i]) > len(listaCuentas[j])
	})
		
	// 2. Propagación jerárquica automática
	for _, cuenta := range listaCuentas {
			partes := strings.Split(cuenta, "-")
			
			// Subimos niveles: de 1-01-01-01 a 1-01-01-00, luego 1-01-00-00, etc.
			for i := len(partes) - 1; i >= 1; i-- {
					idPadre := generarPadre(partes, i)
					
					// Si el padre existe en el TXT y no es la misma cuenta, le pasamos los saldos
					if padre, existe := reporte[idPadre]; existe && idPadre != cuenta {
							padre.Ingreso += reporte[cuenta].Ingreso
							padre.Egreso += reporte[cuenta].Egreso
					}
			}
	}

	// 3. Regla especial: Sumar raíces (1, 2, 3) a la cuenta de Ejercicios Actuales (0-02)
	if destino, ok := reporte["0-02-00-00"]; ok {
			raices := []string{"1-00-00-00", "2-00-00-00", "3-00-00-00"}
			for _, r := range raices {
					if datosRaiz, existe := reporte[r]; existe {
							destino.Ingreso += datosRaiz.Ingreso
							destino.Egreso += datosRaiz.Egreso
					}
			}
	}

	// 4. Volvemos a ordenar la lista alfabéticamente para que el PDF salga prolijo (0, 1, 2, 3...)
	sort.Strings(listaCuentas)

// --- PASO FINAL: Consolidación en la cuenta 0-00-00-00 ---

// 1. Obtenemos las referencias de las 3 cuentas involucradas
c000, ex00 := reporte["0-00-00-00"]
c001, ex01 := reporte["0-01-00-00"]
c002, ex02 := reporte["0-02-00-00"]

if ex00 && ex01 && ex02 {
    // 2. Sumamos ingresos y egresos de las dos cuentas de ejercicio (0-01 y 0-02)
    totalIngresos := c001.Ingreso + c002.Ingreso
    totalEgresos  := c001.Egreso + c002.Egreso
    
    // 3. Calculamos el Resultado Neto
    resultado := totalIngresos - totalEgresos
    
    // 4. Aplicamos la lógica solicitada para la cuenta 0-00-00-00
    if resultado >= 0 {
        c000.Ingreso = resultado
        c000.Egreso  = 0
    } else {
        // Usamos math.Abs o simplemente negamos el valor para que sea positivo en Egresos
        c000.Ingreso = 0
        c000.Egreso  = -resultado 
    }
	}


	// --- PASO D: Dibujar en el PDF ---
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "REPORTE DE SALDOS POR CUENTA", "0", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Cabecera Tabla
	pdf.SetFillColor(230, 230, 230)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(30, 8, "Cuenta", "1", 0, "C", true, 0, "")
	pdf.CellFormat(70, 8, "Nombre", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 8, "Egreso", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 8, "Ingreso", "1", 1, "C", true, 0, "")

	// Filas
	pdf.SetFont("Arial", "", 9)
	for _, c := range listaCuentas {
		v := reporte[c]
		pdf.CellFormat(30, 7, c, "1", 0, "L", false, 0, "")
		pdf.CellFormat(70, 7, v.Nombre, "1", 0, "L", false, 0, "")
		
		// Formatear con separador de miles si lo deseas, aquí usamos simple:
		sEgreso := formatearMoneda(v.Egreso)
		sIngreso := formatearMoneda(v.Ingreso)
		
		pdf.CellFormat(45, 7, sEgreso, "1", 0, "R", false, 0, "")
		pdf.CellFormat(45, 7, sIngreso, "1", 1, "R", false, 0, "")
	}
}

// generarPadre toma las partes y mantiene las primeras 'n', rellenando el resto con "00"
func generarPadre(partes []string, nivelesAMantener int) string {
	resultado := make([]string, len(partes))
	for i := 0; i < len(partes); i++ {
		if i < nivelesAMantener {
			resultado[i] = partes[i]
		} else {
			resultado[i] = "00" // O "0" según tu formato
		}
	}
	return strings.Join(resultado, "-")
}

func formatearMoneda(valor float64) string {
	// 1. Convertimos a string con 2 decimales
	s := fmt.Sprintf("%.2f", valor)
	partes := strings.Split(s, ".")
	entero := partes[0]
	decimal := partes[1]

	// 2. Insertamos los puntos de miles
	var resultado []string
	n := len(entero)
	for i, r := range entero {
		if i > 0 && (n-i)%3 == 0 && entero[i-1] != '-' {
			resultado = append(resultado, ".")
		}
		resultado = append(resultado, string(r))
	}

	// 3. Unimos con la coma decimal argentina
	return strings.Join(resultado, "") + "," + decimal
}
