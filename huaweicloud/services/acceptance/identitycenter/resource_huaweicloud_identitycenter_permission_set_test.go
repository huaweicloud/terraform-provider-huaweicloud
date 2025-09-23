package identitycenter

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

func getPermissionSetResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	getPermissionSetClient, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center client: %s", err)
	}

	getPermissionSetHttpUrl := "v1/instances/{instance_id}/permission-sets/{id}"
	getPermissionSetPath := getPermissionSetClient.Endpoint + getPermissionSetHttpUrl
	getPermissionSetPath = strings.ReplaceAll(getPermissionSetPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getPermissionSetPath = strings.ReplaceAll(getPermissionSetPath, "{id}", state.Primary.ID)

	getPermissionSetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPermissionSetResp, err := getPermissionSetClient.Request("GET", getPermissionSetPath, &getPermissionSetOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving permission set: %s", err)
	}

	return utils.FlattenResponse(getPermissionSetResp)
}

func TestAccPermissionSet_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_identitycenter_permission_set.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPermissionSetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPermissionSet_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "session_duration", "PT8H"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "data.huaweicloud_identitycenter_instance.system", "id"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
				),
			},
			{
				Config: testPermissionSet_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "session_duration", "PT4H"),
					resource.TestCheckResourceAttr(rName, "description", "updated by terraform"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateIdFunc: testPermissionSetImportState(rName),
				ImportStateVerify: true,
			},
			{
				Config: testPermissionSet_without_desc(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "session_duration", "PT4H"),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
		},
	})
}

func testPermissionSetImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		if instanceID == "" {
			return "", fmt.Errorf("attribute (instance_id) of resource (%s) not found: %s", name, rs)
		}

		return instanceID + "/" + rs.Primary.ID, nil
	}
}

func testPermissionSet_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "system" {}

resource "huaweicloud_identitycenter_permission_set" "test" {
  instance_id      = data.huaweicloud_identitycenter_instance.system.id
  name             = "%s"
  session_duration = "PT8H"
  description      = "created by terraform"

  tags = {
    foo = "bar"
  }
}
`, name)
}

func testPermissionSet_update(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "system" {}

resource "huaweicloud_identitycenter_permission_set" "test" {
  instance_id      = data.huaweicloud_identitycenter_instance.system.id
  name             = "%s"
  session_duration = "PT4H"
  description      = "updated by terraform"
}
`, name)
}

func testPermissionSet_without_desc(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "system" {}

resource "huaweicloud_identitycenter_permission_set" "test" {
  instance_id      = data.huaweicloud_identitycenter_instance.system.id
  name             = "%s"
  session_duration = "PT4H"

  tags = {
    foo = "bar_update"
  }
}
`, name)
}
