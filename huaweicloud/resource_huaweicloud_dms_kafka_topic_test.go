package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/dms/v2/kafka/topics"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDmsKafkaTopic_basic(t *testing.T) {
	var kafkaTopic topics.Topic
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_dms_kafka_topic.topic"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsKafkaTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsKafkaTopic_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsKafkaTopicExists(resourceName, &kafkaTopic),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "partitions", "10"),
					resource.TestCheckResourceAttr(resourceName, "replicas", "3"),
					resource.TestCheckResourceAttr(resourceName, "aging_time", "36"),
					resource.TestCheckResourceAttr(resourceName, "sync_replication", "false"),
					resource.TestCheckResourceAttr(resourceName, "sync_flushing", "false"),
				),
			},
			{
				Config: testAccDmsKafkaTopic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsKafkaTopicExists(resourceName, &kafkaTopic),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "partitions", "20"),
					resource.TestCheckResourceAttr(resourceName, "replicas", "3"),
					resource.TestCheckResourceAttr(resourceName, "aging_time", "72"),
					resource.TestCheckResourceAttr(resourceName, "sync_replication", "false"),
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

func testAccCheckDmsKafkaTopicDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dmsClient, err := config.DmsV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud DMS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dms_kafka_topic" {
			continue
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		allTopics, err := topics.List(dmsClient, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("Error listing DMS kafka topics in %s, err: %s", instanceID, err)
		}

		topicID := rs.Primary.ID
		for _, item := range allTopics {
			if item.Name == topicID {
				return fmt.Errorf("The DMS kafka topic %s still exists", topicID)
			}
		}
	}
	return nil
}

func testAccCheckDmsKafkaTopicExists(n string, topic *topics.Topic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		dmsClient, err := config.DmsV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud DMS client: %s", err)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		allTopics, err := topics.List(dmsClient, instanceID).Extract()
		if err != nil {
			return fmt.Errorf("Error listing DMS kafka topics in %s, err: %s", instanceID, err)
		}

		topicID := rs.Primary.ID
		for _, item := range allTopics {
			if item.Name == topicID {
				*topic = item
				return nil
			}
		}

		return fmt.Errorf("The DMS kafka topic %s not found", topicID)
	}
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
`, testAccDmsKafkaInstance_basic(rName), rName)
}

func testAccDmsKafkaTopic_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_kafka_topic" "topic" {
  instance_id   = huaweicloud_dms_kafka_instance.test.id
  name          = "%s"
  partitions    = 20
  aging_time    = 72
  sync_flushing = true
}
`, testAccDmsKafkaInstance_basic(rName), rName)
}
