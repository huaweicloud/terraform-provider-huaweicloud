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

func getResourceTopicFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.NewServiceClient("dms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return kafka.GetTopicByName(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccTopic_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		obj   interface{}
		rName = "huaweicloud_dms_kafka_topic.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getResourceTopicFunc)
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
				Config: testAccTopic_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_DMS_KAFKA_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "partitions", "4"),
					resource.TestCheckResourceAttr(rName, "replicas", "3"),
					resource.TestCheckResourceAttr(rName, "aging_time", "36"),
					resource.TestCheckResourceAttr(rName, "sync_replication", "true"),
					resource.TestCheckResourceAttr(rName, "sync_flushing", "true"),
					resource.TestCheckResourceAttr(rName, "description", "Created by Terraform script"),
					resource.TestCheckResourceAttr(rName, "configs.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "configs.*", map[string]string{
						"name":  "max.message.bytes",
						"value": "10000000",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "configs.*", map[string]string{
						"name":  "message.timestamp.type",
						"value": "LogAppendTime",
					}),
					resource.TestCheckResourceAttrSet(rName, "policies_only"),
					resource.TestCheckResourceAttrSet(rName, "type"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccTopic_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "partitions", "5"),
					resource.TestCheckResourceAttr(rName, "replicas", "3"),
					resource.TestCheckResourceAttr(rName, "description", "Updated by Terraform script"),
					resource.TestCheckResourceAttr(rName, "aging_time", "72"),
					resource.TestCheckResourceAttr(rName, "sync_replication", "false"),
					resource.TestCheckResourceAttr(rName, "sync_flushing", "false"),
					resource.TestCheckResourceAttr(rName, "configs.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "configs.*", map[string]string{
						"name":  "max.message.bytes",
						"value": "10000001",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "configs.*", map[string]string{
						"name":  "message.timestamp.type",
						"value": "CreateTime",
					}),
				),
			},
			{
				Config: testAccTopic_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "partitions", "5"),
					resource.TestCheckResourceAttr(rName, "replicas", "3"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "aging_time", "72"),
					resource.TestCheckResourceAttr(rName, "sync_replication", "true"),
					resource.TestCheckResourceAttr(rName, "sync_flushing", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccTopicImportStateFunc(rName),
			},
		},
	})
}

func testAccTopicImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		topicName := rs.Primary.ID
		if instanceId == "" || topicName == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<name>', but got '%s/%s'", instanceId, topicName)
		}

		return fmt.Sprintf("%s/%s", instanceId, topicName), nil
	}
}

func testAccTopic_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id      = "%[1]s"
  name             = "%[2]s"
  partitions       = 4
  aging_time       = 36
  description      = "Created by Terraform script"
  sync_flushing    = true
  sync_replication = true

  configs {
    name  = "max.message.bytes"
    value = "10000000"
  }
  configs {
    name  = "message.timestamp.type"
    value = "LogAppendTime"
  }
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func testAccTopic_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id      = "%[1]s"
  name             = "%[2]s"
  partitions       = 5
  aging_time       = 72
  description      = "Updated by Terraform script"
  sync_flushing    = false
  sync_replication = false

  configs {
    name  = "max.message.bytes"
    value = "10000001"
  }
  configs {
    name  = "message.timestamp.type"
    value = "CreateTime"
  }
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func testAccTopic_basic_step3(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id      = "%[1]s"
  name             = "%[2]s"
  partitions       = 5
  aging_time       = 72
  sync_flushing    = true
  sync_replication = true

  configs {
    name  = "max.message.bytes"
    value = "10000001"
  }
  configs {
    name  = "message.timestamp.type"
    value = "CreateTime"
  }
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func TestAccTopic_new_partition_brokers(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		obj   interface{}
		rName = "huaweicloud_dms_kafka_topic.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getResourceTopicFunc)
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
				Config: testAccTopic_new_partition_brokers_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "configs.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttr(rName, "partitions", "4"),
				),
			},
			{
				Config: testAccTopic_new_partition_brokers_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "partitions", "5"),
				),
			},
			{
				Config:      testAccTopic_new_partition_brokers_step3(name),
				ExpectError: regexp.MustCompile(`only support to add partitions`),
			},
		},
	})
}

func testAccTopic_new_partition_brokers_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  partitions  = 4
  replicas    = 1
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func testAccTopic_new_partition_brokers_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id           = "%[1]s"
  name                  = "%[2]s"
  partitions            = 5
  replicas              = 1
  new_partition_brokers = [0]
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func testAccTopic_new_partition_brokers_step3(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id           = "%[1]s"
  name                  = "%[2]s"
  partitions            = 3
  replicas              = 1
  new_partition_brokers = [0]
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}
