package types

import "html/template"

type HomeInternetPage struct {
	Hero      HeroSection
	Features  []Feature
	Packages  []InternetPackage
	Technical TechnicalSection
	FAQs      []QAItem
}

type HeroSection struct {
	Badge        string
	Title        string
	TitleItalic  string
	Text         string
	PrimaryCTA   Link
	SecondaryCTA Link
}

type Feature struct {
	Title       string
	Description string
	Icon        template.HTML
}

type InternetPackage struct {
	Name     string
	Speed    string
	Price    string
	Popular  bool
	Features []string
}

type TechnicalSection struct {
	Title        string
	Description  string
	Stats        []Stat
	ImageURL     string
	BulletPoints []string
}

type Stat struct {
	Value string
	Label string
}

type Link struct {
	Text string
	URL  string
}
