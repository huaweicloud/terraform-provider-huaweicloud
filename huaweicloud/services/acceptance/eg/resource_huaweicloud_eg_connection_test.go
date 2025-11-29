package eg

import (
	"fmt"
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

func getConnectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getConnection: Query the EG Connection detail
	var (
		getConnectionHttpUrl = "v1/{project_id}/connections/{id}"
		getConnectionProduct = "eg"
	)
	getConnectionClient, err := cfg.NewServiceClient(getConnectionProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EG client: %s", err)
	}

	getConnectionPath := getConnectionClient.Endpoint + getConnectionHttpUrl
	getConnectionPath = strings.ReplaceAll(getConnectionPath, "{project_id}", getConnectionClient.ProjectID)
	getConnectionPath = strings.ReplaceAll(getConnectionPath, "{id}", state.Primary.ID)

	getConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getConnectionResp, err := getConnectionClient.Request("GET", getConnectionPath, &getConnectionOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Connection: %s", err)
	}

	getConnectionRespBody, err := utils.FlattenResponse(getConnectionResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Connection: %s", err)
	}

	return getConnectionRespBody, nil
}

func TestAccConnection_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_eg_connection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testConnection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "agency"),
				),
			},
			{
				Config: testConnection_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testConnection_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_eg_connection" "test" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  description = "created by terraform"
}
`, common.TestVpc(name), name)
}

func testConnection_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_eg_connection" "test" {
  name      = "%s"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
}
`, common.TestVpc(name), name)
}

func TestAccConnection_kafka(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_eg_connection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getConnectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testConnection_kafka(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(rName, "type", "KAFKA"),
					resource.TestCheckResourceAttrPair(rName, "kafka_detail.0.instance_id",
						"huaweicloud_dms_kafka_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "kafka_detail.0.connect_address",
						"huaweicloud_dms_kafka_instance.test", "connect_address"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "flavor"),
					resource.TestCheckResourceAttrSet(rName, "flavor.0.bandwidth_type"),
					resource.TestCheckResourceAttrSet(rName, "flavor.0.concurrency"),
					resource.TestCheckResourceAttrSet(rName, "flavor.0.concurrency_type"),
					resource.TestCheckResourceAttrSet(rName, "flavor.0.name"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testConnection_kafka_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type = "cluster"
}

locals {
  query_results = data.huaweicloud_dms_kafka_flavors.test

  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%s"
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
  engine_version = element(local.query_results.versions, length(local.query_results.versions)-1)
  storage_space  = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num     = 3

  access_user        = "user"
  password           = "Kafkatest@123"
  manager_user       = "kafka-user"
  manager_password   = "Kafkatest@123"
  security_protocol  = "SASL_PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]

  cross_vpc_accesses {
    advertised_ip = ""
  }
  cross_vpc_accesses {
    advertised_ip = "www.terraform-test.com"
  }
  cross_vpc_accesses {
    advertised_ip = "192.168.0.53"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testConnection_kafka(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_eg_connection" "test" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  description = "created by terraform"
  type        = "KAFKA"

  kafka_detail {
    instance_id     = huaweicloud_dms_kafka_instance.test.id
    connect_address = huaweicloud_dms_kafka_instance.test.connect_address
  }
}
`, testConnection_kafka_base(name), name)
}
