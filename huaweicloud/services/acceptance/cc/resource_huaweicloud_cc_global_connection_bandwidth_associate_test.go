package cc

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

func getGCBAssociateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CC client: %s", err)
	}
	getGCBHttpUrl := "v3/{domain_id}/gcb/gcbandwidths/{id}"
	getGCBPath := client.Endpoint + getGCBHttpUrl
	getGCBPath = strings.ReplaceAll(getGCBPath, "{domain_id}", cfg.DomainID)
	getGCBPath = strings.ReplaceAll(getGCBPath, "{id}", state.Primary.ID)

	getGCBOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getGCBResp, err := client.Request("GET", getGCBPath, &getGCBOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving global connection bandwidth: %s", err)
	}

	getGCBRespBody, err := utils.FlattenResponse(getGCBResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving global connection bandwidth: %s", err)
	}

	instances := utils.PathSearch("globalconnection_bandwidth.instances", getGCBRespBody, make([]interface{}, 0))
	if v, ok := instances.([]interface{}); ok && len(v) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}
	return instances, nil
}

func TestAccGCBAssociate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cc_global_connection_bandwidth_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGCBAssociateFunc,
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGCBAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "gcb_binding_resources.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "gcb_id"),
				),
			},
			{
				Config: testAccGCBAssociate_update_associate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "gcb_binding_resources.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "gcb_id"),
				),
			},
			{
				Config: testAccGCBAssociate_update_disassociate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "gcb_binding_resources.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "gcb_id"),
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

func testAccGCBAssociate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cc_global_connection_bandwidth_associate" test {
  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.test,
    huaweicloud_global_eip_associate.test1,
    huaweicloud_global_eip_associate.test2,
  ]

  gcb_id = huaweicloud_cc_global_connection_bandwidth.test.id

  gcb_binding_resources {
    resource_id   = huaweicloud_global_eip.test1.id
    resource_type = "GEIP"
    region_id     = "global"
    project_id    = "%[2]s"
  }
}
`, testAccGCBAssociate_base(name), acceptance.HW_DEST_PROJECT_ID_TEST)
}

func testAccGCBAssociate_update_associate(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cc_global_connection_bandwidth_associate" test {
  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.test,
    huaweicloud_global_eip_associate.test1,
    huaweicloud_global_eip_associate.test2,
  ]

  gcb_id = huaweicloud_cc_global_connection_bandwidth.test.id

  gcb_binding_resources {
    resource_id   = huaweicloud_global_eip.test1.id
    resource_type = "GEIP"
    region_id     = "global"
    project_id    = "%[2]s"
  }

  gcb_binding_resources {
    resource_id   = huaweicloud_global_eip.test2.id
    resource_type = "GEIP"
    region_id     = "global"
    project_id    = "%[2]s"
  }
}
`, testAccGCBAssociate_base(name), acceptance.HW_DEST_PROJECT_ID_TEST)
}

func testAccGCBAssociate_update_disassociate(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cc_global_connection_bandwidth_associate" test {
  depends_on = [
    huaweicloud_cc_global_connection_bandwidth.test,
    huaweicloud_global_eip_associate.test1,
    huaweicloud_global_eip_associate.test2,
  ]

  gcb_id = huaweicloud_cc_global_connection_bandwidth.test.id

  gcb_binding_resources {
    resource_id   = huaweicloud_global_eip.test2.id
    resource_type = "GEIP"
    region_id     = "global"
    project_id    = "%[2]s"
  }
}
`, testAccGCBAssociate_base(name), acceptance.HW_DEST_PROJECT_ID_TEST)
}

func testAccGCBAssociate_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cc_global_connection_bandwidth" "test" {
  name        = "%[2]s"
  type        = "Region"  
  bordercross = false
  charge_mode = "bwd"
  size        = 300
  description = "test"
  sla_level   = "Ag"
}
`, testAccGCBAssociate_GEIP(name), name)
}

func testAccGCBAssociate_GEIP(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_global_eip_pools" "all" {
  access_site = "cn-north-beijing"
}

resource "huaweicloud_global_internet_bandwidth" "test" {
  access_site = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  charge_mode = "95peak_guar"
  size        = 300
  isp         = data.huaweicloud_global_eip_pools.all.geip_pools[0].isp
  name        = "%[3]s"
  type        = data.huaweicloud_global_eip_pools.all.geip_pools[0].allowed_bandwidth_types[0].type  
}

resource "huaweicloud_global_eip" "test1" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  internet_bandwidth_id = huaweicloud_global_internet_bandwidth.test.id
  name                  = "%[3]s1"
}
	
resource "huaweicloud_global_eip" "test2" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  internet_bandwidth_id = huaweicloud_global_internet_bandwidth.test.id
  name                  = "%[3]s2"
}

resource "huaweicloud_global_eip_associate" "test1" {
  depends_on = [huaweicloud_vpc_internet_gateway.test]

  global_eip_id  = huaweicloud_global_eip.test1.id
  is_reserve_gcb = false
	  
  associate_instance {
    region        = huaweicloud_compute_instance.test1.region
    project_id    = "%[4]s"
    instance_type = "ECS"
    instance_id   = huaweicloud_compute_instance.test1.id
  }

  lifecycle {
    ignore_changes = [gc_bandwidth]
  }
}

resource "huaweicloud_global_eip_associate" "test2" {
  depends_on = [huaweicloud_vpc_internet_gateway.test]
	  
  global_eip_id  = huaweicloud_global_eip.test2.id
  is_reserve_gcb = false
	  
  associate_instance {
    region        = huaweicloud_compute_instance.test2.region
    project_id    = "%[4]s"
    instance_type = "ECS"
    instance_id   = huaweicloud_compute_instance.test2.id
  }

  lifecycle {
    ignore_changes = [gc_bandwidth]
  }
}
`, testAccGCBAssociate_IGW(name), testAccGCBAssociate_ECS(name), name, acceptance.HW_PROJECT_ID)
}

func testAccGCBAssociate_VPC(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = huaweicloud_vpc.test.id
  ipv6_enable = true
}
`, name)
}

func testAccGCBAssociate_IGW(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_internet_gateway" "test" {
  depends_on = [huaweicloud_vpc_subnet.test]

  vpc_id      = huaweicloud_vpc.test.id
  name        = "%s"
  add_route   = true
  enable_ipv6 = false
}
`, testAccGCBAssociate_VPC(name), name)
}

func testAccGCBAssociate_ECS(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 8
  memory_size       = 16
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[1]s"
}

resource "huaweicloud_compute_instance" "test1" {
  name               = "%[1]s_1"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  key_pair           = huaweicloud_kps_keypair.test.name

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_compute_instance" "test2" {
  name               = "%[1]s_2"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  key_pair           = huaweicloud_kps_keypair.test.name
  
  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, name)
}
