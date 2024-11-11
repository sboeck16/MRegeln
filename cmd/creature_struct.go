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
	attrTHead1   = "| k.Attr | g.Attr | s.Attr |"
	attrTHead2   = "|--------|--------|--------|"
	skillHead1   = "| Fertigkeit Kampf | Fertigkeit sonst. |"
	skillHead2   = "|------------------|-------------------|"

	tagList    = map[string]string{}
	attackList = map[string]*Attack{}

	// tags
	readTag        = regexp.MustCompile(`^\s*\*\s*([^:]+)$`)
	readModFeature = regexp.MustCompile(`^\s*\*\s*(.*)\s*:\s*([+-].*)$`)
	readSetFeature = regexp.MustCompile(`^\s*\*\s*(.*)\s*:\s*([^+-]*)$`)
	// attack
	readAttackLine = regexp.MustCompile(`^\s*(.*):\s*(.*?)\s*$`)
)

type Creature struct {
	name           string
	baseAttributes map[string]int
	attributesMod  map[string]int
	attributes     map[string]int
	baseSkills     map[string]int
	skillsMod      map[string]int
	skills         map[string]int
	baseSize       int
	sizeMod        int
	damageClassMod int
	typ            int
	text           []string
	tags           map[string]int
	attacks        []*Attack
	movements      []int
	moveMod        int
	armor          int
	armorMod       int
}

func NewCreature() *Creature {
	ret := new(Creature)
	ret.baseAttributes = map[string]int{}
	ret.attributes = map[string]int{}
	ret.attributesMod = map[string]int{}
	ret.baseSkills = map[string]int{}
	ret.skills = map[string]int{}
	ret.skillsMod = map[string]int{}
	ret.tags = map[string]int{}
	return ret
}

func (cr *Creature) String() string {
	ret := ""

	// name
	ret += "#### " + cr.name + newLine + newLine

	// text
	ret += strings.Join(cr.text, newLine) + newLine + newLine

	// table: Attr / Skills
	ret += attrTHead1 + skillHead1 + newLine
	ret += attrTHead2 + skillHead2 + newLine

	// defences
	sizeName, rows, modMap, defMod := getSizeInfo(cr)
	ret += "Verteidigung: " + strconv.Itoa(10+defMod+cr.attributes[AttrDex]) +
		newLine
	ret += "Rüstung:      " + strconv.Itoa(cr.armor+cr.armorMod) + newLine
	ret += newLine

	// movements
	if len(cr.movements) == 0 {
		cr.movements = append(cr.movements, move4Legs)
	}
	for _, moveTyp := range cr.movements {
		ret += getMovementStr(cr, moveTyp) + newLine
	}
	ret += newLine

	// weapons
	for _, att := range cr.attacks {
		ret += "* " + att.genText(cr) + newLine
	}
	ret += newLine

	// size
	ret += "Größe: " + sizeName + newLine + newLine

	// damage table
	cols := 5 + cr.attributes[AttrCon]
	ret += cr.MDMonitor(rows, cols, modMap)

	return ret
}

func (cr *Creature) MDMonitor(rows, cols int, modMap map[int]int) string {
	ret := ""
	ret += "|M.|" + strings.Repeat("#|", cols) + newLine
	ret += "|--|" + strings.Repeat("-|", cols) + newLine

	lastmod := 0
	for r := 1; r <= rows; r++ {
		if val, ok := modMap[r]; ok {
			lastmod = val
		}
		if lastmod != 0 {
			ret += "|" + strconv.Itoa(lastmod) + "|"
		} else {
			ret += "|  |"
		}
		ret += strings.Repeat(" |", cols) + newLine
	}

	return ret
}

// #############################################################################
// #							Generate/Tag
// #############################################################################

func loadTags() {
	dir, err := os.ReadDir(tagDir)
	if checkErr(err) {
		os.Exit(1)
	}

	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		raw, err := os.ReadFile(tagDir + string(os.PathSeparator) + file.Name())
		if checkErr(err) {
			continue
		}
		tagList[file.Name()] = string(raw)
	}
}

func generateCreatureFromFile(targetFile string) *Creature {

	raw, err := os.ReadFile(creatureDir + string(os.PathSeparator) +
		targetFile)
	if checkErr(err) {
		os.Exit(1)
	}

	cr := NewCreature()
	cr.name = targetFile

	// gather tags and mods
	cr.AddTag(string(raw))

	// calc resulting attributes and skills
	for _, attr := range AttrList {
		cr.attributes[attr] = cr.baseAttributes[attr] + cr.attributesMod[attr]
		if cr.attributes[attr] < 1 {
			cr.attributes[attr] = 1
		}
	}
	for _, skill := range SkillList {
		cr.skills[skill] = cr.baseSkills[skill] + cr.skillsMod[skill]
		if cr.skills[skill] < 0 {
			cr.skills[skill] = 0
		}
	}

	return cr
}

func (cr *Creature) AddTag(tagString string) {
	for _, line := range strings.Split(tagString, newLine) {

		if line == "" {
			continue
		}

		// found another tag -> rekursion
		if m := readTag.FindStringSubmatch(line); len(m) > 0 {
			if _, ok := tagList[m[1]]; !ok {
				logError("tag not found: " + m[1])
				continue
			}
			if _, ok := cr.tags[m[1]]; !ok {
				cr.tags[m[1]] += 1
				cr.AddTag(tagList[m[1]])
			}
			continue
		}

		// Set
		if m := readSetFeature.FindStringSubmatch(line); len(m) > 0 {
			if m[1] == creatureAttack {
				if attack, ok := attackList[m[2]]; ok {
					cr.attacks = append(cr.attacks, attack)
				} else {
					logError("unknown attack: " + m[2])
				}
				continue
			}
			nr, err := strconv.Atoi(m[2])
			if err != nil {
				logError("could not set feature " + m[1] + ":" + m[2])
				continue
			}
			found := false
			for _, attr := range AttrList {
				if m[1] == attr {
					cr.baseAttributes[m[1]] = nr
					found = true
				}
			}
			for _, attr := range SkillList {
				if m[1] == attr {
					cr.baseSkills[m[1]] = nr
					found = true
				}
			}
			if found {
				continue
			}

			switch m[1] {
			case creatureSize:
				cr.baseSize = nr
			case creatureType:
				cr.typ = nr
			case creatureMovement:
				cr.movements = append(cr.movements, nr)
			case creatureArmor:
				cr.armor = nr
			default:
				logError("unknown set feature: " + m[1] + ":" + m[2])
			}
			continue
		}

		// Mod
		if m := readModFeature.FindStringSubmatch(line); len(m) > 0 {
			nr, err := strconv.Atoi(m[2])
			if err != nil {
				logError("could not mod feature " + m[1] + ":" + m[2])
				continue
			}
			found := false
			for _, attr := range AttrList {
				if m[1] == attr {
					cr.attributesMod[m[1]] += nr
					found = true
				}
			}
			for _, attr := range SkillList {
				if m[1] == attr {
					cr.skillsMod[m[1]] += nr
					found = true
				}
			}
			if found {
				continue
			}

			switch m[1] {
			case creatureSize:
				cr.sizeMod += nr
			case creatureArmor:
				cr.armorMod += nr
			case creatureDCMod:
				cr.damageClassMod += nr
			case creatureMovement:
				cr.moveMod += nr
			default:
				logError("unknown mod feature: " + m[1] + ":" + m[2])
			}
			continue
		}
		cr.text = append(cr.text, line)
	}
}

// #############################################################################
// #							Attack
// #############################################################################

type Attack struct {
	// Basic Stats
	Name     string
	AddDescr string
	Skill    string
	Attr     []string
	// attacks may have different boni
	PoolMod     int
	DcModDamage int
	DcModWC     int
	DcModStr    int
}

func (att *Attack) genText(creature *Creature) string {

	_, dc, _, _ := getSizeInfo(creature)
	dc += creature.damageClassMod

	dDice := getDCDice(dc + att.DcModDamage)
	weapVal := strconv.Itoa(getDCWeaponVal(dc + att.DcModWC))
	strBon := getDCStrBonus(dc + att.DcModStr)

	// get best attribut and default to dex
	attr := AttrDex
	for _, attAttr := range att.Attr {
		if creature.attributes[attAttr] > creature.attributes[attr] {
			attr = attAttr
		}
	}
	// generate AttackPool
	pool := creature.attributes[attr] + creature.skills[att.Skill] +
		att.PoolMod

	// generate string
	ret := att.Name + ": " + strconv.Itoa(pool) + ", " + dDice + "+"
	ret += strconv.Itoa((strBon * creature.attributes[AttrStr]) / 2)
	ret += " WV: " + weapVal + "."
	if att.AddDescr != "" {
		ret += " " + att.AddDescr
	}
	return ret
}

func loadAttacks() {
	dir, err := os.ReadDir(attackDir)
	if checkErr(err) {
		os.Exit(1)
	}

	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		raw, err := os.ReadFile(
			attackDir + string(os.PathSeparator) + file.Name())
		if checkErr(err) {
			continue
		}
		newAttack := new(Attack)
		newAttack.Name = file.Name()
		for _, line := range strings.Split(string(raw), newLine) {
			m := readAttackLine.FindStringSubmatch(line)
			if len(m) < 2 {
				continue
			}
			switch m[1] {
			case AttName:
				newAttack.Name = m[2]
			case AttDescr:
				newAttack.AddDescr = m[2]
			case AttSkill:
				newAttack.Skill = m[2]
			case AttAttr:
				newAttack.Attr = append(newAttack.Attr, m[2])
			case AttPoolMod:
				mod, _ := strconv.Atoi(m[2]) // #nosec
				newAttack.PoolMod = mod
			case AttModDam:
				mod, _ := strconv.Atoi(m[2]) // #nosec
				newAttack.DcModDamage = mod
			case AttModWC:
				mod, _ := strconv.Atoi(m[2]) // #nosec
				newAttack.DcModWC = mod
			case AttModStr:
				mod, _ := strconv.Atoi(m[2]) // #nosec
				newAttack.DcModStr = mod
			}
		}
		attackList[newAttack.Name] = newAttack
	}
}
