package cmd

const (
	AttrStr = "Stä"
	AttrDex = "Ges"
	AttrCon = "Kon"
	AttrInt = "Int"
	AttrWil = "Wil"
	AttrPer = "Wah"
	AttrCha = "Per"

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

	creatureIsBeast  = 0
	creatureIsHuman  = 1
	creatureIsUndead = 2
)

var (
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
			4: {2: -1, 3: -2, 4: -3, 5: -4},
			5: {2: -1, 3: -2, 4: -3, 5: -4},
			6: {2: -1, 3: -2, 4: -3, 5: -4},
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
		3:  2,
		5:  2,
		7:  2,
		9:  3,
		12: 3,
		15: 4,
		19: 4,
	}
)
