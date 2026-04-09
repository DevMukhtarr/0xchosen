package main

import (
	"fmt"
	"log"
	"os"

	"github.com/devmukhtarr/0xchosen/internal/ai"
	"github.com/devmukhtarr/0xchosen/internal/config"
	"github.com/devmukhtarr/0xchosen/internal/dependency"
	"github.com/devmukhtarr/0xchosen/internal/extractor"
	"github.com/devmukhtarr/0xchosen/internal/filelist"
	"github.com/devmukhtarr/0xchosen/internal/parser"
	"github.com/devmukhtarr/0xchosen/internal/relationship"
	"github.com/devmukhtarr/0xchosen/internal/report"
	"github.com/devmukhtarr/0xchosen/internal/slither"
	"github.com/devmukhtarr/0xchosen/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	godotenv.Load()

	if len(os.Args) < 2 {
		fmt.Println("Usage: 0xchosen filelist")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "filelist":
		runFilelist()
	case "getrecon":
		runGetRecon()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runGetRecon() {
	var files []string
	files = runFilelist()

	runGetMappingAndPrintReconMd(files)
}

func runFilelist() (files []string) {
	// Step 1: Run grep → write funcs.json
	if config.FileExists("funcs.json") {
		fmt.Println("[1/10] funcs.json already exists, skipping scan...")
	} else {
		fmt.Println("[1/10] Scanning src/ for public/external functions...")
		err := filelist.Generate("src/", "funcs.json")
		if err != nil {
			log.Fatalf("Failed to generate file list: %v", err)
		}
	}

	// Step 2: Parse funcs.json → extract unique file paths
	fmt.Println("[2/10] Parsing funcs.json...")
	files, err := parser.ExtractFilePaths("funcs.json")
	if err != nil {
		log.Fatalf("Failed to parse funcs.json: %v", err)
	}
	fmt.Printf("Found %d files in scope\n", len(files))
	return files
}

func runGetMappingAndPrintReconMd(files []string) {
	// Step 3: Parse imports → build dependency graph
	fmt.Println("[3/10] Building dependency graph...")
	depGraph, err := dependency.Build(files)
	if err != nil {
		log.Fatalf("Failed to build dependency graph: %v", err)
	}

	// Step 4: Run Slither on each file
	fmt.Println("[4/10] Running Slither analysis...")
	projectRoot, _ := os.Getwd()
	slitherResults, err := slither.RunAll(projectRoot, files)
	if err != nil {
		log.Fatalf("Slither failed: %v", err)
	}

	// Step 6: Extract structured facts
	fmt.Println("[6/10] Extracting structured facts...")
	facts := extractor.Extract(slitherResults, depGraph)

	// Step 7: Store results to JSON
	err = store.Save(facts, "facts.json")
	if err != nil {
		log.Fatalf("Failed to save facts: %v", err)
	}

	// Step 8: Build cross-file relationship summary
	fmt.Println("[8/10] Building cross-file relationships...")
	relationships := relationship.Build(facts)

	// Step 9: Send to AI → generate recon notes
	fmt.Println("[9/10] Generating recon notes via AI...")
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatal("GROQ_API_KEY environment variable not set")
	}
	reconNotes, err := ai.GenerateRecon(facts, relationships, apiKey)
	if err != nil {
		log.Fatalf("AI generation failed: %v", err)
	}

	// Step 10: Save recon.md
	err = report.Save(reconNotes, "recon.md")
	if err != nil {
		log.Fatalf("Failed to save recon.md: %v", err)
	}

	fmt.Println("\n [10/10]✅ Done! recon.md has been generated.")
}
