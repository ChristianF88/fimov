package main

// Image organizer
// This program organizes images in a directory by moving them to a new directory
// based on the creation date of the images pictures will be moved to a new directory
// The start date, end date and foldername  can be passed via cli
// The source folder and destination path where images are stored are passed via a configuration file
// multiple source folders can be defined in the configuration file
// The configuration file is a json file with the following structure
// {
//    "camera": {"source": "your-source-path", "your-destination-path"},
//	  "whatsapp": {"source": "your-source-path", "your-destination-path"}
// }
// The cli should parse the main keywords of the config file and should accept the keywords and use corresping source and destination.
// Examples:
// fimov camera --start 2020-01-01 --end 2020-12-31 --name folder-name
// fimov whatsapp --start 2020-01-01 --name folder-name  # if end is not passed end == now
// fimov camera --start 2020-01-01 # if name is not passed name == start_end

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type PathConfig struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func main() {
	start, end, name, keyword := parseCLIArgs()
	if keyword == "" {
		fmt.Println("Usage: fimov <keyword> --start <start-date> [--end <end-date>] [--name <folder-name>]")
		return
	}

	config, err := readConfig(".fimov.json")
	if err != nil {
		fmt.Println("Error reading configuration file:", err)
		return
	}

	conf, exists := config[keyword]
	if !exists {
		fmt.Println("Invalid keyword:", keyword)
		return
	}

	if err := validatePaths(conf); err != nil {
		fmt.Println(err)
		return
	}

	if end == "" {
		end = time.Now().Format("2006-01-02")
	}

	if name == "" {
		name = fmt.Sprintf("%s_%s", start, end)
	}

	startDate, err := parseDate(start)
	if err != nil {
		fmt.Println(err)
		return
	}

	endDate, err := parseDate(end)
	if err != nil {
		fmt.Println(err)
		return
	}

	destPath := filepath.Join(conf.Destination, name)
	if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
		fmt.Println("Error creating destination directory:", err)
		return
	}

	if err := moveImages(conf.Source, destPath, startDate, endDate); err != nil {
		fmt.Println("Error organizing images:", err)
	} else {
		fmt.Println("Images organized successfully.")
	}
}

func parseCLIArgs() (string, string, string, string) {
	cameraCmd := flag.NewFlagSet("camera", flag.ExitOnError)
	whatsappCmd := flag.NewFlagSet("whatsapp", flag.ExitOnError)

	start := cameraCmd.String("start", "", "Start date (format: YYYY-MM-DD)")
	end := cameraCmd.String("end", "", "End date (format: YYYY-MM-DD)")
	name := cameraCmd.String("name", "", "Folder name")

	startWhatsApp := whatsappCmd.String("start", "", "Start date (format: YYYY-MM-DD)")
	endWhatsApp := whatsappCmd.String("end", "", "End date (format: YYYY-MM-DD)")
	nameWhatsApp := whatsappCmd.String("name", "", "Folder name")

	if len(os.Args) < 2 {
		fmt.Println("Error: keyword is required")
		fmt.Println("Usage: fimov <keyword> --start <start-date> [--end <end-date>] [--name <folder-name>]")
		os.Exit(1)
	}

	keyword := os.Args[1]

	switch keyword {
	case "camera":
		cameraCmd.Parse(os.Args[2:])
		if *start == "" {
			fmt.Println("Error: --start flag is required for camera")
			cameraCmd.Usage()
			os.Exit(1)
		}
		return *start, *end, *name, keyword
	case "whatsapp":
		whatsappCmd.Parse(os.Args[2:])
		if *startWhatsApp == "" {
			fmt.Println("Error: --start flag is required for whatsapp")
			whatsappCmd.Usage()
			os.Exit(1)
		}
		return *startWhatsApp, *endWhatsApp, *nameWhatsApp, keyword
	default:
		fmt.Println("Error: unknown keyword", keyword)
		fmt.Println("Usage: fimov <keyword> --start <start-date> [--end <end-date>] [--name <folder-name>]")
		os.Exit(1)
	}

	return "", "", "", ""
}

func readConfig(filename string) (map[string]PathConfig, error) {
	configFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config map[string]PathConfig
	if err := json.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func validatePaths(conf PathConfig) error {
	if _, err := os.Stat(conf.Source); os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", conf.Source)
	}

	if _, err := os.Stat(conf.Destination); os.IsNotExist(err) {
		return fmt.Errorf("destination path does not exist: %s", conf.Destination)
	}

	return nil
}

func parseDate(date string) (time.Time, error) {
	var layout = "2006-01-02"
	parsedDate, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid start date: %v", err)
	}
	return parsedDate, nil
}

func moveImages(source, destPath string, startDate, endDate time.Time) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			modTime := info.ModTime()
			if modTime.After(startDate) && modTime.Before(endDate) {
				destFile := filepath.Join(destPath, info.Name())
				if err := os.Rename(path, destFile); err != nil {
					fmt.Println("Error moving file:", err)
				}
			}
		}
		return nil
	})
}
