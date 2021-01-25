package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccVPCEPPublicServicesDataSource_Basic(t *testing.T) {
	resourceName := "data.huaweicloud_vpcep_public_services.services"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPPublicServicesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "services.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "services.0.service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "services.0.service_type"),
				),
			},
		},
	})
}

var testAccVPCEPPublicServicesDataSourceBasic = `
data "huaweicloud_vpcep_public_services" "services" {
  service_name = "dns"
}
`
