package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/vtigo/news-catcher/fetcher"
	"github.com/vtigo/news-catcher/rss"
)

// TODO: pegar o endpoint do arquivo de config
// TODO: pegar múltiplas fontes

func main() {
	endpoint := "https://revistacult.uol.com.br/feed/"

	fetcher := fetcher.NewClient(
		fetcher.WithTimeout(15 * time.Second),
	)

	// Um feed RSS é um documento XML, temos o xml bruto (bytes) no corpo da resposta
	responseData, err := fetcher.Fetch(endpoint)
	if err != nil {
		log.Fatalln("failed to fetch feed: ", err)
	}

	// Precisamos:
	// 	1.Ler o XML
	// 	2.Validar a estrutura
	// 	3.Mapear para um struct

	// reader são os bytes contendo o XML
	reader := bytes.NewReader(responseData)

	// passamos os bytes para um decoder de XML
	decoder := xml.NewDecoder(reader)

	// contamos para o decoder como ele deve ler os bytes - queremos converter o encoding do input para UTF-8
	decoder.CharsetReader = spawnCharsetReader

	var feed *rss.RSSDocument

	// executamos a leitura do XML, utilizamos o ponteiro para armazenar os dados na estrutura RSSDocument
	err = decoder.Decode(&feed)
	if err != nil {
		log.Fatalln("failed to unmarshall the RSS document: ", err)
	}

	// Agora que temos o feed em memória, podemos decidir o que fazer com os dados

	// Convertemos o feed em um slice de bytes que representa uma string JSON formatada
	jsonData, err := json.MarshalIndent(feed, "", " ")
	if err != nil {
		log.Fatalln("failed to marshal feed data into json: ", err)
	}

	// Gravamos os bytes em um arquivo json
	err = os.WriteFile("feed.json", jsonData, 0644)
	if err != nil {
		log.Fatalln("failed to write json data: ", err)
	}

}

// CharsetReader retorna um reader que converte o encoding do input para UTF-8.
func spawnCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "iso-8859-1", "latin1":
		buf := new(bytes.Buffer)
		for {
			b := make([]byte, 1)
			_, err := input.Read(b)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			if b[0] < 128 {
				// o byte é um ASCII, então podemos escrever direto no buffer
				buf.WriteByte(b[0])
			} else {
				// caracteres especiais, precisamos converter para UTF-8
				// aqui a gente converte para uma rune(unicode) - retrocompatível com ISO ISO-8859-1,
				// depois convertemos para string porque o go faz a conversão rune -> string(utf8)
				buf.WriteString(string(rune(b[0])))
			}
		}
		return buf, nil
	default:
		return input, nil
	}
}
