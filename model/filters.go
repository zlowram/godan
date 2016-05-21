package model

type Filters struct {
	Ip       string
	Ports    []string
	Services []string
	Regexp   string
}
