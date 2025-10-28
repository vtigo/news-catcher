package rss

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

func UnmarshalXML(data []byte) (*RSSDocument, error) {
	// Precisamos:
	// 	1.Ler o XML
	// 	2.Validar a estrutura
	// 	3.Mapear para um struct

	// reader são os bytes contendo o XML
	reader := bytes.NewReader(data)

	// passamos os bytes para um decoder de XML
	decoder := xml.NewDecoder(reader)

	// contamos para o decoder como ele deve ler os bytes - queremos converter o encoding do input para UTF-8
	decoder.CharsetReader = spawnCharsetReader

	var feed *RSSDocument

	// executamos a leitura do XML, utilizamos o ponteiro para armazenar os dados na estrutura RSSDocument
	err := decoder.Decode(&feed)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

// spawnCharsetReader retorna um reader que converte o encoding do input para UTF-8.
func spawnCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "iso-8859-1", "latin1":
		buf := new(bytes.Buffer)
		chunk := make([]byte, 4096)
		for {
			n, err := input.Read(chunk)
			if n > 0 {
				for i := 0; i < n; i++ {
					if chunk[i] < 128 {
						// o byte é um ASCII, então podemos escrever direto no buffer
						buf.WriteByte(chunk[i])
					} else {
						// caracteres especiais, precisamos converter para UTF-8
						// aqui a gente converte para uma rune(unicode) - retrocompatível com ISO ISO-8859-1,
						buf.WriteRune(rune(chunk[i]))
					}
				}
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				return nil, err
			}
		}
		return buf, nil

	default:
		return input, nil
	}
}
