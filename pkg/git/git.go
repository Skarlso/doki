package git

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/Skarlso/doki/pkg/runner"
	"github.com/google/go-github/v35/github"
)

var gitExtractor = regexp.MustCompile("^(https|git)(://|@)([^/:]+)[/:]([^/:]+)/(.+)$")

// Provider is a git functionality provider.
type Provider struct {
	Runner runner.Runner
}

// NewProvider creates a new git functionality provider.
func NewProvider(runner runner.Runner) *Provider {
	return &Provider{
		Runner: runner,
	}
}

// GetLatestRemoteTag gets the latest tag from the given remote git repo.
func (p *Provider) GetLatestRemoteTag(owner, repo string) (string, error) {
	client := github.NewClient(nil)
	release, response, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		p.logGithubResponseBody(response)
		return "", err
	}
	return release.GetTagName(), nil
}

// GetOwnerAndRepoFromLocal returns the owner and the repo name from a local git repository.
func (p *Provider) GetOwnerAndRepoFromLocal() (string, string, error) {
	out, err := p.Runner.Run("git", "config", "--get", "remote.origin.url")
	if err != nil {
		return "", "", err
	}
	u := strings.ReplaceAll(string(out), ".git", "")
	m := gitExtractor.FindAllStringSubmatch(strings.TrimSuffix(u, "\n"), -1)
	if m == nil {
		return "", "", fmt.Errorf("failed to extract repo information from remote url: %s\n", string(out))
	}
	if len(m[0]) < 5 {
		return "", "", fmt.Errorf("did not find repo information in match: %v", m[0])
	}

	return m[0][4], m[0][5], nil
}

// GetCurrentBranch gets the current active branch of a repository.
func (p *Provider) GetCurrentBranch() (string, error) {
	out, err := p.Runner.Run("git", "branch", "--show-current")
	if err != nil {
		return "", err
	}
	return string(bytes.Trim(out, "\n")), nil
}

// logGithubResponseBody logs a response if it's not nil.
func (p *Provider) logGithubResponseBody(response *github.Response) {
	if response == nil {
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read message body.")
		os.Exit(1)
	}
	fmt.Println("Body of the response: ", string(body))
}
