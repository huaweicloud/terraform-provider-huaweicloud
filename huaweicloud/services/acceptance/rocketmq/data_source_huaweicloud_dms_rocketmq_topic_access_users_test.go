package rocketmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRocketmqTopicAccessUsers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rocketmq_topic_access_users.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDmsRocketmqTopicAccessUsers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.admin"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.perm"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.access_key"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.white_remote_address"),
				),
			},
		},
	})
}

func testDataSourceDmsRocketmqTopicAccessUsers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  name        = "%[2]s"
  queue_num   = 3

  brokers {
    name = "broker-0"
  }
}

data "huaweicloud_dms_rocketmq_topic_access_users" "test" {
  depends_on = [huaweicloud_dms_rocketmq_user.test, huaweicloud_dms_rocketmq_topic.test]

  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  topic       = huaweicloud_dms_rocketmq_topic.test.name
}
`, testDmsRocketMQUser_basic(name), name)
}
