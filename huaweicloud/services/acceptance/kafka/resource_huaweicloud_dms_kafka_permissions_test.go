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

// Before running this test, please ensure that the Kafka instance has SSL authentication enabled.
func TestAccPermissions_basic(t *testing.T) {
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
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPermissions_basic_step1(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "policies.0.user_name",
						"huaweicloud_dms_kafka_user.test1", "name"),
					resource.TestCheckResourceAttr(resourceName, "policies.0.access_policy", "all"),
				),
			},
			{
				Config: testAccPermissions_basic_step2(rName, password),
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

func testAccPermissions_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  partitions  = 3
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, rName)
}

func testAccPermissions_basic_step1(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_user" "test1" {
  instance_id = "%[4]s"
  name        = "%[2]s-1"
  password    = "%[3]s"
}

resource "huaweicloud_dms_kafka_permissions" "test" {
  instance_id = "%[4]s"
  topic_name  = huaweicloud_dms_kafka_topic.test.name

  policies {
    user_name     = huaweicloud_dms_kafka_user.test1.name
    access_policy = "all"
  }
}
`, testAccPermissions_base(rName), rName, password, acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}

func testAccPermissions_basic_step2(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_user" "test1" {
  instance_id = "%[4]s"
  name        = "%[2]s-1"
  password    = "%[3]s"
}

resource "huaweicloud_dms_kafka_user" "test2" {
  instance_id = "%[4]s"
  name        = "%[2]s-2"
  password    = "%[3]s"
}

resource "huaweicloud_dms_kafka_permissions" "test" {
  instance_id = "%[4]s"
  topic_name  = huaweicloud_dms_kafka_topic.test.name

  policies {
    user_name     = huaweicloud_dms_kafka_user.test1.name
    access_policy = "pub"
  }

  policies {
    user_name     = huaweicloud_dms_kafka_user.test2.name
    access_policy = "sub"
  }
}
`, testAccPermissions_base(rName), rName, password, acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}
