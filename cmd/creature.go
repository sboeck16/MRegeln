package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	creatureCmd = &cobra.Command{
		Use:   "creature",
		Short: "Erzeugt eine MD Representation eines Monsters",
		Long: "Mit Hilfe der angegebenen Datei wird ein Monster erzeugt. " +
			"Weitere Parameter können angegeben werden um das Monster " +
			"zu modifizieren. Ist kein Zielmonster angegeben so werden alle " +
			"im Kreaturenverzeichnis erstellt und in der Zieldatei abgelegt.",
		Run: CreateCreature,
	}

	targetCreature string
	creatureDir    string
	attackDir      string
	tagDir         string
	addTags        string
	writeTo        string
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
	creatureCmd.PersistentFlags().StringVarP(&addTags, "tags", "t",
		"", "Zusätzliche Tags, kommasepariert (ohne Leerzeichen)")
	creatureCmd.PersistentFlags().StringVarP(&writeTo, "writeTo", "w",
		"Spielleitung/Kreaturen.md", "Zieldatei für Kapitelgeneration")

	rootCmd.AddCommand(creatureCmd)
}

func CreateCreature(cmd *cobra.Command, args []string) {
	// load additional data
	loadTags()
	loadAttacks()

	if targetCreature == "" {
		str := "## Kreaturen" + newLine + newLine
		dir, err := os.ReadDir(creatureDir)
		if checkErr(err) {
			os.Exit(1)
		}
		for _, file := range dir {
			if file.IsDir() {
				continue
			}
			str += genCreature(file.Name()) + newLine
		}
		os.WriteFile(writeTo, []byte(str), 0644)
	} else {
		deb(genCreature(targetCreature))
	}
}

func genCreature(targetFile string) string {
	creature := generateCreatureFromFile(targetFile)
	return creature.String()
}
