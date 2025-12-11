package mrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccClusterNodeBatchShrink_basic(t *testing.T) {
	var (
		rNameWithExpand = "huaweicloud_mapreduce_cluster_node_batch_expand.test"
		rNameWithShrink = "huaweicloud_mapreduce_cluster_node_batch_shrink.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterID(t)
			acceptance.TestAccPreCheckMrsClusterNodeGroupName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time batch action resource and there is no logic in the delete method.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterNodeBatchShrinkConfig_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rNameWithExpand, "cluster_id"),
					resource.TestCheckResourceAttrSet(rNameWithExpand, "node_group_name"),
					resource.TestCheckResourceAttr(rNameWithExpand, "node_count", "2"),
				),
			},
			{
				Config: testAccClusterNodeBatchShrinkConfig_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rNameWithShrink, "cluster_id"),
					resource.TestCheckResourceAttrSet(rNameWithShrink, "node_group_name"),
					resource.TestCheckResourceAttr(rNameWithShrink, "node_count", "2"),
				),
			},
		},
	})
}

func testAccClusterNodeBatchShrinkConfig_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_mapreduce_cluster_node_batch_expand" "test" {
  cluster_id      = "%[1]s"
  node_group_name = "%[2]s"
  node_count      = 2
}
`, acceptance.HW_MRS_CLUSTER_ID, acceptance.HW_MRS_CLUSTER_NODE_GROUP_NAME)
}

func testAccClusterNodeBatchShrinkConfig_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_mapreduce_cluster_node_batch_shrink" "test" {
  cluster_id      = "%[1]s"
  node_group_name = "%[2]s"
  node_count      = 2
}
`, acceptance.HW_MRS_CLUSTER_ID, acceptance.HW_MRS_CLUSTER_NODE_GROUP_NAME)
}
