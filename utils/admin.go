package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// commaf is a function to format number with thousand separator
// and two decimal points
func Commaf(j interface{}) string {
	v, _ := strconv.ParseFloat(fmt.Sprint(j), 64)
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}
	s := fmt.Sprintf("%.2f", v)

	comma := []byte{','}

	parts := strings.Split(s, ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}

func IsLocal(Addr string) bool {
	if strings.Contains(Addr, ":") && !strings.Contains(Addr, ".") {
		Addr = strings.TrimPrefix(Addr, "[")
		if strings.HasPrefix(Addr, "::") || strings.HasPrefix(Addr, "fc") || strings.HasPrefix(Addr, "fd") {
			return true
		}
	}
	p := strings.Split(strings.Split(Addr, ":")[0], ".")
	if len(p) != 4 {
		return false
	}
	_, err := strconv.ParseInt(p[2], 10, 64)
	if err != nil {
		return false
	}
	_, err = strconv.ParseInt(p[3], 10, 64)
	if err != nil {
		return false
	}
	v1, err := strconv.ParseInt(p[0], 10, 64)
	if err != nil {
		return false
	}
	v2, err := strconv.ParseInt(p[1], 10, 64)
	if err != nil {
		return false
	}
	if v1 == 10 {
		return true
	}
	if v1 == 172 {
		if v2 >= 16 && v2 <= 31 {
			return true
		}
	}
	if v1 == 192 && v2 == 168 {
		return true
	}
	if v1 == 127 {
		return true
	}
	return false
}
