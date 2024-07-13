package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	encodeAlphabet, err := parseAlphabetFromJSON("alphabet.json")

	if err != nil {
		fmt.Println(err)
		return
	}

	decodeAlphabet := make(map[rune]string)
	for k, v := range encodeAlphabet {
		decodeAlphabet[rune(v[0])] = string(k)
	}

	toEncode := "Hello World"
	encoded, err := encode(toEncode, encodeAlphabet)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(encoded)

	decoded, err := decode(encoded, decodeAlphabet)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(decoded)
}

func encode(toEncode string, alphabet map[rune]string) (string, error) {
	toEncode = strings.ToLower(toEncode)

	var encodedBuilder strings.Builder

	for _, char := range toEncode {
		_, err := encodedBuilder.WriteString(alphabet[char])

		if err != nil {
			return "", errors.New("Error encoding string " + err.Error())
		}
	}

	return encodedBuilder.String(), nil
}

func decode(toDecode string, alphabet map[rune]string) (string, error) {
	toDecode = strings.ToLower(toDecode)

	var decodedBuilder strings.Builder

	for _, char := range toDecode {
		_, err := decodedBuilder.WriteString(alphabet[char])

		if err != nil {
			return "", errors.New("Error decoding string " + err.Error())
		}
	}

	return decodedBuilder.String(), nil
}

func parseAlphabetFromJSON(filepath string) (map[rune]string, error) {
	jsonFile, err := os.Open(filepath)

	if err != nil {
		return nil, errors.New("Error opening file " + err.Error())
	}

	defer jsonFile.Close()

	jsonByteValues, _ := io.ReadAll(jsonFile)

	var uncastedAlphabet interface{}

	json.Unmarshal(jsonByteValues, &uncastedAlphabet)

	alphabet := make(map[rune]string)
	for k, v := range uncastedAlphabet.(map[string]interface{}) {
		switch vv := v.(type) {
		case string:
			alphabet[rune(k[0])] = string(vv)
		default:
			return nil, errors.New("error parsing JSON file: invalid value type")
		}
	}

	return alphabet, nil
}
