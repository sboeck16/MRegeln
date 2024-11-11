package cmd

import "strconv"

const (
	// needs to be included in array down below
	AttrStr = "Stä"
	AttrDex = "Ges"
	AttrCon = "Kon"
	AttrInt = "Int"
	AttrWil = "Wil"
	AttrPer = "Wah"
	AttrCha = "Per"

	// needs to be included in array down below
	SkillMelee      = "Nahkampf"
	SkillMeleeShort = "Hiebwaffen"
	SkillMeleeLong  = "Stangenwaffen"
	SkillRanged     = "Fernkampf"
	SkillRangedBows = "Bögen"
	SkillRangedProj = "Projektilwaffen"
	SkillStealth    = "Heimlichkeit"
	skillAthletik   = "Sportlichkeit"
	skillSwim       = "Schwimmen"
	SkillPercept    = "Aufmerksamkeit"
	SkillMagic      = "Magie"

	// creature mod fields
	creatureSize     = "Größe"
	creatureDCMod    = "Waffenklasse"
	creatureType     = "Typ"
	creatureAttack   = "Angriff"
	creatureMovement = "Bewegungstyp"
	creatureArmor    = "Panzerung"

	// Attack fields
	AttName    = "Name"
	AttDescr   = "Beschreibung"
	AttSkill   = "Fertigkeit"
	AttAttr    = "Attribut"
	AttPoolMod = "Modifikation"
	AttModDam  = "Würfelmodifikation"
	AttModWC   = "Waffenwertmodifikation"
	AttModStr  = "Stärkemodifikation"

	creatureIsBeast  = 0
	creatureIsHuman  = 1
	creatureIsUndead = 2

	move2Legs = 0
	move4Legs = 1
	moveSwim  = 2
	moveFly   = 3
)

var (
	AttrList = []string{
		AttrStr, AttrDex, AttrCon, AttrInt, AttrWil, AttrPer, AttrCha}
	SkillList = []string{
		SkillMelee, SkillMeleeShort, SkillMeleeLong, SkillRanged, SkillRangedBows,
		SkillRangedProj, SkillStealth, skillAthletik, skillSwim, SkillPercept,
		SkillMagic,
	}

	// SIZE
	// ----
	sizeToName = map[int]string{
		1: "winzig",
		2: "sehr klein",
		3: "klein",
		4: "mickrig",
		5: "normal",
		6: "groß",
		7: "sehr groß",
		8: "riesig",
		9: "gigantisch",
	}
	sizeToDef = map[int]int{
		1: 16,
		2: 15,
		3: 15,
		4: 14,
		5: 14,
		6: 13,
		7: 12,
		8: 11,
		9: 10,
	}
	sizeToDCAndRows = map[int]int{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
		6: 7,
		7: 10,
		8: 15,
		9: 20,
	}
	sizeToRowMod = map[int]map[int]map[int]int{
		creatureIsBeast: {
			2: {2: -2},
			3: {2: -1, 3: -2},
			4: {2: -1, 4: -2},
			5: {3: -1, 4: -2, 5: -3},
			6: {3: -1, 5: -2, 7: -3},
			7: {4: -1, 7: -2, 9: -3},
			8: {5: -1, 9: -2, 13: -3},
			9: {7: -1, 13: -2, 18: -3},
		},
		creatureIsHuman: {
			4: {2: -1, 3: -2, 4: -4},
			5: {2: -1, 3: -2, 4: -3, 5: -4},
			6: {2: -1, 4: -2, 5: -3, 7: -4},
		},
	}

	// DAMAGE
	// ------
	damageClassToDice = map[int]string{
		0:  "W2",
		3:  "W3",
		5:  "W4",
		7:  "W6",
		9:  "W8",
		12: "2W6",
		15: "2W8",
		19: "2W10",
	}
	damageClassToWeaponVal = map[int]int{
		0:  0,
		3:  1,
		5:  1,
		7:  2,
		9:  2,
		12: 3,
		15: 4,
		19: 6,
	}
	// values are divided by 2
	damageClassToStrBonus = map[int]int{
		1:  1,
		4:  2,
		10: 3,
		16: 4,
	}
)

func getDCDice(dc int) string {
	for dc > 0 {
		if val, ok := damageClassToDice[dc]; ok {
			return val
		}
		dc--
	}
	return damageClassToDice[0]
}

func getDCWeaponVal(dc int) int {
	for dc > 0 {
		if val, ok := damageClassToWeaponVal[dc]; ok {
			return val
		}
		dc--
	}
	return damageClassToWeaponVal[0]
}
func getDCStrBonus(dc int) int {
	for dc > 0 {
		if val, ok := damageClassToStrBonus[dc]; ok {
			return val
		}
		dc--
	}
	return damageClassToStrBonus[0]
}

func getSizeInfo(cr *Creature) (string, int, map[int]int, int) {
	size := cr.baseSize + cr.sizeMod
	if size < 1 {
		size = 1
	}
	if size > 9 {
		size = 9
	}
	return sizeToName[size], sizeToDCAndRows[size],
		sizeToRowMod[cr.typ][size], 9 - size
}

func getMovementStr(cr *Creature, moveType int) string {
	ret := "Bewegung "
	if moveType == move2Legs {
		nr := 3 + (cr.attributes[AttrStr]+cr.attributes[AttrDex])/2
		nr += cr.moveMod
		ret += "(Land): " + strconv.Itoa(nr)
	}
	if moveType == move4Legs {
		nr := 4 + (3*(cr.attributes[AttrStr]+cr.attributes[AttrDex]))/4
		nr += cr.moveMod
		ret += "(Land): " + strconv.Itoa(nr)
	}
	if moveType == moveFly {
		nr := 6 + (cr.attributes[AttrStr] + cr.attributes[AttrDex])
		nr += cr.moveMod
		ret += "(Luft): " + strconv.Itoa(nr)
	}
	if moveType == moveSwim {
		nr := 1 + (cr.attributes[AttrStr]+cr.attributes[AttrDex])/4
		nr += cr.moveMod
		ret += "(Wasser): " + strconv.Itoa(nr)
	}

	return ret
}
