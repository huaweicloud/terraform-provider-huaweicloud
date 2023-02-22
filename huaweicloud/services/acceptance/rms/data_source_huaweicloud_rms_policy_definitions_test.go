package rms

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataPolicyDefinitions_basic(t *testing.T) {
	var (
		dName = "data.huaweicloud_rms_policy_definitions.test"
		dc    = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPolicyDefinitions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dName, "definitions.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

const testAccDataPolicyDefinitions_basic = `
data "huaweicloud_rms_policy_definitions" "test" {
  name = "allowed-ecs-flavors"
}
`

func TestAccDataPolicyDefinitions_keywords(t *testing.T) {
	var (
		dName = "data.huaweicloud_rms_policy_definitions.test"
		dc    = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPolicyDefinitions_keywords,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dName, "definitions.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

const testAccDataPolicyDefinitions_keywords = `
data "huaweicloud_rms_policy_definitions" "test" {
  keywords = ["ecs"]
}
`

func TestAccDataPolicyDefinitions_policyRuleType(t *testing.T) {
	var (
		dName = "data.huaweicloud_rms_policy_definitions.test"
		dc    = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPolicyDefinitions_policyRuleType,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dName, "definitions.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

const testAccDataPolicyDefinitions_policyRuleType = `
data "huaweicloud_rms_policy_definitions" "test" {
  policy_rule_type = "dsl"
}
`

func TestAccDataPolicyDefinitions_triggerType(t *testing.T) {
	var (
		dName = "data.huaweicloud_rms_policy_definitions.test"
		dc    = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPolicyDefinitions_triggerType,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dName, "definitions.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

const testAccDataPolicyDefinitions_triggerType = `
data "huaweicloud_rms_policy_definitions" "test" {
  trigger_type = "resource"
}
`
