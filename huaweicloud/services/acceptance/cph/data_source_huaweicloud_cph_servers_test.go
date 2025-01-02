package cph

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCphServers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cph_servers.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCphServers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.server_name"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.server_id"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.network_version"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.phone_flavor"),
					resource.TestCheckResourceAttrSet(dataSource, "servers.0.status"),

					resource.TestCheckOutput("is_server_name_filter_useful", "true"),
					resource.TestCheckOutput("is_server_id_filter_useful", "true"),
					resource.TestCheckOutput("is_network_version_filter_useful", "true"),
					resource.TestCheckOutput("is_phone_flavor_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
			{
				Config: testCphServerBase(rName),
				Check: resource.ComposeTestCheckFunc(
					waitForDeletionCooldownComplete(),
				),
			},
		},
	})
}

func waitForDeletionCooldownComplete() resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		// lintignore:R018
		time.Sleep(10 * time.Minute)
		return nil
	}
}

func testDataSourceCphServers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cph_servers" "test" {
  depends_on = [huaweicloud_cph_server.test]
}

locals {
  server_name     = data.huaweicloud_cph_servers.test.servers[0].server_name
  server_id       = data.huaweicloud_cph_servers.test.servers[0].server_id
  network_version = data.huaweicloud_cph_servers.test.servers[0].network_version
  phone_flavor    = data.huaweicloud_cph_servers.test.servers[0].phone_flavor
  status          = tostring(data.huaweicloud_cph_servers.test.servers[0].status)
}

data "huaweicloud_cph_servers" "filter_by_server_name" {
  server_name = local.server_name
}

data "huaweicloud_cph_servers" "filter_by_server_id" {
  server_id = local.server_id
}

data "huaweicloud_cph_servers" "filter_by_network_version" {
  network_version = local.network_version
}

data "huaweicloud_cph_servers" "filter_by_phone_flavor" {
  phone_flavor = local.phone_flavor
}

data "huaweicloud_cph_servers" "filter_by_status" {
  status = local.status
}

locals {
  list_by_server_name     = data.huaweicloud_cph_servers.filter_by_server_name.servers
  list_by_server_id       = data.huaweicloud_cph_servers.filter_by_server_id.servers
  list_by_network_version = data.huaweicloud_cph_servers.filter_by_network_version.servers
  list_by_phone_flavor    = data.huaweicloud_cph_servers.filter_by_phone_flavor.servers
  list_by_status          = data.huaweicloud_cph_servers.filter_by_status.servers
}

output "is_server_name_filter_useful" {
  value = length(local.list_by_server_name) > 0 && alltrue(
    [for v in local.list_by_server_name[*].server_name : v == local.server_name]
  )
}

output "is_server_id_filter_useful" {
  value = length(local.list_by_server_id) > 0 && alltrue(
    [for v in local.list_by_server_id[*].server_id : v == local.server_id]
  )
}

output "is_network_version_filter_useful" {
  value = length(local.list_by_network_version) > 0 && alltrue(
    [for v in local.list_by_network_version[*].network_version : v == local.network_version]
  )
}

output "is_phone_flavor_filter_useful" {
  value = length(local.list_by_phone_flavor) > 0 && alltrue(
    [for v in local.list_by_phone_flavor[*].phone_flavor : v == local.phone_flavor]
  )
}

output "is_status_filter_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : tostring(v) == local.status]
  )
}
`, testCphServer_basic(name))
}
