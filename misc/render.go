package misc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"text/template"

	"github.com/lovego/xiaomei/release"
	"github.com/spf13/cobra"
)

func renderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `render <env> <template file> [output file]`,
		Short: `render template with envionment config.`,
		RunE: release.Env2Call(func(env string, tmplFile, outputFile string) error {
			return RenderFileWithConfig(env, tmplFile, outputFile)
		}),
	}
	return cmd
}

func RenderFileWithConfig(env, tmplFile, outputFile string) error {
	if outputFile != "" {
		return RenderFileTo(tmplFile, release.Config(env), outputFile)
	}
	if output, err := RenderFile(tmplFile, release.Config(env)); err != nil {
		return err
	} else {
		fmt.Println(output.String())
		return nil
	}
}

func RenderFileTo(tmplFile string, data interface{}, outputFile string) error {
	output, err := RenderFile(tmplFile, data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outputFile, output.Bytes(), 0644)
}

func RenderFile(tmplFile string, data interface{}) (bytes.Buffer, error) {
	var buf bytes.Buffer
	content, err := ioutil.ReadFile(tmplFile)
	if err != nil {
		return buf, err
	}
	tmpl := template.New(``).Funcs(funcsMap)
	if _, err := tmpl.Parse(string(content)); err != nil {
		return buf, err
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return buf, err
	}
	return buf, nil
}

var funcsMap = template.FuncMap{
	`domainAncestor`:  DomainAncestor,
	`regexpQuoteMeta`: regexp.QuoteMeta,
}

// DomainAncestor return the N'th ancestor of domain
func DomainAncestor(domain string, n int) string {
	if n <= 0 {
		return domain
	}
	var index int
	for i, b := range domain {
		if b == '.' {
			index = i
			n--
			if n == 0 {
				break
			}
		}
	}
	return domain[index+1:]
}
