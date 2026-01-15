package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceInstanceDiskUsage_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dds_instance_disk_usage.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceInstanceDiskUsage_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.entity_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.entity_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.group_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.used"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volumes.0.size"),
				),
			},
		},
	})
}

func testAccDatasourceInstanceDiskUsage_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_instance_disk_usage" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_DDS_INSTANCE_ID)
}
