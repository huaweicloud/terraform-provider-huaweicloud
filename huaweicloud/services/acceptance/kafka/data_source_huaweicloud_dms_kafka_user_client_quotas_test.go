package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsKafkaUserClientQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafka_user_client_quotas.all"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsKafkaUserClientQuotas_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckOutput("user_validation", "true"),
					resource.TestCheckOutput("client_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsKafkaUserClientQuotas_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_kafka_user_client_quotas" "all" {
  depends_on = [huaweicloud_dms_kafka_user_client_quota.test]

  instance_id = huaweicloud_dms_kafka_instance.test.id
}

data "huaweicloud_dms_kafka_user_client_quotas" "test" {
  depends_on = [huaweicloud_dms_kafka_user_client_quota.test]

  instance_id = huaweicloud_dms_kafka_instance.test.id
  user        = huaweicloud_dms_kafka_user_client_quota.test.user
  client      = huaweicloud_dms_kafka_user_client_quota.test.client
}

locals {
  test_results = data.huaweicloud_dms_kafka_user_client_quotas.test
}

output "user_validation" {
  value = alltrue([for v in local.test_results.quotas[*].user : strcontains(v, huaweicloud_dms_kafka_user_client_quota.test.user)])
}

output "client_validation" {
  value = alltrue([for v in local.test_results.quotas[*].client : strcontains(v, huaweicloud_dms_kafka_user_client_quota.test.client)])
}
`, testDmsKafkaUserClientQuota_basic(name))
}
