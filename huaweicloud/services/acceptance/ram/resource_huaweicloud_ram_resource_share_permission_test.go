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

func getAssociatedPermissionFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	var (
		region                   = acceptance.HW_REGION_NAME
		resourceShareId          = acceptance.HW_RAM_SHARE_ID
		permissionId             = acceptance.HW_RAM_PERMISSION_ID
		httpUrl                  = "v1/resource-shares/{resource_share_id}/associated-permissions"
		product                  = "ram"
		marker                   string
		allAssociatedPermissions = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RAM client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{resource_share_id}", resourceShareId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := requestPath + buildGetAssociatedPermissionsQueryParams(marker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving RAM shared resource permissions: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		permissionsResp := utils.PathSearch("associated_permissions", respBody, make([]interface{}, 0)).([]interface{})
		allAssociatedPermissions = append(allAssociatedPermissions, permissionsResp...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	jsonPath := fmt.Sprintf("[?permission_id=='%s'&&status=='associated']|[0]", permissionId)
	associatedPermission := utils.PathSearch(jsonPath, allAssociatedPermissions, nil)
	if associatedPermission == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return associatedPermission, nil
}

func buildGetAssociatedPermissionsQueryParams(marker string) string {
	queryParams := "?limit=2000"
	if marker != "" {
		queryParams = fmt.Sprintf("%s&marker=%v", queryParams, marker)
	}

	return queryParams
}

func TestAccResourceSharePermission_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_ram_resource_share_permission.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAssociatedPermissionFunc,
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
				Config: testResourceSharePermssion_basic(acceptance.HW_RAM_SHARE_ID, acceptance.HW_RAM_PERMISSION_ID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "resource_share_id", acceptance.HW_RAM_SHARE_ID),
					resource.TestCheckResourceAttr(rName, "permission_id", acceptance.HW_RAM_PERMISSION_ID),
					resource.TestCheckResourceAttrSet(rName, "permission_name"),
					resource.TestCheckResourceAttrSet(rName, "resource_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
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

func testResourceSharePermssion_basic(resourceId string, permissionId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ram_resource_share_permission" "test" {
  resource_share_id = "%[1]s"
  permission_id     = "%[2]s"
  replace           = true
}
`, resourceId, permissionId)
}
