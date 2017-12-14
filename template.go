package hime

import (
	"html/template"
	"log"
	"path/filepath"
)

// TemplateFuncs adds template funcs
func (app *App) TemplateFuncs(funcs ...template.FuncMap) *App {
	app.templateFuncs = append(app.templateFuncs, funcs...)
	return app
}

// Component adds global template component
func (app *App) Component(filename ...string) *App {
	app.templateComponents = append(app.templateComponents, filename...)
	return app
}

// Template registers new template
func (app *App) Template(name string, filename ...string) *App {
	if _, ok := app.template[name]; ok {
		log.Panicf("hime: template %s already exist", name)
	}

	t := template.New("")

	// register funcs
	for _, fn := range app.templateFuncs {
		t.Funcs(fn)
	}

	// load templates and components
	fn := make([]string, len(filename))
	copy(fn, filename)
	fn = append(fn, app.templateComponents...)
	t = template.Must(t.ParseFiles(joinTemplateDir(app.templateDir, fn...)...))
	t = t.Lookup(app.templateRoot)

	app.template[name] = t

	return app
}

func joinTemplateDir(dir string, filenames ...string) []string {
	xs := make([]string, len(filenames))
	for i, filename := range filenames {
		xs[i] = filepath.Join(dir, filename)
	}
	return xs
}
