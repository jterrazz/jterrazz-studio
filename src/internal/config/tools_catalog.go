package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jterrazz/jterrazz-studio/src/internal/domain/tool"
)

// checkApp returns a CheckFn for a macOS .app bundle, using the plist for version info.
func checkApp(appName string) func() CheckResult {
	return func() CheckResult {
		if _, err := os.Stat("/Applications/" + appName + ".app"); err != nil {
			return CheckResult{}
		}
		version := tool.VersionFromAppPlist(appName)()
		return CheckResult{Installed: true, Version: version}
	}
}

// checkAppWithCask is like checkApp but uses brew cask for version info.
func checkAppWithCask(appName, caskName string) func() CheckResult {
	return func() CheckResult {
		if _, err := os.Stat("/Applications/" + appName + ".app"); err != nil {
			return CheckResult{}
		}
		version := tool.VersionFromBrewCask(caskName)()
		return CheckResult{Installed: true, Version: version}
	}
}

// Tools is the single source of truth for all installable software
var Tools = []Tool{
	// ==========================================================================
	// Package Managers
	// ==========================================================================
	{
		Name:         "uv",
		Description:  "Fast Python package and project manager",
		Command:      "uv",
		Formula:      "uv",
		Method:       InstallBrewFormula,
		Category:     CategoryPackageManager,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("uv", []string{"--version"}, tool.ParseBrewVersion),
	},
	{
		Name:         "cocoapods",
		Description:  "Dependency manager for iOS and macOS projects",
		Command:      "pod",
		Formula:      "cocoapods",
		Method:       InstallBrewFormula,
		Category:     CategoryPackageManager,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("pod", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:        "homebrew",
		Description: "macOS package manager",
		Command:     "brew",
		Method:      InstallManual,
		Category:    CategoryPackageManager,
		CheckFn: func() CheckResult {
			if _, err := exec.LookPath("brew"); err != nil {
				return CheckResult{}
			}
			out, err := exec.Command("brew", "--version").Output()
			if err != nil {
				return Installed()
			}
			version := tool.ParseBrewVersion(string(out))
			formulaeOut, _ := exec.Command("brew", "list", "--formula", "-1").Output() // non-critical
			caskOut, _ := exec.Command("brew", "list", "--cask", "-1").Output()        // non-critical
			formulaeCount := 0
			caskCount := 0
			if len(strings.TrimSpace(string(formulaeOut))) > 0 {
				formulaeCount = len(strings.Split(strings.TrimSpace(string(formulaeOut)), "\n"))
			}
			if len(strings.TrimSpace(string(caskOut))) > 0 {
				caskCount = len(strings.Split(strings.TrimSpace(string(caskOut)), "\n"))
			}
			return CheckResult{
				Installed: true,
				Version:   version,
				Status:    fmt.Sprintf("%d formulae, %d casks", formulaeCount, caskCount),
			}
		},
		InstallFn: func() error {
			cmd := exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		},
	},
	{
		Name:         "npm",
		Description:  "Node.js package manager",
		Command:      "npm",
		Method:       InstallNvm,
		Category:     CategoryPackageManager,
		Dependencies: []string{"node"},
		CheckFn: func() CheckResult {
			if _, err := exec.LookPath("npm"); err != nil {
				return CheckResult{}
			}
			out, _ := exec.Command("npm", "--version").Output()
			version := tool.TrimVersion(string(out))
			npmOut, _ := exec.Command("npm", "list", "-g", "--depth=0", "--parseable").Output()
			npmLines := strings.Split(strings.TrimSpace(string(npmOut)), "\n")
			count := len(npmLines) - 1
			if count < 0 {
				count = 0
			}
			return CheckResult{
				Installed: true,
				Version:   version,
				Status:    fmt.Sprintf("%d global", count),
			}
		},
	},
	{
		Name:         "nvm",
		Description:  "Manage multiple installed Node.js versions",
		Command:      "",
		Formula:      "nvm",
		Method:       InstallBrewFormula,
		Category:     CategoryPackageManager,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"nvm"},
		CheckFn: func() CheckResult {
			nvmDir := os.Getenv("HOME") + "/.nvm"
			if _, err := os.Stat(nvmDir); err != nil {
				return CheckResult{}
			}
			versionsDir := nvmDir + "/versions/node"
			entries, err := os.ReadDir(versionsDir)
			status := ""
			if err == nil {
				count := 0
				for _, e := range entries {
					if e.IsDir() && strings.HasPrefix(e.Name(), "v") {
						count++
					}
				}
				if count > 0 {
					status = fmt.Sprintf("%d versions", count)
				}
			}
			version := tool.VersionFromBrewFormula("nvm")()
			return CheckResult{Installed: true, Version: version, Status: status}
		},
	},
	{
		Name:         "pnpm",
		Description:  "Fast, disk space efficient Node package manager",
		Command:      "pnpm",
		Formula:      "pnpm",
		Method:       InstallBrewFormula,
		Category:     CategoryPackageManager,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("pnpm", []string{"--version"}, tool.TrimVersion),
	},

	// ==========================================================================
	// Runtimes
	// ==========================================================================
	{
		Name:         "bun",
		Description:  "Fast JavaScript runtime and package manager",
		Command:      "bun",
		Formula:      "oven-sh/bun/bun",
		Method:       InstallBrewFormula,
		Category:     CategoryRuntimes,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("bun", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:         "go",
		Description:  "Compiled systems programming language by Google",
		Command:      "go",
		Formula:      "go",
		Method:       InstallBrewFormula,
		Category:     CategoryRuntimes,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("go", []string{"version"}, tool.ParseGoVersion),
	},
	{
		Name:         "node",
		Description:  "JavaScript runtime built on Chrome's V8 engine",
		Command:      "node",
		Method:       InstallNvm,
		Category:     CategoryRuntimes,
		Dependencies: []string{"nvm"},
		VersionFn:    tool.VersionFromCmd("node", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:         "openjdk",
		Description:  "Open-source Java Development Kit",
		Command:      "java",
		Formula:      "openjdk",
		Method:       InstallBrewFormula,
		Category:     CategoryRuntimes,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"java"},
		CheckFn: func() CheckResult {
			brewJava := "/opt/homebrew/opt/openjdk/bin/java"
			if _, err := os.Stat(brewJava); err == nil {
				out, _ := exec.Command(brewJava, "-version").CombinedOutput()
				return CheckResult{Installed: true, Version: tool.ParseJavaVersion(string(out))}
			}
			cmd := exec.Command("/usr/libexec/java_home")
			if err := cmd.Run(); err != nil {
				return CheckResult{}
			}
			out, _ := exec.Command("java", "-version").CombinedOutput()
			return CheckResult{Installed: true, Version: tool.ParseJavaVersion(string(out))}
		},
	},
	{
		Name:         "python",
		Description:  "Python programming language runtime",
		Command:      "python3",
		Method:       InstallManual,
		Category:     CategoryRuntimes,
		Dependencies: []string{"uv"},
		VersionFn:    tool.VersionFromCmd("python3", []string{"--version"}, tool.ParsePythonVersion),
		InstallFn: func() error {
			return ExecCommand("uv", "python", "install")
		},
	},
	{
		Name:         "rust",
		Description:  "Systems programming language for safety and speed",
		Command:      "rustc",
		Formula:      "rust",
		Method:       InstallBrewFormula,
		Category:     CategoryRuntimes,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("rustc", []string{"--version"}, tool.ParseRustVersion),
	},

	// ==========================================================================
	// Shell & Terminal
	// ==========================================================================
	{
		Name:         "ghostty",
		Description:  "GPU-accelerated terminal emulator",
		Formula:      "ghostty",
		Method:       InstallBrewCask,
		Category:     CategoryShellTerminal,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"ghostty"},
		CheckFn:      checkAppWithCask("Ghostty", "ghostty"),
	},
	{
		Name:        "ohmyzsh",
		Description: "Oh My Zsh shell framework",
		Command:     "",
		Method:      InstallManual,
		Category:    CategoryShellTerminal,
		CheckFn: func() CheckResult {
			omzPath := os.Getenv("HOME") + "/.oh-my-zsh"
			if _, err := os.Stat(omzPath); err != nil {
				return CheckResult{}
			}
			cmd := exec.Command("git", "-C", omzPath, "rev-parse", "--short", "HEAD")
			out, err := cmd.Output()
			version := ""
			if err == nil {
				version = strings.TrimSpace(string(out))
			}
			return CheckResult{Installed: true, Version: version}
		},
		InstallFn: func() error {
			cmd := exec.Command("sh", "-c", "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		},
	},
	{
		Name:         "tmux",
		Description:  "Terminal multiplexer for persistent sessions",
		Command:      "tmux",
		Formula:      "tmux",
		Method:       InstallBrewFormula,
		Category:     CategoryShellTerminal,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"tmux"},
		VersionFn:    tool.VersionFromCmd("tmux", []string{"-V"}, tool.ParseTmuxVersion),
	},
	{
		Name:         "starship",
		Description:  "Cross-shell prompt (path, git branch/status, versions)",
		Command:      "starship",
		Formula:      "starship",
		Method:       InstallBrewFormula,
		Category:     CategoryShellTerminal,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"starship"},
		VersionFn:    tool.VersionFromCmd("starship", []string{"--version"}, tool.ParseBrewVersion),
	},
	{
		Name:         "zoxide",
		Description:  "Smarter cd command",
		Replaces:     "cd (jumping around)",
		Command:      "zoxide",
		Formula:      "zoxide",
		Method:       InstallBrewFormula,
		Category:     CategoryShellTerminal,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("zoxide", []string{"--version"}, tool.ParseBrewVersion),
	},

	// ==========================================================================
	// CLI Tools
	// ==========================================================================
	{
		Name:         "yazi",
		Description:  "Terminal file manager",
		Replaces:     "Finder in the terminal",
		Command:      "yazi",
		Formula:      "yazi",
		Method:       InstallBrewFormula,
		Category:     CategoryCLITools,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("yazi", []string{"--version"}, tool.ParseBrewVersion),
	},
	{
		Name:         "bat",
		Description:  "Cat clone with syntax highlighting",
		Replaces:     "cat (for humans)",
		Command:      "bat",
		Formula:      "bat",
		Method:       InstallBrewFormula,
		Category:     CategoryCLITools,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("bat", []string{"--version"}, tool.ParseBrewVersion),
	},
	{
		Name:         "dust",
		Description:  "Intuitive disk usage tool",
		Replaces:     "du -sh",
		Command:      "dust",
		Formula:      "dust",
		Method:       InstallBrewFormula,
		Category:     CategoryCLITools,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("dust", []string{"--version"}, tool.ParseBrewVersion),
	},
	{
		Name:         "eza",
		Description:  "Modern ls replacement",
		Replaces:     "ls",
		Command:      "eza",
		Formula:      "eza",
		Method:       InstallBrewFormula,
		Category:     CategoryCLITools,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("eza", []string{"--version"}, tool.ParseBrewVersion),
	},
	{
		Name:         "fd",
		Description:  "Fast find alternative",
		Replaces:     "find",
		Command:      "fd",
		Formula:      "fd",
		Method:       InstallBrewFormula,
		Category:     CategoryCLITools,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("fd", []string{"--version"}, tool.ParseBrewVersion),
	},
	{
		Name:         "ripgrep",
		Description:  "Fast grep alternative",
		Replaces:     "grep -r",
		Command:      "rg",
		Formula:      "ripgrep",
		Method:       InstallBrewFormula,
		Category:     CategoryCLITools,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("rg", []string{"--version"}, tool.ParseBrewVersion),
	},
	{
		Name:         "sd",
		Description:  "Intuitive sed alternative",
		Replaces:     "sed (simple replaces)",
		Command:      "sd",
		Formula:      "sd",
		Method:       InstallBrewFormula,
		Category:     CategoryCLITools,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("sd", []string{"--version"}, tool.ParseBrewVersion),
	},
	{
		Name:         "difftastic",
		Description:  "Structural diff tool",
		Replaces:     "diff",
		Command:      "difft",
		Formula:      "difftastic",
		Method:       InstallBrewFormula,
		Category:     CategoryCLITools,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("difft", []string{"--version"}, tool.ParseBrewVersion),
	},

	// ==========================================================================
	// Git
	// ==========================================================================
	{
		Name:         "gpg",
		Description:  "GNU Privacy Guard for encryption and signing",
		Command:      "gpg",
		Formula:      "gnupg",
		Method:       InstallBrewFormula,
		Category:     CategoryGit,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"gpg"},
		VersionFn:    tool.VersionFromBrewFormula("gnupg"),
	},
	{
		Name:         "lazygit",
		Description:  "Terminal UI for git",
		Replaces:     "raw git porcelain (interactive)",
		Command:      "lazygit",
		Formula:      "lazygit",
		Method:       InstallBrewFormula,
		Category:     CategoryGit,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("lazygit", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:        "git",
		Description: "Distributed version control system",
		Command:     "git",
		Method:      InstallXcode,
		Category:    CategoryGit,
		VersionFn:   tool.VersionFromCmd("git", []string{"--version"}, tool.ParseGitVersion),
	},
	{
		Name:         "gh",
		Description:  "GitHub CLI for repository management",
		Replaces:     "raw GitHub API / web UI",
		Command:      "gh",
		Formula:      "gh",
		Method:       InstallBrewFormula,
		Category:     CategoryGit,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"gh"},
		VersionFn:    tool.VersionFromCmd("gh", []string{"--version"}, tool.ParseGhVersion),
	},

	// ==========================================================================
	// Editors & IDEs
	// ==========================================================================
	{
		Name:         "zed",
		Description:  "Zed code editor",
		Formula:      "zed",
		Method:       InstallBrewCask,
		Category:     CategoryEditorsIDEs,
		Dependencies: []string{"homebrew"},
		Scripts:      []string{"zed"},
		CheckFn:      checkApp("Zed"),
	},
	{
		Name:         "android-studio",
		Description:  "Android development IDE",
		Formula:      "android-studio",
		Method:       InstallBrewCask,
		Category:     CategoryEditorsIDEs,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("Android Studio"),
	},
	{
		Name:         "cursor",
		Description:  "AI-powered code editor",
		Formula:      "cursor",
		Method:       InstallBrewCask,
		Category:     CategoryEditorsIDEs,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("Cursor"),
	},
	{Name: "xcode", Description: "Apple development IDE", Method: InstallMAS, Category: CategoryEditorsIDEs, CheckFn: checkApp("Xcode")},

	// ==========================================================================
	// Containers & VMs
	// ==========================================================================
	{
		Name:         "multipass",
		Description:  "Lightweight Ubuntu VMs on demand",
		Command:      "multipass",
		Formula:      "multipass",
		Method:       InstallBrewFormula,
		Category:     CategoryContainersVMs,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("multipass", []string{"--version"}, tool.ParseMultipassVersion),
	},
	{
		Name:         "orbstack",
		Description:  "OrbStack container runtime (provides docker CLI)",
		Formula:      "orbstack",
		Method:       InstallBrewCask,
		Category:     CategoryContainersVMs,
		Dependencies: []string{"homebrew"},
		CheckFn: func() CheckResult {
			if _, err := os.Stat("/Applications/OrbStack.app"); err != nil {
				return CheckResult{}
			}
			version := tool.VersionFromAppPlist("OrbStack")()
			status := "stopped"
			if err := exec.Command("docker", "info").Run(); err == nil {
				status = "running"
			}
			return CheckResult{Installed: true, Version: version, Status: status}
		},
	},
	{
		Name:         "lens",
		Description:  "Kubernetes IDE for managing clusters",
		Formula:      "lens",
		Method:       InstallBrewCask,
		Category:     CategoryContainersVMs,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkAppWithCask("Lens", "lens"),
	},

	// ==========================================================================
	// Deploy
	// ==========================================================================
	{
		Name:         "ansible",
		Description:  "Agentless IT automation and configuration management",
		Command:      "ansible",
		Formula:      "ansible",
		Method:       InstallBrewFormula,
		Category:     CategoryDeploy,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("ansible", []string{"--version"}, tool.ParseAnsibleVersion),
	},
	{
		Name:         "eas",
		Description:  "Expo Application Services CLI for app builds",
		Command:      "eas",
		Formula:      "eas-cli",
		Method:       InstallBun,
		Category:     CategoryDeploy,
		Dependencies: []string{"bun"},
		VersionFn:    tool.VersionFromCmd("eas", []string{"--version"}, tool.ParseEasVersion),
	},
	{
		Name:         "pulumi",
		Description:  "Infrastructure as code using real programming languages",
		Command:      "pulumi",
		Formula:      "pulumi/tap/pulumi",
		Method:       InstallBrewFormula,
		Category:     CategoryDeploy,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("pulumi", []string{"version"}, tool.ParsePulumiVersion),
	},
	{
		Name:         "terraform",
		Description:  "Infrastructure as code provisioning tool",
		Command:      "terraform",
		Formula:      "hashicorp/tap/terraform",
		Method:       InstallBrewFormula,
		Category:     CategoryDeploy,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("terraform", []string{"--version"}, tool.ParseTerraformVersion),
	},

	// ==========================================================================
	// AI Agents
	// ==========================================================================
	{
		Name:        "claude",
		Description: "Claude Code CLI for agentic coding",
		Command:     "claude",
		Method:      InstallManual,
		Category:    CategoryAIAgents,
		VersionFn:   tool.VersionFromCmd("claude", []string{"--version"}, tool.ParseClaudeVersion),
		InstallFn: func() error {
			cmd := exec.Command("bash", "-c", "curl -fsSL https://claude.ai/install.sh | bash")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		},
	},
	{
		Name:         "claude-agent-acp",
		Description:  "Claude agent bridge for Zed editor",
		Command:      "claude-agent-acp",
		Formula:      "@zed-industries/claude-agent-acp",
		Method:       InstallBun,
		Category:     CategoryAIAgents,
		Dependencies: []string{"bun"},
	},
	{
		Name:         "codex",
		Description:  "OpenAI's CLI coding agent",
		Command:      "codex",
		Formula:      "codex",
		Method:       InstallBun,
		Category:     CategoryAIAgents,
		Dependencies: []string{"bun"},
		VersionFn:    tool.VersionFromCmd("codex", []string{"--version"}, tool.ParseCodexVersion),
	},
	{
		Name:         "gemini",
		Description:  "Google's CLI coding agent",
		Command:      "gemini",
		Formula:      "gemini-cli",
		Method:       InstallBrewFormula,
		Category:     CategoryAIAgents,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("gemini", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:         "opencode",
		Description:  "Open-source AI coding agent CLI",
		Command:      "opencode",
		Formula:      "opencode",
		Method:       InstallBrewFormula,
		Category:     CategoryAIAgents,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("opencode", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:        "hermes",
		Description: "Self-improving AI agent with persistent memory",
		Command:     "hermes",
		Method:      InstallManual,
		Category:    CategoryAIAgents,
		VersionFn:   tool.VersionFromCmd("hermes", []string{"--version"}, tool.ParseHermesVersion),
		InstallFn: func() error {
			cmd := exec.Command("bash", "-c", "curl -fsSL https://hermes-agent.nousresearch.com/install.sh | bash -s -- --skip-setup")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		},
	},
	{
		Name:         "conductor",
		Description:  "Runs parallel AI coding agents in git worktrees",
		Formula:      "conductor",
		Method:       InstallBrewCask,
		Category:     CategoryAIAgents,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkAppWithCask("Conductor", "conductor"),
	},

	// ==========================================================================
	// AI Tooling
	// ==========================================================================
	{
		Name:         "ollama",
		Description:  "Run large language models locally",
		Command:      "ollama",
		Formula:      "ollama-app",
		Method:       InstallBrewCask,
		Category:     CategoryAITooling,
		Dependencies: []string{"homebrew"},
		CheckFn: func() CheckResult {
			_, appErr := os.Stat("/Applications/Ollama.app")
			if appErr != nil {
				return CheckResult{}
			}
			version := tool.VersionFromBrewCask("ollama-app")()
			status := "stopped"
			if err := exec.Command("pgrep", "-x", "ollama").Run(); err == nil {
				status = "running"
			}
			return CheckResult{Installed: true, Version: version, Status: status}
		},
	},
	{
		Name:         "qmd",
		Description:  "Local semantic search over markdown",
		Replaces:     "grep over markdown notes",
		Command:      "qmd",
		Formula:      "https://github.com/tobi/qmd",
		Method:       InstallBun,
		Category:     CategoryAITooling,
		Dependencies: []string{"bun"},
		VersionFn:    tool.VersionFromCmd("qmd", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:         "rtk",
		Description:  "Token-optimized CLI proxy for agents",
		Replaces:     "raw dev commands (token-heavy output)",
		Command:      "rtk",
		Formula:      "rtk",
		Method:       InstallBrewFormula,
		Category:     CategoryAITooling,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("rtk", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:         "skills",
		Description:  "Agent skills package manager",
		Command:      "skills",
		Formula:      "skills",
		Method:       InstallBun,
		Category:     CategoryAITooling,
		Dependencies: []string{"bun"},
		VersionFn:    tool.VersionFromCmd("skills", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:          "browser-use",
		Description:   "AI agent library for browser automation",
		Command:       "browser-use",
		Formula:       "browser-use",
		Method:        InstallUV,
		PythonVersion: "3.13",
		Category:      CategoryAITooling,
		Dependencies:  []string{"uv"},
		VersionFn:     tool.VersionFromCmd("browser-use", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:          "markitdown",
		Description:   "Convert files to Markdown for LLMs",
		Command:       "markitdown",
		Formula:       "markitdown[all]",
		Method:        InstallUV,
		PythonVersion: "3.13",
		Category:      CategoryAITooling,
		Dependencies:  []string{"uv"},
	},
	{
		Name:         "playwright-mcp",
		Description:  "Browser automation for AI agents via MCP",
		Command:      "npx",
		Formula:      "@playwright/mcp",
		Method:       InstallNpm,
		Category:     CategoryAITooling,
		Dependencies: []string{"node"},
	},
	{
		Name:         "playwright-browsers",
		Description:  "Chromium runtime for website specs (@jterrazz/test)",
		Category:     CategoryAITooling,
		Dependencies: []string{"node"},
		CheckFn: func() CheckResult {
			entries, err := os.ReadDir(os.Getenv("HOME") + "/Library/Caches/ms-playwright")
			if err != nil {
				return CheckResult{}
			}
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), "chromium") {
					return CheckResult{Installed: true, Version: entry.Name()}
				}
			}
			return CheckResult{}
		},
		InstallFn: func() error {
			return ExecCommand("npx", "playwright", "install", "chromium")
		},
	},
	{
		Name:         "inferrs",
		Description:  "TurboQuant LLM inference server",
		Command:      "inferrs",
		Formula:      "ericcurtin/inferrs/inferrs",
		Method:       InstallBrewFormula,
		Category:     CategoryAITooling,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("inferrs", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:         "agent-browser",
		Description:  "Browser automation CLI for AI agents",
		Command:      "agent-browser",
		Formula:      "agent-browser",
		Method:       InstallBun,
		Category:     CategoryAITooling,
		Dependencies: []string{"bun"},
		VersionFn:    tool.VersionFromCmd("agent-browser", []string{"--version"}, tool.TrimVersion),
	},
	{
		Name:        "plannotator",
		Description: "Review AI agent plans and code before committing",
		Command:     "plannotator",
		Method:      InstallManual,
		Category:    CategoryAITooling,
		InstallFn: func() error {
			cmd := exec.Command("bash", "-c", "curl -fsSL https://plannotator.ai/install.sh | bash")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		},
		VersionFn: tool.VersionFromCmd("plannotator", []string{"--version"}, tool.TrimVersion),
	},

	// ==========================================================================
	// AI Apps
	// ==========================================================================
	{
		Name:         "chatgpt",
		Description:  "OpenAI ChatGPT desktop app",
		Formula:      "chatgpt",
		Method:       InstallBrewCask,
		Category:     CategoryAIApps,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("ChatGPT"),
	},
	{
		Name:         "claude-desktop",
		Description:  "Anthropic Claude desktop app",
		Formula:      "claude",
		Method:       InstallBrewCask,
		Category:     CategoryAIApps,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("Claude"),
	},
	{
		Name:        "typewhisper",
		Description: "Whisper-based dictation app",
		Method:      InstallManual,
		Category:    CategoryAIApps,
		CheckFn:     checkApp("TypeWhisper"),
	},

	// ==========================================================================
	// Browsers
	// ==========================================================================
	{
		Name:         "brave",
		Description:  "Privacy-focused web browser",
		Formula:      "brave-browser",
		Method:       InstallBrewCask,
		Category:     CategoryBrowsers,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("Brave Browser"),
	},
	{Name: "dia", Description: "The Browser Company's AI-first browser", Method: InstallMAS, Category: CategoryBrowsers, CheckFn: checkApp("Dia")},

	// ==========================================================================
	// Communication
	// ==========================================================================
	{
		Name:         "discord",
		Description:  "Voice and text chat",
		Formula:      "discord",
		Method:       InstallBrewCask,
		Category:     CategoryCommunication,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("Discord"),
	},
	{
		Name:         "slack",
		Description:  "Team communication",
		Formula:      "slack",
		Method:       InstallBrewCask,
		Category:     CategoryCommunication,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("Slack"),
	},
	{
		Name:         "whatsapp",
		Description:  "Messaging app",
		Formula:      "whatsapp",
		Method:       InstallBrewCask,
		Category:     CategoryCommunication,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("WhatsApp"),
	},
	{Name: "messenger", Description: "Facebook Messenger", Method: InstallMAS, Category: CategoryCommunication, CheckFn: checkApp("Messenger")},

	// ==========================================================================
	// Productivity
	// ==========================================================================
	{
		Name:         "linear",
		Description:  "Project management tool",
		Formula:      "linear-linear",
		Method:       InstallBrewCask,
		Category:     CategoryProductivity,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("Linear"),
	},
	{
		Name:         "notion",
		Description:  "Workspace for notes and docs",
		Formula:      "notion",
		Method:       InstallBrewCask,
		Category:     CategoryProductivity,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("Notion"),
	},
	{
		Name:         "obsidian",
		Description:  "Knowledge base and note-taking",
		Formula:      "obsidian",
		Method:       InstallBrewCask,
		Category:     CategoryProductivity,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("Obsidian"),
	},
	{Name: "keynote", Description: "Apple presentations", Method: InstallMAS, Category: CategoryProductivity, CheckFn: checkApp("Keynote")},
	{Name: "numbers", Description: "Apple spreadsheets", Method: InstallMAS, Category: CategoryProductivity, CheckFn: checkApp("Numbers")},
	{Name: "pages", Description: "Apple word processor", Method: InstallMAS, Category: CategoryProductivity, CheckFn: checkApp("Pages")},
	{Name: "raindrop", Description: "Bookmark manager", Method: InstallMAS, Category: CategoryProductivity, CheckFn: checkApp("Save to Raindrop.io")},

	// ==========================================================================
	// Media
	// ==========================================================================
	{Name: "broadcasts", Description: "Internet radio and stream player", Method: InstallMAS, Category: CategoryMedia, CheckFn: checkApp("Broadcasts")},
	{Name: "compressor", Description: "Apple video compression tool", Method: InstallMAS, Category: CategoryMedia, CheckFn: checkApp("Compressor")},
	{Name: "final-cut-pro", Description: "Professional video editor", Method: InstallMAS, Category: CategoryMedia, CheckFn: checkApp("Final Cut Pro")},
	{Name: "lightroom", Description: "Adobe photo editor", Method: InstallMAS, Category: CategoryMedia, CheckFn: checkApp("Adobe Lightroom")},
	{Name: "logic-pro", Description: "Professional music production", Method: InstallMAS, Category: CategoryMedia, CheckFn: checkApp("Logic Pro")},

	// ==========================================================================
	// Remote Access
	// ==========================================================================
	{
		Name:         "jump-desktop-connect",
		Description:  "Remote GUI recovery host service for headless Macs",
		Formula:      "jump-desktop-connect",
		Method:       InstallBrewCask,
		Category:     CategoryRemoteAccess,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkAppWithCask("Jump Desktop Connect", "jump-desktop-connect"),
	},
	{
		Name:         "jump-desktop",
		Description:  "Jump Desktop viewer/client app",
		Formula:      "jump-desktop",
		Method:       InstallBrewCask,
		Category:     CategoryRemoteAccess,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkAppWithCask("Jump Desktop", "jump-desktop"),
	},
	{
		Name:         "tailscale",
		Description:  "Mesh VPN built on WireGuard",
		Command:      "tailscale",
		Formula:      "tailscale",
		Method:       InstallBrewFormula,
		Category:     CategoryRemoteAccess,
		Dependencies: []string{"homebrew"},
		CheckFn: func() CheckResult {
			if _, err := exec.LookPath("tailscale"); err == nil {
				version := tool.VersionFromCmd("tailscale", []string{"version"}, tool.ParseTailscaleVersion)()
				status := "installed"
				if out, err := exec.Command("tailscale", "status", "--json").Output(); err == nil {
					if strings.Contains(string(out), `"BackendState":"Running"`) {
						status = "running"
					}
				}
				return CheckResult{Installed: true, Version: version, Status: status}
			}
			if _, err := os.Stat("/Applications/Tailscale.app"); err == nil {
				version := tool.VersionFromAppPlist("Tailscale")()
				return CheckResult{Installed: false, Version: version, Status: "app only"}
			}
			return CheckResult{}
		},
	},

	// ==========================================================================
	// Security
	// ==========================================================================
	{
		Name:         "bitwarden",
		Description:  "Password manager",
		Formula:      "bitwarden",
		Method:       InstallBrewCask,
		Category:     CategorySecurity,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkApp("Bitwarden"),
	},
	{Name: "adguard", Description: "Ad blocker for Safari", Method: InstallMAS, Category: CategorySecurity, CheckFn: checkApp("AdGuard for Safari")},
	{Name: "passepartout", Description: "VPN client", Method: InstallMAS, Category: CategorySecurity, CheckFn: checkApp("Passepartout")},

	// ==========================================================================
	// System Utilities
	// ==========================================================================
	{
		Name:         "mole",
		Description:  "Clean, uninstall, analyze, and optimize macOS from the terminal",
		Command:      "mo",
		Formula:      "tw93/tap/mole",
		Method:       InstallBrewFormula,
		Category:     CategorySystemUtilities,
		Dependencies: []string{"homebrew"},
		VersionFn:    tool.VersionFromCmd("mo", []string{"--version"}, tool.ParseMoleVersion),
	},
	{
		Name:         "betterdisplay",
		Description:  "Flexible HiDPI scaling and display management",
		Formula:      "betterdisplay",
		Method:       InstallBrewCask,
		Category:     CategorySystemUtilities,
		Dependencies: []string{"homebrew"},
		CheckFn:      checkAppWithCask("BetterDisplay", "betterdisplay"),
	},
	{Name: "pipifier", Description: "Picture-in-Picture for Safari", Method: InstallMAS, Category: CategorySystemUtilities, CheckFn: checkApp("PiPifier")},
	{Name: "snippety", Description: "Code snippet manager", Method: InstallMAS, Category: CategorySystemUtilities, CheckFn: checkApp("Snippety")},
	{Name: "speedtest", Description: "Internet speed test", Method: InstallMAS, Category: CategorySystemUtilities, CheckFn: checkApp("Speedtest")},
}
