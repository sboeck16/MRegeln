package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	// shell colors. needs moving if ever build for non linux systems
	lnxShellRed   = "\033[91m"
	lnxShellReset = "\033[0m"
	newLine       = "\n"

	errLogStr = "[ERROR]"
)

var (
	// Version to be printed
	Version string

	// root command
	rootCmd = &cobra.Command{
		Use:   "rpgcli",
		Short: "rpgcli ist ein Hilfswerkzeug für MRegeln",
		Long: "rpgcli stellt Hilfsprogramme zum Umgang mit den MRegeln bereit." +
			"Mehr Information mit -h oder --help. Als cli Tool stapelt es" +
			" Funktionen mit `subcommands` und fügt Parameter hinzu (alles" +
			" was mit - anfängt.",
		Run: rootRun,
	}
)

func rootRun(cmd *cobra.Command, args []string) {
	deb("Version", Version)
	checkErr(cmd.Help())
}

func Execute() {
	if checkErr(rootCmd.Execute()) {
		// error is printed now return with non zero
		os.Exit(1)
	}
}

// #############################################################################
// #							Utility
// #############################################################################

func deb(i ...any) {
	fmt.Println(i...)
}

// maybe better error logging and printing?
func checkErr(errs ...error) bool {
	ret := false
	for _, err := range errs {
		if err != nil {
			logError(err.Error())
			ret = true
		}
	}
	return ret
}

// quick and dirty shell printer
func logError(text string) {
	deb(lnxShellRed + errLogStr + lnxShellReset + " " + text)
}

/*
checks if two arrays contain the same elements (are equal). Asserts no elements
are in there twice!
*/
func sameArr[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

/*
MDTable generates a table where row add index 0 is used as head. generates a
human friendly output.
*/
func MDTable(inp [][]string) string {
	if len(inp) == 0 {
		logError("tried to gen a md table with empty input")
		return ""
	}
	// get size
	lens := make([]int, len(inp[0]))
	for rowInd := range inp {
		for colInd := range inp[0] {
			field := inp[rowInd][colInd]
			if len(field) > lens[colInd] {
				lens[colInd] = len(field)
			}
		}
	}
	ret := "|"
	headDiv := "|"
	for colInd := range inp[0] {
		field := inp[0][colInd]
		ret += field + strings.Repeat(" ", lens[colInd]-len(field)) + "|"
		headDiv += strings.Repeat("-", lens[colInd]) + "|"
	}
	ret += newLine + headDiv + newLine
	for rowInd := range inp {
		if rowInd == 0 {
			continue
		}
		ret += "|"
		for colInd := range inp[0] {
			field := inp[rowInd][colInd]
			ret += field + strings.Repeat(" ", lens[colInd]-len(field)) + "|"
		}
		ret += newLine
	}
	return ret
}
