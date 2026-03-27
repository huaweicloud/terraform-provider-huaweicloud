package mrs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, please ensure that the cluster has synchronized IAM users and user groups.
func TestAccDataClusterSynchronizedIamObjects_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_mapreduce_cluster_synchronized_iam_objects.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataClusterSynchronizedIamObjects_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "user_names.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(all, "group_names.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

func testAccDataClusterSynchronizedIamObjects_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_mapreduce_cluster_synchronized_iam_objects" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_MRS_CLUSTER_ID)
}
