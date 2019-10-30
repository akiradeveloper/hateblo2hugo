package hugo

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/catatsuy/movabletype"
)

const pageTpl = `
+++
date = "{{.Date}}"
draft = {{.Draft}}
title = "{{.Title}}"
tags = [{{ StringsJoin .Tags "," }}]

+++
{{.Content}}
`

type HugoPage struct {
	Date    string
	Draft   bool
	Title   string
	Tags    []string
	Content string
}

func CreateHugoPage(entry *movabletype.Entry) HugoPage {

	tags := make([]string, len(entry.Category))
	for i, s := range entry.Category {
		tags[i] = fmt.Sprintf(`"%s"`, s)
	}

	du, _ := time.ParseDuration("-9h")
	d := entry.Date.Add(du)
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	dJST := d.In(jst)

	return HugoPage{
		Date:    dJST.Format(time.RFC822Z),
		Draft:   entry.Status != "Publish",
		Title:   entry.Title,
		Tags:    tags,
		Content: entry.Body,
	}
}

func (p *HugoPage) Render() ([]byte, error) {

	// strings.LastIndex
	tpl := template.Must(template.New("hugo").Funcs(template.FuncMap{"StringsJoin": strings.Join}).Parse(pageTpl))

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	if err := tpl.Execute(writer, *p); err != nil {
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
