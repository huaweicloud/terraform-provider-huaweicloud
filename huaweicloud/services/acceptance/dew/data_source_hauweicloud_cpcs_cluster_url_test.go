package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterUrl_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cpcs_cluster_url.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckCpcsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceClusterUrl_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "uri"),
				),
			},
		},
	})
}

func testAccDataSourceClusterUrl_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cpcs_cluster_url" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_CPCS_CLUSTER_ID)
}
