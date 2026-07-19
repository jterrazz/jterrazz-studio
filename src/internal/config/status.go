package config

// StatusSection represents a section in the status display
type StatusSection struct {
	Title    string
	SubTitle string // Optional subsection title
	RenderFn func() // Function to render this section
}

// StatusSections defines all status sections in display order
var StatusSections = []StatusSection{
	// System — live performance
	{Title: "System", SubTitle: "CPU", RenderFn: nil},
	{Title: "System", SubTitle: "Memory", RenderFn: nil},

	// Environment — network, services, system health
	{Title: "Environment", SubTitle: "Network", RenderFn: nil},
	{Title: "Environment", SubTitle: "Services", RenderFn: nil},
	{Title: "Environment", SubTitle: "Health", RenderFn: nil},

	// Workspace — local project state
	{Title: "Workspace", SubTitle: "Git", RenderFn: nil},
	{Title: "Workspace", SubTitle: "Disk", RenderFn: nil},

	// Config — per-Script-Category subsections + Identity + Network (j remote)
	{Title: "Config", SubTitle: "Terminal", RenderFn: nil},
	{Title: "Config", SubTitle: "Security", RenderFn: nil},
	{Title: "Config", SubTitle: "Editor", RenderFn: nil},
	{Title: "Config", SubTitle: "System", RenderFn: nil},
	{Title: "Config", SubTitle: "Server", RenderFn: nil},
	{Title: "Config", SubTitle: "Network", RenderFn: nil},
	{Title: "Config", SubTitle: "Identity", RenderFn: nil},

	// Tools — inventory (cards)
	{Title: "Tools", SubTitle: "Package Managers", RenderFn: nil},
	{Title: "Tools", SubTitle: "Runtimes", RenderFn: nil},
	{Title: "Tools", SubTitle: "Shell & Terminal", RenderFn: nil},
	{Title: "Tools", SubTitle: "CLI Tools", RenderFn: nil},
	{Title: "Tools", SubTitle: "Git", RenderFn: nil},
	{Title: "Tools", SubTitle: "Editors & IDEs", RenderFn: nil},
	{Title: "Tools", SubTitle: "Containers & VMs", RenderFn: nil},
	{Title: "Tools", SubTitle: "Deploy", RenderFn: nil},
	{Title: "Tools", SubTitle: "AI Agents", RenderFn: nil},
	{Title: "Tools", SubTitle: "AI Tooling", RenderFn: nil},
	{Title: "Tools", SubTitle: "AI Apps", RenderFn: nil},
	{Title: "Tools", SubTitle: "Browsers", RenderFn: nil},
	{Title: "Tools", SubTitle: "Communication", RenderFn: nil},
	{Title: "Tools", SubTitle: "Productivity", RenderFn: nil},
	{Title: "Tools", SubTitle: "Media", RenderFn: nil},
	{Title: "Tools", SubTitle: "Remote Access", RenderFn: nil},
	{Title: "Tools", SubTitle: "Security", RenderFn: nil},
	{Title: "Tools", SubTitle: "System Utilities", RenderFn: nil},
}

// ToolCategories defines the order of tool categories in status display
var ToolCategories = []ToolCategory{
	CategoryPackageManager,
	CategoryRuntimes,
	CategoryShellTerminal,
	CategoryCLITools,
	CategoryGit,
	CategoryEditorsIDEs,
	CategoryContainersVMs,
	CategoryDeploy,
	CategoryAIAgents,
	CategoryAITooling,
	CategoryAIApps,
	CategoryBrowsers,
	CategoryCommunication,
	CategoryProductivity,
	CategoryMedia,
	CategoryRemoteAccess,
	CategorySecurity,
	CategorySystemUtilities,
}
