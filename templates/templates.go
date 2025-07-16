package templates

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

func LoadTemplates() (*template.Template, error) {
	templates := template.New("")

	// Define template functions
	funcMap := template.FuncMap{
		"formatSize": formatFileSize,
		"formatTime": formatTime,
		"split":      strings.Split,
		"hasSuffix":  strings.HasSuffix,
	}
	templates.Funcs(funcMap)

	// Load templates
	_, err := templates.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}
