package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataProviderProtocols_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		byProviderId   = "data.huaweicloud_identity_provider_protocols.filter_by_provider_id"
		dcByProviderId = acceptance.InitDataSourceCheck(byProviderId)

		byProtocolId   = "data.huaweicloud_identity_provider_protocols.filter_by_protocol_id"
		dcByProtocolId = acceptance.InitDataSourceCheck(byProtocolId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataProviderProtocols_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcByProviderId.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_provider_id_result", "true"),
					dcByProtocolId.CheckResourceExists(),
					resource.TestCheckOutput("filter_by_protocol_id_result", "true"),
				),
			},
		},
	})
}

func testAccDataProviderProtocols_base(name string) string {
	return fmt.Sprintf(`
variable "provider_protocols" {
  type    = list(string)
  default = ["oidc", "saml"]
}

resource "huaweicloud_identity_provider" "test" {
  count = length(var.provider_protocols)

  name     = format("%[1]s_%%s", var.provider_protocols[count.index])
  protocol = var.provider_protocols[count.index]
}
`, name)
}

func testAccDataProviderProtocols_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Filter by provider_id
locals {
  provider_id = try(huaweicloud_identity_provider.test[0].id, "NOT_FOUND")
}

data "huaweicloud_identity_provider_protocols" "filter_by_provider_id" {
  provider_id = local.provider_id
}

output "filter_by_provider_id_result" {
  value = length(data.huaweicloud_identity_provider_protocols.filter_by_provider_id.protocols) == 1
}

# Filter by protocol_id
data "huaweicloud_identity_provider_protocols" "filter_by_protocol_id" {
  provider_id = local.provider_id
  protocol_id = "oidc"
}

locals {
  filter_by_protocol_id_result = length(data.huaweicloud_identity_provider_protocols.filter_by_protocol_id.protocols) == 1
}

output "filter_by_protocol_id_result" {
  value = local.filter_by_protocol_id_result
}
`, testAccDataProviderProtocols_base(name))
}
