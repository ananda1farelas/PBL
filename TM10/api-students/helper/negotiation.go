package helper

import (
	"mime"
	"strconv"
	"strings"
)

const (
	StudentFormatJSON = "json"
	StudentFormatCSV  = "csv"
)

// NegotiateStudentFormat memilih format respons berdasarkan header Accept.
// Header kosong dan */* memakai JSON sebagai format bawaan.
func NegotiateStudentFormat(accept string) (string, error) {
	accept = strings.TrimSpace(accept)
	if accept == "" || accept == "*/*" || accept == "*" {
		return StudentFormatJSON, nil
	}

	bestFormat := ""
	bestQuality := -1.0
	for _, item := range strings.Split(accept, ",") {
		mediaType, params, err := mime.ParseMediaType(strings.TrimSpace(item))
		if err != nil {
			continue
		}

		quality := 1.0
		if rawQuality, ok := params["q"]; ok {
			quality, err = strconv.ParseFloat(rawQuality, 64)
			if err != nil || quality <= 0 || quality > 1 {
				continue
			}
		}

		format := ""
		switch strings.ToLower(mediaType) {
		case "application/json", "application/*":
			format = StudentFormatJSON
		case "text/csv", "text/*":
			format = StudentFormatCSV
		case "*/*":
			format = StudentFormatJSON
		}

		if format != "" && quality > bestQuality {
			bestFormat = format
			bestQuality = quality
		}
	}

	if bestFormat == "" {
		return "", NotAcceptable("format yang diminta tidak tersedia, pilih application/json atau text/csv")
	}
	return bestFormat, nil
}
