package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	latex "github.com/sboeck16/MDLatex"

	"github.com/spf13/cobra"
)

var (
	// debug printer
	debugPrint = ""

	ruleBookCmd = &cobra.Command{
		Use:   "rulebook",
		Short: "Erzeugt ein LaTeX - Dokument aus Skriptdatei",
		Long: "Mithilfe der Skriptdatei wird ein LaTex - Dokument erzeugt" +
			" und entweder in die Zieldatei (target-file) oder nach SDTOUT" +
			" geschrieben. Programme wie pdflatex können dann aus diesem " +
			" ein pdf generieren. Als Beispiel können vorhandene .rb Dateien" +
			" im Texteditor betrachtet werden.",
		Run: RunRuleBook,
	}

	// by cli flags
	script      string
	pathToRules string
	designFile  string
	targetFile  string

	// internal
	printTarget = os.Stdout

	// "constants" that will maybe be changable
	MDEnd   = `.md`
	JSONEnd = `.json`

	// subpaths maybe as part of rulebook and/or cli flags?
	IntroPath      = "Allgemein"
	CharactersPath = "Charaktere"
	DicesPath      = "Würfelproben"
	AttributesPath = "Attribute"
	BattlePath     = "Kampf"
	SkillsPath     = "Fertigkeiten"
	PathsPath      = "Wege"
	MasteriesPath  = "MeisterWege"
)

func init() {
	ruleBookCmd.PersistentFlags().StringVarP(&script, "script", "s", "script.json",
		"Pfad zur Zieldatei (oder als Ergänzung zum rule-path)")
	ruleBookCmd.PersistentFlags().StringVarP(&pathToRules, "rule_path", "r", ".",
		"Pfad zum Hauptverzeichnis der Markdown Regeln")
	ruleBookCmd.PersistentFlags().StringVarP(&targetFile, "target_file", "t", "",
		"Pfad für das generierte LaTeX. Wenn nicht gesetzt wird STDOUT"+
			" verwendet")
	ruleBookCmd.PersistentFlags().StringVarP(&designFile, "design-file", "d", "",
		"Ziel zur Datei mit zusätzlichen LaTeX Regeln (kann auch im"+
			" Skript gesetzt werde. JSON Format, siehe latex Modul oder"+
			" Beispiel anschauen")

	rootCmd.AddCommand(ruleBookCmd)
}

func RunRuleBook(cmd *cobra.Command, args []string) {
	rules, err := ReadRuleBook(pathToRules, script)
	if checkErr(err) {
		os.Exit(1)
	}

	// store all lines from read md lines
	mdLines := []string{}

	// holds prepared filenames with all directories from
	rFiles := []string{}

	// Add Chapters
	// ------------
	// Intro
	rFiles = append(rFiles, prepareChapterFlat(
		pathToRules, IntroPath, rules.Introduction, false)...)
	// Characters
	rFiles = append(rFiles, prepareChapterFlat(
		pathToRules, CharactersPath, rules.Characters, false)...)
	// Dices
	rFiles = append(rFiles, prepareChapterFlat(
		pathToRules, DicesPath, rules.Dices, false)...)
	// Battle
	rFiles = append(rFiles, prepareChapterFlat(
		pathToRules, BattlePath, rules.Battle, false)...)
	// Attributes
	rFiles = append(rFiles, prepareChapterFlat(
		pathToRules, AttributesPath, rules.Attributes, false)...)
	// Skills
	rFiles = append(rFiles, prepareChapterMap(
		pathToRules, SkillsPath, rules.Skills)...)
	// Paths
	rFiles = append(rFiles, prepareChapterMap(
		pathToRules, PathsPath, rules.Paths)...)
	//Masteries
	rFiles = append(rFiles, prepareChapterMap(
		pathToRules, MasteriesPath, rules.Masteries)...)

	// add additional chapters
	rFiles = append(rFiles, prepareChapterFlat(
		pathToRules, "", rules.Additional, false)...)

	// read files in this order
	for _, rFile := range rFiles {
		mdLines = append(mdLines, readMDFile(rFile)...)
	}

	// get doc and parse lines
	doc := getRuleBookDocHead()

	// check for title page
	if rules.Title != "" {
		picStr := ""
		if rules.TitlePicture != "" {
			picStr = `\includegraphics[width=\linewidth]{` +
				rules.TitlePicture + `}`

		}
		title := doc.AddTitle(rules.Title, rules.Author, picStr)
		if rules.FirstPageMD != "" {
			title.AddRaw(strings.Join(readMDFile(rules.FirstPageMD), "\n"))
		}
	}

	// parse markdown
	latex.WriteMDToLatexDoc(mdLines, doc)

	// add design option if defined
	if designFile != "" {
		checkErr(doc.GetDefaultsFromFile(designFile))
		checkErr(latex.LoadPreDefined(designFile))
	}

	// add character sheet
	addCharacterSheet(doc, rules)

	// write to target file
	if targetFile != "" {
		fHandle, err := os.OpenFile(targetFile, os.O_CREATE|os.O_RDWR, 0600)
		if err == nil {
			printTarget = fHandle
		}
	}

	// add last page
	if rules.LastPagePic != "" {
		doc.AddRaw(`\newgeometry{top=1mm, bottom=1mm, left=0mm, right=1mm}`)
		doc.AddRaw(`\includegraphics[width=\textwidth,height=\textheight]{` +
			rules.LastPagePic + `}`)
	}

	// debug
	if debugPrint != "" {
		debHandle, _ := os.OpenFile(debugPrint, os.O_CREATE|os.O_RDWR, 0600)
		debHandle.WriteString(strings.Join(mdLines, "\n"))
	}

	printTarget.WriteString(doc.String())
}

/*
EnsureFileEnding takes a path and ending. if path doesnt have that
ending, ending will be added and returned
*/
func EnsureFileEnding(file, ending string) string {
	if len(ending) > len(file) || file[len(file)-len(ending):] != ending {
		return file + ending
	}
	return file
}

/*
Prepares and returns a latex doc with needed packages and parameters
*/
func getRuleBookDocHead() *latex.Doc {

	doc := latex.NewDoc("article")
	// add our packagaes, maybe configurable (design, template, latex pre file?)
	doc.AddPackage("inputenc", []string{"utf8"})
	doc.AddPackage("fontenc", []string{"T1"})
	doc.AddPackage("babel", []string{"ngerman"})
	doc.AddPackage("geometry", []string{
		"a4paper", "left=2cm", "bottom=15mm", "top=2cm", "right=2cm"})
	doc.AddPackage("hyperref", []string{})
	doc.AddPackage("multicol", []string{})
	doc.AddTOC()
	doc.AddPackage("graphicx", []string{})
	doc.AddPackage("xcolor", []string{})
	doc.AddPackage("mdframed", []string{
		"framemethod=tikz"})
	doc.AddPackage("array", []string{})
	return doc
}

/*
Utility method that will join all file names to be read and returns them
*/
func prepareChapterFlat(
	rulePath, chapterPath string, addFiles []string, sortFiles bool) []string {

	ret := []string{}
	if len(addFiles) == 0 {
		// ommitted chapter return nothing
		return ret
	}

	fPath := filepath.Join(rulePath, chapterPath)

	// check if there is a headline file like Fertigkeiten.md in Fertigkeiten/
	if _, err := os.Stat(filepath.Join(fPath,
		EnsureFileEnding(chapterPath, MDEnd))); err == nil {
		ret = append(ret, filepath.Join(fPath, chapterPath))
	}

	// sort files, (better)
	if sortFiles {
		sort.Strings(addFiles)
	}

	// now add all files
	alreadyRead := map[string]bool{}
	for _, file := range addFiles {
		// remove duplications
		if _, ok := alreadyRead[file]; ok {
			continue
		}
		alreadyRead[file] = true

		ret = append(ret, filepath.Join(fPath, file))
	}

	return ret
}

/*
Utility method to read a chapter with sub chapters, and returns files to be read
*/
func prepareChapterMap(
	rulePath, chapterPath string, addFiles map[string][]string) []string {

	ret := []string{}

	fPath := filepath.Join(rulePath, chapterPath)

	// check if there is a headline file like Fertigkeiten.md in Fertigkeiten/
	if _, err := os.Stat(filepath.Join(fPath,
		EnsureFileEnding(chapterPath, MDEnd))); err == nil {
		ret = append(ret, filepath.Join(fPath, chapterPath))
	}

	// we saved to a map -> so we always sort
	addFilesKeys := make([]string, len(addFiles))
	ind := 0
	for key := range addFiles {
		addFilesKeys[ind] = key
		ind++
	}
	sort.Strings(addFilesKeys)

	// iterate through and use it as flat with redirected path
	for _, sub := range addFilesKeys {
		files := addFiles[sub]
		ret = append(ret, prepareChapterFlat(fPath, sub, files, true)...)
	}
	return ret
}

/*
Utility method to read a md file.
*/
func readMDFile(file string) []string {

	file = EnsureFileEnding(file, MDEnd)

	data, err := os.ReadFile(file)
	if checkErr(err) {
		return []string{}
	}

	return strings.Split(string(data), newLine)
}
