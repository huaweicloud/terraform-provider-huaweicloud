package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePublicipInstancesCount_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_eip_publicip_instances_count.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePublicipInstancesCount_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_num"),
				),
			},
		},
	})
}

func testDataSourcePublicipInstancesCount_basic() string {
	return `
data "huaweicloud_eip_publicip_instances_count" "test" {}
`
}
