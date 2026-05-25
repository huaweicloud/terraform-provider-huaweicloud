package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsInstanceTopology_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dcs_instance_topology.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDcsInstanceTopology_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.node_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.max_memory"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.used_memory"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.qps"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.bandwidth.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.bandwidth.0.input"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.bandwidth.0.output"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.db_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dbs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dbs.0.db_idx"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dbs.0.keys"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dbs.0.expires"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dbs.0.avg_ttl"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.relation_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.relation_port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dims.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dims.0.dim_k"),
					resource.TestCheckResourceAttrSet(dataSourceName, "redis_server.0.dims.0.dim_v"),

					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_proxy.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_proxy.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_proxy.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_proxy.0.ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_proxy.0.port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_proxy.0.node_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_proxy.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_proxy.0.dims.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_proxy.0.dims.0.dim_k"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_proxy.0.dims.0.dim_v"),

					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_endpoint.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_endpoint.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_endpoint.0.ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_endpoint.0.port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_endpoint.0.node_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_endpoint.0.status"),

					resource.TestCheckResourceAttrSet(dataSourceName, "elb.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "elb.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "elb.0.ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "elb.0.port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "elb.0.node_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "elb.0.status"),
				),
			},
		},
	})
}

func testAccDataSourceDcsInstanceTopology_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dcs_instance_topology" "test" {
  instance_id = "%s"
}
`, acceptance.HW_DCS_INSTANCE_ID)
}
