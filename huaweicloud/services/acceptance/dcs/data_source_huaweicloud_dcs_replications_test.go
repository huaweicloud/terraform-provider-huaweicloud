package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceReplications_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dcs_replications.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDCSInstanceID(t)
			acceptance.TestAccPreCheckDcsGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceReplications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "replications.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replications.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replications.0.replication_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replications.0.is_replication"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replications.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replications.0.status"),
				),
			},
		},
	})
}

func testAccDatasourceReplications_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dcs_replications" "test" {
  instance_id = "%[1]s"
  group_id    = "%[2]s"
}
`, acceptance.HW_DCS_INSTANCE_ID, acceptance.HW_DCS_GROUP_ID)
}
