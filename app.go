package main

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	nds "github.com/szerookii/desmume-launcher/nds"
	"github.com/szerookii/desmume-launcher/utils"
	"os"
	"path/filepath"
	"runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if !utils.FileExists(filepath.Join(".", "games")) {
		if err := os.Mkdir(filepath.Join(".", "games"), 0755); err != nil {
			fmt.Println("Cannot create games directory")
			os.Exit(0)
		} else {
			fmt.Println("Created games directory")
		}
	} else {
		fmt.Println("Games directory already exists")
	}

	if !utils.FileExists(filepath.Join(".", "desmume")) {
		if err := os.Mkdir(filepath.Join(".", "desmume"), 0755); err != nil {
			fmt.Println("Cannot create desmume directory")
			os.Exit(0)
		}

		fmt.Println("Created desmume directory")
		fmt.Println(utils.DownloadAndExtract(filepath.Join(".", "desmume")))
	} else {
		fmt.Println("Desmume directory already exists")
	}
}

func (a *App) ListGames() []*nds.NDSFile {
	files, err := utils.LoadGameFiles()
	if err != nil {
		log.Error(err)
		return []*nds.NDSFile{}
	}

	return files
}

func (a *App) LaunchGame(game *nds.NDSFile) bool {
	if utils.FileExists(filepath.Join(".", "games", game.Path)) {
		goos := runtime.GOOS
		var desumeExePath string
		var extension string

		if goos == "windows" {
			extension = ".exe"
		} else if goos == "darwin" {
			extension = ".dmg"
		}

		err := filepath.Walk(filepath.Join(".", "desmume"), func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == extension {
				desumeExePath = path
				return nil
			}

			return nil
		})

		if err != nil {
			return false
		}

		if desumeExePath == "" {
			return false
		}

		if err := utils.RunCommand(desumeExePath, []string{filepath.Join(".", "games", game.Path)}); err != nil {
			return false
		}

		return true
	}

	return false
}
