package ram

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

func getRAMAssociatedPermissionFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	resourceShareId := acceptance.HW_RAM_SHARE_ID
	permissionId := acceptance.HW_RAM_PERMISSION_ID

	var (
		getRAMAssociatedPermissionHttpUrl = "v1/resource-shares/{resource_share_id}/associated-permissions"
		getRAMAssociatedPermissionProduct = "ram"
	)
	getRAMAssociatedPermissionClient, err := cfg.NewServiceClient(getRAMAssociatedPermissionProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RAM client: %s", err)
	}

	getRAMAssociatedPermissionPath := getRAMAssociatedPermissionClient.Endpoint + getRAMAssociatedPermissionHttpUrl
	getRAMAssociatedPermissionPath = strings.ReplaceAll(getRAMAssociatedPermissionPath, "{resource_share_id}", resourceShareId)
	getRAMShareOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var associatedPermission interface{}
	var marker string
	var queryPath string

	for {
		queryPath = getRAMAssociatedPermissionPath + buildGetAssociatedPermissionQueryParams(marker)
		getRAMAssociatedPermissionResp, err := getRAMAssociatedPermissionClient.Request("GET", queryPath, &getRAMShareOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving associated permissions: %s", err)
		}

		getRAMAssociatedPermissionRespBody, err := utils.FlattenResponse(getRAMAssociatedPermissionResp)
		if err != nil {
			return nil, err
		}

		associatedPermission = utils.PathSearch(
			fmt.Sprintf("associated_permissions[?permission_id=='%s'&&status=='associated']|[0]", permissionId),
			getRAMAssociatedPermissionRespBody,
			nil,
		)
		if associatedPermission != nil {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", getRAMAssociatedPermissionRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	if associatedPermission == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return associatedPermission, nil
}

func buildGetAssociatedPermissionQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func TestAccRAMSharePermission_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_ram_resource_share_permission.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRAMAssociatedPermissionFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAMResourceSharePermission(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRAMSharePermssion_basic(acceptance.HW_RAM_SHARE_ID, acceptance.HW_RAM_PERMISSION_ID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"replace",
				},
			},
		},
	})
}

func testRAMSharePermssion_basic(resourceId string, permissionId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ram_resource_share_permission" "test" {
  resource_share_id = "%[1]s"
  permission_id     = "%[2]s"
}
`, resourceId, permissionId)
}
