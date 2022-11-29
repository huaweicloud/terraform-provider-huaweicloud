package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccComputeInstancesDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_compute_instances.test"
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
				Config: testAccComputeInstancesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.network.#"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr("data.huaweicloud_compute_instances.byID", "instances.#", "1"),
				),
			},
		},
	})
}

func testAccComputeInstancesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  tags = {
    foo = "bar"
  }
}

data "huaweicloud_compute_instances" "test" {
  name = huaweicloud_compute_instance.test.name
}

data "huaweicloud_compute_instances" "byID" {
  instance_id = huaweicloud_compute_instance.test.id
}
`, testAccCompute_data, rName)
}
