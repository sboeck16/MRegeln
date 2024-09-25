package cmd

import (
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
)

const (
	fltPrec           = 3
	defaultDice       = 6
	defaultBestAmount = 4
	defaultMaxDice    = 20
)

var (
	// regex for parsing
	findSpace  = regexp.MustCompile(`\s`)
	findDice   = regexp.MustCompile(`^(.*?)([0-9.]*)[wdWD]([0-9.]*)(.*?)$`)
	findPar    = regexp.MustCompile(`^(.*?)\(([^\(\)]*)\)(.*?)$`)
	findDivide = regexp.MustCompile(`^(.*?)([0-9.]+)\/([0-9.]+)(.*?)$`)
	findMult   = regexp.MustCompile(`^(.*?)([0-9.]+)\*([0-9.]+)(.*?)$`)
	findAdd    = regexp.MustCompile(`^(.*?)([0-9.]+)\+([0-9.]+)(.*?)$`)
	findSub    = regexp.MustCompile(`^(.*?)([0-9.]+)-([0-9.]+)(.*?)$`)

	rollCmd = &cobra.Command{
		Use:   "roll",
		Short: "Wird verwendet um Würfelwürfe zu erzeugen und zu berechnen",
		Long: "Nimmt eine Zeichenkette entgegen und versucht diese zu" +
			" verarbeiten. Besteht die Zeichenkette nur aus einer Ganzzahl " +
			"so werden entsprechend Würfel geworfen und die besten 4 " +
			"addiert. Weitere Möglichkeiten w4+2w8+5. Grundrechenarten und " +
			"Klammern sind möglich",
		Run: printDiceRollResString,
	}

	genCmd = cobra.Command{
		Use:   "gen",
		Short: "Generiert ein Tabelle von Ergebniswahrscheinlichkeiten",
		Long: "Die Tabelle der Wahrscheinlichkeiten mit X Würfeln ein " +
			"Ergebnis mindestens zu erreichen wird in Markdown generiert. " +
			"Die Angaben sind in Prozent. Es besteht die Möglichkeit die " +
			"Parameter anzupassen.",
		Run: printPercTable,
	}
	// flags for generation
	precTable  bool
	diceSide   int
	bestAmount int
	toPool     int
)

func init() {
	rootCmd.AddCommand(rollCmd)
	rollCmd.AddCommand(&genCmd)
	genCmd.Flags().BoolVarP(&precTable, "prec", "p", false,
		"Wenn gesetzt, Wahrscheinlichkeiten werden so berechnet das "+
			"das Ergebnis genau erreicht wird.")
	genCmd.Flags().IntVarP(&diceSide, "dice", "d", defaultDice,
		"Seitenzahl des verwendeten Würfels")
	genCmd.Flags().IntVarP(&bestAmount, "amount", "a", defaultBestAmount,
		"Die 'X' besten Würfel die gewertet werden.")
	genCmd.Flags().IntVarP(&toPool, "maxDice", "m", defaultMaxDice,
		"Maximale Poolgröße (Spalten in der Tabelle)")

}

/*
eval input and add defaults before add generating table
*/
func printPercTable(_ *cobra.Command, _ []string) {
	fromPool := bestAmount
	fromTarget := bestAmount
	toTarget := bestAmount * diceSide
	deb(generatePercTable(
		diceSide, bestAmount, fromPool, toPool, fromTarget, toTarget, precTable))
}

/*
Generates a md string for a dice result table. If prec is set percentage to
roll nr and nothing else. If fset to false number must be reached
*/
func generatePercTable(
	dice, best, fromPool, toPool, fromTarget, toTarget int, prec bool) string {

	md := `|pool\mw |`
	for targ := fromTarget; targ <= toTarget; targ++ {
		md += ` ` + strconv.Itoa(targ) + ` |`
	}
	md += "\n|"
	for targ := fromTarget - 1; targ <= toTarget; targ++ {
		md += `-|`
	}
	md += "\n|"

	for pool := fromPool; pool <= toPool; pool++ {
		md += `| ` + strconv.Itoa(pool) + ` |`
		all := math.Pow(float64(dice), float64(pool))
		remain := all
		res := calcResMap(calcPercDiceVsRes(dice, best, pool))
		for targ := 0; targ <= toTarget; targ++ {
			perc := strconv.FormatFloat(remain/all*100, 'f', 1, 64)
			val, _ := res[targ]
			if prec {
				perc = strconv.FormatFloat(float64(val)/all*100, 'f', 1, 64)

			}
			remain -= float64(val)
			if targ < fromTarget {
				continue
			}
			md += ` ` + perc + ` |`
		}

		md += "\n"
	}
	return md
}

/*
returns amount of results per set of n d sided dices from k.
*/
func calcPercDiceVsRes(d, n, k int) map[*[]int]int {
	// prepare n rolled d sided dice
	diceResAmount := map[*[]int]int{}
	// initial values so we can iterate over map
	for _, initial := range expandByNumbers([]int{}, d) {
		il := []int{initial[0]}
		diceResAmount[&il] = 1
	}
	// add dices
	for ik := 2; ik <= k; ik++ {
		// store the next iteration over one more dice
		newDiceResAmount := make(map[*[]int]int)
		// iterate over each known result
		for dRes, amount := range diceResAmount {
			// we add each possible dice result to the result
			newResults := expandByNumbers(*dRes, d)
			for _, newRes := range newResults {
				if len(newRes) > n {
					newRes = newRes[:n]
				}
				// see if newRes is any better than old result
				// look if already exist and if so count up
				found := false
				for lookRes := range newDiceResAmount {
					if sameArr(*lookRes, newRes) {
						newDiceResAmount[lookRes] += amount
						found = true
						break
					}
				}
				if !found {
					newDiceResAmount[&newRes] = amount
				}
			}
		}
		diceResAmount = newDiceResAmount
	}
	return diceResAmount
}

/*
utility function that expands array by 1 to d.
Sort the array before returning
*/
func expandByNumbers(arr []int, d int) [][]int {
	ret := [][]int{}
	for i := 1; i <= d; i++ {
		newArr := []int{i}
		for _, elem := range arr {
			newArr = append(newArr, elem)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(newArr)))
		ret = append(ret, newArr)
	}
	return ret
}

/*
utility function for adding up calc res maps
*/
func calcResMap(r map[*[]int]int) map[int]int {
	resValAmount := map[int]int{}
	for k, v := range r {
		res := 0
		for _, dr := range *k {
			res += dr
		}
		resValAmount[res] += v
	}
	return resValAmount
}

/*
debug utility
*/
func printDRes(r map[*[]int]int) {
	resVal := []int{}
	resValAmount := map[int]int{}
	for k, v := range r {
		deb(*k, v)
		res := 0
		for _, dr := range *k {
			res += dr
		}
		resValAmount[res] += v
	}
	deb("----")
	for k := range resValAmount {
		resVal = append(resVal, k)
	}
	sort.Ints(resVal)
	overall := 0
	for _, res := range resVal {
		overall += resValAmount[res]
		deb(res, resValAmount[res])
	}
	deb("----")
	deb(overall)
}

// #############################################################################
// #############################################################################

func printDiceRollResString(_ *cobra.Command, args []string) {
	for _, arg := range args {

		deb(rollString(arg))
	}
}

/*
Returns a string where "k" d-sided dices are thrown and
the best "n" are used. if add is set result will be raised
by one for every max rolled dice that is not part of n.
*/
func bestOf(d, n, k int, add bool) (string, int) {
	rolls := []int{}
	for i := 0; i < k; i++ {
		rolls = append(rolls, roll(d))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(rolls)))
	ret := "["
	res := 0

	for i := 0; i < k; i++ {
		if i < n {
			res += rolls[i]
		} else if rolls[i] == d && add {
			res += 1
		}
		if i != 0 {
			ret += ","
		}
		ret += strconv.Itoa(rolls[i])
	}
	ret += "] -> " + strconv.Itoa(res)
	return ret, res
}

/*
Rolls a single d sided dice. if d is zero default size is used
*/
func roll(d int) int {
	if d == 0 {
		d = defaultDice
	}

	ret := rand.Intn(d) + 1
	return ret
}

/*
Rolls and calculates a complex string.
*/
func rollString(s string) string {
	s = findSpace.ReplaceAllString(s, ``)
	matches := findPar.FindStringSubmatch(s)
	for len(matches) > 1 {
		s = matches[1] + rollString(matches[2]) + matches[3]
		matches = findPar.FindStringSubmatch(s)
	}

	matches = findDice.FindStringSubmatch(s)
	for len(matches) > 1 {
		amount, _ := strconv.ParseFloat(matches[2], 64) // #nosec
		if amount == 0.0 {
			amount = 1.0
		}
		d, _ := strconv.ParseFloat(matches[3], 64)
		res := 0
		for amount > 0 {
			amount--
			res += roll(int(d))
		}
		s = matches[1] + strconv.Itoa(res) + matches[4]
		matches = findDice.FindStringSubmatch(s)
	}

	matches = findDivide.FindStringSubmatch(s)
	for len(matches) > 1 {
		a, _ := strconv.ParseFloat(matches[2], 64) // #nosec
		b, _ := strconv.ParseFloat(matches[3], 64) // #nosec
		s = matches[1] + strconv.FormatFloat(a/b, 'f', fltPrec, 64) + matches[4]
		matches = findDivide.FindStringSubmatch(s)
	}

	matches = findMult.FindStringSubmatch(s)
	for len(matches) > 1 {
		a, _ := strconv.ParseFloat(matches[2], 64) // #nosec
		b, _ := strconv.ParseFloat(matches[3], 64) // #nosec
		s = matches[1] + strconv.FormatFloat(a*b, 'f', fltPrec, 64) + matches[4]
		matches = findMult.FindStringSubmatch(s)
	}

	matches = findAdd.FindStringSubmatch(s)
	for len(matches) > 1 {
		a, _ := strconv.ParseFloat(matches[2], 64) // #nosec
		b, _ := strconv.ParseFloat(matches[3], 64) // #nosec
		s = matches[1] + strconv.FormatFloat(a+b, 'f', fltPrec, 64) + matches[4]
		matches = findAdd.FindStringSubmatch(s)
	}

	matches = findSub.FindStringSubmatch(s)
	for len(matches) > 1 {
		a, _ := strconv.ParseFloat(matches[2], 64) // #nosec
		b, _ := strconv.ParseFloat(matches[3], 64) // #nosec
		s = matches[1] + strconv.FormatFloat(a-b, 'f', fltPrec, 64) + matches[4]
		matches = findSub.FindStringSubmatch(s)
	}

	return s
}
