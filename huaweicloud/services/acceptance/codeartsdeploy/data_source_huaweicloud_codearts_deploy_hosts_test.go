package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsDeployHosts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_deploy_hosts.test"
	name := acceptance.RandomAccResourceName()
	proxyName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsDeployHosts_basic(proxyName, name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.username"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.trusted_type"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.as_proxy"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.lastest_connection_at"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.connection_status"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.permission.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.permission.0.can_view"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.permission.0.can_edit"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.permission.0.can_delete"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.permission.0.can_add_host"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.permission.0.can_copy"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.owner_id"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.owner_name"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.env_count"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.created_at"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_as_proxy_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCodeartsDeployHosts_basic(proxyName, name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_deploy_hosts" "test" {
  depends_on = [huaweicloud_codearts_deploy_host.test]

  group_id = huaweicloud_codearts_deploy_group.test.id
}

// filter by name
data "huaweicloud_codearts_deploy_hosts" "filter_by_name" {
  group_id = huaweicloud_codearts_deploy_group.test.id
  name     = huaweicloud_codearts_deploy_host.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_codearts_deploy_hosts.filter_by_name.hosts[*].name : 
    v == huaweicloud_codearts_deploy_host.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0
}

// filter by as proxy
data "huaweicloud_codearts_deploy_hosts" "filter_by_as_proxy" {
  group_id = huaweicloud_codearts_deploy_group.test.id
  as_proxy = huaweicloud_codearts_deploy_host.test.as_proxy
}

locals {
  filter_result_by_as_proxy = [for v in data.huaweicloud_codearts_deploy_hosts.filter_by_as_proxy.hosts[*].as_proxy : 
    v == huaweicloud_codearts_deploy_host.test.as_proxy]
}

output "is_as_proxy_filter_useful" {
  value = length(local.filter_result_by_as_proxy) > 0 && alltrue(local.filter_result_by_as_proxy) 
}
`, testDeployHost_withProxyMode(proxyName, name))
}
