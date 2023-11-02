package mdcoach

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func visualizeYamlNode(w io.Writer, node interface{}) {
	switch n := node.(type) {
	case []interface{}: // lists
		w.Write([]byte(`<ol>`))
		for _, v := range n {
			w.Write([]byte(`<li>`))
			visualizeYamlNode(w, v)
			w.Write([]byte(`</li>`))
		}
		w.Write([]byte(`</ol>`))
	case []map[string](interface{}): // lists
		w.Write([]byte(`<ol>`))
		for _, v := range n {
			w.Write([]byte(`<li>`))
			visualizeYamlNode(w, v)
			w.Write([]byte(`</li>`))
		}
		w.Write([]byte(`</ol>`))
	case map[string]interface{}: // dictionary
		w.Write([]byte(`<ul>`))
		for k, v := range n {
			w.Write([]byte(`<li><strong>`))
			visualizeYamlNode(w, k)
			w.Write([]byte(`:</strong> `))
			visualizeYamlNode(w, v)
			w.Write([]byte(`</li>`))
		}
		w.Write([]byte(`</ul>`))
	default: // everything else
		// TODO: iterate through each byte instead and force UTF/escape
		escapeHTML(w, []byte(fmt.Sprintf(`%v`, node)))
	}
}

// YamlWithTemplate puts a yaml file through a template.
func YamlWithTemplate(w io.Writer, source string, templates ...string) error {
	temp := make([]map[string]interface{}, 0)
	yamlFile, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &temp)
	if err != nil {
		return fmt.Errorf(`failed to parse YAML from "%s": %s`, source, err.Error())
	}

	t := template.New(``).Funcs(template.FuncMap{
		`dump`: func(node ...interface{}) template.HTML {
			var b bytes.Buffer
			for i, n := range node {
				if i > 0 {
					fmt.Fprint(&b, `<hr />`)
				}
				visualizeYamlNode(&b, n)
			}
			return template.HTML(b.String())
		},
	})
	for _, tt := range templates {
		b, err := ioutil.ReadFile(tt)
		if err != nil {
			return err
		}
		t, err = t.Parse(string(b))
		if err != nil {
			return err
		}
	}
	if len(templates) == 0 { // Load default template.
		// TODO: make a proper template
		t.Parse(`<h1>Default Template</h1>{{ .| dump }}`)
	}
	return t.Execute(w, temp)
}
