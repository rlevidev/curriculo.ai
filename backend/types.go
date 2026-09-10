// Domain types for the resume — mirrors the JSON for POST /generate-pdf.
package main

type Education struct {
	Institution string   `json:"institution"`
	Degree      string   `json:"degree"`
	Period      string   `json:"period"`
	Notes       []string `json:"notes"`
}

type Experience struct {
	Company  string   `json:"company"`
	Location string   `json:"location"`
	Role     string   `json:"role"`
	Period   string   `json:"period"`
	Bullets  []string `json:"bullets"`
}

type Project struct {
	Name      string   `json:"name"`
	LinkURL   string   `json:"link_url"`
	LinkLabel string   `json:"link_label"`
	Bullets   []string `json:"bullets"`
}

type Language struct {
	Language string `json:"language"`
	Level    string `json:"level"`
}

type Skills struct {
	Languages    []string `json:"languages"`
	Technologies []string `json:"technologies"`
}

type ResumeData struct {
	Name            string       `json:"name"`
	Title           string       `json:"title"`
	Email           string       `json:"email"`
	Phone           string       `json:"phone"`
	LinkedIn        string       `json:"linkedin"`
	GitHub          string       `json:"github"`
	Location        string       `json:"location"`
	Education       []Education  `json:"education"`
	Experiences     []Experience `json:"experiences"`
	Projects        []Project    `json:"projects"`
	SpokenLanguages []Language   `json:"spokenLanguages"`
	Certifications  []string     `json:"certifications"`
	Skills          Skills       `json:"skills"`
}
