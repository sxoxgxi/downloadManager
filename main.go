package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Config struct {
	parentDir   string
	directories map[string]string
	extensions  map[string]map[string]struct{}
}

func NewConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	parentDir := filepath.Join(homeDir, "Downloads")
	dirs := map[string]string{
		"audio":      filepath.Join(parentDir, "Audios"),
		"video":      filepath.Join(parentDir, "Videos"),
		"image":      filepath.Join(parentDir, "Images"),
		"document":   filepath.Join(parentDir, "Documents"),
		"datafile":   filepath.Join(parentDir, "Datafiles"),
		"executable": filepath.Join(parentDir, "Executables"),
		"other":      filepath.Join(parentDir, "Others"),
	}

	extensions := make(map[string]map[string]struct{})
	extensions["audio"] = makeExtensionSet([]string{
		".3ga", ".aac", ".ac3", ".aif", ".aiff", ".alac", ".amr", ".ape", ".au", ".dss",
		".flac", ".flv", ".m4a", ".m4b", ".m4p", ".mp3", ".mpga", ".ogg", ".oga", ".mogg",
		".opus", ".qcp", ".tta", ".voc", ".wav", ".wma", ".wv", ".ra", ".mid", ".midi",
		".pcm", ".aiff", ".aifc",
	})
	extensions["video"] = makeExtensionSet([]string{
		".webm", ".mts", ".m2ts", ".ts", ".mov", ".mp4", ".m4p", ".m4v", ".mxf", ".avi",
		".mkv", ".flv", ".wmv", ".rmvb", ".3gp", ".vob", ".ogv", ".3g2",
	})
	extensions["image"] = makeExtensionSet([]string{
		".jpg", ".jpeg", ".jfif", ".pjpeg", ".pjp", ".png", ".gif", ".webp", ".svg", ".apng",
		".avif", ".bmp", ".ico", ".tiff", ".tif", ".heif", ".heic", ".raw", ".nef", ".cr2",
		".orf", ".sr2", ".arw", ".rw2", ".dng", ".eps",
	})
	extensions["document"] = makeExtensionSet([]string{
		".doc", ".docx", ".log", ".msg", ".odt", ".pages", ".rtf", ".tex", ".txt", ".wpd",
		".wps", ".indo", ".pct", ".pdf", ".xls", ".xlsx", ".csv", ".ppt", ".pptx", ".key",
		".odp", ".ods", ".epub", ".mobi", ".ibooks", ".chm", ".xps",
	})
	extensions["datafile"] = makeExtensionSet([]string{
		".csv", ".dat", ".ged", ".key", ".keychain", ".ppt", ".pptx", ".sdf", ".tar",
		".tax2016", ".tax2020", ".tax2021", ".vcf", ".xml", ".zip", ".rar", ".7z", ".gz",
		".bz2", ".xz", ".iso", ".cab", ".dmg", ".tgz", ".pkg", ".rpm", ".deb",
	})
	extensions["executable"] = makeExtensionSet([]string{
		".bat", ".bin", ".cmd", ".com", ".cpl", ".ex_", ".exe", ".gadget", ".inf1", ".ins",
		".inx", ".isu", ".job", ".jse", ".lnk", ".msc", ".msi", ".msp", ".mst", ".paf",
		".pif", ".ps1", ".reg", ".rgs", ".scr", ".sct", ".shb", ".shs", ".ws", ".wsf",
		".wsh", ".py", ".jar", ".app", ".apk", ".run", ".sh", ".command",
	})

	return &Config{
		parentDir:   parentDir,
		directories: dirs,
		extensions:  extensions,
	}, nil
}

func makeExtensionSet(extensions []string) map[string]struct{} {
	set := make(map[string]struct{}, len(extensions))
	for _, ext := range extensions {
		set[strings.ToLower(ext)] = struct{}{}
	}
	return set
}

type FileOrganizer struct {
	config *Config
	logger *log.Logger
}

func NewFileOrganizer(config *Config) *FileOrganizer {
	return &FileOrganizer{
		config: config,
		logger: log.New(os.Stdout, "[FileOrganizer] ", log.LstdFlags),
	}
}

func (fo *FileOrganizer) getFileCategoryDir(ext string) string {
	ext = strings.ToLower(ext)
	for category, extensions := range fo.config.extensions {
		if _, ok := extensions[ext]; ok {
			return fo.config.directories[category]
		}
	}
	return fo.config.directories["other"]
}

func (fo *FileOrganizer) createDirectories() error {
	for _, dir := range fo.config.directories {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

func (fo *FileOrganizer) moveFile(file string) error {
	ext := strings.ToLower(filepath.Ext(file))
	destinationDir := fo.getFileCategoryDir(ext)
	destination := filepath.Join(destinationDir, filepath.Base(file))

	if _, err := os.Stat(destination); err == nil {
		base := filepath.Base(file)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)

		for i := 1; ; i++ {
			newName := fmt.Sprintf("%s_%d%s", name, i, ext)
			destination = filepath.Join(destinationDir, newName)
			if _, err := os.Stat(destination); os.IsNotExist(err) {
				break
			}
		}
	}

	return os.Rename(file, destination)
}

func (fo *FileOrganizer) OrganizeFiles() error {
	if err := fo.createDirectories(); err != nil {
		return err
	}

	files, err := os.ReadDir(fo.config.parentDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(files))

	const maxWorkers = 5
	semaphore := make(chan struct{}, maxWorkers)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(f os.DirEntry) {
			defer wg.Done()
			defer func() { <-semaphore }()

			filePath := filepath.Join(fo.config.parentDir, f.Name())
			if err := fo.moveFile(filePath); err != nil {
				errors <- fmt.Errorf("failed to move file %s: %w", f.Name(), err)
			}
		}(file)
	}

	wg.Wait()
	close(errors)

	var errCount int
	for err := range errors {
		fo.logger.Println(err)
		errCount++
	}

	if errCount > 0 {
		return fmt.Errorf("failed to move %d files", errCount)
	}

	return nil
}

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	organizer := NewFileOrganizer(config)

	fmt.Println("Managing files in the downloads folder...")
	if err := organizer.OrganizeFiles(); err != nil {
		log.Printf("Error organizing files: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Done managing files in %q!\n", config.parentDir)
}
