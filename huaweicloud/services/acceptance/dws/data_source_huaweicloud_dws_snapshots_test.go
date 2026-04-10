package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSnapshots_basic(t *testing.T) {
	var (
		all     = "data.huaweicloud_dws_snapshots.all"
		dcByAll = acceptance.InitDataSourceCheck(all)

		byClusterId         = "data.huaweicloud_dws_snapshots.filter_by_cluster_id"
		dcFilterByClusterId = acceptance.InitDataSourceCheck(byClusterId)

		rName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSnapshots_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dcByAll.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "snapshots.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "snapshots.0.id"),
					resource.TestCheckResourceAttrSet(all, "snapshots.0.name"),
					resource.TestCheckResourceAttrSet(all, "snapshots.0.cluster_id"),
					resource.TestCheckResourceAttrSet(all, "snapshots.0.type"),
					resource.TestCheckResourceAttrSet(all, "snapshots.0.status"),
					resource.TestMatchResourceAttr(all, "snapshots.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "snapshots.0.finished_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_existed_snapshot", "true"),

					// Filter by cluster ID
					dcFilterByClusterId.CheckResourceExists(),
					resource.TestCheckOutput("is_cluster_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSnapshots_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_snapshot" "test" {
  cluster_id  = "%[1]s"
  name        = "%[2]s"
  description = "Created by terraform script"
}

data "huaweicloud_dws_snapshots" "all" {
  depends_on = [
    huaweicloud_dws_snapshot.test
  ] 
}

locals {
  cluster_ids  = data.huaweicloud_dws_snapshots.all.snapshots[*].cluster_id
  snapshot_ids = data.huaweicloud_dws_snapshots.all.snapshots[*].id
}

output "is_existed_snapshot" {
  value = contains(local.cluster_ids, "%[1]s") && contains(local.snapshot_ids, huaweicloud_dws_snapshot.test.id)
}

# Filter by cluster ID
data "huaweicloud_dws_snapshots" "filter_by_cluster_id" {
  cluster_id = "%[1]s"

  depends_on = [
    huaweicloud_dws_snapshot.test
  ]
}

locals {
  snapshot_ids_by_cluster_id = data.huaweicloud_dws_snapshots.filter_by_cluster_id.snapshots[*].id
}

output "is_cluster_id_filter_useful" {
  value = contains(local.snapshot_ids_by_cluster_id, huaweicloud_dws_snapshot.test.id)
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}
