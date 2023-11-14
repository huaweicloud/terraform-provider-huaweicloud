package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceSecurityPolicies_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_security_policies.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSecurityPolicies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "security_policies.#"),
					resource.TestCheckResourceAttrSet(rName, "security_policies.0.name"),
					resource.TestCheckResourceAttrSet(rName, "security_policies.0.type"),
					resource.TestCheckResourceAttrSet(rName, "security_policies.0.protocols.0"),
					resource.TestCheckResourceAttrSet(rName, "security_policies.0.ciphers.0"),
					resource.TestCheckOutput("custom_name_filter_is_useful", "true"),
					resource.TestCheckOutput("custom_id_filter_is_useful", "true"),
					resource.TestCheckOutput("custom_type_filter_is_useful", "true"),
					resource.TestCheckOutput("custom_protocol_filter_is_useful", "true"),
					resource.TestCheckOutput("custom_cipher_filter_is_useful", "true"),
					resource.TestCheckOutput("system_type_filter_is_useful", "true"),
					resource.TestCheckOutput("system_name_filter_is_useful", "true"),
					resource.TestCheckOutput("system_protocol_filter_is_useful", "true"),
					resource.TestCheckOutput("system_cipher_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceSecurityPolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_security_policies" "test" {
  depends_on = [huaweicloud_elb_security_policy.test]
}

data "huaweicloud_elb_security_policies" "custom_name_filter" {
  depends_on = [huaweicloud_elb_security_policy.test]
  name       = "%[2]s"
}
output "custom_name_filter_is_useful" {
  value = length(data.huaweicloud_elb_security_policies.custom_name_filter.security_policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_security_policies.custom_name_filter.security_policies[*].name :v == "%[2]s"]
  )  
}

locals {
  id = huaweicloud_elb_security_policy.test.id
}
data "huaweicloud_elb_security_policies" "custom_id_filter" {
  security_policy_id = huaweicloud_elb_security_policy.test.id
}
output "custom_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_security_policies.custom_id_filter.security_policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_security_policies.custom_id_filter.security_policies[*].id :v == local.id]
)
}

data "huaweicloud_elb_security_policies" "custom_type_filter" {
  depends_on = [huaweicloud_elb_security_policy.test]
  type       = "custom"
}
output "custom_type_filter_is_useful" {
  value = length(data.huaweicloud_elb_security_policies.custom_type_filter.security_policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_security_policies.custom_type_filter.security_policies[*].type :v == "custom"]
)
}

locals {
  protocol = huaweicloud_elb_security_policy.test.protocols.0
}
data "huaweicloud_elb_security_policies" "custom_protocol_filter" {
  protocol = huaweicloud_elb_security_policy.test.protocols.0
}
output "custom_protocol_filter_is_useful" {
  value = length(data.huaweicloud_elb_security_policies.custom_protocol_filter.security_policies) > 0
}

locals {
  cipher = huaweicloud_elb_security_policy.test.ciphers.0
}
data "huaweicloud_elb_security_policies" "custom_cipher_filter" {
  cipher = huaweicloud_elb_security_policy.test.ciphers.0
}
output "custom_cipher_filter_is_useful" {
  value = length(data.huaweicloud_elb_security_policies.custom_cipher_filter.security_policies) > 0 
}

data "huaweicloud_elb_security_policies" "systemTest" {
}

data "huaweicloud_elb_security_policies" "system_type_filter" {
  type = "system"
}
output "system_type_filter_is_useful" {
  value = length(data.huaweicloud_elb_security_policies.system_type_filter.security_policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_security_policies.system_type_filter.security_policies[*].type :v == "system"]
)
}

data "huaweicloud_elb_security_policies" "system_name_filter" {
  name = data.huaweicloud_elb_security_policies.systemTest.security_policies.0.name
}
output "system_name_filter_is_useful" {
  value = length(data.huaweicloud_elb_security_policies.system_name_filter.security_policies) > 0 && alltrue(
  [for v in data.huaweicloud_elb_security_policies.system_name_filter.security_policies[*].name :
  v == data.huaweicloud_elb_security_policies.system_name_filter.name]
  )  
}

data "huaweicloud_elb_security_policies" "system_protocol_filter" {
  protocol = data.huaweicloud_elb_security_policies.systemTest.security_policies.0.protocols.0
}
output "system_protocol_filter_is_useful" {
  value = length(data.huaweicloud_elb_security_policies.system_protocol_filter.security_policies) > 0
}

data "huaweicloud_elb_security_policies" "system_cipher_filter" {
  cipher = data.huaweicloud_elb_security_policies.systemTest.security_policies.0.ciphers.0
}
output "system_cipher_filter_is_useful" {
  value = length(data.huaweicloud_elb_security_policies.system_cipher_filter.security_policies) > 0 
}

`, testSecurityPoliciesV3_basic(name), name)
}
