package rocketmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDmsRocketMQUsers_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_dms_rocketmq_users.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQSearchUsers(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.access_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.white_remote_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.admin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.default_topic_perm"),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.0.default_group_perm"),
					resource.TestCheckOutput("access_key_filter_is_useful", "true"),
					resource.TestCheckOutput("white_remote_address_filter_is_useful", "true"),
					resource.TestCheckOutput("admin_filter_is_useful", "true"),
					resource.TestCheckOutput("default_topic_perm_filter_is_useful", "true"),
					resource.TestCheckOutput("default_group_perm_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDmsRocketMQSearchUsers(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_rocketmq_users" "test" {
  depends_on  = [huaweicloud_dms_rocketmq_user.test]
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
}

data "huaweicloud_dms_rocketmq_users" "access_key_filter" {
  depends_on  = [huaweicloud_dms_rocketmq_user.test]
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  access_key  = huaweicloud_dms_rocketmq_user.test.access_key
}
  
output "access_key_filter_is_useful" {
  value = length(data.huaweicloud_dms_rocketmq_users.access_key_filter.users) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rocketmq_users.access_key_filter.users[*].access_key : v == huaweicloud_dms_rocketmq_user.test.access_key]
  )  
}

data "huaweicloud_dms_rocketmq_users" "white_remote_address_filter" {
  depends_on           = [huaweicloud_dms_rocketmq_user.test]
  instance_id          = huaweicloud_dms_rocketmq_instance.test.id
  white_remote_address = huaweicloud_dms_rocketmq_user.test.white_remote_address
}

locals {
  white_remote_address = huaweicloud_dms_rocketmq_user.test.white_remote_address
}
	
output "white_remote_address_filter_is_useful" {
  value = length(data.huaweicloud_dms_rocketmq_users.white_remote_address_filter.users) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rocketmq_users.white_remote_address_filter.users[*].white_remote_address : v == local.white_remote_address]
  )  
}

data "huaweicloud_dms_rocketmq_users" "admin_filter" {
  depends_on  = [huaweicloud_dms_rocketmq_user.test]
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  admin       = huaweicloud_dms_rocketmq_user.test.admin
}

locals {
  admin = huaweicloud_dms_rocketmq_user.test.admin
}
	
output "admin_filter_is_useful" {
  value = length(data.huaweicloud_dms_rocketmq_users.admin_filter.users) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rocketmq_users.admin_filter.users[*].admin : v == local.admin]
  )
}

data "huaweicloud_dms_rocketmq_users" "default_topic_perm_filter" {
  depends_on         = [huaweicloud_dms_rocketmq_user.test]
  instance_id        = huaweicloud_dms_rocketmq_instance.test.id
  default_topic_perm = huaweicloud_dms_rocketmq_user.test.default_topic_perm
}
  
locals {
  default_topic_perm = huaweicloud_dms_rocketmq_user.test.default_topic_perm
}
	  
output "default_topic_perm_filter_is_useful" {
  value = length(data.huaweicloud_dms_rocketmq_users.default_topic_perm_filter.users) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rocketmq_users.default_topic_perm_filter.users[*].default_topic_perm : v == local.default_topic_perm]
  )
}

data "huaweicloud_dms_rocketmq_users" "default_group_perm_filter" {
  depends_on         = [huaweicloud_dms_rocketmq_user.test]
  instance_id        = huaweicloud_dms_rocketmq_instance.test.id
  default_group_perm = huaweicloud_dms_rocketmq_user.test.default_group_perm
}

locals {
  default_group_perm = huaweicloud_dms_rocketmq_user.test.default_group_perm
}
	
output "default_group_perm_filter_is_useful" {
  value = length(data.huaweicloud_dms_rocketmq_users.default_group_perm_filter.users) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rocketmq_users.default_group_perm_filter.users[*].default_group_perm : v == local.default_group_perm]
  )
}

`, testDmsRocketMQUser_basic(name))
}
