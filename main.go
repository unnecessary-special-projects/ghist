package main

import (
	"embed"

	"github.com/coderstone/ghist/cmd"
)

//go:embed skills/*.md
var skillsFS embed.FS

//go:embed all:web/dist
var webFS embed.FS

func main() {
	cmd.SetSkillsFS(skillsFS)
	cmd.SetWebFS(webFS)
	cmd.Execute()
}
