package help

import "fmt"

func PrintHelp() {
	help := `
 ██████╗ ██╗  ██╗ ██████╗██╗  ██╗ ██████╗ ███████╗███████╗███╗   ██╗
██╔═████╗╚██╗██╔╝██╔════╝██║  ██║██╔═══██╗██╔════╝██╔════╝████╗  ██║
██║██╔██║ ╚███╔╝ ██║     ███████║██║   ██║███████╗█████╗  ██╔██╗ ██║
████╔╝██║ ██╔██╗ ██║     ██╔══██║██║   ██║╚════██║██╔══╝  ██║╚██╗██║
╚██████╔╝██╔╝ ██╗╚██████╗██║  ██║╚██████╔╝███████║███████╗██║ ╚████║
 ╚═════╝ ╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═══╝
                                                
  AI-powered recon tool for smart contract auditors
  
USAGE:
  0xchosen <command>

COMMANDS:
  filelist       Scan src/ for public/external functions → writes funcs.json
                 Edit funcs.json to trim your scope before running getrecon

  getrecon       Run full recon pipeline on current scope:
                   → Parse funcs.json
                   → Build dependency graph
                   → Run Slither static analysis
                   → Extract structured facts → facts.json
                   → Generate AI recon notes  → recon.md

  --help         Show this help message

WORKFLOW:
  1. cd into your project root (where foundry.toml lives)
  2. Run:  0xchosen filelist
  3. Edit funcs.json — delete out-of-scope entries
  4. Run:  0xchosen getrecon
  5. Open recon.md

REQUIREMENTS:
  - foundry.toml must exist in current directory
  - forge     must be installed and in PATH
  - slither   must be installed and in PATH
  - .env      must contain GROQ_API_KEY=your_key

OUTPUT FILES:
  funcs.json    All public/external functions found in src/
  facts.json    Structured facts per file (inheritance, state vars, detectors)
  recon.md      AI-generated recon report

EXAMPLES:
  0xchosen filelist
  0xchosen getrecon
`
	fmt.Println(help)
}
