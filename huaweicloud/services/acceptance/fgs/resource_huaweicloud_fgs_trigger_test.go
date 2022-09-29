package fgs

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/trigger"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getTriggerResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FgsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud FunctionGraph v2 client: %s", err)
	}
	return trigger.Get(c, state.Primary.Attributes["function_urn"], state.Primary.Attributes["type"],
		state.Primary.ID).Extract()
}

func TestAccFunctionGraphTrigger_basic(t *testing.T) {
	var (
		timeTrigger  trigger.Trigger
		randName     = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_fgs_trigger.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&timeTrigger,
		getTriggerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphTimingTrigger_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Rate"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "3d"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "function_urn",
						"${huaweicloud_fgs_function.test.urn}"),
				),
			},
			{
				Config: testAccFunctionGraphTimingTrigger_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Rate"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "3d"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "function_urn",
						"${huaweicloud_fgs_function.test.urn}"),
				),
			},
		},
	})
}

func TestAccFunctionGraphTrigger_cronTimer(t *testing.T) {
	var (
		randName     = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_fgs_trigger.test"
		timeTrigger  trigger.Trigger
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&timeTrigger,
		getTriggerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphTimingTrigger_cron(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Cron"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "@every 1h30m"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "function_urn",
						"${huaweicloud_fgs_function.test.urn}"),
				),
			},
			{
				Config: testAccFunctionGraphTimingTrigger_cronUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "TIMER"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.name", randName),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule_type", "Cron"),
					resource.TestCheckResourceAttr(resourceName, "timer.0.schedule", "@every 1h30m"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "function_urn",
						"${huaweicloud_fgs_function.test.urn}"),
				),
			},
		},
	})
}

func TestAccFunctionGraphTrigger_smn(t *testing.T) {
	var (
		randName     = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_fgs_trigger.test"
		timeTrigger  trigger.Trigger
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&timeTrigger,
		getTriggerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphSmnTrigger_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "SMN"),
					resource.TestCheckResourceAttrSet(resourceName, "smn.0.topic_urn"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "function_urn",
						"${huaweicloud_fgs_function.test.urn}"),
				),
			},
		},
	})
}

func TestAccFunctionGraphTrigger_lts(t *testing.T) {
	var (
		randName     = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_fgs_trigger.test"
		ltsTrigger   trigger.Trigger
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ltsTrigger,
		getTriggerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckFgsTrigger(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphLtsTrigger_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "LTS"),
					resource.TestCheckResourceAttrSet(resourceName, "lts.0.log_group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "lts.0.log_topic_id"),
				),
			},
		},
	})
}

func testAccFunctionGraphTimingTrigger_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 10
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}`, rName)
}

func testAccFunctionGraphTimingTrigger_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"

  timer {
    name          = "%s"
    schedule_type = "Rate"
    schedule      = "3d"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphTimingTrigger_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "DISABLED"

  timer {
	name          = "%s"
	schedule_type = "Rate"
	schedule      = "3d"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphTimingTrigger_cron(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"

  timer {
    name          = "%s"
    schedule_type = "Cron"
    schedule      = "@every 1h30m"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphTimingTrigger_cronUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "DISABLED"

  timer {
	name          = "%s"
	schedule_type = "Cron"
	schedule      = "@every 1h30m"
  }
}
`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphSmnTrigger_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_smn_topic" "test" {
  name = "%s"
}

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "SMN"

  smn {
    topic_urn = huaweicloud_smn_topic.test.topic_urn
  }
}`, testAccFunctionGraphTimingTrigger_base(rName), rName)
}

func testAccFunctionGraphLtsTrigger_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}

resource "huaweicloud_identity_agency" "test" {
  name = "%[1]s"
  delegated_service_name = "%[3]s"

  project_role {
    project = "%[2]s"
    roles = ["LTS FullAccess"]
  }
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 10
  runtime     = "Python2.7"
  code_type   = "inline"
  agency      = huaweicloud_identity_agency.test.name
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "LTS"

  lts {
    log_group_id = huaweicloud_lts_group.test.id
    log_topic_id = huaweicloud_lts_stream.test.id
  }
}`, rName, acceptance.HW_REGION_NAME, acceptance.HW_FGS_TRIGGER_LTS_AGENCY)
}

func testAccNetwork_config(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.128.0/20"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.128.0/24"
  gateway_ip = "192.168.128.1"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9092
  port_range_max    = 9092
  remote_ip_prefix  = "0.0.0.0/0"
}`, rName, rName, rName)
}

func testAccDmsKafka_config(rName, password string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_az" "test" {}

data "huaweicloud_dms_product" "test" {
  engine            = "kafka"
  version           = "1.1.0"
  instance_type     = "cluster"
  partition_num     = 300
  storage           = 600
  storage_spec_code = "dms.physical.storage.high"
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  available_zones   = [data.huaweicloud_dms_az.test.id]
  product_id        = data.huaweicloud_dms_product.test.id
  engine_version    = data.huaweicloud_dms_product.test.version
  bandwidth         = data.huaweicloud_dms_product.test.bandwidth
  storage_space     = data.huaweicloud_dms_product.test.storage
  storage_spec_code = data.huaweicloud_dms_product.test.storage_spec_code
  manager_user      = "%s"
  manager_password  = "%s"
}

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%s"
  partitions  = 20
}`, testAccNetwork_config(rName), rName, rName, password, rName)
}
