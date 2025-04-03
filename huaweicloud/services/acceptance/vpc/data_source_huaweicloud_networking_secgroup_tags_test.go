package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcNetworkingSecgroupTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_networking_secgroup_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcNetworkingSecgroupTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcNetworkingSecgroupTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup_tags" "test" {
  depends_on = [ huaweicloud_networking_secgroup.secgroup_1 ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_networking_secgroup_tags.test.tags) > 0
}
`, testAccSecGroup_basic(name))
}
