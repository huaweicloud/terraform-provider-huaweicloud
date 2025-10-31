package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityServices_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identity_services.test"

	config := testAccIdentityServices_basic()
	configById := testAccIdentityServicesById_basic()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSource, "services.0.enabled", "true"),
				),
			},
			{
				Config: configById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSource, "services.0.enabled", "true"),
				),
			},
		},
	})
}

func testAccIdentityServices_basic() string {
	return `
data "huaweicloud_identity_services" "test" {}
`
}
func testAccIdentityServicesById_basic() string {
	return `
data "huaweicloud_identity_services" "all" {}

data "huaweicloud_identity_services" "test" {
  service_id = data.huaweicloud_identity_services.all.services[0].id
}
`
}
