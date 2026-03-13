package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccV3Services_basic(t *testing.T) {
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
				Config: testAccV3Services_basic(),
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
				Config: testAccV3ServicesById_basic(),
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

func testAccV3Services_basic() string {
	return `
data "huaweicloud_identity_services" "test" {}
`
}
func testAccV3ServicesById_basic() string {
	return `
data "huaweicloud_identity_services" "all" {}

data "huaweicloud_identity_services" "test" {
  service_id = data.huaweicloud_identity_services.all.services[0].id
}
`
}
