package cmd

import (
	"github.com/spf13/cobra"
)

var (
	creatureCmd = &cobra.Command{
		Use:   "creature",
		Short: "Erzeugt eine MD Representation eines Monsters",
		Long: "Mit Hilfe der angegebenen Datei wird ein Monster erzeugt. " +
			"Weitere Parameter können angegeben werden um das Monster " +
			"zu modifizieren",
		Run: CreateCreature,
	}

	targetCreature string
	creatureDir    string
	tagDir         string
)

func init() {
	creatureCmd.PersistentFlags().StringVarP(&tagDir, "tagdir", "s",
		"Kreaturen/Tags", "Verzeichnis mit den Kreaturen - Tags")
	creatureCmd.PersistentFlags().StringVarP(&creatureDir, "creaturefile", "c",
		"Kreaturen", "Verzeichnis mit den Kreaturen")

}

func CreateCreature(cmd *cobra.Command, args []string) {

}
