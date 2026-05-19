package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsClusterInstanceTopology_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dcs_cluster_instance_topology.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDcsClusterInstanceTopology_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.bandwidth.0.input"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.bandwidth.0.output"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.db_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dbs.0.avg_ttl"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dbs.0.db_idx"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dbs.0.expires"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dbs.0.keys"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dims.0.dim_k"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dims.0.dim_v"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.max_memory"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.node_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.qps"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.relation_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.relation_port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.used_memory"),
				),
			},
		},
	})
}

func testAccDataSourceDcsClusterInstanceTopology_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_cluster_instance_topology" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  start       = 0
  limit       = 100
}
`, testAccDcsInstance_base(name))
}
