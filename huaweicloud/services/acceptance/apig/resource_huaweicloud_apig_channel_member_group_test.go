package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getChannelMemberGroupsFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		instanceId    = state.Primary.Attributes["instance_id"]
		vpcChannelId  = state.Primary.Attributes["vpc_channel_id"]
		memberGroupId = state.Primary.ID
	)

	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	return apig.GetChannelMemberGroupById(client, instanceId, vpcChannelId, memberGroupId)
}

func TestAccChannelMemberGroups_basic(t *testing.T) {
	var (
		memberGroup interface{}

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		memberGroupName = "huaweicloud_apig_channel_member_group.test"
		rcMemberGroup   = acceptance.InitResourceCheck(memberGroupName, &memberGroup, getChannelMemberGroupsFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcMemberGroup.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccChannelMemberGroups_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcMemberGroup.CheckResourceExists(),
					resource.TestCheckResourceAttr(memberGroupName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(memberGroupName, "name", name),
					resource.TestCheckResourceAttrSet(memberGroupName, "vpc_channel_id"),
					resource.TestCheckResourceAttrSet(memberGroupName, "create_time"),
					resource.TestCheckResourceAttrSet(memberGroupName, "update_time"),
				),
			},
			{
				Config: testAccChannelMemberGroups_basic_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rcMemberGroup.CheckResourceExists(),
					resource.TestCheckResourceAttr(memberGroupName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(memberGroupName, "name", updateName),
					resource.TestCheckResourceAttrSet(memberGroupName, "vpc_channel_id"),
					resource.TestCheckResourceAttrSet(memberGroupName, "description"),
					resource.TestCheckResourceAttrSet(memberGroupName, "weight"),
					resource.TestCheckResourceAttrSet(memberGroupName, "create_time"),
					resource.TestCheckResourceAttrSet(memberGroupName, "update_time"),
				),
			},
			{
				ResourceName:      memberGroupName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccChannelMemberGroupsImportStateFunc(memberGroupName),
			},
		},
	})
}

func testAccChannelMemberGroupsImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", rsName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["vpc_channel_id"] == "" || rs.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", rs.Primary.Attributes["instance_id"],
				rs.Primary.Attributes["vpc_channel_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["vpc_channel_id"], rs.Primary.ID), nil
	}
}

func testAccChannelMemberGroups_basic_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

resource "huaweicloud_apig_channel" "test" {
  instance_id      = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name             = "%[1]s"
  port             = 80
  balance_strategy = 1
}`, name, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccChannelMemberGroups_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_channel_member_group" "test" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.test.id
  name           = "%[2]s"
}`, testAccChannelMemberGroups_basic_base(name), name)
}

func testAccChannelMemberGroups_basic_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_channel_member_group" "test" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.test.id
  name           = "%[2]s"
  description    = "terraform script test."
  weight         = 20
}`, testAccChannelMemberGroups_basic_base(name), updateName)
}
