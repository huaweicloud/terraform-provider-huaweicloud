package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/deprecated/dedicated/v2/channels"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getVpcChannelFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return channels.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccVpcChannel_basic(t *testing.T) {
	var (
		channel channels.VpcChannel

		// Only letters, digits and underscores (_) are allowed in the environment name and dedicated instance name.
		rName      = "huaweicloud_apig_vpc_channel.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&channel,
		getVpcChannelFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcChannel_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "port", "80"),
					resource.TestCheckResourceAttr(rName, "member_type", "ECS"),
					resource.TestCheckResourceAttr(rName, "algorithm", "WRR"),
					resource.TestCheckResourceAttr(rName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "path", "/"),
					resource.TestCheckResourceAttr(rName, "members.#", "1"),
				),
			},
			{
				Config: testAccVpcChannel_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "port", "8080"),
					resource.TestCheckResourceAttr(rName, "member_type", "ECS"),
					resource.TestCheckResourceAttr(rName, "algorithm", "WLC"),
					resource.TestCheckResourceAttr(rName, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(rName, "path", "/terraform/"),
					resource.TestCheckResourceAttr(rName, "members.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVpcChannelImportStateFunc(),
			},
		},
	})
}

func TestAccVpcChannel_withEipMembers(t *testing.T) {
	var (
		channel channels.VpcChannel

		// Only letters, digits and underscores (_) are allowed in the environment name and dedicated instance name.
		rName = "huaweicloud_apig_vpc_channel.test"
		name  = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&channel,
		getVpcChannelFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcChannel_withEipMembers(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "port", "80"),
					resource.TestCheckResourceAttr(rName, "member_type", "EIP"),
					resource.TestCheckResourceAttr(rName, "algorithm", "WRR"),
					resource.TestCheckResourceAttr(rName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "path", "/"),
					resource.TestCheckResourceAttr(rName, "members.#", "1"),
				),
			},
			{
				Config: testAccVpcChannel_withEipMembersUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "members.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccVpcChannelImportStateFunc(),
			},
		},
	})
}

func testAccVpcChannelImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "huaweicloud_apig_vpc_channel.test"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["instance_id"],
				rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccVpcChannel_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, common.TestBaseComputeResources(name), name)
}

func testAccVpcChannel_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_vpc_channel" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
  port        = 80
  algorithm   = "WRR"
  protocol    = "HTTP"
  path        = "/"
  http_code   = "201"

  members {
    id = huaweicloud_compute_instance.test.id
  }
}
`, testAccVpcChannel_base(name), name)
}

func testAccVpcChannel_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "newone" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_apig_vpc_channel" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
  port        = 8080
  algorithm   = "WLC"
  protocol    = "HTTPS"
  path        = "/terraform/"
  http_code   = "201,202,203"

  members {
    id     = huaweicloud_compute_instance.test.id
    weight = 30
  }
  members {
    id     = huaweicloud_compute_instance.newone.id
    weight = 70
  }
}
`, testAccVpcChannel_base(name), name)
}

func testAccVpcChannel_eipBase(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_vpc_eip" "newone" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s_newone"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccVpcChannel_withEipMembers(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_vpc_channel" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
  port        = 80
  member_type = "EIP"
  algorithm   = "WRR"
  protocol    = "HTTP"
  path        = "/"
  http_code   = "201"

  members {
    ip_address = huaweicloud_vpc_eip.test.address
  }
}
`, testAccVpcChannel_eipBase(rName), rName)
}

func testAccVpcChannel_withEipMembersUpdate(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_vpc_channel" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
  port        = 80
  member_type = "EIP"
  algorithm   = "WRR"
  protocol    = "HTTP"
  path        = "/"
  http_code   = "201"

  members {
    ip_address = huaweicloud_vpc_eip.test.address
    weight     = 30
  }
  members {
    ip_address = huaweicloud_vpc_eip.newone.address
    weight     = 70
  }
}
`, testAccVpcChannel_eipBase(rName), rName, rName)
}
