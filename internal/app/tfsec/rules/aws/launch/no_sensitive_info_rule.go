package launch

// ATTENTION!
// This rule was autogenerated!
// Before making changes, consider updating the generator.

// generator-locked

import (
	"encoding/base64"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/debug"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/hclcontext"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/pkg/provider"
	"github.com/aquasecurity/tfsec/pkg/result"
	"github.com/aquasecurity/tfsec/pkg/rule"
	"github.com/aquasecurity/tfsec/pkg/severity"
	"github.com/owenrumney/squealer/pkg/squealer"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		Provider:  provider.AWSProvider,
		Service:   "launch",
		ShortCode: "no-sensitive-info",
		Documentation: rule.RuleDocumentation{
			Summary:     "Ensure all data stored in the Launch configuration EBS is securely encrypted",
			Explanation: `When creating Launch Configurations, user data can be used for the initial configuration of the instance. User data must not contain any sensitive data.`,
			Impact:      "Sensitive credentials in user data can be leaked",
			Resolution:  "Don't use sensitive data in user data",
			BadExample: []string{`
resource "aws_launch_configuration" "as_conf" {
  name          = "web_config"
  image_id      = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
  user_data     = <<EOF
export DATABASE_PASSWORD=\"SomeSortOfPassword\"
EOF
}
`, `
resource "aws_launch_configuration" "as_conf" {
  name             = "web_config"
  image_id         = data.aws_ami.ubuntu.id
  instance_type    = "t2.micro"
  user_data_base64 = "ZXhwb3J0IERBVEFCQVNFX1BBU1NXT1JEPSJTb21lU29ydE9mUGFzc3dvcmQi"
}
`},
			GoodExample: []string{`
resource "aws_launch_configuration" "as_conf" {
  name          = "web_config"
  image_id      = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
  user_data     = <<EOF
export GREETING="Hello there"
EOF
}
`, `
resource "aws_launch_configuration" "as_conf" {
	name             = "web_config"
	image_id         = data.aws_ami.ubuntu.id
	instance_type    = "t2.micro"
	user_data_base64 = "ZXhwb3J0IEVESVRPUj12aW1hY3M="
  }
  `,
			},
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/launch_configuration#user_data,user_data_base64",
			},
		},
		RequiredTypes: []string{
			"resource",
		},
		RequiredLabels: []string{
			"aws_launch_configuration",
		},
		DefaultSeverity: severity.High,
		CheckFunc: func(set result.Set, resourceBlock block.Block, _ *hclcontext.Context) {
			// TODO: code goes here

			if resourceBlock.MissingChild("user_data") && resourceBlock.MissingChild("user_data_base64") {
				return
			}

			customDataAttr := resourceBlock.GetAttribute("user_data")
			customDataBase64Attr := resourceBlock.GetAttribute("user_data_base64")

			if customDataAttr.IsNotNil() {
				for _, str := range customDataAttr.ValueAsStrings() {
					if checkStringForSensitive(str) {
						set.AddResult().
							WithDescription("Resource '%s' has user_data with sensitive data.", resourceBlock.FullName()).
							WithAttribute(customDataAttr)
					}
				}
			}

			if customDataBase64Attr.IsNotNil() && customDataBase64Attr.IsString() {
				encoded, err := base64.StdEncoding.DecodeString(customDataBase64Attr.Value().AsString())
				if err != nil {
					debug.Log("could not decode the base64 string in the terraform, trying with the string verbatim")
					encoded = []byte(customDataAttr.Value().AsString())
				}
				if checkStringForSensitive(string(encoded)) {
					set.AddResult().
						WithDescription("Resource '%s' has user_data_base64 with sensitive data.", resourceBlock.FullName()).
						WithAttribute(customDataAttr)
				}
			}
		},
	})
}

func checkStringForSensitive(stringToCheck string) bool {
	scanResult := squealer.NewStringScanner().Scan(stringToCheck)
	return scanResult.TransgressionFound
}
