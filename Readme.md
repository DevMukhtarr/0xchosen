# 0xchosen

AI-powered recon tool for smart contract security auditors. Scan a Solidity codebase, extract structured facts, and generate a full recon report in seconds.

---

## What It Does

`0xchosen` automates the tedious early phase of a smart contract audit:

1. Scans your `src/` directory for all public/external functions
2. Builds a dependency graph of imports across files
3. Runs [Slither](https://github.com/crytic/slither) static analysis on the entire project
4. Extracts structured facts per file — inheritance, entry points, access control, state variables, detectors
5. Sends everything to an AI model (Groq) which writes a full recon report covering:
   - What the protocol does
   - Core modules and their roles
   - Trust assumptions and privileged roles
   - Upgrade/admin risks
   - Interesting attack surfaces

**Output:** A `recon.md` file ready to drop into your audit notes.

---

## Prerequisites

Before running `0xchosen`, make sure the following are installed and available in your `PATH`:

### 1. Foundry

Required for Slither to resolve imports and remappings in Foundry projects.

```bash
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

Verify:
```bash
forge --version
```

### 2. Python + Slither

Slither requires Python 3.8+.

```bash
pip3 install slither-analyzer
```

Verify:
```bash
slither --version
```

> **Windows users:** If `slither` is not found after install, add your Python Scripts directory to PATH:
> `C:\Users\<you>\AppData\Local\Programs\Python\Python3x\Scripts`

### 3. Groq API Key

`0xchosen` uses [Groq](https://console.groq.com) to generate recon notes. Get a free API key at https://console.groq.com.

---

## Installation

Download the latest binary for your platform from the [Releases](../../releases) page.

### Linux / macOS

```bash
# Download
wget https://github.com/devmukhtarr/0xchosen/releases/download/v1.0.0/0xchosen-linux-amd64
chmod +x 0xchosen-linux-amd64
mv 0xchosen-linux-amd64 /usr/local/bin/0xchosen
```

Verify:
```bash
0xchosen --help
```

### Windows

1. Download `0xchosen-windows-amd64.exe` from the [Releases](../../releases) page
2. Rename it to `0xchosen.exe`
3. Move it to a folder that is in your `PATH`, e.g. `C:\tools\`
4. Or run it directly from the folder containing your project

---

## Setup

### 1. Set Your Groq API Key

Create a `.env` file in your project root (the same directory you will run `0xchosen` from):

```bash
# .env
GROQ_API_KEY=your_groq_api_key_here
```

> **Never commit your `.env` file.** Add it to `.gitignore`.

### 2. Confirm Your Project Structure

`0xchosen` expects a standard Foundry project layout:

```
your-project/
├── foundry.toml        ← must exist at project root
├── .env                ← your Groq API key goes here
├── src/
│   ├── Token.sol
│   ├── Vault.sol
│   └── ...
└── lib/
```

> You **must run `0xchosen` from your project root** — the directory that contains `foundry.toml` and `src/`.

---

## Usage

### Step 1 — Generate the function list

```bash
0xchosen filelist
```

This scans `src/` for all public and external functions and writes them to `funcs.json`.

**Example output:**
```
[1/10] Scanning src/ for public/external functions...
Found 18 files in scope
```

**Edit `funcs.json` to define your scope.**
Delete any entries for files or functions that are out of scope for your audit. The remaining entries are what `0xchosen` will analyze.

```json
[
  "src/Vault.sol:function deposit(uint256 amount) external {",
  "src/Vault.sol:function withdraw(uint256 amount) external nonReentrant {",
  "src/Token.sol:function transfer(address to, uint256 amount) public returns (bool) {"
]
```

### Step 2 — Run full recon

```bash
0xchosen getrecon
```

This runs all remaining steps automatically:

```
[1/10] funcs.json already exists, skipping scan...
[2/10] Parsing funcs.json...
Found 18 files in scope
[3/10] Building dependency graph...
[4/10] Running Slither analysis...
[6/10] Extracting structured facts...
[8/10] Building cross-file relationships...
[9/10] Generating recon notes via AI...
[10/10] ✅ Done! recon.md has been generated.
```

Open `recon.md` for your AI-generated recon report.

---

## Output Files

| File | Description |
|------|-------------|
| `funcs.json` | All public/external functions found in `src/`. Edit this to define scope. |
| `facts.json` | Structured facts extracted per file — inheritance, functions, state vars, Slither detectors. |
| `recon.md` | Final AI-generated recon report. |

---

## Example `recon.md` Output

```markdown
## Protocol Overview
This is a lending protocol that allows users to deposit ERC20 collateral...

## Core Modules
- **Vault.sol** — Handles deposits and withdrawals. Entry point for all user funds.
- **PriceOracle.sol** — Fetches asset prices. Trusted by Vault for liquidation logic.

## Trust Assumptions & Roles
- `owner` can update the oracle address — critical trust assumption
- `LIQUIDATOR_ROLE` can trigger liquidations without user consent

## Upgrade / Admin Risks
- Contract uses a transparent proxy pattern. The `ProxyAdmin` owner can upgrade implementation at any time.

## Interesting Attack Surfaces
- Oracle manipulation: if the price feed is stale or manipulated, undercollateralized borrows are possible
- Reentrancy in `withdraw()` — Slither flagged a potential issue (medium confidence)
```

---

## Troubleshooting

### `GetFileAttributesEx src/: The system cannot find the file specified`

You are not running `0xchosen` from your project root. `cd` to the directory containing `foundry.toml` and `src/` first:

```bash
cd path/to/your/project
0xchosen filelist
```

### Slither fails or produces no output

**Check forge is in PATH:**
```bash
forge --version
```

If that fails, run `foundryup` and restart your terminal.

**Check slither is in PATH:**
```bash
slither --version
```

**Run slither manually to see the raw error:**
```bash
slither . --json slither_test.json
```

### `GROQ_API_KEY environment variable not set`

Make sure your `.env` file is in the same directory you are running `0xchosen` from, and that it contains:
```
GROQ_API_KEY=your_key_here
```

### `funcs.json` already exists but you want to re-scan

Delete it and re-run:
```bash
rm funcs.json
0xchosen filelist
```

---

## Notes for Auditors

- **`funcs.json` is your scope file.** The filelist step is intentionally separate so you can review and trim it before analysis runs. Delete anything out of scope.
- Slither runs once on the entire project and results are filtered to your scope — this ensures cross-contract relationships are resolved correctly.
- The dependency graph resolves relative imports, so files imported by in-scope contracts are understood even if they are not directly in scope.
- `facts.json` is human-readable — you can inspect it directly if you want to see raw extracted data before the AI step.

---

## License

MIT
