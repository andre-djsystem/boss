package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// BossPackage is a model of boss.json
type BossPackage struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Version      string      `json:"version"`
	Homepage     string      `json:"homepage"`
	MainSrc      string      `json:"mainsrc"`
	Projects     []string    `json:"projects"`
	Scripts      []string    `json:"scripts"`
	Dependencies interface{} `json:"dependencies"`
}

// MakeBossPackage create a new instance of BossPackage
func MakeBossPackage() *BossPackage {
	return &BossPackage{
		Scripts:  []string{},
		Projects: []string{},
	}
}

// LoadPackage open a boss.json
func LoadPackage(bossPath string) (*BossPackage, error) {
	buf, err := ioutil.ReadFile(bossPath)
	if err != nil {
		return nil, err
	}
	var bossPackage = MakeBossPackage()
	err = json.Unmarshal(buf, &bossPackage)
	if err != nil {
		return nil, err
	}
	return bossPackage, nil
}

// SaveToFile save changes of boss.json file
func (b *BossPackage) SaveToFile(bossPath string) ([]byte, error) {
	buf, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return buf, err
	}
	return buf, ioutil.WriteFile(bossPath, buf, os.ModePerm)
}

func (b *BossPackage) UninstallDependency(dependency string) {
	if b.Dependencies != nil {
		deps := b.Dependencies.(map[string]interface{})
		for key := range deps {
			if strings.ToLower(key) == strings.ToLower(dependency) {
				delete(deps, key)
			}
		}
		b.Dependencies = deps
	}
}

func (b *BossPackage) AddDependency(dependency string, version string) {
	if b.Dependencies == nil {
		b.Dependencies = make(map[string]interface{})
	}
	deps := b.Dependencies.(map[string]interface{})
	for key := range deps {
		if strings.ToLower(key) == strings.ToLower(dependency) {
			deps[key] = version
			return
		}
	}
	deps[dependency] = version
}
