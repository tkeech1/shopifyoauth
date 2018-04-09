package main

type EnvInterface interface {
	Getenv(string) string
}

type EnvHandler struct {
	Env EnvInterface
}
