package main

import (
	"fmt"
	"strings"
)

type Perfdata struct {
	Name  string
	Value string
	Warn  string
	Crit  string
	Min   string
	Max   string
}

type PerfdataList []Perfdata

func (p Perfdata) String() (s string) {
	s = fmt.Sprintf("'%s'=%s;%s;%s;%s;%s", p.Name, p.Value, p.Warn, p.Crit, p.Min, p.Max)
	s = strings.TrimRight(s, ";")

	return
}

func (pl PerfdataList) String() (s string) {
	for _, p := range pl {
		s += p.String() + " "
	}

	s = strings.Trim(s, " ")

	return
}
