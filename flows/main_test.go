package main

import (
	"testing"

	"github.com/matryer/is"
)

func TestToEnc(t *testing.T) {
	is := is.New(t)
	is.Equal(toEnc("../flows.secret/team.md"), "../flows.secret/team.enc.md")
}

func TestUnEnc(t *testing.T) {
	is := is.New(t)
	is.Equal(unEnc("../flows.secret/team.enc.md"), "../flows.secret/team.md")
}

func TestUnsecretDir(t *testing.T) {
	is := is.New(t)
	is.Equal(unsecretDir("../flows.secret/team.enc.md"), "../flows/team.enc.md")
}

func TestUnsecretPath(t *testing.T) {
	is := is.New(t)
	is.Equal(unsecretDir(toEnc("../flows.secret/team.md")), "../flows/team.enc.md")
}

func TestSecretDir(t *testing.T) {
	is := is.New(t)
	is.Equal(secretDir("../flows/team.md"), "../flows.secret/team.md")
}
