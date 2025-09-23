package sms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmsAgentConfigs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sms_agent_configs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSmsAgentConfigs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "config.disktype"),
					resource.TestCheckResourceAttrSet(dataSource, "config.mainRegion"),
					resource.TestCheckResourceAttrSet(dataSource, "config.obs_domain"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.cloud_type"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.create_linux_image_id"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.desc"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.dns_domain"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.ecs_domain"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.evs_domain"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.iam_domain"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.ims_domain"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.kms_address"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.linux_image_id"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.linux_uefi_image_id"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.obs_domain"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.region_name"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.sms_domain"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.vpc_domain"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.windows_image_id"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.windows_ssh_image_id"),
					resource.TestCheckResourceAttrSet(dataSource, "regions.0.windows_uefi_image_id"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSmsAgentConfigs_basic() string {
	return `
data "huaweicloud_sms_agent_configs" "test" {}
`
}
