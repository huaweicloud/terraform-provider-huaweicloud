package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAsGroupQuotas_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_as_group_quotas.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare AS group ID in advance.
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceAsGroupQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.max"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.min"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.quota"),
				),
			},
		},
	})
}

func testDataSourceDataSourceAsGroupQuotas_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_as_group_quotas" "test" {
  scaling_group_id = "%s"
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}
