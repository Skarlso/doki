package gomod_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Skarlso/doki/pkg/git/fakes"
	"github.com/Skarlso/doki/pkg/gomod"
)

var _ = Describe("GoMod", func() {
	Context("GetModReplacements", func() {
		It("can list replacements for modules", func() {
			args := []string{"module1", "module2"}
			p := gomod.NewProvider(gomod.Config{})
			result, err := p.GetModReplacements(args)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(ConsistOf("-replace module1", "-replace module2"))
		})
	})
	Context("GetLatestModuleList", func() {
		var fakeVcsProvider *fakes.FakeVCSProvider
		BeforeEach(func() {
			fakeVcsProvider = &fakes.FakeVCSProvider{}
		})
		It("returns a list of modules with their latest tags", func() {
			fakeVcsProvider.GetLatestRemoteTagReturnsOnCall(0, "v0.0.1", nil)
			fakeVcsProvider.GetLatestRemoteTagReturnsOnCall(1, "v0.0.2", nil)
			p := gomod.NewProvider(gomod.Config{
				VCS: fakeVcsProvider,
			})
			args := []string{"github.com/weaveworks/pctl", "github.com/weaveworks/profiles"}
			result, err := p.GetLatestModuleList(args)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(ConsistOf("github.com/weaveworks/pctl@v0.0.1", "github.com/weaveworks/profiles@v0.0.2"))
			Expect(fakeVcsProvider.GetLatestRemoteTagCallCount()).To(Equal(2))
		})
		When("getting the latest remote tag returns an error", func() {
			It("returns a sensible error message", func() {
				fakeVcsProvider.GetLatestRemoteTagReturns("", errors.New("nope"))
				p := gomod.NewProvider(gomod.Config{
					VCS: fakeVcsProvider,
				})
				args := []string{"github.com/weaveworks/pctl"}
				result, err := p.GetLatestModuleList(args)
				Expect(err).To(MatchError("failed to get latest version for owner weaveworks and repo pctl import github.com/weaveworks/pctl: nope"))
				Expect(result).To(BeNil())
			})
		})
		When("the modules is a novel domain", func() {
			It("returns that anything other than github is not supported right now", func() {
				p := gomod.NewProvider(gomod.Config{
					VCS: fakeVcsProvider,
				})
				args := []string{"k8s.io/api"}
				result, err := p.GetLatestModuleList(args)
				Expect(err).To(MatchError("novel domains are not supported yet k8s.io/api"))
				Expect(result).To(BeNil())
				Expect(fakeVcsProvider.GetLatestRemoteTagCallCount()).To(Equal(0))
			})
		})
		When("the modules already contains a version peg", func() {
			It("returns that unmodified", func() {
				p := gomod.NewProvider(gomod.Config{
					VCS: fakeVcsProvider,
				})
				args := []string{"github.com/weaveworks/pctl@v0.0.5-dev"}
				result, err := p.GetLatestModuleList(args)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(ConsistOf("github.com/weaveworks/pctl@v0.0.5-dev"))
				Expect(fakeVcsProvider.GetLatestRemoteTagCallCount()).To(Equal(0))
			})
		})
	})
})
