package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/kafka"
)

func getConsumerGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return kafka.GetConsumerGroupByName(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["name"])
}

func TestAccConsumerGroup_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		obj   interface{}
		rName = "huaweicloud_dms_kafka_consumer_group.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getConsumerGroupResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConsumerGroup_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "state", "EMPTY"),
					resource.TestCheckResourceAttrSet(rName, "coordinator_id"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccConsumerGroup_basic_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", rName),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConsumerGroup_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_consumer_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Created by terraform script"
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func testAccConsumerGroup_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_consumer_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}
