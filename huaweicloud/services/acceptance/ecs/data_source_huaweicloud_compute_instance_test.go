package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getEcsInstanceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.ComputeV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating compute v1 client: %s", err)
	}

	resourceID := state.Primary.ID
	found, err := cloudservers.Get(client, resourceID).Extract()
	if err == nil && found.Status == "DELETED" {
		return nil, fmt.Errorf("the resource %s has been deleted", resourceID)
	}

	return found, err
}

func TestAccComputeInstanceDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.huaweicloud_compute_instance.this"
	var instance cloudservers.CloudServer

	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rc := acceptance.InitResourceCheck(
		"huaweicloud_compute_instance.test",
		&instance,
		getEcsInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "system_disk_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume_attached.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_compute_instance.byID", "status"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_compute_instance.byTags", "status"),
					resource.TestCheckResourceAttrPair("data.huaweicloud_compute_instance.byIP", "network.0.fixed_ip_v4",
						"huaweicloud_compute_instance.test", "network.0.fixed_ip_v4"),
				),
			},
		},
	})
}

func testAccComputeInstanceDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  tags = {
    server_name = "%[2]s"
  }
}

data "huaweicloud_compute_instance" "this" {
  name = huaweicloud_compute_instance.test.name
}

data "huaweicloud_compute_instance" "byID" {
  instance_id = huaweicloud_compute_instance.test.id
}

data "huaweicloud_compute_instance" "byIP" {
  fixed_ip_v4 = huaweicloud_compute_instance.test.network[0].fixed_ip_v4
}

data "huaweicloud_compute_instance" "byTags" {
  tags = {
    server_name = "%[2]s"
  }

  depends_on = [huaweicloud_compute_instance.test]
}
`, testAccCompute_data, rName)
}
