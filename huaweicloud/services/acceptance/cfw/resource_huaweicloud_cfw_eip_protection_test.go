package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
)

func getEipProtectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		queryHttpUrl            = "v1/{project_id}/eips/protect"
		getProtectedEipsProduct = "cfw"
	)
	client, err := cfg.NewServiceClient(getProtectedEipsProduct, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW Client: %s", err)
	}

	resp, err := cfw.QuerySyncedEips(client, queryHttpUrl, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting protected EIPs: %s", err)
	}
	if !cfw.ProtectedEipExist(resp, state.Primary.ID) {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, nil
}

func TestAccEipProtection_basic(t *testing.T) {
	var (
		obj interface{}

		rName       = "huaweicloud_cfw_eip_protection.test"
		rName1      = "huaweicloud_cfw_eip_protection.test.0"
		rName2      = "huaweicloud_cfw_eip_protection.test.1"
		basicConfig = testEipProtection_base()

		rc = acceptance.InitResourceCheck(
			rName,
			&obj,
			getEipProtectionResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting at least two different firewall instance IDs.
			acceptance.TestAccPreCheckCfwIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: append([]resource.TestStep{
			{
				Config: testEipProtection_basic_step1(basicConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(2),
					resource.TestCheckResourceAttrPair(rName1, "object_id",
						"data.huaweicloud_cfw_firewalls.test.0", "records.0.protect_objects.0.object_id"),
					resource.TestCheckResourceAttr(rName1, "protected_eip.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(rName1, "protected_eip.*.id", "huaweicloud_vpc_eip.test.0", "id"),
					resource.TestCheckTypeSetElemAttrPair(rName1, "protected_eip.*.id", "huaweicloud_vpc_eip.test.1", "id"),
					resource.TestCheckResourceAttrPair(rName2, "object_id",
						"data.huaweicloud_cfw_firewalls.test.1", "records.0.protect_objects.0.object_id"),
					resource.TestCheckResourceAttr(rName2, "protected_eip.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(rName2, "protected_eip.*.id", "huaweicloud_vpc_eip.test.2", "id"),
					resource.TestCheckTypeSetElemAttrPair(rName2, "protected_eip.*.id", "huaweicloud_vpc_eip.test.3", "id"),
				),
			},
			{
				Config: testEipProtection_basic_step2(basicConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(2),
					resource.TestCheckResourceAttr(rName1, "protected_eip.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(rName1, "protected_eip.*.id", "huaweicloud_vpc_eip.test.1", "id"),
					resource.TestCheckResourceAttr(rName2, "protected_eip.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(rName2, "protected_eip.*.id", "huaweicloud_vpc_eip.test.3", "id"),
				),
			},
		}, rc.MultiResourcesImportSteps(2)...),
	})
}

func testEipProtection_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
locals {
  fw_instance_ids = slice(split(",", "%[1]s"), 0, 2)
}

data "huaweicloud_cfw_firewalls" "test" {
  count          = length(local.fw_instance_ids)
  fw_instance_id = local.fw_instance_ids[count.index]
}

resource "huaweicloud_vpc_eip" "test" {
  count = 4

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[2]s_${count.index}"
    size        = 10
    charge_mode = "traffic"
  }
}
`, acceptance.HW_CFW_INSTANCE_IDS, name)
}

func testEipProtection_basic_step1(basicConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_eip_protection" "test" {
  count     = length(local.fw_instance_ids)
  object_id = data.huaweicloud_cfw_firewalls.test[count.index].records[0].protect_objects[0].object_id

  dynamic "protected_eip" {
    for_each = slice(huaweicloud_vpc_eip.test[*], count.index * 2, count.index * 2 + 2)
    content {
      id          = protected_eip.value.id
      public_ipv4 = protected_eip.value.address
    }
  }
}
`, basicConfig)
}

func testEipProtection_basic_step2(basicConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_eip_protection" "test" {
  count     = length(local.fw_instance_ids)
  object_id = data.huaweicloud_cfw_firewalls.test[count.index].records[0].protect_objects[0].object_id

  dynamic "protected_eip" {
    for_each = slice(huaweicloud_vpc_eip.test[*], count.index * 2 + 1, count.index * 2 + 2)
    content {
      id          = protected_eip.value.id
      public_ipv4 = protected_eip.value.address
    }
  }
}
`, basicConfig)
}
