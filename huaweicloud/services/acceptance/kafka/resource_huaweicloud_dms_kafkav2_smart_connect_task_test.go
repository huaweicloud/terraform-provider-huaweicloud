package kafka

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDmsKafkav2SmartConnectTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	getTaskHttpUrl := "v2/{project_id}/instances/{instance_id}/connector/tasks/{task_id}"
	getTaskPath := client.Endpoint + getTaskHttpUrl
	getTaskPath = strings.ReplaceAll(getTaskPath, "{project_id}", client.ProjectID)
	getTaskPath = strings.ReplaceAll(getTaskPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getTaskPath = strings.ReplaceAll(getTaskPath, "{task_id}", state.Primary.ID)
	getTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTaskResp, err := client.Request("GET", getTaskPath, &getTaskOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DMS kafka smart connect task: %s", err)
	}

	return utils.FlattenResponse(getTaskResp)
}

func TestAccDmsKafkav2SmartConnectTask_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafkav2_smart_connect_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsKafkav2SmartConnectTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsKafkav2SmartConnectTask_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "task_name", rName),
					resource.TestCheckResourceAttr(resourceName, "destination_type", "OBS_SINK"),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.consumer_strategy", "latest"),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.destination_file_type", "TEXT"),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.record_delimiter", ";"),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.deliver_time_interval", "300"),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.obs_bucket_name", rName),
					resource.TestCheckResourceAttr(resourceName, "destination_task.0.partition_format", "yyyy/MM/dd/HH/mm"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"destination_task.0.access_key", "destination_task.0.secret_key",
				},
				ImportStateIdFunc: testKafkav2SmartConnectTaskResourceImportState(resourceName),
			},
		},
	})
}

func TestAccDmsKafkav2SmartConnectTask_KafkaToKafka(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafkav2_smart_connect_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsKafkav2SmartConnectTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsKafkav2SmartConnectTask_kafkaToKafka(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "task_name"),
					resource.TestCheckResourceAttr(resourceName, "source_type", "KAFKA_REPLICATOR_SOURCE"),
					resource.TestCheckResourceAttr(resourceName, "task_name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "source_task.0.peer_instance_id",
						"huaweicloud_dms_kafka_instance.test2", "id"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.direction", "push"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.replication_factor", "3"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.task_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.rename_topic_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.provenance_header_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.sync_consumer_offsets_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.consumer_strategy", "latest"),
					resource.TestCheckResourceAttr(resourceName, "source_task.0.compression_type", "snappy"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testKafkav2SmartConnectTaskResourceImportState(resourceName),
			},
		},
	})
}

func testDmsKafkav2SmartConnectTask_basic(rName string) string {
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

resource "huaweicloud_dms_kafkav2_smart_connect_task" "test" {
  depends_on = [huaweicloud_dms_kafka_smart_connect.test, huaweicloud_dms_kafka_topic.test]

  instance_id      = huaweicloud_dms_kafka_instance.test.id
  task_name        = "%[2]s"
  destination_type = "OBS_SINK"
  topics           = [huaweicloud_dms_kafka_topic.test.name]

  destination_task {
    consumer_strategy     = "latest"
    destination_file_type = "TEXT"
    access_key            = "%[3]s"
    secret_key            = "%[4]s"
    obs_bucket_name       = huaweicloud_obs_bucket.test.bucket
    partition_format      = "yyyy/MM/dd/HH/mm"
    record_delimiter      = ";"
    deliver_time_interval = 300
  }
}`, testAccKafkaInstance_newFormat(rName), rName, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}

func testDmsKafkav2SmartConnectTask_kafkaToKafka(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafkav2_smart_connect_task" "test" {
  depends_on = [
	huaweicloud_dms_kafka_smart_connect.test1,
	huaweicloud_dms_kafka_smart_connect.test2,
	huaweicloud_dms_kafka_topic.test
  ]

  instance_id = huaweicloud_dms_kafka_instance.test1.id
  task_name   = "%[2]s"
  topics      = [huaweicloud_dms_kafka_topic.test.name]
  source_type = "KAFKA_REPLICATOR_SOURCE"

  source_task {
    peer_instance_id              = huaweicloud_dms_kafka_instance.test2.id
    direction                     = "push"
    replication_factor            = 3
    task_num                      = 2
    provenance_header_enabled     = true
    sync_consumer_offsets_enabled = true
    rename_topic_enabled          = true
    consumer_strategy             = "latest"
    compression_type              = "snappy"
  }
}`, testAccKafkav2SmartConnectTaskKafkaToKafKaBase(rName), rName)
}

func testAccKafkav2SmartConnectTaskKafkaToKafKaBase(rName string) string {
	kafka1 := testAccKafkav2SmartConnectTaskKafkaToKafKaInstanceBase(rName, 1)
	kafka2 := testAccKafkav2SmartConnectTaskKafkaToKafKaInstanceBase(rName, 2)
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "rule" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  ports             = "9092"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "out_v4_all" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "egress"
  remote_ip_prefix  = "0.0.0.0/0"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

%s

%s

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test1.id
  name        = "%s"
  partitions  = 10
  aging_time  = 36
}`, common.TestBaseNetwork(rName), kafka1, kafka2, rName)
}

func testAccKafkav2SmartConnectTaskKafkaToKafKaInstanceBase(rName string, num int) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_instance" "test%[2]v" {
  name              = "%[1]s-%[2]v"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor_id          = local.flavor.id
  storage_spec_code  = local.flavor.ios[0].storage_spec_code
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2]
  ]

  engine_version     = "2.7"
  storage_space      = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num         = 3
  arch_type          = "X86"
}

resource "huaweicloud_dms_kafka_smart_connect" "test%[2]v" {
  instance_id       = huaweicloud_dms_kafka_instance.test%[2]v.id
  storage_spec_code = "dms.physical.storage.high.v2"
  node_count        = 2
  bandwidth         = "100MB"
}`, rName, num)
}

// testKafkaSmartConnectTaskResourceImportState is used to return an import ID with format <instance_id>/<task_id>
func testKafkav2SmartConnectTaskResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceID, rs.Primary.ID), nil
	}
}
