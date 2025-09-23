package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKafkaInstanceRestart_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstanceRestart_basic(rName),
			},
		},
	})
}

func testAccKafkaInstanceRestart_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_instance_restart" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
}`, testAccKafkaInstance_newFormat(rName))
}
