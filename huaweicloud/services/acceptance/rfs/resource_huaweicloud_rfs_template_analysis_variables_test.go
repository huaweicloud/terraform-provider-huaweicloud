package rfs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTemplateAnalysisVariables_basic(t *testing.T) {
	rName := "huaweicloud_rfs_template_analysis_variables.test"
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateAnalysisVariables_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "variables.0.name"),
					resource.TestCheckResourceAttrSet(rName, "variables.0.type"),
					resource.TestCheckResourceAttrSet(rName, "variables.0.description"),
					resource.TestCheckResourceAttrSet(rName, "variables.0.default"),
					resource.TestCheckResourceAttrSet(rName, "variables.0.sensitive"),
					resource.TestCheckResourceAttrSet(rName, "variables.0.nullable"),
				),
			},
		},
	})
}

const testAccTemplateAnalysisVariables_basic = `
resource "huaweicloud_rfs_template_analysis_variables" "test" {
  template_body = <<EOT
variable "instance_name" {
  type        = string
  description = "The name of the instance"
  default     = "test-instance"
}

variable "instance_count" {
  type        = number
  description = "The number of instances to create"
  default     = 1
  validation {
    condition     = var.instance_count > 0
    error_message = "Instance count must be greater than 0."
  }
}

variable "enable_monitoring" {
  type        = bool
  description = "Whether to enable monitoring"
  default     = false
  sensitive   = false
  nullable    = false
}

variable "tags" {
  type        = map(string)
  description = "Tags to apply to the instances"
  default     = {
    Environment = "test"
    Project     = "demo"
  }
}
EOT
}
`
