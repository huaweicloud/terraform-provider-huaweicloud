package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsSecondaryDimMonitoredObjects_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dcs_secondary_dim_monitored_objects.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDcsSecondaryDimMonitoredObjects_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "dim_name", "dcs_instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "router.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.dim_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "children.0.dim_route"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "dcs_cluster_redis_node.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "dcs_cluster_redis_node.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "dcs_cluster_redis_node.0.status"),
				),
			},
		},
	})
}

func testAccDataSourceDcsSecondaryDimMonitoredObjects_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_secondary_dim_monitored_objects" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  dim_name    = "dcs_instance_id"
}
`, testAccDcsInstance_base(name))
}
