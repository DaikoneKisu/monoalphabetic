package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"strings"
)

func main() {
	alphabetNext, _ := parseAlphabetFromJSON("alphabet-next.json")

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

	bruteForceDecodeWithStopWord("WJRAKGT", "CIPHERS", alphabetNext)
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

func bruteForceDecodeWithStopWord(toDecode string, stopWord string, alphabetNext map[rune]string) string {
	baseAlphabet, _ := parseAlphabetFromJSON("base-alphabet.json")

	alphabetPrevious, _ := parseAlphabetFromJSON("alphabet-previous.json")

	toDecode = strings.ToLower(toDecode)
	stopWord = strings.ToLower(stopWord)

	changingAlphabet := map[rune]string{}

	maps.Copy(changingAlphabet, baseAlphabet)

	var iterations uint64 = 0

	for {
		iterations += 1
		fmt.Printf("Iteration: %d\n", iterations)
		fmt.Println()

		decoded, _ := decode(toDecode, changingAlphabet)

		fmt.Println("Alphabet: " + printAlphabetAsJSONString(changingAlphabet))
		fmt.Println("String decoded: " + decoded)

		if strings.Contains(decoded, stopWord) {
			return decoded
		}

		changingAlphabet['a'] = alphabetNext[rune(changingAlphabet['a'][0])]

		if changingAlphabet['a'] == baseAlphabet['a'] {
			for k := range changingAlphabet {
				previous := alphabetPrevious[rune(changingAlphabet[k][0])]
				previousInBaseAlphabet := alphabetPrevious[rune(baseAlphabet[k][0])]
				if k != 'a' {
					fmt.Println(string(k), previous)
					fmt.Println(previousInBaseAlphabet)
					if previous == previousInBaseAlphabet {
						if k == 'b' {
							fmt.Println("b:" + alphabetNext[rune(changingAlphabet[k][0])])
						}
						changingAlphabet[k] = alphabetNext[rune(changingAlphabet[k][0])]
					}
					previous = changingAlphabet[k]
					previousInBaseAlphabet = baseAlphabet[k]
				}
			}
		}
	}
}

func frequencyAnalysisDecode(toDecode string) string {
	return toDecode
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

func printAlphabetAsJSONString(alphabet map[rune]string) string {
	stringMap := make(map[string]string)
	for k, v := range alphabet {
		stringMap[string(k)] = v
	}

	jsonString, _ := json.Marshal(stringMap)

	return string(jsonString)
}
