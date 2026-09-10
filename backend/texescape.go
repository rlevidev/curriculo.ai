package main

import "strings"

var texReplacer = strings.NewReplacer(
	"\\", "\\textbackslash{}",
	"&", "\\&",
	"%", "\\%",
	"$", "\\$",
	"#", "\\#",
	"_", "\\_",
	"{", "\\{",
	"}", "\\}",
	"~", "\\textasciitilde{}",
	"^", "\\textasciicircum{}",
)

// texEscape escapes LaTeX special characters using a precompiled global replacer.
func texEscape(s string) string {
	return texReplacer.Replace(s)
}

// escapeResumeData applies texEscape to every field of ResumeData (value copy).
func escapeResumeData(d ResumeData) ResumeData {
	d.Name = texEscape(d.Name)
	d.Title = texEscape(d.Title)
	d.Email = texEscape(d.Email)
	d.Phone = texEscape(d.Phone)
	d.LinkedIn = texEscape(d.LinkedIn)
	d.GitHub = texEscape(d.GitHub)
	d.Location = texEscape(d.Location)
	for i := range d.Education {
		d.Education[i].Institution = texEscape(d.Education[i].Institution)
		d.Education[i].Degree = texEscape(d.Education[i].Degree)
		d.Education[i].Period = texEscape(d.Education[i].Period)
		for j := range d.Education[i].Notes {
			d.Education[i].Notes[j] = texEscape(d.Education[i].Notes[j])
		}
	}
	for i := range d.Experiences {
		d.Experiences[i].Company = texEscape(d.Experiences[i].Company)
		d.Experiences[i].Location = texEscape(d.Experiences[i].Location)
		d.Experiences[i].Role = texEscape(d.Experiences[i].Role)
		d.Experiences[i].Period = texEscape(d.Experiences[i].Period)
		for j := range d.Experiences[i].Bullets {
			d.Experiences[i].Bullets[j] = texEscape(d.Experiences[i].Bullets[j])
		}
	}
	for i := range d.Projects {
		d.Projects[i].Name = texEscape(d.Projects[i].Name)
		d.Projects[i].LinkURL = texEscape(d.Projects[i].LinkURL)
		d.Projects[i].LinkLabel = texEscape(d.Projects[i].LinkLabel)
		for j := range d.Projects[i].Bullets {
			d.Projects[i].Bullets[j] = texEscape(d.Projects[i].Bullets[j])
		}
	}
	for i := range d.SpokenLanguages {
		d.SpokenLanguages[i].Language = texEscape(d.SpokenLanguages[i].Language)
		d.SpokenLanguages[i].Level = texEscape(d.SpokenLanguages[i].Level)
	}
	for i := range d.Certifications {
		d.Certifications[i] = texEscape(d.Certifications[i])
	}
	for i := range d.Skills.Languages {
		d.Skills.Languages[i] = texEscape(d.Skills.Languages[i])
	}
	for i := range d.Skills.Technologies {
		d.Skills.Technologies[i] = texEscape(d.Skills.Technologies[i])
	}
	return d
}
