package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getLogstashClusterFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	return cluster.Get(client, state.Primary.ID)
}

func TestAccLogstashCluster_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_logstash_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getLogstashClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLogstashCluster_basic(rName, 1, "bar"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "node_config.0.instance_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "logstash"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccLogstashCluster_basic(rName+"-update", 2, "bar_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "node_config.0.instance_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
				),
			},
			{
				Config: testAccLogstashCluster_basic(rName+"-update", 1, "bar_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "node_config.0.instance_number", "1"),
				),
			},
		},
	})
}

func TestAccLogstashCluster_prePaid(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_logstash_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getLogstashClusterFunc,
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
				Config: testAccLogstashCluster_prePaid(rName, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "logstash"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccLogstashCluster_prePaid(rName, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "logstash"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func TestAccLogstashCluster_updateWithEpsId(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_logstash_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getLogstashClusterFunc,
	)
	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLogstashCluster_withEpsId(rName, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccLogstashCluster_withEpsId(rName, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func TestAccLogstashCluster_route(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_logstash_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getLogstashClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLogstashCluster_route_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "routes.0.ip_address", "192.168.0.0"),
					resource.TestCheckResourceAttr(resourceName, "routes.0.ip_net_mask", "255.255.255.0"),
				),
			},
			{
				Config: testAccLogstashCluster_route_add(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "routes.#", "2"),
				),
			},
			{
				Config: testAccLogstashCluster_route_del(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "routes.0.ip_address", "192.168.10.0"),
					resource.TestCheckResourceAttr(resourceName, "routes.0.ip_net_mask", "255.255.255.0"),
				),
			},
		},
	})
}

func testAcclogstashBase(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}
`, common.TestBaseNetwork(rName))
}

func testAccLogstashCluster_basic(rName string, nodeNum int, tag string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.0"

  node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = %[3]d
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  tags = {
    foo = "%[4]s"
  }
}
`, testAcclogstashBase(rName), rName, nodeNum, tag)
}

func testAccLogstashCluster_prePaid(rName string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.0"

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[3]v"

  node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }
}
`, testAcclogstashBase(rName), rName, isAutoRenew)
}

func testAccLogstashCluster_withEpsId(rName string, epsId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.0"

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1
  enterprise_project_id = "%[3]s"

  node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }
}
`, testAcclogstashBase(rName), rName, epsId)
}

func testAccLogstashCluster_route_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.0"

  node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  routes {
    ip_address  = "192.168.0.0"
    ip_net_mask = "255.255.255.0"
  }
}
`, testAcclogstashBase(rName), rName)
}

func testAccLogstashCluster_route_add(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.0"

  node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  routes {
    ip_address  = "192.168.0.0"
    ip_net_mask = "255.255.255.0"
  }

  routes {
    ip_address  = "192.168.10.0"
    ip_net_mask = "255.255.255.0"
  }
}
`, testAcclogstashBase(rName), rName)
}

func testAccLogstashCluster_route_del(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_logstash_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.0"

  node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  routes {
    ip_address  = "192.168.10.0"
    ip_net_mask = "255.255.255.0"
  }
}
`, testAcclogstashBase(rName), rName)
}
