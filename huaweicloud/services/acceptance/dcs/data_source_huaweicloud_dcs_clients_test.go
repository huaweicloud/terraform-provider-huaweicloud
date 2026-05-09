package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceClients_basic(t *testing.T) {
	rName := "data.huaweicloud_dcs_clients.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceClients_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "clients.#"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.id"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.addr"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.fd"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.name"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.cmd"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.age"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.idle"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.db"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.flags"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.sub"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.psub"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.multi"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.qbuf"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.qbuf_free"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.obl"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.oll"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.omem"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.events"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.network"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.peer"),
					resource.TestCheckResourceAttrSet(rName, "clients.0.user"),
					resource.TestCheckOutput("addr_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_order_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceClients_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  replication_list = data.huaweicloud_dcs_instance_shards.test.group_list[0].replication_list
  node_id          = [for v in local.replication_list : v.node_id if v.replication_role == "slave"][0]
}

resource "huaweicloud_dcs_sessions_query" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
  node_id     = local.node_id
  clean_cache = false
}

data "huaweicloud_dcs_clients" "test" {
  instance_id = huaweicloud_dcs_instance.test.id
  node_id     = local.node_id 

  depends_on  = [huaweicloud_dcs_sessions_query.test]
}

locals {
  test_addr = data.huaweicloud_dcs_clients.test.clients[0].addr
}

data "huaweicloud_dcs_clients" "addr_filter" {
  instance_id = huaweicloud_dcs_instance.test.id
  node_id     = local.node_id
  addr        = local.test_addr

  depends_on  = [data.huaweicloud_dcs_clients.test]
}

output "addr_filter_is_useful" {  
  value = length(data.huaweicloud_dcs_clients.addr_filter.clients) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_clients.addr_filter.clients[*].addr : v == local.test_addr]
  )
}

data "huaweicloud_dcs_clients" "sort_order_filter" {
  instance_id = huaweicloud_dcs_instance.test.id
  node_id     = local.node_id
  sort        = "age"
  order       = "desc"

  depends_on  = [data.huaweicloud_dcs_clients.test]
}

locals {
  has_enough_clients = length(data.huaweicloud_dcs_clients.sort_age_asc.clients) > 1
  asc_first_age      = data.huaweicloud_dcs_clients.sort_age_asc.clients[0].age
  asc_last_age       = data.huaweicloud_dcs_clients.sort_age_asc.clients[length(data.huaweicloud_dcs_clients.sort_age_asc.clients)-1].age
  desc_first_age     = data.huaweicloud_dcs_clients.sort_age_desc.clients[0].age
}

output "sort_order_filter_is_useful" {  
  value = local.has_enough_clients && (local.asc_first_age <= local.asc_last_age) && (local.desc_first_age == local.asc_last_age)
}
`, testAccDcsSessionsQuery_base(name))
}
