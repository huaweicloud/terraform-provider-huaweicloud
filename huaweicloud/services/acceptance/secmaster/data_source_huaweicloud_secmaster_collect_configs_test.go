package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCollectConfigs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_collect_configs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCollectConfigs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "workspace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "dataspace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "dataspace_name"),
					resource.TestCheckResourceAttrSet(dataSource, "datasets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "all_vendors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "config_statistics.#"),
				),
			},
		},
	})
}

func testAccDataSourceCollectConfigs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_collect_configs" "test" {
  region_id        = "%[1]s"
  domain_id        = "%[2]s"
  query_statistics = true
}
`, acceptance.HW_REGION_NAME, acceptance.HW_DOMAIN_ID)
}
