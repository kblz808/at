package main

type theme struct {
	name     string
	location string
}

func (t theme) FilterValue() string { return t.name }
