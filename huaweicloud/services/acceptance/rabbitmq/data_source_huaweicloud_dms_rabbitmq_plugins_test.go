package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDmsRabbitmqPlugins_basic(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.huaweicloud_dms_rabbitmq_plugins.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDmsRabbitmqPlugins_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "plugins.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "plugins.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "plugins.0.enable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "plugins.0.running"),
					resource.TestCheckResourceAttrSet(dataSourceName, "plugins.0.version"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("enable_filter_is_useful", "true"),
					resource.TestCheckOutput("running_filter_is_useful", "true"),
					resource.TestCheckOutput("version_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDmsRabbitmqPlugins_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_rabbitmq_plugins" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
}

locals {
  name    = data.huaweicloud_dms_rabbitmq_plugins.test.plugins[0].name
  enable  = data.huaweicloud_dms_rabbitmq_plugins.test.plugins[0].enable
  running = data.huaweicloud_dms_rabbitmq_plugins.test.plugins[0].running
  version = data.huaweicloud_dms_rabbitmq_plugins.test.plugins[0].version
}

data "huaweicloud_dms_rabbitmq_plugins" "name_filter" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  name        = local.name
}
  
output "name_filter_is_useful" {
  value = length(data.huaweicloud_dms_rabbitmq_plugins.name_filter.plugins) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rabbitmq_plugins.name_filter.plugins[*].name : v == local.name]
  )
}

data "huaweicloud_dms_rabbitmq_plugins" "enable_filter" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  enable      = local.enable
}
	
output "enable_filter_is_useful" {
  value = length(data.huaweicloud_dms_rabbitmq_plugins.enable_filter.plugins) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rabbitmq_plugins.enable_filter.plugins[*].enable : v == local.enable]
  ) 
}

data "huaweicloud_dms_rabbitmq_plugins" "running_filter" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  running     = local.running
}
	  
output "running_filter_is_useful" {
  value = length(data.huaweicloud_dms_rabbitmq_plugins.running_filter.plugins) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rabbitmq_plugins.running_filter.plugins[*].running : v == local.running]
  ) 
}

data "huaweicloud_dms_rabbitmq_plugins" "version_filter" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  version     = local.version
}
		
output "version_filter_is_useful" {
  value = length(data.huaweicloud_dms_rabbitmq_plugins.version_filter.plugins) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rabbitmq_plugins.version_filter.plugins[*].version : v == local.version]
  ) 
}

`, testAccDmsRabbitmqInstance_newFormat_single(name))
}
