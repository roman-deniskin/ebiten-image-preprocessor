package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/urfave/cli/v2"
	"image"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	app := &cli.App{
		Name:  "image-preprocessor",
		Usage: "Preprocess images using ebiten",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Required: true,
				Usage:    "Path to the folder with images",
			},
		},
		Action: func(c *cli.Context) error {
			path := c.String("path")
			return preprocessImages(path)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func preprocessImages(path string) error {
	err := filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		_, img, err := ebitenutil.NewImageFromFile(file)
		if err != nil {
			return nil // пропускаем файл, если не удалось загрузить изображение
		}

		err = saveImageAsDat(img, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func saveImageAsDat(img image.Image, originalPath string) error {
	preparedDataPath := "prepared_data"
	err := os.MkdirAll(preparedDataPath, 0755)
	if err != nil {
		return err
	}

	newPath := filepath.Join(preparedDataPath, filepath.Base(originalPath)+".dat")
	file, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(img)
	if err != nil {
		return err
	}

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
