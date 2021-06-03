package gomod_test

import (
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
		It("return a list of modules with their latest tags", func() {
			fakeVcsProvider.GetLatestRemoteTagReturns("v0.0.1", nil)
			p := gomod.NewProvider(gomod.Config{
				VCS: fakeVcsProvider,
			})
			args := []string{"github.com/weaveworks/pctl"}
			result, err := p.GetLatestModuleList(args)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(ConsistOf("github.com/weaveworks/pctl@v0.0.1"))
		})
	})
})
