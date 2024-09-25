package cmd

import (
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const (
	ingredientsFile = "Alchemie/Zutaten.md"
	switchMode      = "#### Additive"
	alchBaseMode    = "Basis"
	alchAddMode     = "Addititiv"

	// MAYBE move to a constant file to be shared with magic?
	// start words
	Obsek = "OBSEK"
	Ort   = "ORT"
	In    = "IN"
	An    = "AN"
	// "real" form words"
	Scien = "SCIEN"
	Syma  = "SYMA"
	Aum   = "AUM"
	Ur    = "UR"
	Tel   = "TEL"
	Lep   = "LEP"
	Mutar = "MUTAR"
	// form manipulation
	Leko   = "LEKO"
	Konfar = "KONFAR"
	Koniu  = "KONIU"
	// colors
	Magi  = "MAGI"
	Prix  = "PRIX"
	Nox   = "NOX"
	Flam  = "FLAM"
	Waku  = "WAKU"
	Litax = "LITAX"
	Ebor  = "EBOR"
	Mani  = "MANI"
	Mort  = "MORT"
	Sicr  = "SICR"
	Fera  = "FERA"
	Flora = "FLORA"
	Humi  = "HUMI"
	Anima = "ANIMA"
	Tora  = "TORA"
	Hora  = "HORA"

	minAddsWithoutCrit = 4
)

var (
	// Maybe move them together with string const
	FormWords = []string{Obsek, Scien, Syma, Aum, Ur, Tel, Lep, Leko, Konfar,
		Koniu, Mutar}
	ColorWords = []string{Magi, Prix, Nox, Flam, Waku, Litax, Ebor, Mani, Mort,
		Sicr, Fera, Flora, Humi, Anima, Tora, Hora}
	StartWords     = []string{Obsek, Ort, In}
	MiddleWords    = []string{Scien, Syma, Aum, Ur, Tel, Lep, Mutar}
	EndWords       = []string{Leko, Konfar, Koniu}
	colorOpposites = map[string]string{
		Flam:  Waku,
		Litax: Ebor,
		Prix:  Nox,
		Waku:  Flam,
		Ebor:  Litax,
		Nox:   Prix,
	}

	alchData = readAlchData()

	// parse md file into 4 section (name, diff, cost, string for words)
	parseRex = regexp.MustCompile(
		`^\|\s*(\S.*?\S)\s*\|\s*(\d+)\s*\|\s*(\d+)\s*\|\s*(\S.*?\S)\s*\|`)
	// read next magic word
	readEntity = regexp.MustCompile(`^,?\s*([^\+\-,]+)([\+\-]*)`)

	// special rules
	alchConducMode bool

	// max search depth
	maxSearchDepth = 2
	maxDifficulty  = 26

	// sheet mode
	emptyRecipes = 20

	// verbose Mode
	verboseMode    bool
	alchSearchMode bool

	alchCmd = &cobra.Command{
		Use:   "alch",
		Short: "Werkzeuge zum Umgang mit Alchemie",
		Long:  "Stellt Kommandos für den Umgang mit Alchemie bereit.",
	}

	alchCalcCmd = &cobra.Command{
		Use:   "calc",
		Short: "Berechnet das Ergebnis einer Alchemistischen Prozedur.",
		Long: "Berechnet das Ergebnis einer Alchemisten Prozedur. " +
			"Die Substanzen werden als kommaseparierte Liste übergeben." +
			"Bei unbekannten Substanzen wird abgebrochen.",
		Run: printAlchCalc,
	}

	alchSearchCmd = &cobra.Command{
		Use:   "search",
		Short: "Versucht die besten Rezepte für den Effekt zu finden.",
		Long: "Versucht die besten Rezepte für den angegebenen Effekt zu " +
			"finden. Effekt als kommaseparierte Liste angeben!",
		Run: printSearchRes,
	}

	alchPrintDataCmd = &cobra.Command{
		Use:   "printData",
		Short: "JSON Dump der Alchemiedaten.",
		Long:  "JSON Dump der Alchemiedaten.",
		Run:   printAlchData,
	}

	alchSheetCmd = &cobra.Command{
		Use:   "sheet",
		Short: "Generiert die leere Alchemie-Zutaten-Tabelle.",
		Long:  "Generiert die leere Alchemie-Zutaten-Tabelle.",
		Run:   printAlchSheet,
	}
)

func init() {
	rootCmd.AddCommand(alchCmd)
	alchCmd.AddCommand(alchCalcCmd)
	alchCmd.AddCommand(alchSearchCmd)
	alchCmd.AddCommand(alchPrintDataCmd)
	alchCmd.AddCommand(alchSheetCmd)
	alchCmd.PersistentFlags().BoolVarP(&alchConducMode, "conduc", "o", false,
		"Ermöglicht den CONDUC Modus aus alchemistischen Wegen.")
	alchCmd.PersistentFlags().IntVarP(&maxSearchDepth, "maxDepth", "m", 2,
		"Maximale Suchtiefe für Rezepte")
	alchCmd.PersistentFlags().IntVarP(&maxDifficulty, "maxDifficulty", "d", 26,
		"maximale Schwierigkeit für Rezepte")
	alchCmd.PersistentFlags().IntVarP(&emptyRecipes, "emptyRecipes", "", 20,
		"Anzahl Einträge für das Alchemieblatt.")
	alchCmd.PersistentFlags().BoolVarP(&verboseMode, "verbose", "", false,
		"Zeigt zusätzliche Informationen an")
}

func printSearchRes(_ *cobra.Command, args []string) {
	for _, arg := range args {
		printRecipes(searchRecipes(splitOpts(arg)))
	}
}

func printAlchCalc(_ *cobra.Command, args []string) {
	alchSearchMode = true
	for _, arg := range args {
		deb(calculateResult((splitOpts(arg))))
	}
}

// #############################################################################
// #							Ingredients Struct								   #
// #############################################################################

type Ingredient struct {
	CWords     map[string]int `json:"Farbworte"`
	FWords     []string       `json:"Formwort"`
	Difficulty int            `json:"Schwierigkeit"`
	Price      int            `json:"Preis"`
}

// #############################################################################
// #							Recipe Struct								   #
// #############################################################################

type Recipe struct {
	items         string
	formWords     []string
	colorWords    []string
	strength      int
	difficulty    int
	criticality   int
	priceSumItems int
}

func (r *Recipe) String() string {
	invalid := r.difficulty == 0
	ret := "--------------\n"
	ret += "Rezept für    (" + r.items + ")"
	if invalid {
		ret += " FEHLSCHLAG!"
	}
	ret += "\n"
	ret += "--------------\n"
	ret += "Formwörter:    " + strings.Join(r.formWords, ", ") + "\n"
	ret += "Farbwörter:    " + strings.Join(r.colorWords, ", ") + "\n"
	ret += "Effektstärke:  " + strconv.Itoa(r.strength) + "\n"
	ret += "--------------\n"
	ret += "Schwierigkeit: " + strconv.Itoa(r.difficulty) + "\n"
	ret += "Kritikalität:  " + strconv.Itoa(r.criticality) + " / "
	ret += strconv.Itoa(r.criticality*2+10) + "\n"
	return ret
}

func (r *Recipe) getTableArr() []string {
	return []string{
		r.items,
		strconv.Itoa(r.difficulty),
		strconv.Itoa(r.criticality) + " / " + strconv.Itoa(r.criticality*2+10),
		strings.Join(r.formWords, ", "),
		strings.Join(r.colorWords, ", "),
		strconv.Itoa(r.strength),
		strconv.Itoa(r.priceSumItems),
	}
}

// #############################################################################
// #							INIT										   #
// #############################################################################

// reads md file into global alchData map
func readAlchData() map[string]map[string]*Ingredient {
	raw, _ := os.ReadFile(ingredientsFile)
	parseMode := alchBaseMode

	alchData := make(map[string]map[string]*Ingredient)
	alchData[parseMode] = make(map[string]*Ingredient)

	for _, line := range strings.Split(string(raw), "\n") {
		if line == switchMode {
			parseMode = alchAddMode
			alchData[parseMode] = make(map[string]*Ingredient)
		}
		matches := parseRex.FindStringSubmatch(line)
		if len(matches) < 4 {
			continue
		}

		newIngredient := new(Ingredient)
		newIngredient.CWords = make(map[string]int)
		newIngredient.FWords = make([]string, 0)
		alchData[parseMode][matches[1]] = newIngredient

		newIngredient.Difficulty, _ = strconv.Atoi(matches[2])
		newIngredient.Price, _ = strconv.Atoi(matches[3])

		info := matches[4]
		for len(info) > 0 {
			wordMatch := readEntity.FindStringSubmatch(info)
			if len(wordMatch) == 0 {
				break
			}
			if isCWord(wordMatch[1]) {
				if len(wordMatch[2]) == 0 {
					newIngredient.CWords[wordMatch[1]] = 0
				} else if strings.Contains(wordMatch[2], "-") {
					newIngredient.CWords[wordMatch[1]] = -1 * len(wordMatch[2])
				} else {
					newIngredient.CWords[wordMatch[1]] = len(wordMatch[2])
				}
			} else if isFWord(wordMatch[1]) {
				newIngredient.FWords = append(newIngredient.FWords, wordMatch[1])
			} else {
				deb("unbekannt: |"+wordMatch[1]+"|", line)
			}
			info = readEntity.ReplaceAllString(info, "")
		}
	}
	return alchData
}

// prints alchData as formated and sorted json
func printAlchData(_ *cobra.Command, _ []string) {
	alchDataBytes, err := json.MarshalIndent(alchData, "", "\t")
	if err != nil {
		deb(err)
	} else {
		deb(string(alchDataBytes))
	}
}

// #############################################################################
// #							SEARCH										   #
// #############################################################################

/*
Searches for matching Recipes
*/
func searchRecipes(opts []string) []*Recipe {

	recipes := []*Recipe{}

	adds := []string{}
	for a := range alchData[alchAddMode] {
		adds = append(adds, a)
	}
	base := []string{}
	for b := range alchData[alchBaseMode] {
		base = append(base, b)
	}
	addsRange := permForward(adds, maxSearchDepth)
	baseRange := permForward(base, maxSearchDepth)
	amount := len(addsRange) * len(baseRange)
	sNr := 0
	for _, a := range addsRange {
		for _, b := range baseRange {
			items := append(a, b...)
			sNr++
			r := calculateResultSearch(items)
			if verboseMode {
				deb(sNr, amount, items, r.colorWords, r.formWords)
			}
			if r.difficulty >= 10 &&
				r.difficulty <= maxDifficulty &&
				sameElemInArr(opts, append(r.colorWords, r.formWords...)) {
				recipes = append(recipes, r)
			}

		}
	}
	return recipes
}

/*
Prints a list of Recipes
*/
func printRecipes(recipes []*Recipe) {
	for _, r := range recipes {
		deb(r.getTableArr())
	}

	deb("\nANZAHL REZEPTE", len(recipes))
	// some infos about best recipes
	bestPrice := make(map[int]*Recipe)
	bestDiff := make(map[int]*Recipe)
	for _, r := range recipes {
		if val, ok := bestDiff[r.strength]; ok {
			if val.difficulty > r.difficulty {
				bestDiff[r.strength] = r
			}
			if val.priceSumItems > r.priceSumItems {
				bestPrice[r.strength] = r
			}
		} else {
			bestPrice[r.strength] = r
			bestDiff[r.strength] = r
		}
	}
	for i := 0; i < 10; i++ {
		deb("---", i, "---")
		if _, ok := bestDiff[i]; ok {
			deb("DIFF ", bestDiff[i].difficulty, bestDiff[i].getTableArr())
			deb("PRICE", bestPrice[i].priceSumItems, bestPrice[i].getTableArr())
		}
	}
}

// #############################################################################
// #							CALC										   #
// #############################################################################

func calculateResult(opts []string) *Recipe {

	// filter both levels (rule 1)
	fwordsBase := make(map[string]bool)
	fwords := make(map[string]bool)
	cwordsStrengthBase := make(map[string]int)
	cwordsStrength := make(map[string]int)
	diff := 10
	price := 0
	// check base items
	for _, item := range opts {
		itemType, ingredient := itemTyp(item)
		if itemType == alchBaseMode {
			for k, v := range ingredient.CWords {
				cwordsStrengthBase[k] += v
			}
			for _, k := range ingredient.FWords {
				fwordsBase[k] = true
			}
		} else if itemType == "" {
			deb("unknown substance " + item + ", aborting!")
			return nil
		}
		diff += ingredient.Difficulty
		price += ingredient.Price
	}
	// filter with adds
	for _, item := range opts {
		itemType, ingredient := itemTyp(item)
		if itemType == alchAddMode {
			for k, v := range ingredient.CWords {
				if v2, ok := cwordsStrengthBase[k]; ok {
					if _, ok2 := cwordsStrength[k]; !ok2 {
						cwordsStrength[k] = v2
					}
					cwordsStrength[k] += v
				}
			}
			for _, k := range ingredient.FWords {
				if v, ok := fwordsBase[k]; ok && v {
					fwords[k] = true
				}
			}
		}
	}

	// get only defined words
	fWordsSorted := make([]string, 0)
	for k := range fwords {
		for _, isThis := range [][]string{StartWords, MiddleWords, EndWords} {
			if isIn(k, isThis) {
				fWordsSorted = append(fWordsSorted, k)
			}
		}
	}

	// make rule checks and
	ret := new(Recipe)
	ret.items = strings.Join(opts, ", ")
	cwordsCleared, criticality, strength, valid := checkCwords(cwordsStrength)
	ret.criticality = criticality
	if !IsValidForm(fWordsSorted) {
		ret.criticality += len(fWordsSorted)
		return ret
	}
	if !valid {
		return ret
	}
	extraCritical := minAddsWithoutCrit - len(cwordsCleared) - len(fWordsSorted)
	if extraCritical > 0 {
		ret.criticality += extraCritical
	}

	ret.colorWords = cwordsCleared
	ret.formWords = fWordsSorted
	ret.strength = strength
	ret.difficulty = diff
	ret.priceSumItems = price
	return ret
}

func calculateResultSearch(opts []string) *Recipe {

	// to compare reduced ingredients against full list we get the full,
	// validation checks are made later in calc method so we can
	r := calculateResult(opts)
	if r.difficulty < 10 {
		return r
	}
	for _, itemtoRemove := range opts {
		newItems := []string{}
		for _, item := range opts {
			if itemtoRemove != item {
				newItems = append(newItems, item)
			}
		}
		smallerR := calculateResult(newItems)
		if sameElemInArr(smallerR.colorWords, r.colorWords) &&
			sameElemInArr(smallerR.formWords, r.formWords) &&
			smallerR.strength >= r.strength {

			// recipe can be reduced and is still the same or stronger we
			// invalidate the original and brea
			r.difficulty = 0
			break
		}
	}
	return r
}

// #############################################################################
// #							Subrules    								   #
// #############################################################################

func checkCwords(words map[string]int) ([]string, int, int, bool) {
	crit := 0
	maxStr := 0
	cwords := []string{}
	valid := true

	// reduce opposite colors
	newStrength := make(map[string]int)
	for word, str := range words {
		if opposite, okOp := colorOpposites[word]; okOp {
			if strOpp, oppIsHere := words[opposite]; oppIsHere {
				if str > 0 && strOpp > 0 {
					if str >= strOpp {
						newStrength[word] = str - strOpp
					} else {
						crit += str
						newStrength[word] = 0
					}
				} else if str < 0 && strOpp < 0 {
					if str <= strOpp {
						newStrength[word] = str - strOpp
					} else {
						crit += str
						newStrength[word] = 0
					}
				}
			}
		}
	}
	// validate check
	for k, v := range newStrength {
		if !alchConducMode {
			words[k] = v
		}
	}
	// get greates strength
	overAllNeg := false
	for _, str := range words {
		negativ := false
		if str < 0 {
			str *= -1
			negativ = true
		}
		if str > maxStr {
			overAllNeg = negativ
			maxStr = str
		}
	}
	for word, str := range words {
		negativ := false
		if str < 0 {
			negativ = true
			word = An + " " + word
			str *= -1
		}
		if str == maxStr {
			if negativ != overAllNeg {
				valid = false
			}
			cwords = append(cwords, word)
		} else {
			crit += str
		}
	}

	if len(cwords) > 1 && !alchConducMode {
		crit += len(cwords) * maxStr
		valid = false
	}
	if verboseMode && alchSearchMode && !valid {
		deb("color words", cwords)
	}

	return cwords, crit, maxStr, valid
}

/*
checks if magic words do make sense
*/
func IsValidForm(words []string) bool {
	aumInMiddle := false
	konfarEff := false
	sWords := []string{}
	mWords := []string{}
	eWords := []string{}
	for _, word := range words {
		if isIn(word, StartWords) {
			sWords = append(sWords, word)
		}
		if isIn(word, MiddleWords) {
			sWords = append(mWords, word)
		}
		if isIn(word, EndWords) {
			sWords = append(eWords, word)
		}
		aumInMiddle = word == Aum || aumInMiddle
		konfarEff = word == Konfar || konfarEff

	}
	invalid := len(mWords) > 2 ||
		(len(mWords) > 1 && !aumInMiddle) ||
		len(sWords) > 1 ||
		len(eWords) > 1 ||
		(konfarEff && len(words) > 1)

	if verboseMode && alchSearchMode && invalid {
		deb("start words", len(sWords), sWords)
		deb("form words", len(mWords), mWords)
		deb("end words", len(eWords), eWords)
	}
	return !invalid
}

// #############################################################################
// #							Sheet Mode									   #
// #############################################################################

func printAlchSheet(_ *cobra.Command, _ []string) {

	out := "#### Alchemieblatt\n"
	out += "Rezeptetabelle und Besitz fertiger alchemistischer Erzeugnisse "
	out += "gestafflet nach Herstellungserfolg und resultierender Willenskraft"
	out += "\"W\"\n\n"

	// recipes
	out += "| Rezept | Zutaten (Basis und Additive) | MW | krit. MW. | mag. Effekt | W0 | W1 | W2 | W3 | W4 |\n"
	out += "|--------|------------------------------|----|-----------|-------------|----|----|----|----|----|\n"
	for i := 0; i < emptyRecipes; i++ {
		out += "| | | | | | | | | | |\n"
	}
	out += "\n"
	for mode := range alchData {
		out += "| Zutat | Schwierigkeit | Preisvorschlag | Effekte oder sonstige Notizen |\n"
		out += "|-------|---------------|----------------|-------------------------------|\n"
		for _, itemName := range getSortedKeys(alchData[mode]) {
			item := alchData[mode][itemName]
			out += "| " + itemName + " | " + strconv.Itoa(item.Difficulty) + " | "
			out += strconv.Itoa(item.Price) + " | |\n"
		}
		out += "\n"
	}
	deb(out)
}

// #############################################################################
// #							UTIL										   #
// #############################################################################

// splits optional words and removes empty lines
func splitOpts(opts string) []string {
	ret := []string{}

	for len(opts) > 0 {
		wordMatch := readEntity.FindStringSubmatch(opts)
		if len(wordMatch) == 0 {
			break
		}
		ret = append(ret, wordMatch[1])
		opts = readEntity.ReplaceAllString(opts, "")
	}
	return ret
}

// searches alchData if item is base or add type (empty if unknown item)
func itemTyp(w string) (string, *Ingredient) {
	for itemType, itemMap := range alchData {
		for item, ingredient := range itemMap {
			if item == w {
				return itemType, ingredient
			}
		}
	}
	return "", nil
}

// checks if a give word (string) is part of the color words
func isCWord(w string) bool {
	return isIn(w, ColorWords)
}

// checks if a give word (string) is part of the form words
func isFWord(w string) bool {
	return isIn(w, FormWords)
}

// checks if word is in array
func isIn(w string, a []string) bool {
	for _, check := range a {
		if check == w {
			return true
		}
	}
	return false
}

/*
returns words depending on length, order doesnt matter so some shortcuts
*/
func permForward(words []string, depth int) [][]string {
	ret := [][]string{}
	if depth > 1 {
		for i, word := range words {
			for _, inner := range permForward(words[i:], depth-1) {
				if !isIn(word, inner) {
					ret = append(ret, append(inner, word))
				}
			}
		}
	}
	for _, w := range words {
		ret = append(ret, []string{w})
	}
	return ret
}

func getSortedKeys(m map[string]*Ingredient) []string {
	ret := make([]string, len(m))
	i := 0
	for key := range m {
		ret[i] = key
		i++
	}
	sort.Strings(ret)
	return ret
}

/*
checks if two arrays contain the same elements (are equal). Asserts no elements
are in there twice!
*/
func sameElemInArr(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, w := range a {
		if !isIn(w, b) {
			return false
		}
	}
	return true
}
