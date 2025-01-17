package main

const checkTemplate = `package {{ .Package}}

// ATTENTION!
// This rule was autogenerated!
// Before making changes, consider updating the generator.

import (
	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/hclcontext"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/pkg/provider"
	"github.com/aquasecurity/tfsec/pkg/result"
	"github.com/aquasecurity/tfsec/pkg/rule"
	"github.com/aquasecurity/tfsec/pkg/severity"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		Provider:       provider.{{.Provider.ConstName}}Provider,
		Service:   "{{ .Service}}",
		ShortCode: "{{ .ShortCode}}",
		Documentation: rule.RuleDocumentation{
			Summary:     "{{.Summary}}",
			Explanation: ` + "`" + `{{.Explanation}}` + "`" + `,
			Impact:      "{{.Impact}}",
			Resolution:  "{{.Resolution}}",
			BadExample: []string{  ` + "`" + `
{{.BadExampleCode}}
` + "`" + `},
			GoodExample: []string{ ` + "`" + `
{{.GoodExampleCode}}
` + "`" + `},
			Links: []string{
				"{{.FirstLink}}",
			},
		},
		RequiredTypes:  []string{ {{range .RequiredTypes}}
			"{{.}}",{{end}}
		},
		RequiredLabels: []string{ {{range .RequiredLabels}}
			"{{.}}",{{end}}
		},
		DefaultSeverity: severity.{{.Severity}}, 
		CheckFunc: func(set result.Set, resourceBlock block.Block, _ *hclcontext.Context){
			{{.RuleCode}}
		},
	})
}
`

const checkTestTemplate = `package {{ .Package}}

import (
	"strings"
	"testing"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/testutil"
)

func Test_{{.TestName}}_FailureExamples(t *testing.T) {
	expectedCode := "{{ .Provider}}-{{ .Service }}-{{.ShortCode}}"

	rule, err := scanner.GetRuleById(expectedCode)
	if err != nil {
		t.Fatalf("Rule not found: %s", expectedCode)
	}
	for i, badExample := range rule.Documentation.BadExample {
		t.Logf("Running bad example for '%s' #%d", expectedCode, i+1)
		if strings.TrimSpace(badExample) == "" {
			t.Fatalf("bad example code not provided for %s", rule.ID())
		}
		defer func() {
			if err := recover(); err != nil {
				t.Fatalf("Scan (bad) failed: %s", err)
			}
		}()
		results := testutil.ScanHCL(badExample, t)
		testutil.AssertCheckCode(t, rule.ID(), "", results)
	}
}

func Test_{{.TestName}}_SuccessExamples(t *testing.T) {
	expectedCode := "{{ .Provider}}-{{ .Service }}-{{.ShortCode}}"

	rule, err := scanner.GetRuleById(expectedCode)
	if err != nil {
		t.Fatalf("Rule not found: %s", expectedCode)
	}
	for i, example := range rule.Documentation.GoodExample {
		t.Logf("Running good example for '%s' #%d", expectedCode, i+1)
		if strings.TrimSpace(example) == "" {
			t.Fatalf("good example code not provided for %s", rule.ID())
		}
		defer func() {
			if err := recover(); err != nil {
				t.Fatalf("Scan (good) failed: %s", err)
			}
		}()
		results := testutil.ScanHCL(example, t)
		testutil.AssertCheckCode(t, "", rule.ID(), results)
	}
}
`
