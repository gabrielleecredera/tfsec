package test

import (
	"testing"

	"github.com/aquasecurity/tfsec/internal/pkg/testutil/filesystem"

	"github.com/aquasecurity/defsec/provider"
	"github.com/aquasecurity/defsec/rules"

	"github.com/aquasecurity/tfsec/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aquasecurity/defsec/severity"
	"github.com/aquasecurity/tfsec/internal/pkg/block"
	"github.com/aquasecurity/tfsec/internal/pkg/parser"
	"github.com/aquasecurity/tfsec/internal/pkg/scanner"
	"github.com/aquasecurity/tfsec/pkg/rule"
)

var panicRule = rule.Rule{
	Base: rules.Register(
		rules.Rule{
			Provider:  provider.AWSProvider,
			Service:   "service",
			ShortCode: "abc",
			Severity:  severity.High,
		},
		nil,
	),
	RequiredTypes:  []string{"resource"},
	RequiredLabels: []string{"problem"},
	CheckTerraform: func(resourceBlock *block.Block, _ *block.Module) (results rules.Results) {
		if resourceBlock.GetAttribute("panic").IsTrue() {
			panic("This is fine")
		}
		return
	},
}

func Test_PanicInCheckNotAllowed(t *testing.T) {

	scanner.RegisterCheckRule(panicRule)
	defer scanner.DeregisterCheckRule(panicRule)

	fs, err := filesystem.New()
	require.NoError(t, err)
	defer fs.Close()

	require.NoError(t, fs.WriteTextFile("project/main.tf", `
resource "problem" "this" {
	panic = true
}
`))

	blocks, err := parser.New(fs.RealPath("project/"), parser.OptionStopOnHCLError()).ParseDirectory()
	require.NoError(t, err)
	results, _ := scanner.New().Scan(blocks)
	testutil.AssertRuleNotFound(t, panicRule.ID(), results, "")
}

func Test_PanicInCheckAllowed(t *testing.T) {

	scanner.RegisterCheckRule(panicRule)
	defer scanner.DeregisterCheckRule(panicRule)

	fs, err := filesystem.New()
	require.NoError(t, err)
	defer fs.Close()

	require.NoError(t, fs.WriteTextFile("project/main.tf", `
resource "problem" "this" {
	panic = true
}
`))

	blocks, err := parser.New(fs.RealPath("project/"), parser.OptionStopOnHCLError()).ParseDirectory()
	require.NoError(t, err)
	_, err = scanner.New(scanner.OptionStopOnErrors()).Scan(blocks)
	assert.Error(t, err)
}

func Test_PanicNotInCheckNotIncludePassed(t *testing.T) {

	scanner.RegisterCheckRule(panicRule)
	defer scanner.DeregisterCheckRule(panicRule)

	fs, err := filesystem.New()
	require.NoError(t, err)
	defer fs.Close()

	require.NoError(t, fs.WriteTextFile("project/main.tf", `
resource "problem" "this" {
	panic = true
}
`))

	blocks, err := parser.New(fs.RealPath("project/"), parser.OptionStopOnHCLError()).ParseDirectory()
	require.NoError(t, err)
	results, _ := scanner.New().Scan(blocks)
	testutil.AssertRuleNotFound(t, panicRule.ID(), results, "")
}

func Test_PanicNotInCheckNotIncludePassedStopOnError(t *testing.T) {

	scanner.RegisterCheckRule(panicRule)
	defer scanner.DeregisterCheckRule(panicRule)

	fs, err := filesystem.New()
	require.NoError(t, err)
	defer fs.Close()

	require.NoError(t, fs.WriteTextFile("project/main.tf", `
resource "problem" "this" {
	panic = true
}
`))

	blocks, err := parser.New(fs.RealPath("project/"), parser.OptionStopOnHCLError()).ParseDirectory()
	require.NoError(t, err)

	_, err = scanner.New(scanner.OptionStopOnErrors()).Scan(blocks)
	assert.Error(t, err)
}
