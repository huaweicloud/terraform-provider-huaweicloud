package as

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getASGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	asClient, err := cfg.AutoscalingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating autoscaling client: %s", err)
	}

	return groups.Get(asClient, state.Primary.ID).Extract()
}

func TestAccASGroup_basic(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_group.acc_as_group"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testASGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "this is a basic AS group"),
					resource.TestCheckResourceAttr(resourceName, "desire_instance_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "min_instance_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_instance_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "delete_publicip", "false"),
					resource.TestCheckResourceAttr(resourceName, "delete_volume", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "lbaas_listeners.0.pool_id",
						"huaweicloud_lb_pool.pool_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lbaas_listeners.0.protocol_port",
						"huaweicloud_lb_listener.listener_1", "protocol_port"),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.weight", "20"),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.protocol_version", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "networks.0.source_dest_check", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "multi_az_scaling_policy", "EQUILIBRIUM_DISTRIBUTE"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "300"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_time", "5"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_grace_period", "600"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zones.#"),
				),
			},
			{
				Config: testASGroup_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scaling_group_name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "this is an updated AS group"),
					resource.TestCheckResourceAttr(resourceName, "min_instance_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_instance_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "delete_publicip", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_volume", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "lbaas_listeners.0.pool_id",
						"huaweicloud_lb_pool.pool_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lbaas_listeners.0.protocol_port",
						"huaweicloud_lb_listener.listener_1", "protocol_port"),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.weight", "30"),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.protocol_version", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "agency_name", "ims_admin"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_instances",
				},
			},
			{
				Config: testASGroup_basic_disable(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "PAUSED"),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.weight", "1"),
				),
			},
			{
				Config: testASGroup_basic_enable(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "multi_az_scaling_policy", "PICK_FIRST"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "600"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_time", "15"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_grace_period", "900"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
				),
			},
		},
	})
}

func TestAccASGroup_withEpsId(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_group.acc_as_group"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testASGroup_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.weight", "1"),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.protocol_version", "ipv4"),
				),
			},
		},
	})
}

func TestAccASGroup_forceDelete(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_group.acc_as_group"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testASGroup_forceDelete(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "desire_instance_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "min_instance_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_instance_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "instances.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
					resource.TestCheckResourceAttr(resourceName, "delete_publicip", "true"),
					resource.TestCheckResourceAttr(resourceName, "delete_volume", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "scaling_configuration_name"),
					resource.TestCheckResourceAttrSet(resourceName, "is_scaling"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
		},
	})
}

func TestAccASGroup_sourceDestCheck(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_group.acc_as_group"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testASGroup_sourceDestCheck(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "networks.0.source_dest_check", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
				),
			},
			{
				Config: testASGroup_sourceDestErrorCheck(rName),
				ExpectError: regexp.MustCompile("invalid parameters: it should be min_instance_number <=" +
					" desire_instance_number <= max_instance_number"),
			},
		},
	})
}

//nolint:revive
func testASGroup_Base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "acc_key" {
  name       = "%[2]s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[2]s"
  vip_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "%[2]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name        = "%[2]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.listener_1.id
}

resource "huaweicloud_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%[2]s"
  instance_config {
	image    = data.huaweicloud_images_image.test.id
	flavor   = data.huaweicloud_compute_flavors.test.ids[0]
    key_name = huaweicloud_kps_keypair.acc_key.id
    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }
  }
}`, common.TestBaseComputeResources(rName), rName)
}

func testASGroup_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  vpc_id                   = huaweicloud_vpc.test.id
  max_instance_number      = 5
  description              = "this is a basic AS group"
  delete_publicip          = false
  delete_volume            = true

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
  lbaas_listeners {
    pool_id          = huaweicloud_lb_pool.pool_1.id
    protocol_port    = huaweicloud_lb_listener.listener_1.protocol_port
    weight           = 20
    protocol_version = "ipv4"
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName)
}

// update the following fields:
// scaling_group_name, description, agency_name, tags
func testASGroup_update(rName, newName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  vpc_id                   = huaweicloud_vpc.test.id
  max_instance_number      = 5
  description              = "this is an updated AS group"
  agency_name              = "ims_admin"
  delete_publicip          = true
  delete_volume            = false

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
  lbaas_listeners {
    pool_id          = huaweicloud_lb_pool.pool_1.id
    protocol_port    = huaweicloud_lb_listener.listener_1.protocol_port
    weight           = 30
    protocol_version = "ipv4"
  }
  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, testASGroup_Base(rName), newName)
}

func testASGroup_basic_disable(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  vpc_id                   = huaweicloud_vpc.test.id
  max_instance_number      = 5
  enable                   = false

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
  lbaas_listeners {
    pool_id       = huaweicloud_lb_pool.pool_1.id
    protocol_port = huaweicloud_lb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName)
}

func testASGroup_basic_enable(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  vpc_id                   = huaweicloud_vpc.test.id
  max_instance_number      = 5
  enable                   = true

  multi_az_scaling_policy            = "PICK_FIRST"
  cool_down_time                     = 600
  health_periodic_audit_time         = 15
  health_periodic_audit_grace_period = 900

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
  lbaas_listeners {
    pool_id       = huaweicloud_lb_pool.pool_1.id
    protocol_port = huaweicloud_lb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName)
}

func testASGroup_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  vpc_id                   = huaweicloud_vpc.test.id
  enterprise_project_id    = "%s"

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
  lbaas_listeners {
    pool_id       = huaweicloud_lb_pool.pool_1.id
    protocol_port = huaweicloud_lb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testASGroup_forceDelete(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  min_instance_number      = 2
  max_instance_number      = 5
  force_delete             = true
  vpc_id                   = huaweicloud_vpc.test.id
  delete_publicip          = true
  delete_volume            = true

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
}
`, testASGroup_Base(rName), rName)
}

func testASGroup_sourceDestCheck(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  min_instance_number      = 0
  max_instance_number      = 2
  desire_instance_number   = 2
  vpc_id                   = huaweicloud_vpc.test.id
  delete_instances         = "yes"

  networks {
    id                = huaweicloud_vpc_subnet.test.id
    source_dest_check = false
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
}
`, testASGroup_Base(rName), rName)
}

func testASGroup_sourceDestErrorCheck(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  min_instance_number      = 0
  max_instance_number      = 0
  desire_instance_number   = 2
  vpc_id                   = huaweicloud_vpc.test.id
  delete_instances         = "yes"

  networks {
    id                = huaweicloud_vpc_subnet.test.id
    source_dest_check = false
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
}
`, testASGroup_Base(rName), rName)
}
