package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"sort"
	"strings"
)

func main() {
	// bruteForceDecodeWithStopWord("BOLC AODO ESTCS", "como")
	// frequencyAnalysisDecode("BOLC AODO ESTCS")
	// frequencyAnalysisDecode("WJRAKGT YJMK OAJT, NAJWA HTK XPGK OAQZ PZK WJRAKG QYRAQVKO QGK MZPNZ QT RPYLQYRAQVKOJW WJRAKGT. OAKTK WQZ VK JZWGKCJVYL CJSSJWHYO OP CKWJRAKG, VKWQHTK PS OAKJG GKTJTOQZWK OP YKOOKG SGKIHKZWL QZQYLTJT")
	// frequencyAnalysisDecode(`KPTY HFXY SNQWY EPZNYHYG TPXXYG LWI'Z. VQHILIE. KHFLZKFX SLZNQFZ. NYHA NYHA SPW SPW WPOLIE. GPOW OQF, NYPJYI, ZHYY ZNYO'HY ZSQ FW THYPZFHY GHO P QFH VPDY. ZNY ZNYV KPTY EHYYI ELJY ZNYV AHQFENZ KHFLZKFX QK GQVLILQI KQH PXX ELJY XPIG YPHZN AYPWZ NPG, WQ VPXY L QK FW FIZQ XLENZ QJYH. XLJLIE KYVPXY. XLENZ LI FIZQ KHFLZ GLJLGYG GPOW. SYHY ZNYHY LZ WMLHLZ XLKY, NPG FIZQ NPZN KQHV SPW AYPHLIE. XLDYIYWW, AYPWZ. XLENZW NLV GYYM LW WYZ VPDY SPZYHW XPIG VPO FW GLJLGYG TPI'Z GPO SYHY SNPXYW VPXY ZNPZ JQLG. YJYHO XLJLIE, WYPW PXWQ WFAGFY OQF'XX KYVPXY EPZNYHYG SNQWY THYYMLIE WNY'G SPW PXWQ KYVPXY. GQYWI'Z WMLHLZ KLWN PLH GLJLGYG. AYELIILIE GPOW WYTQIG NPZN TPZZXY SPW VYPZ WYTQIG XLENZW EQG JYHO ZNYLH EQG HYMXYILWN KPTY SQI'Z EHYYI ZSQ ZNPZ, WNY'G VQJYG XLJLIE NPJY, WYZ WNPXX WPOLIE WYTQIG AXYWWYG OQF WYZ TPZZXY. PXWQ WLUZN NLW GQI'Z. AYNQXG. JQLG GPO. VPO XPIG DLIG VPO VPXY KQHV VQJYZN WQ ZQ. AXYWWYG. PKZYH ZNPZ WYTQIG THYPZFHY OQF. ZSQ EQG PIG YPHZN KHFLZKFX OYPHW THYPZYG QMYI XYWWYH VPI QK GLJLGYG KQH. KPTY WYYG XLENZW ZQ AXYWWYG SQI'Z, GLJLGYG ZNYO'HY AY JQLG XYZ ZNYHY LVPEY KHQV VPXY SNYHYLI ZNYO'HY WMLHLZ NY OQF'HY SPZYHW WLEIW AYPHLIE AY KHFLZ.`)
}

func encode(toEncode string, alphabet map[rune]string) (string, error) {
	toEncode = strings.ToLower(toEncode)

	var encodedBuilder strings.Builder

	for _, char := range toEncode {
		if char >= 'a' && char <= 'z' {
			_, err := encodedBuilder.WriteString(alphabet[char])

			if err != nil {
				return "", errors.New("Error encoding string " + err.Error())
			}
		} else {
			encodedBuilder.WriteRune(char)
		}
	}

	return encodedBuilder.String(), nil
}

func decode(toDecode string, alphabet map[rune]string) (string, error) {
	toDecode = strings.ToLower(toDecode)

	var decodedBuilder strings.Builder

	for _, char := range toDecode {
		if char >= 'a' && char <= 'z' {
			_, err := decodedBuilder.WriteString(alphabet[char])

			if err != nil {
				return "", errors.New("Error decoding string " + err.Error())
			}
		} else {
			decodedBuilder.WriteString(string(char))
		}
	}

	fmt.Println(decodedBuilder.String())

	return decodedBuilder.String(), nil
}

func bruteForceDecodeWithStopWord(toDecode string, stopWord string) string {
	orderedEnglishAlphabet := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

	baseAlphabet, _ := parseAlphabetFromJSON("base-alphabet.json")

	alphabetPrevious, _ := parseAlphabetFromJSON("alphabet-previous.json")
	alphabetNext, _ := parseAlphabetFromJSON("alphabet-next.json")

	toDecode = strings.ToLower(toDecode)
	stopWord = strings.ToLower(stopWord)

	changingAlphabet := map[rune]string{}

	maps.Copy(changingAlphabet, baseAlphabet)

	var iterations uint64 = 0

	for {
		iterations += 1

		decoded, _ := decode(toDecode, changingAlphabet)

		fmt.Printf("Iteration: %d\n", iterations)
		fmt.Println("Alphabet: " + printAlphabetAsJSONString(changingAlphabet))
		fmt.Println("String decoded: " + decoded)
		fmt.Println()

		if strings.Contains(decoded, stopWord) {
			return decoded
		}

		changingAlphabet['a'] = alphabetNext[rune(changingAlphabet['a'][0])]

		if changingAlphabet['a'] == baseAlphabet['a'] {
			for _, k := range orderedEnglishAlphabet {
				previous := changingAlphabet[rune(alphabetPrevious[k][0])]
				previousInBaseAlphabet := baseAlphabet[rune(alphabetPrevious[k][0])]
				if k != 'a' {
					if previous == previousInBaseAlphabet {
						changingAlphabet[k] = alphabetNext[rune(changingAlphabet[k][0])]
					} else {
						break
					}
				}
			}
		}
	}
}

func frequencyAnalysisDecode(toDecode string) string {
	var englishLetterFrequencies = map[rune]float32{
		'a': 8.2, 'b': 1.5, 'c': 2.8, 'd': 4.3,
		'e': 12.7, 'f': 2.2, 'g': 2.0, 'h': 6.1,
		'i': 7.0, 'j': 0.2, 'k': 0.8, 'l': 4.0,
		'm': 2.4, 'n': 6.7, 'o': 7.5, 'p': 1.9,
		'q': 0.1, 'r': 6.0, 's': 6.3, 't': 9.1,
		'u': 2.8, 'v': 1.0, 'w': 2.4, 'x': 0.2,
		'y': 2.0, 'z': 0.1,
	}

	var frequencyToChar = make(map[float32]rune)

	for char, freq := range englishLetterFrequencies {
		frequencyToChar[freq] = char
	}

	toDecode = strings.ToLower(toDecode)

	frequencies := make(map[rune]float32)
	totalChars := 0

	for _, char := range toDecode {
		if char >= 'a' && char <= 'z' {
			frequencies[char]++
			totalChars++
		}
	}

	for char := 'a'; char <= 'z'; char++ {
		if count, found := frequencies[char]; found {
			frequencies[char] = (count / float32(totalChars)) * 100.0
		} else {
			frequencies[char] = 0
		}
	}

	foundLetters := make(map[rune]bool)
	uniqueLetters := []rune{}
	for _, char := range toDecode {
		if char >= 'a' && char <= 'z' && !foundLetters[char] {
			uniqueLetters = append(uniqueLetters, char)
			foundLetters[char] = true
		}
	}

	sort.Slice(uniqueLetters, func(i, j int) bool {
		return frequencies[uniqueLetters[i]] > frequencies[uniqueLetters[j]]
	})

	decoded := strings.Builder{}
	guessedLetters := make(map[rune]rune)
	remainingLetters := make(map[rune]bool)

	for char := 'a'; char <= 'z'; char++ {
		guessedLetters[char] = char
		remainingLetters[char] = true
	}

	for _, char := range uniqueLetters {
		if char >= 'a' && char <= 'z' {
			minDiff := float32(100)
			var bestChar rune

			for englishChar, freq := range englishLetterFrequencies {
				if remainingLetters[englishChar] {
					diff := abs(frequencies[char] - freq)
					if diff < minDiff {
						minDiff = diff
						bestChar = englishChar
					}
				}
			}

			guessedLetters[char] = bestChar
			remainingLetters[bestChar] = false
		} else {
			decoded.WriteRune(char)
		}
	}

	for _, char := range toDecode {
		if char >= 'a' && char <= 'z' {
			decoded.WriteRune(guessedLetters[char])
		} else {
			decoded.WriteRune(char)
		}
	}

	fmt.Println("Original: " + toDecode)
	fmt.Println(" Decoded: " + decoded.String())

	return decoded.String()
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
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
