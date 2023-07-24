package ram

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRAMShareResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getRAMShare: Query the RAM share.
	var (
		getRAMShareHttpUrl = "v1/resource-shares/search"
		getRAMShareProduct = "ram"
	)
	getRAMShareClient, err := cfg.NewServiceClient(getRAMShareProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RAM Client: %s", err)
	}
	getRAMSharePath := getRAMShareClient.Endpoint + getRAMShareHttpUrl
	getRAMShareOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"resource_share_ids": []string{state.Primary.ID},
			"resource_owner":     "self",
		},
	}
	getRAMShareResp, err := getRAMShareClient.Request("POST", getRAMSharePath, &getRAMShareOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RAM share: %s", err)
	}
	getRAMShareRespBody, err := utils.FlattenResponse(getRAMShareResp)
	if err != nil {
		return nil, err
	}

	curJson := utils.PathSearch("resource_shares", getRAMShareRespBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil, fmt.Errorf("the target resource share is not exist")
	}
	if len(curArray) > 1 {
		return nil, fmt.Errorf("except retrieving one RAM share, but got %d", len(curArray))
	}

	resourceShare := curArray[0]
	status := utils.PathSearch("status", resourceShare, "")
	if status == "deleted" {
		// The deleted resource share will exist 48 hours with "deleted" status. And will be removed after 48 hours.
		return nil, fmt.Errorf("the RAM share has been deleted")
	}
	return resourceShare, nil
}

func TestAccRAMShare_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ram_resource_share.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRAMShareResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAM(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRAMShare_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description information"),
					resource.TestCheckResourceAttr(rName, "principals.0", acceptance.HW_RAM_SHARE_ACCOUNT_ID),
					resource.TestCheckResourceAttr(rName, "resource_urns.0", acceptance.HW_RAM_SHARE_RESOURCE_URN),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "associated_permissions.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "owning_account_id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				// update principal
				Config: testRAMShare_basic_update(name, acceptance.HW_RAM_SHARE_UPDATE_ACCOUNT_ID,
					acceptance.HW_RAM_SHARE_RESOURCE_URN),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description",
						"test description information update"),
					resource.TestCheckResourceAttr(rName, "principals.0",
						acceptance.HW_RAM_SHARE_UPDATE_ACCOUNT_ID),
					resource.TestCheckResourceAttr(rName, "resource_urns.0", acceptance.HW_RAM_SHARE_RESOURCE_URN),
					resource.TestCheckResourceAttr(rName, "tags.foo_update", "bar_update"),
				),
			},
			{
				// update resource urn
				Config: testRAMShare_basic_update(name, acceptance.HW_RAM_SHARE_UPDATE_ACCOUNT_ID,
					acceptance.HW_RAM_SHARE_UPDATE_RESOURCE_URN),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "principals.0",
						acceptance.HW_RAM_SHARE_UPDATE_ACCOUNT_ID),
					resource.TestCheckResourceAttr(rName, "resource_urns.0",
						acceptance.HW_RAM_SHARE_UPDATE_RESOURCE_URN),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"permission_ids",
				},
			},
		},
	})
}

func testRAMShare_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ram_resource_share" "test" {
  name        = "%[1]s"
  description = "test description information"

  principals    = ["%[2]s"]
  resource_urns = ["%[3]s"]

  tags = {
    foo = "bar"
  }
}
`, name, acceptance.HW_RAM_SHARE_ACCOUNT_ID, acceptance.HW_RAM_SHARE_RESOURCE_URN)
}

func testRAMShare_basic_update(name, accountID, resourceUrn string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ram_resource_share" "test" {
  name        = "%[1]s_update"
  description = "test description information update"

  principals    = ["%[2]s"]
  resource_urns = ["%[3]s"]

  tags = {
    foo_update = "bar_update"
  }
}
`, name, accountID, resourceUrn)
}
