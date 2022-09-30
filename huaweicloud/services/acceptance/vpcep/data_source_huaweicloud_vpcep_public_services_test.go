package vpcep

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVPCEPPublicServicesDataSource_Basic(t *testing.T) {
	resourceName := "data.huaweicloud_vpcep_public_services.services"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPPublicServicesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
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
