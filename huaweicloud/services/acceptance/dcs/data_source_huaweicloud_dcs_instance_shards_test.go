package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDcsInstanceShard_basic(t *testing.T) {
	rName := "data.huaweicloud_dcs_instance_shards.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDcsInstanceShard_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "shards.#", "1"),
					resource.TestCheckResourceAttr(rName, "shards.0.shard_name", "group-0"),
					resource.TestCheckResourceAttr(rName, "shards.0.replicas.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "shards.0.shard_id"),
					resource.TestCheckResourceAttrSet(rName, "shards.0.shard_name"),
					resource.TestCheckResourceAttrSet(rName, "shards.0.replicas.0.id"),
					resource.TestCheckResourceAttrSet(rName, "shards.0.replicas.0.ip"),
					resource.TestCheckResourceAttrSet(rName, "shards.0.replicas.0.role"),
					resource.TestCheckResourceAttrSet(rName, "shards.0.replicas.0.node_id"),
					resource.TestCheckResourceAttrSet(rName, "shards.0.replicas.0.status"),
				),
			},
		},
	})
}

func testAccDatasourceDcsInstanceShard_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_instance_shards" "test" {
  instance_id  = huaweicloud_dcs_instance.instance_1.id
  replica_role = "slave"
  shard_names  = ["group-0"]
}
`, testAccDcsV1Instance_basic(name))
}
