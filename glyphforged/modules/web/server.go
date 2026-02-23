package web

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// TODO: Page Data

// Server holds shared state for the web layer.
// - template: the parsed template tree
// - rootDir: absolute path to the project root
type Server struct {
	template *template.Template
	rootDir  string
}

// NewServer constructs a server instance.
// - Defines template helper functions
// - Figure out the project root dir
// - Parse template files
func NewServer() (*Server, error) {
	templateFuncMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"nowYear": func() int { return time.Now().Year() },
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get cwd: %w", err)
	}
	rootDir, err := discoverProjectRoot(workingDir)
	if err != nil {
		return nil, err
	}

	// Define the glob patterns for templates
	// Layouts, then partials, then pages
	templateGlobs := []string{
		filepath.Join(rootDir, "templates/layouts/*.gohtml"),
		filepath.Join(rootDir, "templates/partials/*gohtml"),
		filepath.Join(rootDir, "templates/pages/*.gohtml"),
	}

	// Create root template and attach helper functions
	template, err := template.
		New("site").
		Funcs(templateFuncMap).
		ParseGlob(templateGlobs[0])
	if err != nil {
		return nil, fmt.Errorf("parse layouts: %w", err)
	}

	// Parse remaining template groups into the tree
	// Each ParseGlob() appends a template
	for _, globPattern := range templateGlobs[1:] {
		if _, err := template.ParseGlob(globPattern); err != nil {
			return nil, fmt.Errorf("parse %s: %w", globPattern, err)
		}
	}

	return &Server{
		template: template,
		rootDir:  rootDir,
	}, nil
}

func (s *Server) Routes() http.Handler {
	router := http.NewServeMux()

	staticRoot := filepath.Join(s.rootDir, "static")
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticRoot))))

	// router.HandleFunc("/", s.home)
	// TODO: Handle remaining path handlers

	return router
}

// TODO: Create renderers (full and partial)
// TODO: Create path handlers

// Walk upward from a starting directory to find project root
// Looks for:
// - At least one file matching templates/layouts/*.gohtml
// AND
// - A static/ directory
//
// Stops after 5 layers or when reaching filesystem root
func discoverProjectRoot(start string) (string, error) {
	currentDir := filepath.Clean(start)

	for range 5 {
		layoutsGlob := filepath.Join(currentDir, "templates/layouts/*.gohtml")
		staticDir := filepath.Join(currentDir, "static")
		layoutMatches, _ := filepath.Glob(layoutsGlob)

		// Check if at least one layout template exists
		if len(layoutMatches) > 0 {
			// Then check for static directory
			if staticInfo, err := os.Stat(staticDir); err == nil && staticInfo.IsDir() {
				// Found a directory, return that
				return currentDir, nil
			}
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("project root not found from %q: expected templates/ and static/", start)
}
