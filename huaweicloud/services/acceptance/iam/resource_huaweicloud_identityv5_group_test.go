package iam

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getIdentityGroupResourceFuncV5(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam_no_version", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	getGroupHttpUrl := "v5/groups/{group_id}"
	getGroupPath := client.Endpoint + getGroupHttpUrl
	getGroupPath = strings.ReplaceAll(getGroupPath, "{group_id}", state.Primary.ID)
	getGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getGroupResp, err := client.Request("GET", getGroupPath, &getGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM group: %s", err)
	}
	return utils.FlattenResponse(getGroupResp)
}

func TestAccIdentityV5Group_basic(t *testing.T) {
	var group groups.Group

	resourceName := "huaweicloud_identityv5_group.group"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getIdentityGroupResourceFuncV5,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV5Group_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "description", "kkkkkstring"),
					resource.TestCheckResourceAttr(resourceName, "group_name", "kkkkk4"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityV5Group_update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "description", "kkkkkstringupdate"),
					resource.TestCheckResourceAttr(resourceName, "group_name", "kkkkkupdate"),
				),
			},
			{
				Config: testAccIdentityV5Group_noDescription,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "group_name", "kkkkkupdate1"),
				),
			},
		},
	})
}

const testAccIdentityV5Group_basic = `
resource "huaweicloud_identityv5_group" "group" {
  group_name  = "kkkkk4"
  description = "kkkkkstring"
}
`

const testAccIdentityV5Group_update = `
resource "huaweicloud_identityv5_group" "group" {
  group_name  = "kkkkkupdate"
  description = "kkkkkstringupdate"
}
`

const testAccIdentityV5Group_noDescription = `
resource "huaweicloud_identityv5_group" "group" {
  group_name  = "kkkkkupdate1"
}
`
