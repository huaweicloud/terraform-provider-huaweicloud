package apig

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/channels"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApigVpcChannelV2_basic(t *testing.T) {
	var (
		// The dedicated instance name only allow letters, digits and underscores (_).
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_vpc_channel.test"
		channel      channels.VpcChannel
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigVpcChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigVpcChannel_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigVpcChannelExists(resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "member_type", "ECS"),
					resource.TestCheckResourceAttr(resourceName, "algorithm", "WRR"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
				),
			},
			{
				Config: testAccApigVpcChannel_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigVpcChannelExists(resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "member_type", "ECS"),
					resource.TestCheckResourceAttr(resourceName, "algorithm", "WLC"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "path", "/terraform/"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigSubResNameImportStateFunc(resourceName),
			},
		},
	})
}

func TestAccApigVpcChannelV2_withEipMembers(t *testing.T) {
	var (
		// The dedicated instance name only allow letters, digits and underscores (_).
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_vpc_channel.test"
		channel      channels.VpcChannel
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigVpcChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigVpcChannel_withEipMembers(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigVpcChannelExists(resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "member_type", "EIP"),
					resource.TestCheckResourceAttr(resourceName, "algorithm", "WRR"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "path", "/"),
					resource.TestCheckResourceAttr(resourceName, "members.#", "1"),
				),
			},
			{
				Config: testAccApigVpcChannel_withEipMembersUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigVpcChannelExists(resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "members.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigSubResNameImportStateFunc(resourceName),
			},
		},
	})
}

func testAccCheckApigVpcChannelDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_vpc_channel" {
			continue
		}
		_, err := channels.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("APIG v2 Vpc Channel (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigVpcChannelExists(n string, app *channels.VpcChannel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no vpc channel id")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
		}
		found, err := channels.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*app = *found
		return nil
	}
}

func testAccApigSubResNameImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["instance_id"],
				rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccApigVpcChannel_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, testAccApigApplication_base(rName), rName)
}

func testAccApigVpcChannel_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_vpc_channel" "test" {
  name        = "%s"
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
`, testAccApigVpcChannel_base(rName), rName)
}

func testAccApigVpcChannel_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "newone" {
  name               = "%s"
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
  name        = "%s_update"
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
`, testAccApigVpcChannel_base(rName), rName, rName)
}

func testAccApigVpcChannel_eipBase(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}
`, rName)
}

func testAccApigVpcChannel_withEipMembers(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_apig_vpc_channel" "test" {
  name        = "%s"
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
`, testAccApigApplication_base(rName), testAccApigVpcChannel_eipBase(rName), rName)
}

func testAccApigVpcChannel_withEipMembersUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_vpc_eip" "newone" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%s_newone"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_apig_vpc_channel" "test" {
  name        = "%s"
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
`, testAccApigApplication_base(rName), testAccApigVpcChannel_eipBase(rName), rName, rName)
}
