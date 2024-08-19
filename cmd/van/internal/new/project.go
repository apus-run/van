package new

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/charmbracelet/huh"

	"github.com/apus-run/van/cmd/van/config"
)

// Project is a project template.
type Project struct {
	Name string
}

func NewProject(name string) *Project {
	return &Project{
		Name: name,
	}
}

func (p *Project) New(ctx context.Context, dir string, layout string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", p.Name)
		override := false
		e := huh.NewConfirm().
			Title("📂 Do you want to override the folder ?").
			Description("Delete the existing folder and create the project.").
			Affirmative("Yes!").
			Negative("No.").
			Value(&override).Run()
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		e = os.RemoveAll(to)
		if e != nil {
			fmt.Println("remove old project error: ", err)
			return e
		}
	}

	repo := ""
	if layout == "" {
		selected := ""
		err := huh.NewSelect[string]().
			Title("Please select a layout:").
			Options(
				huh.NewOptions("Basic", "Advanced", "Multiple")...,
			).Value(&selected).Run()
		if err != nil {
			return err
		}

		switch selected {
		case "Basic":
			repo = config.RepoBase
		case "Advanced":
			repo = config.RepoAdvanced
		case "Multiple":
			repo = config.RepoMultiple
		default:
			repo = config.RepoBase
		}

		err = os.RemoveAll(p.Name)
		if err != nil {
			fmt.Println("remove old project error: ", err)
			return err
		}
	} else {
		repo = layout
	}

	fmt.Printf("🚀 Creating service %s, layout repo is %s, please wait a moment.\n\n", p.Name, repo)

	fmt.Printf("git clone %s\n", repo)
	cmd := exec.Command("git", "clone", repo, p.Name)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("git clone %s error: %s\n", repo, err)
		return err
	}
	return true, nil
}
