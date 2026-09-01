package engine

import (
	"strings"
)

// FormatCompactTable, dar ekranlar veya mobil terminal görünümleri için 
// uzun metinleri ve tabloları sıkıştırarak taşmayı önleyen bir formatlayıcıdır.
func FormatCompactTable(header string, rows []string, maxLineWidth int) string {
	if maxLineWidth <= 0 {
		maxLineWidth = 60 // Mobil veya dar terminal için varsayılan sınır
	}

	var sb strings.Builder
	sb.WriteString("=== " + header + " ===\n")

	for _, row := range rows {
		if len(row) > maxLineWidth {
			// Satır çok uzunsa dar ekran için kes ve üç nokta ekle
			sb.WriteString(row[:maxLineWidth-3] + "...\n")
		} else {
			sb.WriteString(row + "\n")
		}
	}

	sb.WriteString(strings.Repeat("-", maxLineWidth) + "\n")
	return sb.String()
}
