package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcAuthorizations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_authorizations.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCAuth(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcAuthorizations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.cloud_connection_id"),
					resource.TestCheckResourceAttrSet(dataSource, "authorizations.0.cloud_connection_domain_id"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),
					resource.TestCheckOutput("is_cloud_connection_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcAuthorizations_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  id                  = huaweicloud_cc_authorization.test.id
  name                = huaweicloud_cc_authorization.test.name
  instance_id         = huaweicloud_cc_authorization.test.instance_id
  cloud_connection_id = huaweicloud_cc_authorization.test.cloud_connection_id
}

data "huaweicloud_cc_authorizations" "test" {
  depends_on = [
    huaweicloud_cc_authorization.test,
  ]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_cc_authorizations.test.authorizations) >= 1
}

data "huaweicloud_cc_authorizations" "filter_by_id" {
  authorization_id = local.id
  
  depends_on = [
    huaweicloud_cc_authorization.test,
  ]
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_cc_authorizations.filter_by_id.authorizations) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_authorizations.filter_by_id.authorizations[*] : v.id == local.id]
  )
}

data "huaweicloud_cc_authorizations" "filter_by_name" {
  name = local.name
  
  depends_on = [
    huaweicloud_cc_authorization.test,
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_cc_authorizations.filter_by_name.authorizations) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_authorizations.filter_by_name.authorizations[*] : v.name == local.name]
  )
}

data "huaweicloud_cc_authorizations" "filter_by_instance_id" {
  instance_id = local.instance_id
  
  depends_on = [
    huaweicloud_cc_authorization.test,
  ]
}

output "is_instance_id_filter_useful" {
  value = length(data.huaweicloud_cc_authorizations.filter_by_instance_id.authorizations) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_authorizations.filter_by_instance_id.authorizations[*] : v.instance_id == local.instance_id]
  )
}

data "huaweicloud_cc_authorizations" "filter_by_cloud_connection_id" {
  cloud_connection_id = local.cloud_connection_id
  
  depends_on = [
    huaweicloud_cc_authorization.test,
  ]
}

output "is_cloud_connection_id_filter_useful" {
  value = length(data.huaweicloud_cc_authorizations.filter_by_cloud_connection_id.authorizations) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_authorizations.filter_by_cloud_connection_id.authorizations[*] : v.cloud_connection_id == local.cloud_connection_id]
  )
}
`, testDataSourceCcAuthorizations_base(name))
}

func testDataSourceCcAuthorizations_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_cc_authorization" "test" {
  name                       = "%[1]s"
  instance_type              = "vpc"
  instance_id                = huaweicloud_vpc.test.id
  cloud_connection_domain_id = "%[2]s"
  cloud_connection_id        = "%[3]s"
  description                = "This is a test"
}
`, name, acceptance.HW_CC_PEER_DOMAIN_ID, acceptance.HW_CC_PEER_CONNECTION_ID)
}
