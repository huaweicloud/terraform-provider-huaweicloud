package kafka

import (
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

func getDmsKafkaSmartConnectTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getKafkaSmartConnectTask: query DMS kafka smart connect task
	var (
		getKafkaSmartConnectTaskHttpUrl = "v2/{project_id}/connectors/{connector_id}/sink-tasks/{task_id}"
		getKafkaSmartConnectTaskProduct = "dms"
	)
	getKafkaSmartConnectTaskClient, err := cfg.NewServiceClient(getKafkaSmartConnectTaskProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS Client: %s", err)
	}

	connectorID := state.Primary.Attributes["connector_id"]
	getKafkaSmartConnectTaskPath := getKafkaSmartConnectTaskClient.Endpoint + getKafkaSmartConnectTaskHttpUrl
	getKafkaSmartConnectTaskPath = strings.ReplaceAll(getKafkaSmartConnectTaskPath, "{project_id}",
		getKafkaSmartConnectTaskClient.ProjectID)
	getKafkaSmartConnectTaskPath = strings.ReplaceAll(getKafkaSmartConnectTaskPath, "{connector_id}", connectorID)
	getKafkaSmartConnectTaskPath = strings.ReplaceAll(getKafkaSmartConnectTaskPath, "{task_id}", state.Primary.ID)

	getKafkaSmartConnectTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getKafkaSmartConnectTaskResp, err := getKafkaSmartConnectTaskClient.Request("GET", getKafkaSmartConnectTaskPath,
		&getKafkaSmartConnectTaskOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DMS kafka smart connect task: %s", err)
	}

	getKafkaSmartConnectTaskRespBody, err := utils.FlattenResponse(getKafkaSmartConnectTaskResp)
	if err != nil {
		return nil, err
	}

	return getKafkaSmartConnectTaskRespBody, nil
}

func TestAccDmsKafkaSmartConnectTask_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_smart_connect_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsKafkaSmartConnectTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsKafkaSmartConnectTask_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "source_type"),
					resource.TestCheckResourceAttrSet(resourceName, "task_name"),
					resource.TestCheckResourceAttrSet(resourceName, "destination_type"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "source_type", "BLOB"),
					resource.TestCheckResourceAttr(resourceName, "task_name", rName),
					resource.TestCheckResourceAttr(resourceName, "destination_type", "OBS"),
					resource.TestCheckResourceAttr(resourceName, "topics", rName),
					resource.TestCheckResourceAttr(resourceName, "consumer_strategy", "latest"),
					resource.TestCheckResourceAttr(resourceName, "destination_file_type", "TEXT"),
					resource.TestCheckResourceAttr(resourceName, "record_delimiter", ";"),
					resource.TestCheckResourceAttr(resourceName, "deliver_time_interval", "300"),
					resource.TestCheckResourceAttr(resourceName, "obs_bucket_name", rName),
					resource.TestCheckResourceAttr(resourceName, "partition_format", "yyyy/MM/dd/HH/mm"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"source_type", "access_key", "secret_key",
				},
				ImportStateIdFunc: testKafkaSmartConnectTaskResourceImportState(resourceName),
			},
		},
	})
}

func testDmsKafkaSmartConnectTask_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[2]s"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_dms_kafka_smart_connect" "test" {
  instance_id       = huaweicloud_dms_kafka_instance.test.id
  storage_spec_code = "dms.physical.storage.high.v2"
  node_count        = 2
  bandwidth         = "100MB"
}

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%[2]s"
  partitions  = 10
  aging_time  = 36
}

resource "huaweicloud_dms_kafka_smart_connect_task" "test" {
  connector_id          = huaweicloud_dms_kafka_smart_connect.test.id
  source_type           = "BLOB"
  task_name             = "%[2]s"
  destination_type      = "OBS"
  topics                = huaweicloud_dms_kafka_topic.test.name
  consumer_strategy     = "latest"
  destination_file_type = "TEXT"
  access_key            = "%[3]s"
  secret_key            = "%[4]s"
  obs_bucket_name       = huaweicloud_obs_bucket.test.bucket
  partition_format      = "yyyy/MM/dd/HH/mm"
  record_delimiter      = ";"
  deliver_time_interval = 300
}

`, testAccKafkaInstance_newFormat(rName), rName, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}

// testKafkaSmartConnectTaskResourceImportState is used to return an import id with format <connector_id>/<id>
func testKafkaSmartConnectTaskResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		connector_id := rs.Primary.Attributes["connector_id"]
		return fmt.Sprintf("%s/%s", connector_id, rs.Primary.ID), nil
	}
}
