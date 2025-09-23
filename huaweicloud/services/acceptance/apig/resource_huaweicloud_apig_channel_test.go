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
		channel interface{}

		// Only letters, digits and underscores (_) are allowed in the environment name and dedicated instance name.
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		typeServer       = "huaweicloud_apig_channel.type_server"
		typeServerLegacy = "huaweicloud_apig_channel.type_server_legacy"
		typeReference    = "huaweicloud_apig_channel.type_reference"

		rcTypeServer       = acceptance.InitResourceCheck(typeServer, &channel, getChannelFunc)
		rcTypeServerLegacy = acceptance.InitResourceCheck(typeServerLegacy, &channel, getChannelFunc)
		rcTypeReference    = acceptance.InitResourceCheck(typeReference, &channel, getChannelFunc)

		baseConfig = testAccChannel_base(name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckApigChannelRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcTypeServer.CheckResourceDestroy(),
			rcTypeServerLegacy.CheckResourceDestroy(),
			rcTypeReference.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccChannel_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rcTypeServer.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeServer, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeServer, "name", name+"_type_server"),
					resource.TestCheckResourceAttr(typeServer, "port", "80"),
					resource.TestCheckResourceAttr(typeServer, "balance_strategy", "1"),
					resource.TestCheckResourceAttr(typeServer, "member_type", "ecs"),
					resource.TestCheckResourceAttr(typeServer, "type", "builtin"),
					resource.TestCheckResourceAttr(typeServer, "health_check.#", "1"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.threshold_normal", "1"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.threshold_abnormal", "1"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.interval", "1"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.timeout", "1"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.path", ""),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.method", ""),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.port", "0"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.http_codes", ""),
					resource.TestCheckResourceAttr(typeServer, "member.#", "1"),
					resource.TestCheckResourceAttrPair(typeServer, "member.0.id", "huaweicloud_compute_instance.test.0", "id"),
					resource.TestCheckResourceAttrPair(typeServer, "member.0.name", "huaweicloud_compute_instance.test.0", "name"),
					// If `member.port` is not specified, this value is the channel's `port`.
					resource.TestCheckResourceAttr(typeServer, "member.0.port", "80"),
					resource.TestCheckResourceAttr(typeServer, "member.0.weight", "0"),
					resource.TestCheckResourceAttr(typeServer, "member.0.is_backup", "false"),
					resource.TestCheckResourceAttr(typeServer, "member.0.status", "1"),
					rcTypeServerLegacy.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeServerLegacy, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeServerLegacy, "name", name+"_type_server_legacy"),
					resource.TestCheckResourceAttr(typeServerLegacy, "port", "81"),
					resource.TestCheckResourceAttr(typeServerLegacy, "balance_strategy", "1"),
					resource.TestCheckResourceAttr(typeServerLegacy, "member_type", "ecs"),
					resource.TestCheckResourceAttr(typeServerLegacy, "type", "builtin"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.#", "1"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.threshold_normal", "1"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.threshold_abnormal", "1"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.interval", "1"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.timeout", "1"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.path", ""),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.method", ""),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.port", "0"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.http_codes", ""),
					resource.TestCheckResourceAttr(typeServerLegacy, "member.#", "1"),
					rcTypeReference.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeReference, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeReference, "name", name+"_type_reference"),
					resource.TestCheckResourceAttr(typeReference, "port", "82"),
					resource.TestCheckResourceAttr(typeReference, "balance_strategy", "1"),
					resource.TestCheckResourceAttr(typeReference, "member_type", "ecs"),
					resource.TestCheckResourceAttr(typeReference, "type", "reference"),
					resource.TestCheckResourceAttr(typeReference, "member_group.#", "1"),
					resource.TestCheckResourceAttr(typeReference, "member_group.0.name", name),
					resource.TestCheckResourceAttr(typeReference, "member_group.0.description", "Created by terraform script"),
					resource.TestCheckResourceAttr(typeReference, "member_group.0.weight", "1"),
					resource.TestCheckResourceAttrPair(typeReference, "member_group.0.reference_vpc_channel_id",
						"huaweicloud_apig_channel.type_server", "id"),
					resource.TestCheckResourceAttr(typeReference, "health_check.#", "1"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.threshold_normal", "1"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.threshold_abnormal", "1"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.interval", "1"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.timeout", "1"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.path", ""),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.method", ""),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.port", "0"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.http_codes", ""),
					resource.TestCheckResourceAttr(typeReference, "member.#", "0"),
				),
			},
			{
				Config: testAccChannel_basic_step2(baseConfig, name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rcTypeServer.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeServer, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeServer, "name", updateName+"_type_server"),
					resource.TestCheckResourceAttr(typeServer, "port", "8000"),
					resource.TestCheckResourceAttr(typeServer, "balance_strategy", "2"),
					resource.TestCheckResourceAttr(typeServer, "member_type", "ecs"),
					resource.TestCheckResourceAttr(typeServer, "type", "builtin"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.threshold_normal", "10"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.threshold_abnormal", "10"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.interval", "300"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.timeout", "30"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.path", "/terraform/"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.method", "HEAD"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.port", "8080"),
					resource.TestCheckResourceAttr(typeServer, "health_check.0.http_codes", "201,202,303-404"),
					resource.TestCheckResourceAttr(typeServer, "member_group.#", "2"),
					resource.TestCheckResourceAttrSet(typeServer, "member_group.0.name"),
					resource.TestCheckResourceAttr(typeServer, "member.#", "2"),
					resource.TestCheckResourceAttrSet(typeServer, "member.0.id"),
					resource.TestCheckResourceAttrSet(typeServer, "member.0.name"),
					resource.TestCheckResourceAttrSet(typeServer, "member.0.group_name"),
					resource.TestCheckResourceAttrSet(typeServer, "member.1.id"),
					resource.TestCheckResourceAttrSet(typeServer, "member.1.name"),
					resource.TestCheckResourceAttrSet(typeServer, "member.1.group_name"),
					rcTypeServerLegacy.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeServerLegacy, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeServerLegacy, "name", updateName+"_type_server_legacy"),
					resource.TestCheckResourceAttr(typeServerLegacy, "port", "8001"),
					resource.TestCheckResourceAttr(typeServerLegacy, "balance_strategy", "2"),
					resource.TestCheckResourceAttr(typeServerLegacy, "member_type", "ecs"),
					resource.TestCheckResourceAttr(typeServerLegacy, "type", "builtin"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.threshold_normal", "10"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.threshold_abnormal", "10"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.interval", "300"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.timeout", "30"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.path", "/terraform/"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.method", "HEAD"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.port", "8080"),
					resource.TestCheckResourceAttr(typeServerLegacy, "health_check.0.http_codes", "201,202,303-404"),
					resource.TestCheckResourceAttr(typeServerLegacy, "member.#", "2"),
					resource.TestCheckResourceAttr(typeServerLegacy, "member.0.status", "2"),
					rcTypeReference.CheckResourceExists(),
					resource.TestCheckResourceAttr(typeReference, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(typeReference, "name", updateName+"_type_reference"),
					resource.TestCheckResourceAttr(typeReference, "port", "8002"),
					resource.TestCheckResourceAttr(typeReference, "balance_strategy", "2"),
					resource.TestCheckResourceAttr(typeReference, "member_type", "ecs"),
					resource.TestCheckResourceAttr(typeReference, "type", "reference"),
					resource.TestCheckResourceAttr(typeReference, "member_group.#", "1"),
					resource.TestCheckResourceAttr(typeReference, "member_group.0.name", updateName),
					resource.TestCheckResourceAttr(typeReference, "member_group.0.description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(typeReference, "member_group.0.weight", "2"),
					resource.TestCheckResourceAttrPair(typeReference, "member_group.0.reference_vpc_channel_id",
						"huaweicloud_apig_channel.type_server_legacy", "id"),
					resource.TestCheckResourceAttr(typeReference, "health_check.#", "1"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.threshold_normal", "2"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.threshold_abnormal", "5"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.interval", "10"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.timeout", "5"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.path", "/terraform/"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.method", "GET"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.port", "50"),
					resource.TestCheckResourceAttr(typeReference, "health_check.0.http_codes", "500"),
					resource.TestCheckResourceAttr(typeReference, "member.#", "0"),
				),
			},
			{
				ResourceName:      typeServer,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccChannelImportStateFunc(typeServer),
			},
			{
				ResourceName:      typeServerLegacy,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccChannelImportStateFunc(typeServerLegacy),
			},
			{
				ResourceName:      typeReference,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccChannelImportStateFunc(typeReference),
			},
		},
	})
}

func testAccChannelImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", rsName, rs)
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

resource "huaweicloud_apig_channel" "type_server" {
  instance_id        = local.instance_id
  name               = "%[2]s_type_server"
  port               = 80
  balance_strategy   = 1
  member_type        = "ecs"
  type               = "builtin"

  health_check {
    protocol           = "TCP"
    threshold_normal   = 1 # minimum value
    threshold_abnormal = 1 # minimum value
    interval           = 1 # minimum value
    timeout            = 1 # minimum value
  }

  member_group {
    name   = "%[2]s"
    weight = 1
  }

  dynamic "member" {
    for_each = huaweicloud_compute_instance.test[*]

    content {
      id   = member.value.id
      name = member.value.name
    }
  }
}


resource "huaweicloud_apig_channel" "type_server_legacy" {
  instance_id        = local.instance_id
  name               = "%[2]s_type_server_legacy"
  port               = 81
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

resource "huaweicloud_apig_channel" "type_reference" {
  instance_id      = local.instance_id
  name             = "%[2]s_type_reference"
  port             = 82
  balance_strategy = 1
  member_type      = "ecs"
  type             = "reference"

  member_group {
    name                     = "%[2]s"
    description              = "Created by terraform script"
    weight                   = 1
    reference_vpc_channel_id = huaweicloud_apig_channel.type_server.id
  }

  health_check {
    protocol           = "TCP"
    threshold_normal   = 1 # minimum value
    threshold_abnormal = 1 # minimum value
    interval           = 1 # minimum value
    timeout            = 1 # minimum value
  }
}
`, baseConfig, name, acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID)
}

func testAccChannel_basic_step2(baseConfig, name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name               = format("%[4]s-%%d", count.index)
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = "%[3]s"
  }
}

resource "huaweicloud_apig_channel" "type_server" {
  instance_id      = local.instance_id
  name             = "%[2]s_type_server"
  port             = 8000
  balance_strategy = 2
  member_type      = "ecs"
  type             = "builtin"

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

  member_group {
    name   = "%[2]s"
    weight = 1
  }
  member_group {
    name = "%[2]s_2"
  }

  member {
    id         = huaweicloud_compute_instance.test[0].id
    name       = huaweicloud_compute_instance.test[0].name
	group_name = "%[2]s"
  }
  member {
    id         = huaweicloud_compute_instance.test[1].id
    name       = huaweicloud_compute_instance.test[1].name
    port       = 8009
    status     = 2
    is_backup  = true
	group_name = "%[2]s_2"
  }
}

resource "huaweicloud_apig_channel" "type_server_legacy" {
  instance_id      = local.instance_id
  name             = "%[2]s_type_server_legacy"
  port             = 8001
  balance_strategy = 2
  member_type      = "ecs"
  type             = 2

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
      id     = member.value.id
      name   = member.value.name
      status = 2
    }
  }
}

resource "huaweicloud_apig_channel" "type_reference" {
  instance_id      = local.instance_id
  name             = "%[2]s_type_reference"
  port             = 8002
  balance_strategy = 2
  member_type      = "ecs"
  type             = "reference"

  member_group {
    name                     = "%[2]s"
    description              = "Updated by terraform script"
    weight                   = 2
    reference_vpc_channel_id = huaweicloud_apig_channel.type_server_legacy.id
  }

  health_check {
    protocol           = "HTTPS"
    threshold_normal   = 2  # default value
    threshold_abnormal = 5  # default value
    interval           = 10 # default value
    timeout            = 5  # default value
    path               = "/terraform/"
    method             = "GET"
    port               = "50"
    http_codes         = "500"
  }
}
`, baseConfig, updateName, acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID, name)
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
					resource.TestCheckResourceAttr(rName, "type", "builtin"),
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
				ImportStateIdFunc: testAccChannelImportStateFunc(rName),
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
