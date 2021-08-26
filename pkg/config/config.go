/*
Package config provides application configuration; it should not import any other
packages outside those required from the standard library
*/
package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

// AppConfig holds the application configuration
// It defines an element TemplateCache that is of type map of strings
// pointing to the Template type from the text/template package in the standard library
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}
