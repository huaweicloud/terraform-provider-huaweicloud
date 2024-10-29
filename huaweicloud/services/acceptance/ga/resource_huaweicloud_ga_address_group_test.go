package ga

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

func getIpAddressGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/ip-groups/{ip_group_id}"
		product = "ga"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{ip_group_id}", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GA IP address group: %s", err)
	}
	return utils.FlattenResponse(resp)
}

func TestAccIpAddressGroup_basic(t *testing.T) {
	var (
		obj        interface{}
		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()
		rName      = "huaweicloud_ga_address_group.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIpAddressGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testIpAddressGroup_basic_step_1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testIpAddressGroup_basic_step_2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "ip_addresses.#", "23"),
					resource.TestCheckResourceAttrSet(rName, "status"),
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

func testIpAddressGroup_basic_step_1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ga_address_group" "test" {
  name        = "%[1]s"
  description = "Created by terraform"

  ip_addresses {
    cidr        = "192.168.0.0/24"
    description = "The first CIDR block in the IP address group"
  }
}
`, name)
}

func testIpAddressGroup_basic_step_2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ga_address_group" "test" {
  name = "%[1]s"

  dynamic "ip_addresses" {
    for_each = range(0, 23)
    content {
      cidr        = "192.168.${ip_addresses.value + 1}.0/24"
      description = "The ${ip_addresses.value + 1}th CIDR block in the IP address group"
    }
  }
}`, name)
}

func TestAccIpAddressGroup_associateListener(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceNameWithDash()
		rName = "huaweicloud_ga_address_group.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIpAddressGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testIpAddressGroup_associateListener(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform create"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.0.cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttrPair(rName, "listeners.0.id", "huaweicloud_ga_listener.test", "id"),
					resource.TestCheckResourceAttr(rName, "listeners.0.type", "WHITE"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testIpAddressGroup_update_associateListener(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform create"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.0.cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttrPair(rName, "listeners.0.id", "huaweicloud_ga_listener.listener", "id"),
					resource.TestCheckResourceAttr(rName, "listeners.0.type", "BLACK"),
					resource.TestCheckResourceAttrSet(rName, "status"),
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

func associateListener_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ga_accelerator" "test" {
  name        = "%[1]s"
  description = "terraform test"

  ip_sets {
    area = "CM"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_ga_listener" "test" {
  accelerator_id = huaweicloud_ga_accelerator.test.id
  name           = "%[1]s"
  protocol       = "TCP"
  description    = "terraform test"

  port_ranges {
    from_port = 60
    to_port   = 80
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_ga_listener" "listener" {
  accelerator_id = huaweicloud_ga_accelerator.test.id
  name           = "%[1]s"
  protocol       = "TCP"
  description    = "terraform test"

  port_ranges {
    from_port = 90
    to_port   = 100
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testIpAddressGroup_associateListener(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ga_address_group" "test" {
  name        = "%[2]s"
  description = "terraform create"

  ip_addresses {
    cidr        = "192.168.1.0/24"
    description = "The IP addresses included in the address group"
  }

  listeners {
    id   = huaweicloud_ga_listener.test.id
    type = "WHITE"
  }
}
`, associateListener_base(name), name)
}

func testIpAddressGroup_update_associateListener(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ga_address_group" "test" {
  name        = "%[2]s"
  description = "terraform create"

  ip_addresses {
    cidr        = "192.168.1.0/24"
    description = "The IP addresses included in the address group"
  }

  listeners {
    id   = huaweicloud_ga_listener.listener.id
    type = "BLACK"
  }
}
`, associateListener_base(name), name)
}
