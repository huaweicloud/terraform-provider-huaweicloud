package css

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getCssClusterFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	return cluster.Get(client, state.Primary.ID)
}

func TestAccCssCluster_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_basic(rName, 1, 7, "bar"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccCssCluster_basic(rName, 2, 8, "bar_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccCssCluster_prePaid(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_prePaid(rName, 1, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "security_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.instance_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "client_node_config.0.instance_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.instance_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpcep_endpoint.0.endpoint_with_dns_name", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcep_endpoint_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcep_ip"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccCssCluster_prePaid(rName, 1, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func testAccCssBase(rName string) string {
	bucketName := acceptance.RandomAccResourceNameWithDash()
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  vpc_id     = huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s"
  description = "terraform security group acceptance test"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_obs_bucket" "cssObs" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true
}
`, rName, rName, rName, bucketName)
}

func testAccCssCluster_basic(rName string, nodeNum int, keepDays int, tag string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_css_cluster" "test" {
  name            = "%s"
  engine_version  = "7.10.2"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = %d
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  backup_strategy {
    keep_days  = %d
    start_time = "00:00 GMT+08:00"
    prefix     = "snapshot"
    bucket      = huaweicloud_obs_bucket.cssObs.bucket
    agency      = "css_obs_agency"
    backup_path = "css_repository/acctest"
  }

  tags = {
    foo = "%s"
    key = "value"
  }
}
`, testAccCssBase(rName), rName, nodeNum, keepDays, tag)
}

func testAccCssCluster_prePaid(rName string, nodeNum int, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"
  security_mode  = true
  password       = "Test@passw0rd"

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = %[3]d
  auto_renew    = "%[4]v"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  master_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 3
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  client_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  cold_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  vpcep_endpoint {
    endpoint_with_dns_name = true
  }
}
`, testAccCssBase(rName), rName, nodeNum, isAutoRenew)
}
