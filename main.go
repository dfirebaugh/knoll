package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dfirebaugh/knoll/web"
	"gopkg.in/yaml.v2"
)

type Attribute struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Default string `yaml:"default"`
}

type Property struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Default string `yaml:"default"`
}

type Event struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Element struct {
	Name        string            `yaml:"name"`
	Tag         string            `yaml:"tag"`
	Attributes  []Attribute       `yaml:"attributes"`
	Properties  []Property        `yaml:"properties"`
	Events      []Event           `yaml:"events"`
	ExampleData map[string]string `yaml:"exampleData"`
	Script      string            `yaml:"script"`
}

type Config struct {
	Scripts  []string  `yaml:"scripts"`
	Elements []Element `yaml:"elements"`
}

func renderElement(element Element) template.HTML {
	openTag := fmt.Sprintf("<%s", element.Tag)
	for attr, value := range element.ExampleData {
		openTag += fmt.Sprintf(` %s="%s"`, attr, value)
	}
	openTag += ">"
	closeTag := fmt.Sprintf("</%s>", element.Tag)
	return template.HTML(openTag + closeTag)
}

func toJSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func copyFile(srcFS embed.FS, src, dstDir, relPath string) error {
	input, err := srcFS.ReadFile(src)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(dstDir, relPath)
	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(outputPath, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func copyDir(srcFS embed.FS, srcDir, dstDir, relPath string) error {
	entries, err := srcFS.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstRelPath := filepath.Join(relPath, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcFS, srcPath, dstDir, dstRelPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcFS, srcPath, dstDir, dstRelPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyScript(src, dstDir string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", src, err)
	}

	relPath, err := filepath.Rel(".", src)
	if err != nil {
		return fmt.Errorf("error getting relative path for %s: %v", src, err)
	}

	outputPath := filepath.Join(dstDir, relPath)
	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return fmt.Errorf("error creating directories for %s: %v", outputPath, err)
	}

	err = os.WriteFile(outputPath, input, 0644)
	if err != nil {
		return fmt.Errorf("error writing %s to output directory: %v", outputPath, err)
	}

	return nil
}

func serve(outputDir string, port int) {
	fs := http.FileServer(http.Dir(outputDir))

	http.Handle("/", fs)

	fmt.Printf("Serving files from the '%s' directory on http://localhost:%d\n", outputDir, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func main() {
	// Parse command-line flags
	serveFlag := flag.Bool("serve", false, "Serve the output directory with a web server")
	outputDirFlag := flag.String("output", "output", "Directory to write the output files")
	portFlag := flag.Int("port", 8080, "Port to serve the web server on")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("Please provide a YAML configuration file.")
	}

	configFile := flag.Arg(0)
	yamlData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}

	funcMap := template.FuncMap{
		"renderElement": renderElement,
		"toJSON":        toJSON,
	}

	tmpl, err := template.New("gallery").Funcs(funcMap).ParseFS(web.WebFS, "templates/gallery.tmpl", "templates/element.tmpl")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	outputDir := *outputDirFlag
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	// Copy the contents of the web directory to the output directory
	err = copyDir(web.WebFS, "css", outputDir, "")
	if err != nil {
		log.Fatalf("Error copying web directory: %v", err)
	}

	outputFile, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	err = tmpl.ExecuteTemplate(outputFile, "gallery", config)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	for _, script := range config.Scripts {
		err = copyScript(script, outputDir)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	for _, element := range config.Elements {
		if element.Script != "" {
			err = copyScript(element.Script, outputDir)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
	}

	fmt.Println("Static site generated successfully!")

	if *serveFlag {
		serve(outputDir, *portFlag)
	}
}
