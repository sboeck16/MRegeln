package cmd

import (
	"fmt"
	"os"

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
			deb(lnxShellRed + errLogStr + lnxShellReset + " " + err.Error())
			ret = true
		}
	}
	return ret
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
