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

func getIdentityCenterGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getIdentityCenterGroup: query Identity Center group
	var (
		getIdentityCenterGroupHttpUrl = "v1/identity-stores/{identity_store_id}/groups/{group_id}"
		getIdentityCenterGroupProduct = "identitystore"
	)
	getIdentityCenterGroupClient, err := cfg.NewServiceClient(getIdentityCenterGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	getIdentityCenterGroupPath := getIdentityCenterGroupClient.Endpoint + getIdentityCenterGroupHttpUrl
	getIdentityCenterGroupPath = strings.ReplaceAll(getIdentityCenterGroupPath, "{identity_store_id}",
		state.Primary.Attributes["identity_store_id"])
	getIdentityCenterGroupPath = strings.ReplaceAll(getIdentityCenterGroupPath, "{group_id}", state.Primary.ID)

	getIdentityCenterGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getIdentityCenterGroupResp, err := getIdentityCenterGroupClient.Request("GET", getIdentityCenterGroupPath,
		&getIdentityCenterGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center Group: %s", err)
	}
	return utils.FlattenResponse(getIdentityCenterGroupResp)
}

func TestAccIdentityCenterGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_identitycenter_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIdentityCenterGroupResourceFunc,
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
				Config: testIdentityCenterGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testIdentityCenterGroup_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testIdentityCenterGroupImportState(rName),
			},
		},
	})
}

func testIdentityCenterGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_group" "test" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  name              = "%s"
  description       = "test description"
}
`, testAccDatasourceIdentityCenter_basic(), name)
}

func testIdentityCenterGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_group" "test" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  name              = "%s"
  description       = "test description update"
}
`, testAccDatasourceIdentityCenter_basic(), name)
}

// testIdentityCenterGroupImportState use to return an id with format <identity_store_id>/<group_id>
func testIdentityCenterGroupImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		identityStoreId := rs.Primary.Attributes["identity_store_id"]
		if identityStoreId == "" {
			return "", fmt.Errorf("attribute (identity_store_id) of Resource (%s) not found: %s", name, rs)
		}
		return identityStoreId + "/" + rs.Primary.ID, nil
	}
}
