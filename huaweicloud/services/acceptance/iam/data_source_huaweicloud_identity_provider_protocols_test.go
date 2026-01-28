package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceProviderProtocols_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dcName = "data.huaweicloud_identity_provider_protocols.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		dcNameByProtocolId = "data.huaweicloud_identity_provider_protocols.filter_by_protocol_id"
		dcByProtocolId     = acceptance.InitDataSourceCheck(dcNameByProtocolId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProviderProtocolsBasicStep1(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "protocols.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "protocols.0.id"),
				),
			},
			{
				Config: testAccDataSourceProviderProtocolsBasicStep2(name),
				Check: resource.ComposeTestCheckFunc(
					dcByProtocolId.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcNameByProtocolId, "protocols.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_protocol_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceProviderProtocolsBase(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "test" {
  name     = "%[1]s"
  protocol = "saml"
}
`, name)
}

func testAccDataSourceProviderProtocolsBasicStep1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_provider_protocols" "test" {
  provider_id = huaweicloud_identity_provider.test.id
}
`, testAccDataSourceProviderProtocolsBase(name))
}

func testAccDataSourceProviderProtocolsBasicStep2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_provider_protocols" "filter_by_protocol_id" {
  provider_id = huaweicloud_identity_provider.test.id
  protocol_id = "saml"
}

locals {
  protocol_id_filter_result = [
    for o in data.huaweicloud_identity_provider_protocols.filter_by_protocol_id.protocols : o.id == "saml"
  ]
}

output "is_protocol_id_filter_useful" {
  value = length(local.protocol_id_filter_result) >= 1 && alltrue(local.protocol_id_filter_result)
}
`, testAccDataSourceProviderProtocolsBase(name))
}
