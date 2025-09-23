package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsInstanceShards_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_instance_shards.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcsInstanceShards_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.replication_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.replication_list.0.replication_role"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.replication_list.0.replication_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.replication_list.0.is_replication"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.replication_list.0.replication_id"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.replication_list.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.replication_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.replication_list.0.az_code"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.replication_list.0.dimensions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.replication_list.0.dimensions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "group_list.0.replication_list.0.dimensions.0.value"),
				),
			},
		},
	})
}

func testDataSourceDcsInstanceShards_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_instance_shards" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
}
`, testAccDcsV1Instance_basic(name))
}
