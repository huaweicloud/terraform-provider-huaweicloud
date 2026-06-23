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

func getLayoutWizardResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		product     = "secmaster"
		region      = acceptance.HW_REGION_NAME
		workspaceId = state.Primary.Attributes["workspace_id"]
		id          = state.Primary.ID
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/layouts/wizards/{wizard_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{wizard_id}", id)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code",
			"SecMaster.20097006")
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

func TestAccLayoutWizard_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceNameWithDash()
		rName = "huaweicloud_secmaster_layout_wizard.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLayoutWizardResourceFunc,
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
				Config: testLayoutWizard_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "layout_id", "huaweicloud_secmaster_layout.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "wizard_json", "test wizard json"),
					resource.TestCheckResourceAttr(rName, "is_binding", "true"),
					resource.TestCheckResourceAttr(rName, "boa_version", "v3"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				Config: testLayoutWizard_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
					resource.TestCheckResourceAttr(rName, "wizard_json", "test wizard json update"),
					resource.TestCheckResourceAttr(rName, "is_binding", "false"),
					resource.TestCheckResourceAttr(rName, "boa_version", "v5"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testLayoutWizardImportState(rName),
				ImportStateVerifyIgnore: []string{"layout_id"},
			},
		},
	})
}

func testLayoutWizard_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_layout" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s_layout"
  used_by      = "DATACLASS"
  layout_type  = "List"
  binding_code = "Alert"
  boa_version  = "v3"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testLayoutWizard_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_secmaster_layout_wizard" "test" {
  workspace_id = "%[2]s"
  layout_id    = huaweicloud_secmaster_layout.test.id
  name         = "%[3]s"
  description  = "test description"

  wizard_json = "test wizard json"

  is_binding  = "true"
  boa_version = "v3"
}
`, testLayoutWizard_base(name), acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testLayoutWizard_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_secmaster_layout_wizard" "test" {
  workspace_id = "%[2]s"
  layout_id    = huaweicloud_secmaster_layout.test.id
  name         = "%[3]s_update"
  description  = "test description update"

  wizard_json = "test wizard json update"

  is_binding  = "false"
  boa_version = "v5"
}
`, testLayoutWizard_base(name), acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testLayoutWizardImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		workspaceId := rs.Primary.Attributes["workspace_id"]
		id := rs.Primary.ID
		if workspaceId == "" || id == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>',"+
				" but got '%s/%s'", workspaceId, id)
		}

		return fmt.Sprintf("%s/%s", workspaceId, id), nil
	}
}
