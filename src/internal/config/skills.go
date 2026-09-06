package config

// SkillRepo represents a repository containing AI agent skills
type SkillRepo struct {
	Name        string // Repository name (owner/repo format)
	Description string
}

// Skill represents a favorite/installed skill
type Skill struct {
	Repo  string // Repository (owner/repo format)
	Skill string // Skill name within the repo
}

// StudioSkills are @jterrazz skills — the foundation of every project
var StudioSkills = []Skill{
	{"jterrazz/jterrazz-studio", "jterrazz-toolbelt"},
	{"jterrazz/jterrazz-studio", "jterrazz-stack"},
	{"jterrazz/jterrazz-studio", "jterrazz-new-project"},
	{"jterrazz/jterrazz-studio", "jterrazz-repo-structure"},
	{"jterrazz/package-reach", "jterrazz-reach"},
	{"jterrazz/package-reach", "jterrazz-seo"},
	{"jterrazz/package-reach", "jterrazz-geo"},
	{"jterrazz/package-reach", "jterrazz-structured-data"},
	{"jterrazz/package-reach", "jterrazz-content-reach"},
	{"jterrazz/jterrazz-infrastructure", "jterrazz-infra"},
	{"jterrazz/package-typescript", "jterrazz-typescript"},
	{"jterrazz/package-attestation", "jterrazz-attestation"},
	{"jterrazz/package-broadcast", "jterrazz-broadcast"},
	{"jterrazz/package-test", "jterrazz-test"},
	{"jterrazz/jterrazz-actions", "jterrazz-workflows"},
}

// CommunitySkills are third-party skills worth having
var CommunitySkills = []Skill{
	{"anthropics/skills", "frontend-design"},
	{"expo/skills", "upgrading-expo"},
	{"GoogleChrome/modern-web-guidance", "modern-web-guidance"},
	{"GoogleChrome/modern-web-guidance", "chrome-extensions"},
	{"jakubkrehel/skills", "better-interface"},
	{"jakubkrehel/skills", "better-accessibility"},
	{"jakubkrehel/skills", "better-layout"},
	{"jakubkrehel/skills", "better-writing"},
	{"jakubkrehel/skills", "better-typography"},
	{"jakubkrehel/skills", "better-colors"},
	{"jakubkrehel/skills", "better-ui"},
	{"shadcn-ui/ui", "shadcn"},
	{"sickn33/antigravity-awesome-skills", "last30days"},
	{"tobi/qmd", "qmd"},
	{"vercel-labs/agent-browser", "agent-browser"},
	{"vercel-labs/agent-skills", "vercel-react-best-practices"},
	{"vercel-labs/agent-skills", "vercel-react-native-skills"},
}

// FavoriteSkills is all pinned skills combined (studio + community)
var FavoriteSkills = append(StudioSkills, CommunitySkills...)

// StudioRepos are @jterrazz skill repositories
var StudioRepos = []SkillRepo{
	{"jterrazz/jterrazz-studio", "The studio — machine CLI, stack conventions, repo doctrine"},
	{"jterrazz/jterrazz-infrastructure", "Infrastructure and deployment (K3s, Helm, Traefik)"},
	{"jterrazz/jterrazz-actions", "Shared CI/CD workflows (validate, release)"},
	{"jterrazz/package-typescript", "TypeScript toolchain — build, lint, format, docs"},
	{"jterrazz/package-attestation", "Article attestation — EIP-712 + OpenTimestamps"},
	{"jterrazz/package-broadcast", "Multi-channel announcements (App Store, push)"},
	{"jterrazz/package-reach", "Site reach — SEO/GEO surfaces, projections, conformance"},
	{"jterrazz/package-test", "Testing utilities (vitest mocks)"},
}

// CommunityRepos are third-party skill repositories
var CommunityRepos = []SkillRepo{
	{"anthropics/skills", "Official Anthropic skills for Claude"},
	{"better-auth/skills", "Authentication best practices"},
	{"code-with-beto/skills", "Beto's development skills"},
	{"coreyhaines31/marketingskills", "Marketing and SEO skills"},
	{"expo/skills", "Expo and React Native mobile development"},
	{"firecrawl/cli", "Web content extraction for AI agents"},
	{"GoogleChrome/modern-web-guidance", "Chrome's current web platform guidance, and extension authoring"},
	{"jakubkrehel/skills", "Interface craft — UI polish, typography, color, a11y, layout, copy"},
	{"shadcn-ui/ui", "Official shadcn/ui components and patterns"},
	{"obra/superpowers", "Development workflow and productivity skills"},
	{"remotion-dev/skills", "Remotion video creation skills"},
	{"resend/email-best-practices", "Email development best practices"},
	{"supabase/agent-skills", "Supabase database and backend skills"},
	{"tobi/qmd", "Local search engine for docs and knowledge bases"},
	{"vercel-labs/agent-browser", "Browser automation CLI for AI agents"},
	{"vercel-labs/agent-skills", "Vercel React and web development skills"},
}

// SkillRepos is all repositories combined (studio + community)
var SkillRepos = append(StudioRepos, CommunityRepos...)

// GetAllSkillRepos returns all skill repositories
func GetAllSkillRepos() []SkillRepo {
	return SkillRepos
}

// GetStudioSkills returns @jterrazz skills
func GetStudioSkills() []Skill {
	return StudioSkills
}

// GetCommunitySkills returns third-party skills
func GetCommunitySkills() []Skill {
	return CommunitySkills
}

// GetStudioRepos returns @jterrazz skill repositories
func GetStudioRepos() []SkillRepo {
	return StudioRepos
}

// GetCommunityRepos returns third-party skill repositories
func GetCommunityRepos() []SkillRepo {
	return CommunityRepos
}

// GetSkillRepoByName returns a skill repo by name
func GetSkillRepoByName(name string) *SkillRepo {
	for i := range SkillRepos {
		if SkillRepos[i].Name == name {
			return &SkillRepos[i]
		}
	}
	return nil
}

// GetFavoriteSkills returns all favorite skills
func GetFavoriteSkills() []Skill {
	return FavoriteSkills
}

// IsFavoriteSkill checks if a skill is in the favorites list
// If repo is empty, only the skill name is checked
func IsFavoriteSkill(repo, skill string) bool {
	for _, fav := range FavoriteSkills {
		if fav.Skill == skill {
			if repo == "" || fav.Repo == repo {
				return true
			}
		}
	}
	return false
}
