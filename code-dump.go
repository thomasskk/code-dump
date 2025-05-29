package main

import (
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	allowedExtensions    map[string]struct{}
	allowedFilenames     map[string]struct{}
	skipDirectories      map[string]struct{}
	skipDirectoryRegexes []*regexp.Regexp
	skipFileRegexes      []*regexp.Regexp
)

type FileInfo struct {
	Size         int64
	LastModified string
}

func init() {

	allowedExtensionsList := []string{
		".txt", ".md", ".markdown", ".json", ".xml", ".yaml", ".yml", ".toml",
		".ini", ".cfg", ".conf", ".sql", ".graphql", ".proto",
		".py", ".pyx", ".pyd", ".pyo", ".pyc", ".pyw", ".pyi",
		".c", ".h", ".i", ".cpp", ".hpp", ".cc", ".hh", ".cxx", ".hxx",
		".jl",
		".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs",
		".html", ".htm", ".css", ".scss", ".sass", ".less",
		".java", ".kt", ".kts", ".groovy", ".scala", ".clj", ".cljs",
		".cs", ".fs", ".vb",
		".rb", ".rake", ".gemspec",
		".php", ".phtml", ".php3", ".php4", ".php5", ".phps",
		".go",
		".rs",
		".swift",
		".sh", ".bash", ".zsh", ".fish",
		".ps1", ".psm1", ".psd1",
		".pl", ".pm",
		".lua",
		".hs", ".lhs",
		".r", ".R", ".Rmd",
		".dart",
		".m", ".mm",
		".elm",
		".ex", ".exs",
		".erl", ".hrl",
		".lisp", ".cl", ".el",
		".f", ".for", ".f90", ".f95", ".f03", ".f08",
		".tf", ".tfvars",
		".tex", ".sty", ".cls",
	}
	allowedExtensions = make(map[string]struct{}, len(allowedExtensionsList))
	for _, ext := range allowedExtensionsList {
		allowedExtensions[ext] = struct{}{}
	}

	allowedExtensions[".mat"] = struct{}{}

	allowedFilenamesList := []string{
		"readme", "license", "dockerfile", "makefile", ".gitignore", ".dockerignore",
		".editorconfig", ".env", "requirements.txt", "package.json", "tsconfig.json",
		"setup.py", "setup.cfg", "pyproject.toml", "pipfile", "manifest.in",
		".pylintrc", ".flake8", "pytest.ini", "tox.ini",
		"cmakelists.txt",
		"project.toml", "manifest.toml", "juliaconfig.toml",
		".npmignore", ".babelrc", ".eslintrc", ".prettierrc",
		"tslint.json", "webpack.config.js", "package-lock.json", "yarn.lock",
		"gemfile", "rakefile",
		"composer.json", "composer.lock",
		"go.mod", "go.sum",
		"cargo.toml", "cargo.lock",
		"packages.config", "nuget.config",
		"pom.xml", "build.gradle", "build.gradle.kts", "settings.gradle", "settings.gradle.kts",
		"docker-compose.yml", "docker-compose.yaml",
		".gitattributes",
		".travis.yml", ".gitlab-ci.yml", "jenkinsfile",
		"azure-pipelines.yml",
		"elm.json",
		"paket.dependencies", "paket.lock",
		"mix.exs", "mix.lock",
		"rebar.config",
		".octaverc",
		"build.sbt",
		".terraform.lock.hcl",
		"ansible.cfg", "hosts",
		"latexmkrc",
	}

	allowedFilenames = make(map[string]struct{}, len(allowedFilenamesList))
	for _, name := range allowedFilenamesList {
		allowedFilenames[strings.ToLower(name)] = struct{}{}
	}

	skipDirectoriesList := []string{
		"__pycache__", "node_modules", "venv", "env", ".venv", ".env",
		"build", "dist", "target", "out", "bin", "obj",
		".git", ".svn", ".hg",
		".idea", ".vscode",
		"logs",
		"output",
	}
	skipDirectories = make(map[string]struct{}, len(skipDirectoriesList))
	for _, dir := range skipDirectoriesList {
		skipDirectories[dir] = struct{}{}
	}

	skipDirectoryPatternsStr := []string{
		`\.egg-info$`,
	}
	skipDirectoryRegexes = make([]*regexp.Regexp, len(skipDirectoryPatternsStr))
	for i, pattern := range skipDirectoryPatternsStr {
		re, err := regexp.Compile(pattern)
		if err != nil {
			log.Fatalf("Invalid skip directory regex pattern %s: %v", pattern, err)
		}
		skipDirectoryRegexes[i] = re
	}

	skipFilePatternsStr := []string{
		`\.log(\.[0-9]+)?$`,
		`^log\.`,
		`\.bak$`,
		`\.tmp$`,
		`\.temp$`,
		`\.swp$`,
		`~$`,
	}
	skipFileRegexes = make([]*regexp.Regexp, len(skipFilePatternsStr))
	for i, pattern := range skipFilePatternsStr {
		re, err := regexp.Compile(pattern)
		if err != nil {
			log.Fatalf("Invalid skip file regex pattern %s: %v", pattern, err)
		}
		skipFileRegexes[i] = re
	}
}

func getFileInfo(filePath string) (FileInfo, error) {
	stats, err := os.Stat(filePath)
	if err != nil {
		return FileInfo{}, err
	}
	return FileInfo{
		Size:         stats.Size(),
		LastModified: stats.ModTime().Format("2006-01-02 15:04:05"),
	}, nil
}

func shouldSkip(path string, isDir bool) bool {
	name := filepath.Base(path)
	nameLower := strings.ToLower(name)

	if isDir {
		if strings.HasPrefix(name, ".") {
			if _, found := allowedFilenames[nameLower]; !found {
				return true
			}
		}
		if _, found := skipDirectories[name]; found {
			return true
		}
		for _, pattern := range skipDirectoryRegexes {
			if pattern.MatchString(name) {
				return true
			}
		}
		return false
	}

	ext := strings.ToLower(filepath.Ext(name))

	for _, pattern := range skipFileRegexes {
		if pattern.MatchString(nameLower) {
			return true
		}
	}

	if strings.HasPrefix(name, ".") {
		if _, found := allowedFilenames[nameLower]; !found {
			return true
		}
	}

	_, extAllowed := allowedExtensions[ext]
	_, nameAllowed := allowedFilenames[nameLower]

	if ext != "" {
		if !extAllowed && !nameAllowed {
			return true
		}
	} else {
		if !nameAllowed {
			return true
		}
	}

	return false
}

func concatenateFiles(directory string, listOnly bool) (string, error) {
	var output strings.Builder

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {

			fmt.Fprintf(os.Stderr, "Error accessing path %s: %v\n", path, err)
			if d == nil {
				return err
			}
			return nil
		}

		if path == directory && d.IsDir() {
			return nil
		}

		if shouldSkip(path, d.IsDir()) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			info, err := getFileInfo(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting file info for %s: %v\n", path, err)
				return nil
			}

			displayPath := filepath.ToSlash(path)

			if listOnly {
				output.WriteString(displayPath + "\n")
			} else {
				separator := strings.Repeat("=", 80)
				output.WriteString("\n\n" + separator + "\n")
				output.WriteString(fmt.Sprintf("File: %s\n", displayPath))
				output.WriteString(fmt.Sprintf("Size: %d bytes\n", info.Size))
				output.WriteString(fmt.Sprintf("Last Modified: %s\n", info.LastModified))
				output.WriteString(separator + "\n\n")

				content, err := os.ReadFile(path)
				if err != nil {
					output.WriteString(fmt.Sprintf("Error reading file: %v\n", err))
				} else {
					if strings.ToValidUTF8(string(content), "") != string(content) {
						output.WriteString(fmt.Sprintf("Content (binary or non-UTF-8, %d bytes)\n", len(content)))
					} else {
						output.WriteString(string(content))
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error walking directory %s: %w", directory, err)
	}

	return output.String(), nil
}

func main() {

	listOnly := flag.Bool("l", false, "Only list file paths without content")

	flag.BoolVar(listOnly, "list-only", false, "Only list file paths without content (long form)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [directory]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Concatenate files in a directory.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nIf [directory] is not provided, it defaults to the current directory ('.').\n")
	}

	flag.Parse()

	rootDir := "."
	if flag.NArg() > 0 {
		rootDir = flag.Arg(0)
	}

	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		log.Fatalf("Error getting absolute path for %s: %v", rootDir, err)
	}

	result, err := concatenateFiles(absRootDir, *listOnly)
	if err != nil {
		log.Fatalf("Error processing directory: %v", err)
	}

	if result != "" {
		fmt.Print(result)
		err = clipboard.WriteAll(result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n\nError copying to clipboard: %v\n", err)
			fmt.Fprintf(os.Stderr, "You might need to install a clipboard tool (e.g., xclip/xsel on Linux, pbcopy on macOS).\n")
		} else {
			fmt.Fprintf(os.Stderr, "\n\nOutput for directory '%s' has been copied to clipboard.\n", absRootDir)
		}
	} else {
		fmt.Fprintf(os.Stderr, "No files found or processed in directory '%s' based on criteria.\n", absRootDir)
	}
}
