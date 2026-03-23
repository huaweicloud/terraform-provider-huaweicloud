package mrs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataClusterFiles_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_mapreduce_cluster_files.test"
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
				Config: testAccDataClusterFiles_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "files.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "files.0.path_suffix"),
					resource.TestCheckResourceAttrSet(all, "files.0.owner"),
					resource.TestCheckResourceAttrSet(all, "files.0.group"),
					resource.TestCheckResourceAttrSet(all, "files.0.permission"),
					resource.TestCheckResourceAttrSet(all, "files.0.replication"),
					resource.TestCheckResourceAttrSet(all, "files.0.type"),
					// The block_size, length, children_num, access_time, and modification_time attributes are not guaranteed
					// to be non-empty, so they are not validated.
				),
			},
		},
	})
}

func testAccDataClusterFiles_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_mapreduce_cluster_files" "test" {
  cluster_id = "%s"
  path       = "user"
}
`, acceptance.HW_MRS_CLUSTER_ID)
}
