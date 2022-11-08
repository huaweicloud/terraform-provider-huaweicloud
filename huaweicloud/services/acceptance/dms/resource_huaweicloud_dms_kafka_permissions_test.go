package dms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kafka/v2/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDmsKafkaPermissionsFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.HcDmsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	// Split instance_id and topic_name from resource id
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<topic_name>")
	}
	instanceId := parts[0]
	topicName := parts[1]

	request := &model.ShowTopicAccessPolicyRequest{
		InstanceId: instanceId,
		TopicName:  topicName,
	}

	response, err := client.ShowTopicAccessPolicy(request)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DMS kafka permissions: %s", err)
	}

	if response.Policies != nil && len(*response.Policies) != 0 {
		policies := *response.Policies
		return policies, nil
	}

	return nil, fmt.Errorf("can not found DMS kafka user")
}

func TestAccDmsKafkaPermissions_basic(t *testing.T) {
	var policies []model.PolicyEntity
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
