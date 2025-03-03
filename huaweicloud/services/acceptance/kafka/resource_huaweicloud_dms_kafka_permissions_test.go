package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/kafka"
)

func getDmsKafkaPermissionsFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.DmsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return kafka.GetDmsKafkaPermissions(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["topic_name"])
}

func TestAccDmsKafkaPermissions_basic(t *testing.T) {
	var policies interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_permissions.test"
	password := acceptance.RandomPassword()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&policies,
		getDmsKafkaPermissionsFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaPermissions_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "policies.0.user_name",
						"huaweicloud_dms_kafka_user.test1", "name"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.access_policy", "all"),
				),
			},
			{
				Config: testAccDmsKafkaPermissions_update(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "policies.0.user_name",
						"huaweicloud_dms_kafka_user.test1", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "policies.1.user_name",
						"huaweicloud_dms_kafka_user.test2", "name"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.access_policy", "pub"),
					resource.TestCheckResourceAttr(resourceName, "policies.1.access_policy", "sub"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDmsKafkaPermissions_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_user" "test1" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%[2]s-1"
  password    = "%[3]s"
}

resource "huaweicloud_dms_kafka_permissions" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  topic_name  = huaweicloud_dms_kafka_topic.topic.name

  policies {
    user_name     = huaweicloud_dms_kafka_user.test1.name
    access_policy = "all"
  }
}
`, testAccDmsKafkaTopic_basic(rName), rName, password)
}

func testAccDmsKafkaPermissions_update(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_user" "test1" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%[2]s-1"
  password    = "%[3]s"
}

resource "huaweicloud_dms_kafka_user" "test2" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%[2]s-2"
  password    = "%[3]s"
}

resource "huaweicloud_dms_kafka_permissions" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  topic_name  = huaweicloud_dms_kafka_topic.topic.name

  policies {
    user_name     = huaweicloud_dms_kafka_user.test1.name
    access_policy = "pub"
  }

  policies {
    user_name     = huaweicloud_dms_kafka_user.test2.name
    access_policy = "sub"
  }
}
`, testAccDmsKafkaTopic_basic(rName), rName, password)
}
