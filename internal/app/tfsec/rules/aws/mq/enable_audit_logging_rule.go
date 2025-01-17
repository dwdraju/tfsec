package mq

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
		Provider:       provider.AWSProvider,
		Service:   "mq",
		ShortCode: "enable-audit-logging",
		Documentation: rule.RuleDocumentation{
			Summary:     "MQ Broker should have audit logging enabled",
			Explanation: `Logging should be enabled to allow tracing of issues and activity to be investigated more fully. Logs provide additional information and context which is often invalauble during investigation`,
			Impact:      "Without audit logging it is difficult to trace activity in the MQ broker",
			Resolution:  "Enable audit logging",
			BadExample: []string{  `
resource "aws_mq_broker" "bad_example" {
  broker_name = "example"

  configuration {
    id       = aws_mq_configuration.test.id
    revision = aws_mq_configuration.test.latest_revision
  }

  engine_type        = "ActiveMQ"
  engine_version     = "5.15.0"
  host_instance_type = "mq.t2.micro"
  security_groups    = [aws_security_group.test.id]

  user {
    username = "ExampleUser"
    password = "MindTheGap"
  }
  logs {
    audit = false
  }
}
`},
			GoodExample: []string{ `
resource "aws_mq_broker" "good_example" {
  broker_name = "example"

  configuration {
    id       = aws_mq_configuration.test.id
    revision = aws_mq_configuration.test.latest_revision
  }

  engine_type        = "ActiveMQ"
  engine_version     = "5.15.0"
  host_instance_type = "mq.t2.micro"
  security_groups    = [aws_security_group.test.id]

  user {
    username = "ExampleUser"
    password = "MindTheGap"
  }
  logs {
    audit = true
  }
}
`},
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/mq_broker#audit",
			},
		},
		RequiredTypes:  []string{ 
			"resource",
		},
		RequiredLabels: []string{ 
			"aws_mq_broker",
		},
		DefaultSeverity: severity.Medium, 
		CheckFunc: func(set result.Set, resourceBlock block.Block, _ *hclcontext.Context){
			if auditAttr := resourceBlock.GetBlock("logs").GetAttribute("audit"); auditAttr.IsNil() { // alert on use of default value
				set.AddResult().
					WithDescription("Resource '%s' uses default value for logs.audit", resourceBlock.FullName())
			} else if auditAttr.IsFalse() {
				set.AddResult().
					WithDescription("Resource '%s' does not have logs.audit set to true", resourceBlock.FullName()).
					WithAttribute(auditAttr)
			}
		},
	})
}
