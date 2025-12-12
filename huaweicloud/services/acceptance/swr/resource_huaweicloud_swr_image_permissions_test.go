package swr

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

func getSwrImagePermissionsResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getSwrImagePermissions: Query SWR image permissions
	var (
		getSwrImagePermissionsHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/access"
		getSwrImagePermissionsProduct = "swr"
	)
	getSwrImagePermissionsClient, err := cfg.NewServiceClient(getSwrImagePermissionsProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	organization := state.Primary.Attributes["organization"]
	repository := strings.ReplaceAll(state.Primary.Attributes["repository"], "/", "$")

	getSwrImagePermissionsPath := getSwrImagePermissionsClient.Endpoint + getSwrImagePermissionsHttpUrl
	getSwrImagePermissionsPath = strings.ReplaceAll(getSwrImagePermissionsPath, "{namespace}",
		organization)
	getSwrImagePermissionsPath = strings.ReplaceAll(getSwrImagePermissionsPath, "{repository}",
		repository)

	getSwrImagePermissionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSwrImagePermissionsResp, err := getSwrImagePermissionsClient.Request("GET",
		getSwrImagePermissionsPath, &getSwrImagePermissionsOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SWR image permissions: %s", err)
	}

	return utils.FlattenResponse(getSwrImagePermissionsResp)
}

func TestAccSwrImagePermissions_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_swr_image_permissions.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSwrImagePermissionsResourceFunc,
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
				Config: testSwrImagePermissions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"huaweicloud_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"huaweicloud_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "users.#", "3"),
					resource.TestCheckResourceAttrPair(rName, "users.0.user_id",
						"huaweicloud_identity_user.user_1", "id"),
					resource.TestCheckResourceAttrPair(rName, "users.0.user_name",
						"huaweicloud_identity_user.user_1", "name"),
					resource.TestCheckResourceAttr(rName, "users.0.permission", "Read"),
					resource.TestCheckResourceAttrPair(rName, "users.1.user_id",
						"huaweicloud_identity_user.user_2", "id"),
					resource.TestCheckResourceAttrPair(rName, "users.1.user_name",
						"huaweicloud_identity_user.user_2", "name"),
					resource.TestCheckResourceAttr(rName, "users.1.permission", "Write"),
					resource.TestCheckResourceAttrPair(rName, "users.2.user_id",
						"huaweicloud_identity_user.user_3", "id"),
					resource.TestCheckResourceAttrPair(rName, "users.2.user_name",
						"huaweicloud_identity_user.user_3", "name"),
					resource.TestCheckResourceAttr(rName, "users.2.permission", "Manage"),
					resource.TestCheckResourceAttr(rName, "self_permission.#", "1"),
				),
			},
			{
				Config: testSwrImagePermissions_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"huaweicloud_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"huaweicloud_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "users.#", "4"),
					resource.TestCheckResourceAttrPair(rName, "users.0.user_id",
						"huaweicloud_identity_user.user_1", "id"),
					resource.TestCheckResourceAttrPair(rName, "users.0.user_name",
						"huaweicloud_identity_user.user_1", "name"),
					resource.TestCheckResourceAttr(rName, "users.0.permission", "Write"),
					resource.TestCheckResourceAttrPair(rName, "users.1.user_id",
						"huaweicloud_identity_user.user_2", "id"),
					resource.TestCheckResourceAttrPair(rName, "users.1.user_name",
						"huaweicloud_identity_user.user_2", "name"),
					resource.TestCheckResourceAttr(rName, "users.1.permission", "Read"),
					resource.TestCheckResourceAttrPair(rName, "users.2.user_id",
						"huaweicloud_identity_user.user_4", "id"),
					resource.TestCheckResourceAttrPair(rName, "users.2.user_name",
						"huaweicloud_identity_user.user_4", "name"),
					resource.TestCheckResourceAttr(rName, "users.2.permission", "Manage"),
					resource.TestCheckResourceAttrPair(rName, "users.3.user_id",
						"huaweicloud_identity_user.user_5", "id"),
					resource.TestCheckResourceAttrPair(rName, "users.3.user_name",
						"huaweicloud_identity_user.user_5", "name"),
					resource.TestCheckResourceAttr(rName, "users.3.permission", "Write"),
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

func testSwrImagePermissions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user" "user_1" {
  name     = "%[2]s_1"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_user" "user_2" {
  name     = "%[2]s_2"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_user" "user_3" {
  name     = "%[2]s_3"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_swr_image_permissions" "test" {
  organization = huaweicloud_swr_organization.test.name
  repository   = huaweicloud_swr_repository.test.name

  users {
    user_id    = huaweicloud_identity_user.user_1.id
    user_name  = huaweicloud_identity_user.user_1.name
    permission = "Read"
  }

  users {
    user_id    = huaweicloud_identity_user.user_2.id
    permission = "Write"
  }

  users {
    user_id    = huaweicloud_identity_user.user_3.id
    permission = "Manage"
  }
}
`, testAccSWRRepository_basic(name), name)
}

func testSwrImagePermissions_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_user" "user_1" {
  name     = "%[2]s_1"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_user" "user_2" {
  name     = "%[2]s_2"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_user" "user_4" {
  name     = "%[2]s_4"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_user" "user_5" {
  name     = "%[2]s_5"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_swr_image_permissions" "test" {
  organization = huaweicloud_swr_organization.test.name
  repository   = huaweicloud_swr_repository.test.name

  users {
    user_id    = huaweicloud_identity_user.user_1.id
    user_name  = huaweicloud_identity_user.user_1.name
    permission = "Write"
  }

  users {
    user_id    = huaweicloud_identity_user.user_2.id
    permission = "Read"
  }

  users {
    user_id    = huaweicloud_identity_user.user_4.id
    permission = "Manage"
  }

  users {
    user_id    = huaweicloud_identity_user.user_5.id
    permission = "Write"
  }
} 
`, testAccSWRRepository_basic(name), name)
}
