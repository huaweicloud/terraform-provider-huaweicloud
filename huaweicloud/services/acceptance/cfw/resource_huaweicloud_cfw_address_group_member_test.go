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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAddressGroupMemberResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAddressGroupMember: Query the CFW IP address group member detail
	product := "cfw"

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW Client: %s", err)
	}

	addressGroupMembers, _, err := cfw.ReadAddressGroupMembers(state.Primary.Attributes["group_id"], client)
	if err != nil {
		return nil, fmt.Errorf("error retrieving address group member: %s", err)
	}

	findAddressGroupMemberExpr := fmt.Sprintf("[?item_id == '%s']|[0]", state.Primary.ID)
	addressGroupMember := utils.PathSearch(findAddressGroupMemberExpr, addressGroupMembers, nil)
	if addressGroupMember == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return addressGroupMember, nil
}

func testAddressGroupMemberImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["group_id"] == "" {
			return "", fmt.Errorf("attribute (group_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (id) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["group_id"] + "/" + rs.Primary.ID, nil
	}
}

func TestAccAddressGroupMember_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_address_group_member.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAddressGroupMemberResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAddressGroupMember_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "group_id", "huaweicloud_cfw_address_group.test", "id"),
					resource.TestCheckResourceAttr(rName, "address", "192.168.0.1"),
					resource.TestCheckResourceAttr(rName, "address_type", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAddressGroupMemberImportState(rName),
			},
		},
	})
}

func TestAccAddressGroupMember_addressRange(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_address_group_member.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAddressGroupMemberResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAddressGroupMember_addressRange(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "group_id", "huaweicloud_cfw_address_group.test", "id"),
					resource.TestCheckResourceAttr(rName, "address", "192.168.0.1-192.168.0.16"),
					resource.TestCheckResourceAttr(rName, "address_type", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAddressGroupMemberImportState(rName),
			},
		},
	})
}

func testAddressGroupMember_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_address_group_member" "test" {
  group_id = huaweicloud_cfw_address_group.test.id
  address  = "192.168.0.1"
}
`, testAddressGroup_basic(name))
}

func testAddressGroupMember_addressRange(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_address_group_member" "test" {
  group_id = huaweicloud_cfw_address_group.test.id
  address  = "192.168.0.1-192.168.0.16"
}
`, testAddressGroup_basic(name))
}
