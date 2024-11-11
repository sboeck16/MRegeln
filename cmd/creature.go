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
	attackDir      string
	tagDir         string
)

func init() {
	creatureCmd.PersistentFlags().StringVarP(&tagDir, "tagdir", "s",
		"Kreaturen/Tags", "Verzeichnis mit den Kreaturen - Tags")
	creatureCmd.PersistentFlags().StringVarP(&creatureDir, "creaturedir", "d",
		"Kreaturen", "Verzeichnis mit den Kreaturen")
	creatureCmd.PersistentFlags().StringVarP(&attackDir, "attackdir", "a",
		"Kreaturen/Angriffe", "Verzeichnis mit den Angriffen")
	creatureCmd.PersistentFlags().StringVarP(&targetCreature, "creaturefile", "c",
		"", "Kreaturdatei, es werden alle generiert wenn leer.")

	rootCmd.AddCommand(creatureCmd)
}

func CreateCreature(cmd *cobra.Command, args []string) {
	// load additional data
	loadTags()
	loadAttacks()

	if targetCreature == "" {
		deb("TODO")
	} else {
		deb(genCreature(targetCreature))
	}
}

func genCreature(targetFile string) string {
	creature := generateCreatureFromFile(targetFile)
	return creature.String()
}
