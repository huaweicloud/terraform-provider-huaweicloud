package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcNetworkAclTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_network_acl_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcNetworkAclTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcNetworkAclTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_network_acl_tags" "test" {
  depends_on = [ huaweicloud_vpc_network_acl.test ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_vpc_network_acl_tags.test.tags) > 0
}
`, testAccNetworkAcl_basic(name))
}
