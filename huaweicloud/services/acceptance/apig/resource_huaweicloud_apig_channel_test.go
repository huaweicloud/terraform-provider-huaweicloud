package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/channels"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getChannelFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return channels.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccChannel_basic(t *testing.T) {
	var (
		channel channels.Channel

		// Only letters, digits and underscores (_) are allowed in the environment name and dedicated instance name.
		rName      = "huaweicloud_apig_channel.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		baseConfig = testAccChannel_base(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&channel,
		getChannelFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckApigChannelRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccChannel_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "port", "80"),
					resource.TestCheckResourceAttr(rName, "balance_strategy", "1"),
					resource.TestCheckResourceAttr(rName, "member_type", "ecs"),
					resource.TestCheckResourceAttr(rName, "type", "2"),
					resource.TestCheckResourceAttr(rName, "health_check.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(rName, "health_check.0.threshold_normal", "1"),
					resource.TestCheckResourceAttr(rName, "health_check.0.threshold_abnormal", "1"),
					resource.TestCheckResourceAttr(rName, "health_check.0.interval", "1"),
					resource.TestCheckResourceAttr(rName, "health_check.0.timeout", "1"),
					resource.TestCheckResourceAttr(rName, "health_check.0.path", ""),
					resource.TestCheckResourceAttr(rName, "health_check.0.method", ""),
					resource.TestCheckResourceAttr(rName, "health_check.0.port", "0"),
					resource.TestCheckResourceAttr(rName, "health_check.0.http_codes", ""),
					resource.TestCheckResourceAttr(rName, "member.#", "1"),
				),
			},
			{
				Config: testAccChannel_basic_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "port", "8000"),
					resource.TestCheckResourceAttr(rName, "balance_strategy", "2"),
					resource.TestCheckResourceAttr(rName, "member_type", "ecs"),
					resource.TestCheckResourceAttr(rName, "type", "2"),
					resource.TestCheckResourceAttr(rName, "health_check.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(rName, "health_check.0.threshold_normal", "10"),
					resource.TestCheckResourceAttr(rName, "health_check.0.threshold_abnormal", "10"),
					resource.TestCheckResourceAttr(rName, "health_check.0.interval", "300"),
					resource.TestCheckResourceAttr(rName, "health_check.0.timeout", "30"),
					resource.TestCheckResourceAttr(rName, "health_check.0.path", "/terraform/"),
					resource.TestCheckResourceAttr(rName, "health_check.0.method", "HEAD"),
					resource.TestCheckResourceAttr(rName, "health_check.0.port", "8080"),
					resource.TestCheckResourceAttr(rName, "health_check.0.http_codes", "201,202,303-404"),
					resource.TestCheckResourceAttr(rName, "member.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccChannelImportStateFunc(),
			},
		},
	})
}

func testAccChannelImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "huaweicloud_apig_channel.test"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["instance_id"],
				rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccChannel_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}
`, common.TestBaseComputeResources(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccChannel_basic_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  count = 1

  name               = format("%[2]s-%%d", count.index)
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = "%[3]s"
  }
}

resource "huaweicloud_apig_channel" "test" {
  instance_id        = local.instance_id
  name               = "%[2]s"
  port               = 80
  balance_strategy   = 1
  member_type        = "ecs"
  type               = 2

  health_check {
    protocol           = "TCP"
    threshold_normal   = 1 # minimum value
    threshold_abnormal = 1 # minimum value
    interval           = 1 # minimum value
    timeout            = 1 # minimum value
  }

  dynamic "member" {
    for_each = huaweicloud_compute_instance.test[*]

    content {
      id   = member.value.id
      name = member.value.name
    }
  }
}
`, baseConfig, name, acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID)
}

func testAccChannel_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name               = format("%[2]s-%%d", count.index)
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = "%[3]s"
  }
}

resource "huaweicloud_apig_channel" "test" {
  instance_id      = local.instance_id
  name             = "%[2]s"
  port             = 8000
  balance_strategy = 2
  member_type        = "ecs"
  type               = 2

  health_check {
    protocol           = "HTTPS"
    threshold_normal   = 10  # maximum value
    threshold_abnormal = 10  # maximum value
    interval           = 300 # maximum value
    timeout            = 30  # maximum value
    path               = "/terraform/"
    method             = "HEAD"
    port               = 8080
    http_codes         = "201,202,303-404"
  }

  dynamic "member" {
    for_each = huaweicloud_compute_instance.test[*]

    content {
      id   = member.value.id
      name = member.value.name
    }
  }
}
`, baseConfig, name, acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID)
}

func TestAccChannel_eipMembers(t *testing.T) {
	var (
		channel channels.Channel

		// Only letters, digits and underscores (_) are allowed in the environment name and dedicated instance name.
		rName      = "huaweicloud_apig_channel.test"
		name       = acceptance.RandomAccResourceName()
		baseConfig = testAccChannel_eipBase(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&channel,
		getChannelFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccChannel_eipMembers_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "port", "80"),
					resource.TestCheckResourceAttr(rName, "balance_strategy", "2"),
					resource.TestCheckResourceAttr(rName, "member_type", "ip"),
					resource.TestCheckResourceAttr(rName, "type", "2"),
					resource.TestCheckResourceAttr(rName, "health_check.0.protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "health_check.0.threshold_normal", "2"),
					resource.TestCheckResourceAttr(rName, "health_check.0.threshold_abnormal", "2"),
					resource.TestCheckResourceAttr(rName, "health_check.0.interval", "60"),
					resource.TestCheckResourceAttr(rName, "health_check.0.timeout", "10"),
					resource.TestCheckResourceAttr(rName, "health_check.0.path", "/"),
					resource.TestCheckResourceAttr(rName, "health_check.0.method", "HEAD"),
					resource.TestCheckResourceAttr(rName, "health_check.0.port", "8080"),
					resource.TestCheckResourceAttr(rName, "health_check.0.http_codes", "201,202,303-404"),
					resource.TestCheckResourceAttr(rName, "member.#", "1"),
				),
			},
			{
				Config: testAccChannel_eipMembers_step2(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "member.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccChannelImportStateFunc(),
			},
		},
	})
}

func testAccChannel_eipBase(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccChannel_eipMembers_step1(baseConfig, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  count = 1

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = format("%[2]s-%%d", count.index)
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_apig_channel" "test" {
  instance_id      = local.instance_id
  name             = "%[2]s"
  port             = 80
  balance_strategy = 2
  member_type      = "ip"
  type             = 2

  health_check {
    protocol           = "HTTP"
    threshold_normal   = 2
    threshold_abnormal = 2
    interval           = 60
    timeout            = 10
    path               = "/"
    method             = "HEAD"
    port               = 8080
    http_codes         = "201,202,303-404"
  }

  dynamic "member" {
    for_each = huaweicloud_vpc_eip.test[*].address

    content {
      host = member.value
    }
  }
}
`, baseConfig, rName)
}

func testAccChannel_eipMembers_step2(baseConfig, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  count = 2

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = format("%[2]s-%%d", count.index)
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_apig_channel" "test" {
  instance_id      = local.instance_id
  name             = "%[2]s"
  port             = 80
  balance_strategy = 2
  member_type      = "ip"
  type             = 2

  health_check {
    protocol           = "HTTP"
    threshold_normal   = 2
    threshold_abnormal = 2
    interval           = 60
    timeout            = 10
    path               = "/"
    method             = "HEAD"
    port               = 8080
    http_codes         = "201,202,303-404"
  }

  dynamic "member" {
    for_each = huaweicloud_vpc_eip.test[*].address

    content {
      host = member.value
    }
  }
}
`, baseConfig, rName)
}
