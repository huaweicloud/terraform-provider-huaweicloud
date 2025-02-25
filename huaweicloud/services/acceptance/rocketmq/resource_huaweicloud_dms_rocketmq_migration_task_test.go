package rocketmq

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRocketmqMigrationTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getRocketmqMigrationTask: query RocketMQ migration task
	var (
		getRocketmqMigrationTaskHttpUrl = "v2/{project_id}/instances/{instance_id}/metadata"
		getRocketmqMigrationTaskProduct = "dms"
	)
	getRocketmqMigrationTaskClient, err := cfg.NewServiceClient(getRocketmqMigrationTaskProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	getRocketmqMigrationTaskPath := getRocketmqMigrationTaskClient.Endpoint + getRocketmqMigrationTaskHttpUrl
	getRocketmqMigrationTaskPath = strings.ReplaceAll(getRocketmqMigrationTaskPath, "{project_id}", getRocketmqMigrationTaskClient.ProjectID)
	getRocketmqMigrationTaskPath = strings.ReplaceAll(getRocketmqMigrationTaskPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getRocketmqMigrationTaskPath += fmt.Sprintf("?id=%s&type=vhost", state.Primary.ID)

	getRocketmqMigrationTaskOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getRocketmqMigrationTaskResp, err := getRocketmqMigrationTaskClient.Request("GET", getRocketmqMigrationTaskPath, &getRocketmqMigrationTaskOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving RocketMQ migration task: %s", err)
	}

	getRocketmqMigrationTaskRespBody, err := utils.FlattenResponse(getRocketmqMigrationTaskResp)
	if err != nil {
		return nil, err
	}
	jsonContent := utils.PathSearch("json_content", getRocketmqMigrationTaskRespBody, "")
	var jsonData interface{}
	err = json.Unmarshal([]byte(jsonContent.(string)), &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func TestAccRockemqMigrationTask_rocketmq(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rocketmq_migration_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRocketmqMigrationTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRockemqMigrationTask_rocketmq(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_dms_rocketmq_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "overwrite", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "rocketmq"),
					resource.TestCheckResourceAttr(resourceName, "topic_configs.0.order", "false"),
					resource.TestCheckResourceAttr(resourceName, "topic_configs.0.perm", "6"),
					resource.TestCheckResourceAttr(resourceName, "topic_configs.0.read_queue_num", "16"),
					resource.TestCheckResourceAttr(resourceName, "topic_configs.0.topic_filter_type", "SINGLE_TAG"),
					resource.TestCheckResourceAttr(resourceName, "topic_configs.0.topic_name", rName+"-topic"),
					resource.TestCheckResourceAttr(resourceName, "topic_configs.0.topic_sys_flag", "0"),
					resource.TestCheckResourceAttr(resourceName, "topic_configs.0.write_queue_num", "16"),
					resource.TestCheckResourceAttr(resourceName, "subscription_groups.0.consume_broadcast_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "subscription_groups.0.consume_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "subscription_groups.0.consume_from_min_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "subscription_groups.0.group_name", rName+"-group"),
					resource.TestCheckResourceAttr(resourceName, "subscription_groups.0.notify_consumerids_changed_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "subscription_groups.0.retry_max_times", "16"),
					resource.TestCheckResourceAttr(resourceName, "subscription_groups.0.retry_queue_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "subscription_groups.0.which_broker_when_consume_slow", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "start_date"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"overwrite",
				},
				ImportStateIdFunc: testRocketmqMigrationTaskImportState(resourceName),
			},
		},
	})
}

func TestAccRockemqMigrationTask_rabbitToRocket(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rocketmq_migration_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRocketmqMigrationTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRockemqMigrationTask_rabbitToRocket(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_dms_rocketmq_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "overwrite", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "rabbitToRocket"),
					resource.TestCheckResourceAttr(resourceName, "vhosts.0.name", rName+"-vhost"),
					resource.TestCheckResourceAttr(resourceName, "queues.0.name", rName+"-queue"),
					resource.TestCheckResourceAttr(resourceName, "queues.0.vhost", rName+"-vhost"),
					resource.TestCheckResourceAttr(resourceName, "queues.0.durable", "false"),
					resource.TestCheckResourceAttr(resourceName, "exchanges.0.name", rName+"-direct"),
					resource.TestCheckResourceAttr(resourceName, "exchanges.0.vhost", rName+"-vhost"),
					resource.TestCheckResourceAttr(resourceName, "exchanges.0.type", "topic"),
					resource.TestCheckResourceAttr(resourceName, "exchanges.0.durable", "false"),
					resource.TestCheckResourceAttr(resourceName, "bindings.0.source", rName+"-direct"),
					resource.TestCheckResourceAttr(resourceName, "bindings.0.vhost", rName+"-vhost"),
					resource.TestCheckResourceAttr(resourceName, "bindings.0.destination", rName+"-queue"),
					resource.TestCheckResourceAttr(resourceName, "bindings.0.destination_type", "queue"),
					resource.TestCheckResourceAttr(resourceName, "bindings.0.routing_key", rName+"-queue"),
					resource.TestCheckResourceAttrSet(resourceName, "start_date"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"overwrite",
				},
				ImportStateIdFunc: testRocketmqMigrationTaskImportState(resourceName),
			},
		},
	})
}

func testAccRockemqMigrationTask_rocketmq(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_rocketmq_migration_task" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  overwrite   = "true"
  name        = "%[2]s"
  type        = "rocketmq"

  topic_configs {
    order             = false
    perm              = 6
    read_queue_num    = 16
    topic_filter_type = "SINGLE_TAG"
    topic_name        = "%[2]s-topic"
    topic_sys_flag    = 0
    write_queue_num   = 16
  }

  subscription_groups  {     
    consume_broadcast_enable          = true
    consume_enable                    = true
    consume_from_min_enable           = true
    group_name                        = "%[2]s-group"
    notify_consumerids_changed_enable = true
    retry_max_times                   = 16
    retry_queue_num                   = 1
    which_broker_when_consume_slow    = 1        
  }
}
`, testDmsRocketMQInstance_basic(name), name)
}

func testAccRockemqMigrationTask_rabbitToRocket(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_rocketmq_migration_task" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  overwrite   = "true"
  name        = "%[2]s"
  type        = "rabbitToRocket"

  vhosts {
    name = "%[2]s-vhost"
  }

  queues {
    name  = "%[2]s-queue"
    vhost = "%[2]s-vhost"
    durable = false
  }
  
  exchanges {
    name    = "%[2]s-direct"
    vhost   = "%[2]s-vhost"
    type    = "topic"
    durable = false
  }

  bindings {
    source           = "%[2]s-direct"
    vhost            = "%[2]s-vhost"
    destination      = "%[2]s-queue"
    destination_type = "queue"
    routing_key      = "%[2]s-queue"
  }
}
`, testDmsRocketMQInstance_basic(name), name)
}

// testRocketmqMigrationTaskImportState is used to return an import id with format <instance_id>/<id>
func testRocketmqMigrationTaskImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("the resource (%s) not found: %s", name, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceId, rs.Primary.ID), nil
	}
}
