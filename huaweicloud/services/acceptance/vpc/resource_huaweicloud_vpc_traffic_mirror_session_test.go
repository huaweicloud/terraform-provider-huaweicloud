package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getTrafficMirrorSessionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return "", fmt.Errorf("error creating VPC v3 client: %s", err)
	}

	getTrafficMirrorSessionHttpUrl := "vpc/traffic-mirror-sessions/" + state.Primary.ID
	getTrafficMirrorSessionPath := client.ResourceBaseURL() + getTrafficMirrorSessionHttpUrl
	getTrafficMirrorSessionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getTrafficMirrorSessionResp, err := client.Request("GET", getTrafficMirrorSessionPath, &getTrafficMirrorSessionOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving traffic mirror session: %s", err)
	}

	return utils.FlattenResponse(getTrafficMirrorSessionResp)
}

func TestAccTrafficMirrorSession_basic(t *testing.T) {
	var (
		trafficMirrorSession interface{}
		name                 = acceptance.RandomAccResourceNameWithDash()
		resourceName         = "huaweicloud_vpc_traffic_mirror_session.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&trafficMirrorSession,
			getTrafficMirrorSessionResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTrafficMirrorSession_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by Terraform"),
					resource.TestCheckResourceAttr(resourceName, "priority", "10"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mirror_target_type", "eni"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "traffic_mirror_target_id",
						"huaweicloud_compute_instance.test.0", "network.0.port"),
					resource.TestCheckResourceAttrPair(resourceName, "traffic_mirror_sources.0",
						"huaweicloud_compute_instance.test.1", "network.0.port"),
					resource.TestCheckResourceAttrPair(resourceName, "traffic_mirror_filter_id",
						"huaweicloud_vpc_traffic_mirror_filter.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccTrafficMirrorSession_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "priority", "20"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mirror_target_type", "eni"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "traffic_mirror_target_id",
						"huaweicloud_compute_instance.test.3", "network.0.port"),
					resource.TestCheckResourceAttrPair(resourceName, "traffic_mirror_sources.0",
						"huaweicloud_compute_instance.test.1", "network.0.port"),
					resource.TestCheckResourceAttrPair(resourceName, "traffic_mirror_sources.1",
						"huaweicloud_compute_instance.test.2", "network.0.port"),
					resource.TestCheckResourceAttrPair(resourceName, "traffic_mirror_filter_id",
						"huaweicloud_vpc_traffic_mirror_filter.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTrafficMirrorSession_basic(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_compute_instance" "test" {
  count              = 2
  name               = "%s-${count.index}"
  description        = "terraform test"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = "c7t.large.2"
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]

  network {
    uuid              = data.huaweicloud_vpc_subnet.test.id
    source_dest_check = false
  }

  system_disk_type = "SSD"
  system_disk_size = 50

  data_disks {
    type = "SSD"
    size = "10"
  }
}

resource "huaweicloud_vpc_traffic_mirror_session" "test" {
  name                       = "%s"
  description                = "Created by Terraform"
  traffic_mirror_filter_id   = huaweicloud_vpc_traffic_mirror_filter.test.id
  traffic_mirror_sources     = [huaweicloud_compute_instance.test[1].network[0].port]
  traffic_mirror_target_id   = huaweicloud_compute_instance.test[0].network[0].port
  traffic_mirror_target_type = "eni"
  priority                   = 10
}
`, testAccCompute_data, testAccTrafficMirrorFilter_base(name, ""), name, name)
}

func testAccTrafficMirrorSession_update(name string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_compute_instance" "test" {
  count              = 4
  name               = "%s-${count.index}"
  description        = "terraform test"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = "c7t.large.2"
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]

  network {
    uuid              = data.huaweicloud_vpc_subnet.test.id
    source_dest_check = false
  }

  system_disk_type = "SSD"
  system_disk_size = 50

  data_disks {
    type = "SSD"
    size = "10"
  }
}

resource "huaweicloud_vpc_traffic_mirror_session" "test" {
  name                       = "%s-update"
  description                = ""
  traffic_mirror_filter_id   = huaweicloud_vpc_traffic_mirror_filter.test.id
  traffic_mirror_sources     = [
    huaweicloud_compute_instance.test[1].network[0].port,
    huaweicloud_compute_instance.test[2].network[0].port
  ]
  traffic_mirror_target_id   = huaweicloud_compute_instance.test[3].network[0].port
  traffic_mirror_target_type = "eni"
  priority                   = 20
  enabled                    = false
}
`, testAccCompute_data, testAccTrafficMirrorFilter_base(name, ""), name, name)
}
