package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSnapshots_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_snapshots.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSnapshots_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "snapshots.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshots.0.status"),
					resource.TestMatchResourceAttr(dataSource, "snapshots.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "snapshots.0.finished_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_existed_snapshot", "true"),
				),
			},
		},
	})
}

func testDataSourceSnapshots_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_snapshot" "test" {
  cluster_id  = "%[1]s"
  name        = "%[2]s"
  description = "Created by terraform script"
}

data "huaweicloud_dws_snapshots" "test" {
  depends_on = [
    huaweicloud_dws_snapshot.test
  ]
}

locals {
  cluster_ids  = data.huaweicloud_dws_snapshots.test.snapshots[*].cluster_id
  snapshot_ids = data.huaweicloud_dws_snapshots.test.snapshots[*].id
}

output "is_existed_snapshot" {
  value = contains(local.cluster_ids, "%[1]s") && contains(local.snapshot_ids, huaweicloud_dws_snapshot.test.id)
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}
