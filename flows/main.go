package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/orsenkucher/nothing/encio"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()
	sugar := logger.Sugar()

	if err := run(sugar); err != nil {
		sugar.Fatalf("Application failed: %v", err)
	}
}

func run(sugar *zap.SugaredLogger) error {
	var s = flag.String("s", "", "provide encio password")
	flag.Parse()
	if *s == "" {
		return errors.New("[-s] -> encio must be handled")
	}

	key := encio.NewEncIO(*s)
	err := classifySecrets(sugar, key)
	if err != nil {
		return err
	}

	err = revealSecrets(sugar, key)
	if err != nil {
		return err
	}

	return nil
}

func revealSecrets(sugar *zap.SugaredLogger, key encio.EncIO) error {
	dir := "../flows"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			sugar.Infow("Skipped directory", "dir", file.Name())
			continue
		}

		if !strings.Contains(file.Name(), ".enc.md") {
			sugar.Infow("Skipped file", "file", file.Name())
			continue
		}

		err = reveal(sugar, key, path.Join(dir, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

func reveal(sugar *zap.SugaredLogger, key encio.EncIO, path string) error {
	file, err := key.ReadFile(unEnc(path))
	if err != nil {
		return err
	}

	sec := unEnc(secretDir(path))
	err = ioutil.WriteFile(sec, file, 0644)
	if err != nil {
		return err
	}

	sugar.Infow("Revealing",
		"file", path,
		"onto", sec,
		"size", len(file),
	)

	return nil
}

func classifySecrets(sugar *zap.SugaredLogger, key encio.EncIO) error {
	dir := "../flows.secret"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModeDir)
		if err != nil {
			return err
		}
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			sugar.Infow("Skipped directory", "dir", file.Name())
			continue
		}

		if path.Ext(file.Name()) != ".md" {
			sugar.Infow("Skipped file", "file", file.Name())
			continue
		}

		err = classify(sugar, key, path.Join(dir, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

func classify(sugar *zap.SugaredLogger, key encio.EncIO, path string) error {
	file, err := key.ReadFile(path)
	if err != nil {
		return err
	}

	moved := unsecretDir(toEnc(path))
	err = os.Rename(toEnc(path), moved)
	if err != nil {
		return err
	}

	sugar.Infow("Classifying",
		"file", path,
		"onto", moved,
		"size", len(file),
	)

	return nil
}

func toEnc(name string) string {
	ext := path.Ext(name)
	result := strings.Replace(name, ext, ".enc"+ext, 1)
	return result
}

func secretDir(name string) string {
	dir := path.Dir(name) + ".secret"
	return path.Join(dir, path.Base(name))
}

func unsecretDir(name string) string {
	dir := path.Dir(name)
	dir = strings.Replace(dir, ".secret", "", 1)
	return path.Join(dir, path.Base(name))
}

func unEnc(name string) string {
	base := strings.Replace(path.Base(name), ".enc", "", 1)
	return path.Join(path.Dir(name), base)
}
