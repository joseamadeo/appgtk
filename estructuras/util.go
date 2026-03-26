package estructuras

func CompararStrings(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// para las combos
type CuentaInfo struct {
	NumeroCuenta string
	NombreCuenta string
}

