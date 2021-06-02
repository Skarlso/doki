package gomod

import (
	"fmt"
	"os"
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

// GetLatestModuleList returns a list of modules with their respective latest tags.
func (p *Provider) GetLatestModuleList(args []string) ([]string, error) {
	var allTags []string
	for _, imp := range args {
		if !strings.Contains(imp, "github") {
			fmt.Println("This tool does not support novel domains at the moment.")
			os.Exit(1)
		}
		if strings.Contains(imp, "@") {
			allTags = append(allTags, imp)
			continue
		}
		m := gitExtract.FindAllStringSubmatch(imp, -1)
		if len(m) == 0 {
			fmt.Printf("Failed to get latest version for import %s\n", imp)
			os.Exit(1)
		}
		if len(m[0]) < 3 {
			fmt.Println("Match does not contain repo and owner: ", m[0])
			os.Exit(1)
		}
		owner := m[0][1]
		repo := m[0][2]
		latestTag, err := p.VCS.GetLatestRemoteTag(owner, repo)
		if err != nil {
			fmt.Printf("Failed to get latest version for import %s: %s\n", imp, err.Error())
			os.Exit(1)
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
