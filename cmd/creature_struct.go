package cmd

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	AttrOrder = []string{AttrStr, AttrDex, AttrCon, AttrInt,
		AttrWil, AttrPer, AttrCha}

	// to be put into a table they need to be ordered
	attributesPrint = [][]string{
		{AttrStr, AttrPer, AttrCha},
		{AttrDex, AttrInt, ""},
		{AttrCon, AttrWil, ""},
	}

	// for printing into md table
	skillsPrint = [][]string{
		{SkillMelee, skillAthletik},
		{SkillRanged, SkillStealth},
		{SkillMagic, SkillPercept},
	}

	crTableSizes = []int{8, 19}
	attrTHead1   = "| k.Attr | g.Attr | s.Attr |" + newLine
	attrTHead2   = "|--------|--------|--------|" + newLine
	skillHead1   = "| Fertigkeit Kampf | Fertigkeit sonst. |"
	skillHead2   = "|------------------|-------------------|"

	tagList    map[string]string
	attackList = new(AttackList)

	readTag       = regexp.MustCompile(`^\*\s+(.*)$`)
	readTagAction = regexp.MustCompile(`^[^:]+:\s*(.*)$`)
)

type Creature struct {
	name        string
	attributes  map[string]int
	skills      map[string]int
	baseSize    int
	sizeMod     int
	text        []string
	tags        []string
	attacks     []*Attack
	defence     int
	armor       int
	speed       int
	speedAir    int
	speedWater  int
	damageClass int
}

func (cr *Creature) String() string {
	ret := ""

	// name
	ret += "#### " + cr.name + newLine + newLine

	// text
	ret += strings.Join(cr.text, newLine) + newLine + newLine

	// table: Attr / Skills
	eMaps := []map[string]int{cr.attributes, cr.skills}
	tPrints := [][][]string{attributesPrint, skillsPrint}
	t1Head := []string{attrTHead1, skillHead1}
	t2Head := []string{attrTHead2, skillHead2}

	for ind, eMap := range eMaps {
		size := crTableSizes[ind]
		for _, line := range tPrints[ind] {
			ret += t1Head[ind]
			ret += t2Head[ind]
			for _, entry := range line {
				if entry == "" {
					ret += "|" + strings.Repeat(" ", size)
					continue
				}
				ret += "|" + tabStrI(entry, eMap[entry], size)
			}
			ret += "|" + newLine
		}
		ret += newLine
	}

	// defences
	ret += "Verteidigung: " + strconv.Itoa(cr.defence) + newLine
	ret += "Rüstung:      " + strconv.Itoa(cr.armor) + newLine
	ret += newLine

	// movements
	ret += "Geschw. Land:   " + strconv.Itoa(cr.speed) + newLine
	ret += "Geschw. Luft:   " + strconv.Itoa(cr.speedAir) + newLine
	ret += "Geschw. Wasser: " + strconv.Itoa(cr.speedWater) + newLine
	ret += newLine

	// weapons
	for _, att := range cr.attacks {
		ret += "* " + att.String() + newLine
	}
	ret += newLine

	// size
	size := cr.baseSize + cr.sizeMod
	if size < 1 {
		size = 1
	}
	if size > 9 {
		size = 9
	}
	ret += "Größe: " + sizeToName[size] + newLine + newLine

	// damage table
	rows := sizeToDCAndRows[size]
	cols := 5 + cr.attributes[AttrCon]
	modMap := sizeToRowMod[size]
	ret += MDMonitor(rows, cols, modMap)

	return ret
}

func MDMonitor(rows, cols int, modMap map[int]int) string {
	ret := ""
	ret += "|M.|"
	for i := 0; i < cols; i++ {
		ret += "|#"
	}
	ret += "|" + newLine
	ret += "|--"
	for i := 0; i < cols; i++ {
		ret += "|-"
	}
	ret += "|" + newLine

	lastmod := 0
	for r := 1; r <= rows; r++ {
		if val, ok := modMap[r]; ok {
			lastmod = val
		}
		if lastmod > 0 {
			ret += "|" + strconv.Itoa(lastmod)
		} else {
			ret += "|  "
		}
		for c := 0; c < cols; c++ {
			ret += "| "
		}
		ret += "|" + newLine
	}

	return ret
}

/*
adds l and r filling the middle with enough whitespaces to reach size in length
*/
func tabStr(l, r string, size int) string {
	return l + strings.Repeat(" ", size-len(l)-len(r)) + r
}

func tabStrI(l string, r int, size int) string {
	return tabStr(l, strconv.Itoa(r), size)
}

// #############################################################################
// #							Generate/Tag
// #############################################################################

func loadTags() {
	dir, err := os.ReadDir(tagDir)
	if checkErr(err) {
		os.Exit(1)
	}

	tagList = make(map[string]string)

	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		raw, err := os.ReadFile(tagDir + file.Name())
		if checkErr(err) {
			continue
		}
		tagList[file.Name()] = string(raw)
	}
}

func generateCreatureFromFile(targetFile string) *Creature {

	if tagList == nil {
		loadTags()
	}

	raw, err := os.ReadFile(creatureDir + string(os.PathSeparator) + targetCreature)
	if checkErr(err) {
		os.Exit(1)
	}

	ret := new(Creature)
	ret.name = targetFile
	ret.attributes = make(map[string]int)
	for _, attr := range AttrOrder {
		ret.attributes[attr] = 3
	}
	ret.skills = make(map[string]int)

	for _, line := range strings.Split(string(raw), newLine) {
		if len(line) == 0 {
			continue
		}

		if m := readTag.FindStringSubmatch(line); len(m) > 1 {
			ret.tags = append(ret.tags, m[1])
		} else {
			ret.text = append(ret.text, line)
		}
	}
	/*
		ret.linkupTags()
		ret.calcTags()
	*/
	return ret
}

// #############################################################################
// #							Attack
// #############################################################################

type AttackList struct {
	Attacks map[string]*Attack `json:"attacks"`
}

type Attack struct {
	// defined by data
	Name     string
	AddDescr string `json:"descr"`
	Skill    string `json:"skill"`
	Attr     string `json:"attribute"`
	PoolMod  int    `json:"pool_mod"`
	// adds to base damge class from creatur
	DcMod int `json:"dc_mod"`
	// defined after creature creation

	pool        int
	damageClass int
	strength    int
}

func (att *Attack) String() string {

	// safeguard
	dc := att.damageClass
	if dc < 0 {
		dc = 0
	}

	ret := att.Name + ": " + strconv.Itoa(att.pool) + ", "
	for {
		if val, ok := damageClassToDice[dc]; ok {
			ret += val + "+"
			ret += strconv.Itoa(damageClassToStrBonus[dc] * att.strength / 2)
			ret += " WV: " + strconv.Itoa(damageClassToWeaponVal[dc]) + "."
			if att.AddDescr != "" {
				ret += " " + att.AddDescr
			}
			break
		}
	}
	return ret
}

/*
Clone creates a copy of a given attack. needs name.
*/
func (att *Attack) Clone(name string) *Attack {
	ret := new(Attack)
	ret.Name = name
	ret.AddDescr = att.AddDescr
	ret.Skill = att.Skill
	ret.Attr = att.Attr
	ret.PoolMod = att.PoolMod
	ret.DcMod = att.DcMod
	return ret
}

/*

 */
