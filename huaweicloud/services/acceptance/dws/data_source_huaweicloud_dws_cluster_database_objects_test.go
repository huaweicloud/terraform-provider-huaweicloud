package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataClusterDatabaseObjects_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dws_cluster_database_objects.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dws_cluster_database_objects.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataClusterDatabaseObjects_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "objects.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "objects.0.obj_name"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataClusterDatabaseObjects_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_database_objects" "test" {
  cluster_id = "%s"
  type       = "DATABASE"
}

locals {
  first_object_name = data.huaweicloud_dws_cluster_database_objects.test.objects[0].obj_name
}

data "huaweicloud_dws_cluster_database_objects" "filter_by_name" {
  cluster_id = "%s"
  type       = "DATABASE"
  name       = local.first_object_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dws_cluster_database_objects.filter_by_name.objects[*].obj_name :
    v == local.first_object_name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}
`, acceptance.HW_DWS_CLUSTER_ID, acceptance.HW_DWS_CLUSTER_ID)
}
