package secmaster

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getLayoutResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		product     = "secmaster"
		region      = acceptance.HW_REGION_NAME
		workspaceId = state.Primary.Attributes["workspace_id"]
		id          = state.Primary.ID
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts/{layout_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{layout_id}", id)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code",
			"SecMaster.20041303")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	dataResp := utils.PathSearch("data", respBody, nil)
	if dataResp == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return dataResp, nil
}

func TestAccLayout_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceNameWithDash()
		rName = "huaweicloud_secmaster_layout.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLayoutResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLayout_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "used_by", "DATACLASS"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "thumbnail", "test"),
					resource.TestCheckResourceAttr(rName, "layout_type", "List"),
					resource.TestCheckResourceAttr(rName, "binding_code", "Alert"),
					resource.TestCheckResourceAttr(rName, "sections_sum", "10"),
					resource.TestCheckResourceAttr(rName, "tabs_sum", "10"),
					resource.TestCheckResourceAttr(rName, "boa_version", "v3"),
					resource.TestCheckResourceAttrSet(rName, "creator_id"),
					resource.TestCheckResourceAttrSet(rName, "project_id"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				Config: testLayout_update(fmt.Sprintf("%s_update", name)),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testLayoutImportState(rName),
				ImportStateVerifyIgnore: []string{"is_delete", "fields_sum", "wizards_sum"},
			},
		},
	})
}

func testLayout_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_layout" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  used_by      = "DATACLASS"
  description  = "test description"
  thumbnail    = "test"
  layout_type  = "List"
  binding_code = "Alert"
  fields_sum   = 10
  wizards_sum  = 10
  sections_sum = 10
  tabs_sum     = 10
  boa_version  = "v3"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testLayout_update(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_layout" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  used_by      = "DATACLASS"
  description  = "test description update"
  thumbnail    = "test"
  layout_type  = "List"
  binding_code = "Alert"
  fields_sum   = 20
  wizards_sum  = 20
  sections_sum = 20
  tabs_sum     = 20
  boa_version  = "v5"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, updateName)
}

func testLayoutImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		id := rs.Primary.ID
		workspaceId := rs.Primary.Attributes["workspace_id"]

		if workspaceId == "" || id == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<mapping_id>',"+
				" but got '%s/%s'",
				workspaceId, id)
		}

		return fmt.Sprintf("%s/%s", workspaceId, id), nil
	}
}
