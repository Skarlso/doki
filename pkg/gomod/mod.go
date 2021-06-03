package gomod

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Skarlso/doki/pkg/git"
)

var gitExtract = regexp.MustCompile(`github.com/([a-zA-Z0-9\-]+)/([a-zA-Z0-9\-]+)`)

// Config provides configuration options for the github provider.
type Config struct {
	VCS git.VCSProvider
}

// Provider is a git functionality provider.
type Provider struct {
	Config
}

// GetModReplacements returns a list of replacement directives.
func (p *Provider) GetModReplacements(args []string) ([]string, error) {
	result := make([]string, 0)
	for _, replacement := range args {
		result = append(result, fmt.Sprintf("-replace %s", replacement))
	}
	return result, nil
}

// GetLatestModuleList returns a list of modules with their respective latest tags.
func (p *Provider) GetLatestModuleList(args []string) ([]string, error) {
	var allTags []string
	for _, imp := range args {
		if !strings.Contains(imp, "github") {
			return nil, fmt.Errorf("novel domains are not supported yet %s", imp)
		}
		if strings.Contains(imp, "@") {
			allTags = append(allTags, imp)
			continue
		}
		m := gitExtract.FindAllStringSubmatch(imp, -1)
		if len(m) == 0 {
			return nil, fmt.Errorf("failed to get latest version for import %s", imp)
		}
		if len(m[0]) < 3 {
			return nil, fmt.Errorf("match does not contain repo and owner: %s", m[0])
		}
		owner := m[0][1]
		repo := m[0][2]
		latestTag, err := p.VCS.GetLatestRemoteTag(owner, repo)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest version for owner %s and repo %s import %s: %w", owner, repo, imp, err)
		}
		allTags = append(allTags, fmt.Sprintf("%s@%s", imp, latestTag))
	}
	sort.Strings(allTags)
	return allTags, nil
}

// NewProvider creates a new git functionality provider.
func NewProvider(cfg Config) *Provider {
	return &Provider{
		Config: cfg,
	}
}
