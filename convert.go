package main

import (
	"encoding/xml"
	"io"
	"sort"
	"strings"

	"github.com/gertd/go-pluralize"
)

func convert(in io.ReadCloser) string {
	var sql struct {
		Table []struct {
			Name string `xml:"name,attr"`
			Row  []struct {
				Name     string `xml:"name,attr"`
				Datatype string `xml:"datatype"`
				Relation []struct {
					Table string `xml:"table,attr"`
				} `xml:"relation"`
			} `xml:"row"`
		} `xml:"table"`
	}
	if err := xml.NewDecoder(in).Decode(&sql); err != nil {
		return ""
	}

	gs := []string{}
	for _, table := range sql.Table {
		g := "rails g model " + modelName(table.Name)
		for _, row := range table.Row {
			if row.Name == "id" {
				continue
			}

			if len(row.Relation) == 0 {
				g += " " + rowName(row.Name) + ":" + typeAttr(row.Datatype)
			} else {
				for _, r := range row.Relation {
					g += " " + referenceName(r.Table) + ":references"
				}
			}
		}
		gs = append(gs, g)
	}

	commandMap := make(map[string]string)
	dependenciesMap := make(map[string][]string)
	models := []string{}

	for _, command := range gs {
		model, dependencies := parseCommand(command)
		commandMap[model] = command
		dependenciesMap[model] = dependencies
		models = append(models, model)
	}

	sort.Slice(models, func(i, j int) bool {
		for _, dep := range dependenciesMap[models[j]] {
			if models[i] == dep {
				return true
			}
		}
		return false
	})

	res := ""
	for _, model := range models {
		res += commandMap[model] + "\n"
	}
	return res
}

func typeAttr(t string) string {
	switch t {
	case "INTEGER":
		return "integer"
	case "VARCHAR":
		return "string"
	case "DECIMAL":
		return "float"
	case "TEXT":
		return "text"
	case "DATE":
		return "date"
	default:
		return "string"
	}
}

func modelName(m string) string {
	client := pluralize.NewClient()
	singularName := client.Singular(m)
	singularNameWithUnderscores := strings.ReplaceAll(singularName, " ", "_")
	parts := strings.Split(singularNameWithUnderscores, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func referenceName(m string) string {
	client := pluralize.NewClient()
	singularName := client.Singular(m)
	singularName = strings.ReplaceAll(singularName, " ", "_")
	return singularName
}

func rowName(r string) string {
	return strings.ReplaceAll(r, " ", "_")
}

func parseCommand(command string) (string, []string) {
	parts := strings.Split(command, " ")
	model := strings.ToLower(parts[3])
	dependencies := []string{}
	for _, part := range parts[4:] {
		if strings.Contains(part, ":references") {
			dependencies = append(dependencies, strings.Split(part, ":")[0])
		}
	}
	return model, dependencies
}
