package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityServices_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identity_services.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityServices_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "services.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(all, "services.0.name"),
					resource.TestCheckResourceAttrSet(all, "services.0.id"),
					resource.TestCheckResourceAttrSet(all, "services.0.type"),
					resource.TestCheckResourceAttrSet(all, "services.0.link"),
				),
			},
			{
				Config: testAccIdentityServicesById_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "services.0.enabled", "true"),
					resource.TestCheckResourceAttrSet(all, "services.0.name"),
					resource.TestCheckResourceAttrSet(all, "services.0.id"),
					resource.TestCheckResourceAttrSet(all, "services.0.type"),
					resource.TestCheckResourceAttrSet(all, "services.0.link"),
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
