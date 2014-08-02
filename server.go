package main

import (
	"fmt"
	"strings"
	"net/http"
	"github.com/go-martini/martini"
	"github.com/codegangsta/martini-contrib/render"
	
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == a {
			return true
		}
	}
	return false
}

func sourcesForm(r render.Render, params martini.Params) {
	type CheckBox struct {
		Name string
		Checked bool
	}
	type Form struct {
		Releases []CheckBox
		Sources  []CheckBox
	}
	formdata := Form{
		[]CheckBox{
			{"Stable", true},
			{"Testing", false},
			{"Unstable", false},
			{"Experimental", false},
		},
		[]CheckBox{
			{"Main", true},
			{"Contrib", true},
			{"Non-Free", false},
		},
	}
	r.HTML(200, "debform", formdata)
}

func processForm(r render.Render, req *http.Request, params martini.Params)(string) {
	req.ParseForm()
	x := req.Form
	loc := x["location"][0]
	ref_lines := []string{}
	source_lines := []string{}
	line := ""

	package_sources := strings.Join(x["source"], " ")
	package_sources = strings.ToLower(package_sources)

	for _, release := range x["release"] {
		release = strings.ToLower(release)

		if StringInSlice("updates", x["updates"]) {
			line = fmt.Sprintf("http://ftp.%s.debian.org/debian %s-updates %s", loc, release, package_sources)
			ref_lines = append(ref_lines, line)

			if StringInSlice("source_pkgs", x["source_pkgs"]){
				line = fmt.Sprintf("http://ftp.%s.debian.org/debian %s-updates %s", loc, release, package_sources)
				ref_lines = append(ref_lines, line)
				
			}
		}
		
		line = fmt.Sprintf("http://ftp.%s.debian.org/debian %s %s", loc, release, package_sources)
		ref_lines = append(ref_lines, line)
		
		if release != "experimental" {
			line = fmt.Sprintf("http://security.debian.org/ %s/updates %s", release, package_sources)
			ref_lines = append(ref_lines, line)
		}
		
	}
	for _, line := range ref_lines {
		if StringInSlice("source_pkgs", x["source_pkgs"]){
			source_lines = append(source_lines, "deb-src " + line)
		}
		source_lines = append(source_lines, "deb " + line)
	}

	v := string(strings.Join(source_lines, "\n"))
	return v
	
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	// Show default menu
	m.Get("/", sourcesForm)

	m.Post("/", processForm)

	// go
	m.Run()
}
