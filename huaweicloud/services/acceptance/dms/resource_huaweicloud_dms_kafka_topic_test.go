package dms

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/topics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDmsKafkaTopicFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud DMS client(V2): %s", err)
	}
	instanceID := state.Primary.Attributes["instance_id"]
	allTopics, err := topics.List(client, instanceID).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error listing DMS kafka topics in %s, error: %s", instanceID, err)
	}

	topicID := state.Primary.ID
	for _, item := range allTopics {
		if item.Name == topicID {
			return item, nil
		}
	}

	return nil, fmt.Errorf("can not found dms kafka topic instance")
}

func TestAccDmsKafkaTopic_basic(t *testing.T) {
	var kafkaTopic topics.Topic
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_topic.topic"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&kafkaTopic,
		getDmsKafkaTopicFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaTopic_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "partitions", "10"),
					resource.TestCheckResourceAttr(resourceName, "replicas", "3"),
					resource.TestCheckResourceAttr(resourceName, "aging_time", "36"),
					resource.TestCheckResourceAttr(resourceName, "sync_replication", "false"),
					resource.TestCheckResourceAttr(resourceName, "sync_flushing", "false"),
				),
			},
			{
				Config: testAccDmsKafkaTopic_update_part1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "partitions", "20"),
					resource.TestCheckResourceAttr(resourceName, "replicas", "3"),
					resource.TestCheckResourceAttr(resourceName, "aging_time", "72"),
				),
			},
			{
				Config: testAccDmsKafkaTopic_update_part2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "partitions", "20"),
					resource.TestCheckResourceAttr(resourceName, "replicas", "3"),
					resource.TestCheckResourceAttr(resourceName, "aging_time", "72"),
					resource.TestCheckResourceAttr(resourceName, "sync_replication", "true"),
					resource.TestCheckResourceAttr(resourceName, "sync_flushing", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccKafkaTopicImportStateFunc(resourceName),
			},
		},
	})
}

// testAccKafkaTopicImportStateFunc is used to import the resource
func testAccKafkaTopicImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		instance, ok := s.RootModule().Resources["huaweicloud_dms_kafka_instance.test"]
		if !ok {
			return "", fmt.Errorf("DMS kafka instance not found")
		}
		topic, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("DMS kafka topic not found")
		}

		return fmt.Sprintf("%s/%s", instance.Primary.ID, topic.Primary.ID), nil
	}
}

func testAccDmsKafkaTopic_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_topic" "topic" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%s"
  partitions  = 10
  aging_time  = 36
}
`, testAccKafkaInstance_basic(rName), rName)
}

func testAccDmsKafkaTopic_update_part1(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_topic" "topic" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%s"
  partitions  = 20
  aging_time  = 72
}
`, testAccKafkaInstance_basic(rName), rName)
}

func testAccDmsKafkaTopic_update_part2(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_topic" "topic" {
  instance_id      = huaweicloud_dms_kafka_instance.test.id
  name             = "%s"
  partitions       = 20
  aging_time       = 72
  sync_flushing    = true
  sync_replication = true
}
`, testAccKafkaInstance_basic(rName), rName)
}
