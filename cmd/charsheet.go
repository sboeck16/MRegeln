package cmd

import (
	"os"
	"strings"

	latex "github.com/sboeck16/MDLatex"

	"github.com/spf13/cobra"
)

var (
	charsheetCmd = &cobra.Command{
		Use:   "charsheet",
		Short: "Generiert ein Charakterblatt aus einem Rulebook",
		Long: "Erstellt ein Charakterblatt als Latexdokument. Es muss ein " +
			"rulebook angegeben werden damit das Charakterblatt generiert " +
			"werden kann.",
		Run: NewDocCharSheet,
	}
)

func init() {
	ruleBookCmd.AddCommand(charsheetCmd)
}

func NewDocCharSheet(cmd *cobra.Command, args []string) {
	// read and prepare rulebook, code is doubled... (Maybe fix it later?)
	doc := getRuleBookDocHead()
	rules, err := ReadRuleBook(pathToRules, script)
	if checkErr(err) {
		os.Exit(1)
	}
	// only write char sheet
	addCharacterSheet(doc, rules)

	// print it to shell, character sheet is part of book anyway and normal user
	// should use it from there.
	deb(doc.String())
}

func addCharacterSheet(doc *latex.Doc, rules *RuleBook) {
	// page break
	doc.AddRaw(`\newpage`)

	// description
	sheet := doc.AddSection(latex.NewText("Charakterblatt"))

	// GENERAL
	tDesc := sheet.AddTableNoBorder("description", "m{6cm} m{6cm} m{6cm}")
	tDesc.AddRowStr("Name:", "Herkunft:", "Beruf:")
	tDesc.AddRowStr("Beschreibung:", "", "")
	tDesc.AddRowStr("", "", "")
	tDesc.AddRowStr("", "", "")
	tDesc.AddRowStr("", "", "")
	sheet.AddRaw(`\hline`)

	// ATTRIBUTS
	attr := sheet.AddChild(latex.NewText("Attribute"))
	tAttr := attr.AddTableNoBorder(
		"char-sheet-attributes", "m{6cm} m{6cm} m{6cm}")
	for _, r := range flipToRows(rules.Attributes, 3) {
		tAttr.AddRowStr(r...)
	}

	attr.AddRaw(`\newline\hline`)

	// BATTLESTATS
	battleStats := sheet.AddChild(latex.NewText("Kampfdaten"))
	tBattle := battleStats.AddTableNoBorder("battle-overview", "m{10cm} m{10cm}")
	tBattle.AddRowStr("Verteidigung (14+Ges+Boni-Behinderung):", "Waffe:")
	tBattle.AddRowStr("Geschwindigkeit (3 + (Stärke+Geschick)/2):", "Waffe:")
	tBattle.AddRowStr("Rüstungsschutz:", "Waffe:")
	tBattle.AddRowStr("", "")

	// MONITORS
	tMon1 := latex.NewTable("damage-monitor",
		"|"+strings.Repeat(" c |", 14))
	tMon1.NoDoubleHLine = true
	tMon1.AddRowStr("0", "", "", "", "", "", "", "", "", "", "", "", "", "")
	tMon1.AddRowStr("-1", "", "", "", "", "", "", "", "", "", "", "", "", "")
	tMon1.AddRowStr("-2", "", "", "", "", "", "", "", "", "", "", "", "", "")
	tMon1.AddRowStr("-3", "", "", "", "", "", "", "", "", "", "", "", "", "")
	tMon1.AddRowStr("-4", "", "", "", "", "", "", "", "", "", "", "", "", "")
	tMon2 := latex.NewTable("damage-monitor",
		"|"+strings.Repeat(" c |", 14))
	tMon2.NoDoubleHLine = true
	tMon2.AddRowStr("0", "", "", "", "", "", "", "", "", "", "", "", "", "")
	tMon2.AddRowStr("-1", "", "", "", "", "", "", "", "", "", "", "", "", "")
	tMon2.AddRowStr("-2", "", "", "", "", "", "", "", "", "", "", "", "", "")
	tMon2.AddRowStr("-3", "", "", "", "", "", "", "", "", "", "", "", "", "")
	tMon2.AddRowStr("-4", "", "", "", "", "", "", "", "", "", "", "", "", "")
	tMons := battleStats.AddTableNoBorder("damage-monitors", " m{9cm} m{9cm} ")
	tMons.AddRow([]latex.LatexStr{
		latex.BoldText("Körperlicher Monitor (5+Konstitution)"),
		latex.BoldText("Geistiger Monitor (5+Willenskraft)")})
	tMons.AddRow([]latex.LatexStr{tMon1, tMon2})
	battleStats.AddRaw(`\newline`)
	battleStats.AddText("")
	battleStats.AddText("Zustände:")
	sheet.AddRaw(`\newline\hline`)

	// SKILLS
	skills := sheet.AddChild(latex.NewText("Fertigkeiten"))
	skillTab := makeSkillTable(rules)
	skills.AppendContent(skillTab)
	sheet.AddText("")
	sheet.AddRaw(`\newline\hline`)

	// ETC.
	etc := sheet.AddChild(latex.NewText("Sonstiges"))
	eTab := etc.AddTableNoBorder("etc", " m{10cm} m{10cm}")
	eTab.AddRow([]latex.LatexStr{latex.NewText("Erfahrungspunkte / Gesamt:"),
		latex.BoldText("Ausrüstung:")})
	eTab.AddRowStr("Anzahl gekaufter Attributspunkte:", "")
	eTab.AddRowStr("Anzahl erlernter Wege / Meisterwege:", "")
	eTab.AddRow([]latex.LatexStr{latex.BoldText("Weg / Stufe"),
		latex.Text{}})
}

func flipToRows(inp []string, cols int) [][]string {
	ret := [][]string{}
	count := 0
	for colInd := 0; colInd < cols; colInd++ {
		for rowInd := 0; rowInd < len(inp)/cols+1; rowInd++ {
			if count == len(inp) {
				break
			}
			if colInd == 0 {
				newRow := make([]string, cols)
				newRow[0] = inp[count]
				ret = append(ret, newRow)
			} else {
				ret[rowInd][colInd] = inp[count]
			}
			count++
		}
	}
	return ret

}

func makeSkillTable(rules *RuleBook) *latex.Table {
	all := []string{}
	cats := []int{}
	count := 0
	for _, category := range rules.getSortedSkillCategories() {
		cats = append(cats, count)
		all = append(all, category)
		count++
		for _, skill := range rules.Skills[category] {
			all = append(all, skill+":")
			count++
		}
		all = append(all, "______")
		all = append(all, "______")
		all = append(all, "")
		count += 3
	}
	cols := 4
	colLen := len(all) / cols
	if len(all)%cols != 0 {
		colLen += 1
	}
	tabRows := flipToRows(all, cols)
	rows := [][]latex.LatexStr{}
	for ri, r := range tabRows {
		row := []latex.LatexStr{}
		for ci, c := range r {
			pos := ri + colLen*ci
			lText := latex.NewText(c)
			for _, cat := range cats {
				if pos == cat {
					lText = latex.BoldText(c)
				}
			}
			row = append(row, lText)
		}
		rows = append(rows, row)
	}
	tab := latex.NewTable("skills-table", strings.Repeat(" m{4cm} ", cols))
	tab.NoInnerHLine = true
	tab.NoOuterHLine = true
	tab.NoHeadBodyHLine = true
	for _, row := range rows {
		tab.AddRow(row)
	}
	return tab
}
