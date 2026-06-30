package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataClusterSnapshotFlavors_basic(t *testing.T) {
	var (
		name   = acceptance.RandomAccResourceName()
		dcName = "data.huaweicloud_dws_cluster_snapshot_flavors.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// A non-existent snapshot ID returns HTTP 200 with an empty flavors list (no API error).
				Config: testAccDataClusterSnapshotFlavors_nonexistentSnapshot(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "flavors.#", "0"),
				),
			},
			{
				Config: testAccDataClusterSnapshotFlavors_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "flavors.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.code"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.classify"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.scenario"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.version"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.status"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.default_capacity"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.duplicate"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.default_node"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.min_node"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.max_node"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.flavor_code"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.volume_num"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.attribute.#"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.product_version_list.#"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.volume_used.#"),
				),
			},
		},
	})
}

func testAccDataClusterSnapshotFlavors_nonexistentSnapshot() string {
	snapshotId, _ := uuid.NewRandom()
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_snapshot_flavors" "test" {
  snapshot_id = "%[1]s"
}
`, snapshotId.String())
}

func testAccDataClusterSnapshotFlavors_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_snapshot" "test" {
  cluster_id = "%[1]s"
  name       = "%[2]s"
}

data "huaweicloud_dws_cluster_snapshot_flavors" "test" {
  snapshot_id = huaweicloud_dws_snapshot.test.id

  depends_on  = [huaweicloud_dws_snapshot.test]
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}
