package git_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Skarlso/doki/pkg/git"
	"github.com/Skarlso/doki/pkg/git/fakes"
	runnerfakes "github.com/Skarlso/doki/pkg/runner/fakes"
)

var _ = Describe("VCSProvider", func() {
	var fakeHttpRoundTripper *fakes.FakeHTTPRoundTripper
	var fakeRunner *runnerfakes.FakeRunner
	BeforeEach(func() {
		fakeHttpRoundTripper = &fakes.FakeHTTPRoundTripper{}
		fakeRunner = &runnerfakes.FakeRunner{}
	})

	Context("GetLatestRemoteTag", func() {
		It("can get the dev tag for a repository", func() {
			content, err := ioutil.ReadFile(filepath.Join("testdata", "latest_release.json"))
			Expect(err).ToNot(HaveOccurred())
			fakeHttpRoundTripper.RoundTripReturns(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(content)),
			}, nil)
			p := git.Provider{
				Config: git.Config{},
				Client: &http.Client{
					Transport: fakeHttpRoundTripper,
				},
			}
			tag, err := p.GetLatestRemoteTag("octocat", "Hello-World")
			Expect(err).ToNot(HaveOccurred())
			Expect(tag).To(Equal("v1.0.0"))
		})
		When("there is an error retrieving the version", func() {
			It("will provide a sensible error message", func() {
				u, err := url.Parse("https://api.github.com/repos/octocat/hello-world/releases/latest")
				Expect(err).ToNot(HaveOccurred())
				fakeHttpRoundTripper.RoundTripReturns(&http.Response{
					Status:     http.StatusText(http.StatusBadRequest),
					StatusCode: http.StatusBadRequest,
					Request: &http.Request{
						URL: u,
					},
				}, nil)
				p := git.Provider{
					Config: git.Config{},
					Client: &http.Client{
						Transport: fakeHttpRoundTripper,
					},
				}
				tag, err := p.GetLatestRemoteTag("octocat", "Hello-World")
				Expect(err).To(MatchError("failed to get latest version:  https://api.github.com/repos/octocat/hello-world/releases/latest: 400  []"))
				Expect(tag).To(BeEmpty())
			})
		})
	})
	Context("GetOwnerAndRepoFromLocal", func() {
		It("can get the owner and repo from a local repository", func() {
			By("git@github formatted remote")
			fakeRunner.RunReturns([]byte("git@github.com:weaveworks/pctl.git"), nil)
			p := git.Provider{
				Config: git.Config{
					Runner: fakeRunner,
				},
			}
			owner, repo, err := p.GetOwnerAndRepoFromLocal()
			Expect(err).ToNot(HaveOccurred())
			Expect(owner).To(Equal("weaveworks"))
			Expect(repo).To(Equal("pctl"))
			arg, args := fakeRunner.RunArgsForCall(0)
			Expect(arg).To(Equal("git"))
			Expect(args).To(ConsistOf("config", "--get", "remote.origin.url"))
			By("https://github.com formatted remote")
			fakeRunner.RunReturns([]byte("https://github.com/weaveworks/pctl.git"), nil)
			owner, repo, err = p.GetOwnerAndRepoFromLocal()
			Expect(err).ToNot(HaveOccurred())
			Expect(owner).To(Equal("weaveworks"))
			Expect(repo).To(Equal("pctl"))
			arg, args = fakeRunner.RunArgsForCall(1)
			Expect(arg).To(Equal("git"))
			Expect(args).To(ConsistOf("config", "--get", "remote.origin.url"))
		})
		When("there is an error in detecting the repo", func() {
			It("returns a sensible error", func() {
				fakeRunner.RunReturns([]byte("this is the error"), errors.New("nope"))
				p := git.Provider{
					Config: git.Config{
						Runner: fakeRunner,
					},
				}
				_, _, err := p.GetOwnerAndRepoFromLocal()
				Expect(err).To(MatchError("failed to run git command: nope"))
				arg, args := fakeRunner.RunArgsForCall(0)
				Expect(arg).To(Equal("git"))
				Expect(args).To(ConsistOf("config", "--get", "remote.origin.url"))
			})
		})
		When("owner and repo cannot be determined", func() {
			It("returns a sensible error", func() {
				fakeRunner.RunReturns([]byte("git@github.com:weaveworks"), nil)
				p := git.Provider{
					Config: git.Config{
						Runner: fakeRunner,
					},
				}
				_, _, err := p.GetOwnerAndRepoFromLocal()
				Expect(err).To(MatchError("failed to extract repo information from remote url: git@github.com:weaveworks"))
				arg, args := fakeRunner.RunArgsForCall(0)
				Expect(arg).To(Equal("git"))
				Expect(args).To(ConsistOf("config", "--get", "remote.origin.url"))
			})
		})
		When("the remote is invalid", func() {
			It("returns a sensible error", func() {
				fakeRunner.RunReturns([]byte("invalid"), nil)
				p := git.Provider{
					Config: git.Config{
						Runner: fakeRunner,
					},
				}
				_, _, err := p.GetOwnerAndRepoFromLocal()
				Expect(err).To(MatchError("failed to extract repo information from remote url: invalid"))
				arg, args := fakeRunner.RunArgsForCall(0)
				Expect(arg).To(Equal("git"))
				Expect(args).To(ConsistOf("config", "--get", "remote.origin.url"))
			})
		})
	})
	Context("GetCurrentBranch", func() {
		It("can get the current branch", func() {
			fakeRunner.RunReturns([]byte("test-branch"), nil)
			p := git.Provider{
				Config: git.Config{
					Runner: fakeRunner,
				},
			}
			branch, err := p.GetCurrentBranch()
			Expect(err).ToNot(HaveOccurred())
			Expect(branch).To(Equal("test-branch"))
			Expect(fakeRunner.RunCallCount()).To(Equal(1))
			arg, args := fakeRunner.RunArgsForCall(0)
			Expect(arg).To(Equal("git"))
			Expect(args).To(ConsistOf("branch", "--show-current"))

		})
		When("when calling git fails", func() {
			It("returns a sensible error", func() {
				fakeRunner.RunReturns(nil, errors.New("nope"))
				p := git.Provider{
					Config: git.Config{
						Runner: fakeRunner,
					},
				}
				branch, err := p.GetCurrentBranch()
				Expect(err).To(MatchError("failed to run git command: nope"))
				Expect(branch).To(BeEmpty())
				Expect(fakeRunner.RunCallCount()).To(Equal(1))
				arg, args := fakeRunner.RunArgsForCall(0)
				Expect(arg).To(Equal("git"))
				Expect(args).To(ConsistOf("branch", "--show-current"))

			})
		})
	})
	Context("GetDevTag", func() {
		It("returns a dev tag for a repository", func() {
			content, err := ioutil.ReadFile(filepath.Join("testdata", "latest_release.json"))
			Expect(err).ToNot(HaveOccurred())
			fakeHttpRoundTripper.RoundTripReturns(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(content)),
			}, nil)
			fakeRunner.RunReturnsOnCall(0, []byte("git@github.com:octocat/Hello-World.git"), nil)
			fakeRunner.RunReturnsOnCall(1, []byte("test-branch"), nil)
			p := git.Provider{
				Config: git.Config{
					Runner: fakeRunner,
				},
				Client: &http.Client{
					Transport: fakeHttpRoundTripper,
				},
			}
			tag, err := p.GetDevTag()
			Expect(err).NotTo(HaveOccurred())
			Expect(tag).To(Equal("v1.0.0-test-branch"))
		})
		When("there is an error retrieving latest version", func() {
			It("it returns a sensible error", func() {
				u, err := url.Parse("https://github.com/octocat/Hello-World")
				Expect(err).ToNot(HaveOccurred())
				fakeHttpRoundTripper.RoundTripReturns(&http.Response{
					StatusCode: http.StatusBadRequest,
					Request: &http.Request{
						URL: u,
					},
				}, errors.New("nope"))
				p := git.Provider{
					Config: git.Config{
						Runner: fakeRunner,
					},
					Client: &http.Client{
						Transport: fakeHttpRoundTripper,
					},
				}
				fakeRunner.RunReturnsOnCall(0, []byte("git@github.com:octocat/Hello-World.git"), nil)
				tag, err := p.GetDevTag()
				Expect(err).To(MatchError("failed to get latest version: Get \"https://api.github.com/repos/octocat/Hello-World/releases/latest\": nope"))
				Expect(tag).To(BeEmpty())
				Expect(fakeRunner.RunCallCount()).To(Equal(1))
			})
		})
		When("there is an error getting the repo information", func() {
			It("it returns a sensible error", func() {
				fakeHttpRoundTripper.RoundTripReturns(&http.Response{
					StatusCode: http.StatusOK,
				}, nil)
				p := git.Provider{
					Config: git.Config{
						Runner: fakeRunner,
					},
					Client: &http.Client{
						Transport: fakeHttpRoundTripper,
					},
				}
				fakeRunner.RunReturnsOnCall(0, nil, errors.New("nope"))
				tag, err := p.GetDevTag()
				Expect(err).To(MatchError("failed to run git command: nope"))
				Expect(tag).To(BeEmpty())
				Expect(fakeRunner.RunCallCount()).To(Equal(1))
				Expect(fakeHttpRoundTripper.RoundTripCallCount()).To(Equal(0))
			})
		})
		When("there is an error getting the current branch", func() {
			It("it returns a sensible error", func() {
				content, err := ioutil.ReadFile(filepath.Join("testdata", "latest_release.json"))
				Expect(err).ToNot(HaveOccurred())
				fakeHttpRoundTripper.RoundTripReturns(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(content)),
				}, nil)
				Expect(err).ToNot(HaveOccurred())
				p := git.Provider{
					Config: git.Config{
						Runner: fakeRunner,
					},
					Client: &http.Client{
						Transport: fakeHttpRoundTripper,
					},
				}
				fakeRunner.RunReturnsOnCall(0, []byte("git@github.com:octocat/Hello-World.git"), nil)
				fakeRunner.RunReturnsOnCall(1, nil, errors.New("nope"))
				tag, err := p.GetDevTag()
				Expect(err).To(MatchError("failed to run git command: nope"))
				Expect(tag).To(BeEmpty())
				Expect(fakeRunner.RunCallCount()).To(Equal(2))
				Expect(fakeHttpRoundTripper.RoundTripCallCount()).To(Equal(1))
			})
		})
	})
})
