package rms

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_resource_histories.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRmsHistories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "items.0.resource.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSource, "items.0.resource.0.tags.key", "value"),
					resource.TestCheckResourceAttrPair(dataSource, "items.0.resource_id", "huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.capture_time"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.relations.0.from_resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.relations.0.from_resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.relations.0.relation_type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.relations.0.to_resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.relations.0.to_resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.checksum"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.created"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.ep_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.ep_name"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.provider"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.provisioning_state"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.updated"),
				),
			},
			{
				Config: testAccDataSourceRmsHistories_withTimeRange(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "items.0.resource.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSource, "items.0.resource.0.tags.key", "value"),
					resource.TestCheckResourceAttrPair(dataSource, "items.0.resource_id", "huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.capture_time"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.relations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.resource.0.name"),
				),
			},
		},
	})
}

func testDataSourceRmsHistories_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_histories" "test" {
  resource_id = huaweicloud_compute_instance.test.id
}
`, testDataSourceRmsHistories_base())
}

func testAccDataSourceRmsHistories_withTimeRange() string {
	currentTime := time.Now().UTC()
	tenMinutesAgo := currentTime.Add(-10 * time.Minute)
	tenMinutesLater := currentTime.Add(10 * time.Minute)
	earlyTimeString := tenMinutesAgo.Format("2006-01-02 15:04:05")
	laterTimeString := tenMinutesLater.Format("2006-01-02 15:04:05")

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_histories" "test" {
  resource_id  = huaweicloud_compute_instance.test.id
  earlier_time = "%[2]s"
  later_time   = "%[3]s"
}
`, testDataSourceRmsHistories_base(), earlyTimeString, laterTimeString)
}

func testDataSourceRmsHistories_base() string {
	name := acceptance.RandomAccResourceName()
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
	  
data "huaweicloud_availability_zones" "test" {}
	  
data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 8
  memory_size       = 16
}
	  
data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 22.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "egress"
  ethertype         = "IPv4"
  action            = "allow"
  priority          = 1
  remote_ip_prefix  = "0.0.0.0/0"
}
	  
resource "huaweicloud_compute_instance" "test" {
  name                = "%[1]s"
  image_id            = data.huaweicloud_images_image.test.id
  flavor_id           = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  security_group_ids  = [huaweicloud_networking_secgroup.test.id]
  stop_before_destroy = true
  agent_list          = "ces"
	  
  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  tags = {
    foo = "bar"
    key = "value"
  }
	  
  provisioner "local-exec" {
    command = "sleep 60"
  }
}`, name)
}
