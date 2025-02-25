package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsKafkaUsers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafka_users.all"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsKafkaUsers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "users.#"),
					resource.TestCheckOutput("name_validation", "true"),
					resource.TestCheckOutput("description_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsKafkaUsers_basic(name string) string {
	password := acceptance.RandomPassword()
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_kafka_users" "all" {
  depends_on = [huaweicloud_dms_kafka_user.test]

  instance_id = huaweicloud_dms_kafka_instance.test.id
}

data "huaweicloud_dms_kafka_users" "test" {
  depends_on = [huaweicloud_dms_kafka_user.test]

  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = huaweicloud_dms_kafka_user.test.name
  description = huaweicloud_dms_kafka_user.test.description
}

locals {
  test_results = data.huaweicloud_dms_kafka_users.test
}

output "name_validation" {
  value = alltrue([for v in local.test_results.users[*].name : strcontains(v, huaweicloud_dms_kafka_user.test.name)])
}

output "description_validation" {
  value = alltrue([for v in local.test_results.users[*].description : strcontains(v, huaweicloud_dms_kafka_user.test.description)])
}
`, testAccDmsKafkaUser_basic(name, password, "test"))
}
