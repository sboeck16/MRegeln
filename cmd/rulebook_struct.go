package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

/*
Anythin added to rulebook is either anther json or a md file.
File ending are added automatically!
*/
type RuleBook struct {
	// paths
	DesignFile string `json:"design_file"`
	RulePath   string `json:"rule_path"`
	TargetFile string `json:"target_file"`

	// modules get .json file extension
	Modules []string `json:"modules"`

	// preScripts executes command before loading md files
	PreScript []string

	// markdown files
	Introduction []string `json:"intro"`
	Characters   []string `json:"characters"`
	Dices        []string `json:"dices"`
	Battle       []string `json:"battle_rules"`

	// hold names. That names are used to get corresponding
	// MD files! Stärke -> Attribute/Stärke.md
	Attributes []string            `json:"attributes"`
	Skills     map[string][]string `json:"skills"`
	Paths      map[string][]string `json:"paths"`
	Masteries  map[string][]string `json:"masteries"`

	// additional md files after. need full path
	Additional []string `json:"additional"`
}

func ReadRuleBook(path, file string) (*RuleBook, error) {
	// ensure .json ending
	file = EnsureFileEnding(file, JSONEnd)
	// try reading scriptfile
	data, err1 := os.ReadFile(file)
	if err1 != nil {
		var err2 error
		data, err2 = os.ReadFile(filepath.Join(path, file))
		if err2 != nil {
			return nil, fmt.Errorf(err1.Error() + " - " + err2.Error())
		}
	}
	ret := &RuleBook{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}

	// no ensure maps are none nil, empty slices are ok
	if ret.Skills == nil {
		ret.Skills = make(map[string][]string)
	}
	if ret.Paths == nil {
		ret.Paths = make(map[string][]string)
	}
	if ret.Masteries == nil {
		ret.Masteries = make(map[string][]string)
	}

	// set rulepath (unecessary?)
	ret.RulePath = path

	// run prescripts
	for _, cmdStr := range ret.PreScript {
		checkErr(exec.Command(cmdStr).Run())
	}

	// merge sub modules
	for _, subModule := range ret.Modules {
		sub, errS := ReadRuleBook(path, subModule)
		if !checkErr(errS) {
			ret.mergeSub(sub)
		}
	}

	return ret, nil
}

func (r *RuleBook) mergeSub(sub *RuleBook) {
	// sub pre script and subsub modules have been run already
	// Only set Design and Target if not already set (prefer main)
	if r.DesignFile == "" {
		r.DesignFile = sub.DesignFile
	}
	if r.TargetFile == "" {
		r.TargetFile = sub.TargetFile
	}

	// those should be defined in main module but there could be additions
	r.Introduction = append(r.Introduction, sub.Introduction...)
	r.Characters = append(r.Characters, sub.Characters...)
	r.Dices = append(r.Dices, sub.Dices...)
	r.Battle = append(r.Battle, sub.Battle...)
	r.Additional = append(r.Additional, sub.Additional...)

	// attributes should also be in main module
	r.Attributes = append(r.Attributes, sub.Attributes...)

	// join the real modular part
	for skill, mdfiles := range sub.Skills {
		r.Skills[skill] = append(r.Skills[skill], mdfiles...)
	}
	for path, mdfiles := range sub.Paths {
		r.Paths[path] = append(r.Paths[path], mdfiles...)
	}
	for mast, mdfiles := range sub.Masteries {
		r.Masteries[mast] = append(r.Masteries[mast], mdfiles...)
	}
}

func (r *RuleBook) getSortedSkillCategories() []string {
	ret := []string{}
	for skillCat := range r.Skills {
		ret = append(ret, skillCat)
	}
	sort.Strings(ret)
	return ret
}
