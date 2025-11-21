package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5RegisteredServices_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_registered_services.services"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5RegisteredServices_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "service_codes.0"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5RegisteredServices_basic() string {
	return `
data "huaweicloud_identityv5_registered_services" "services" {}
`
}
