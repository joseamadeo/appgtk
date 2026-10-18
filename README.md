# AppGtk 📒

Aplicación de escritorio para **contabilidad** escrita en **Go** con interfaz gráfica **GTK 4**, que funciona con archivos de texto en disco (TXT) y genera reportes en **PDF**.

Incluye además una colección de pantallas de ejemplo para explorar los widgets de GTK 4.

> Stack y dependencias detalladas en [STACK.md](STACK.md).

---

## ✨ Características

### Módulo contable

| Pantalla | Descripción |
|---|---|
| **Cuentas** | Plan de cuentas jerárquico (`00-00-00-00`) con cuentas imputables y no imputables, persistido en `cuentas.txt` |
| **Movimiento Contable** | Alta de movimientos de ingreso/egreso con cuenta, fecha e importe, persistidos en `movimientos.txt` |
| **Informe PDF** | Generación de informe de movimientos en PDF (gofpdf) |
| **Vista PDF** | Visor de PDF con paginación e impresión (cairo) |

### Widgets personalizados (`entradas/`)

- **`Fecha`** — entrada de fecha con calendario emergente, autoformato `DD/MM/AAAA` (también acepta los 8 dígitos seguidos: `DDMMAAAA`), validación de rango Mín/Máx, menú contextual (copiar/pegar formatos) y estilos según validación (neutro / válido / inválido).
- **`Numerico`** — entrada numérica con formato de miles y decimales (`#.###,00`), separador decimal con el teclado numérico, normalización inteligente del portapapeles (formato argentino `1.234,56`, formato internacional `1,234.56` y negativos contables `(1.234,56)`), y menú contextual.
- **`MiTabla[T]`** — tabla genérica para mostrar colecciones de cualquier tipo.
- Otros: validación de **CUIT**, cuentas contables, nombres propios.

### Ejemplos de GTK 4 (`pantalla/`)

Grilla (Grid), Panel (Paned), Caja (Box), ListBox, Tabs, Contador, Diálogos de mensaje, Cuadro de texto, Estilo CSS, Visor PDF y Validación de widgets.

> 🔑 Al iniciar aparece un **login de demostración** (validación simulada): usuario `admin`, contraseña `1234`.

---

## 🛠 Requisitos

**Para compilar:**

- [Go 1.25+](https://go.dev/dl/)
- GTK 4 y GLib de desarrollo con `pkg-config` (necesario por cgo/gotk4).
  En Windows la vía habitual es [MSYS2](https://www.msys2.org/):

  ```bash
  pacman -S mingw-w64-x86_64-gtk4 mingw-w64-x86_64-pkgconf
  ```

**Para ejecutar:** no hace falta instalar GTK. El runtime está empaquetado en el repositorio:

- `lib/` → DLLs de GTK 4
- `share/` → temas, iconos y datos de GTK

---

## ▶️ Compilar y ejecutar

```bash
# Desarrollo (con consola)
go run .

# Ejecutable de distribución (sin ventana de consola)
go build -o appgtk.exe -ldflags="-H=windowsgui" main.go
```

El ejecutable debe correr **desde la raíz del proyecto** (o junto a `lib/`, `share/` y los archivos `.txt`), porque busca el runtime y los datos relativos al directorio de trabajo actual.

### Tests

```bash
go test ./test/...
```

---

## 📁 Estructura del proyecto

```
appgtk/
├── main.go               # Punto de entrada: app GTK, menú principal, CSS embebido
├── style.css             # Hoja de estilos (embebida con go:embed)
├── go.mod / go.sum
├── entradas/             # Widgets de entrada reutilizables
│   ├── fecha.go          #   Campo de fecha con calendario
│   ├── numero.go         #   Campo numérico con formato
│   ├── cuit.go           #   Validación de CUIT
│   ├── cuentasContables.go
│   ├── nombrePropios.go
│   └── tabla.go          #   Tabla genérica MiTabla[T]
├── pantalla/             # Pantallas y ventanas
│   ├── login.go          #   Login (demo)
│   ├── movimientoContable.go
│   ├── cuentas.go
│   ├── informePDF.go     #   Generación del PDF
│   ├── PDFVista.go       #   Visor/impresión de PDF
│   └── ...               #   Ejemplos de widgets GTK
├── estructuras/          # Tipos compartidos y utilidades
├── test/                 # Tests unitarios
├── cuentas.txt           # Datos: plan de cuentas
├── movimientos.txt       # Datos: movimientos contables
├── lib/  share/          # Runtime de GTK 4 para Windows
└── InstaladorAppGtk/     # Instalador (Inno Setup)
```

---

## 💾 Formato de los archivos de datos

**`cuentas.txt`** — separador `|`:

```
Cuenta|Nombre|Imputable
1-01-01-01|Caja imputable|Sí
```

**`movimientos.txt`** — separador `|`:

```
Cuenta|Nombre|Fecha|Importe|Tipo
1-01-01-01|Caja|2026-03-10|33,77|I
```

- `Tipo`: `I` = ingreso, `E` = egreso
- `Fecha`: `AAAA-MM-DD` (los registros con fecha vacía se omiten al cargar)

> ⚠️ Estos archivos son datos de ejemplo. La aplicación trabaja sobre los que se encuentren en su directorio de trabajo.

---

## 📦 Instalador

En `InstaladorAppGtk/` se encuentra el script de [Inno Setup](https://jrsoftware.org/isinfo.php) (`instalador.iss`) y el instalador compilado (`InsAppgtk.exe`).

---

## 📄 Licencia

Proyecto personal. Consultar al autor respecto al uso y la distribución.
