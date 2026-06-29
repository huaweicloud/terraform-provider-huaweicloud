package secmaster

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

func getSecurityReportResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		product     = "secmaster"
		region      = acceptance.HW_REGION_NAME
		workspaceId = state.Primary.Attributes["workspace_id"]
		id          = state.Primary.ID
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/sa/reports/{report_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = strings.ReplaceAll(requestPath, "{report_id}", id)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	if utils.PathSearch("id", respBody, "").(string) == "" {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func TestAccSecurityReport_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceNameWithDash()
		rName = "huaweicloud_secmaster_security_report.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSecurityReportResourceFunc,
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
				Config: testSecurityReport_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "report_name", name),
					resource.TestCheckResourceAttr(rName, "report_period", "daily"),
					resource.TestCheckResourceAttrPair(rName, "layout_id", "huaweicloud_secmaster_layout.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "disable"),
				),
			},
			{
				Config: testSecurityReport_updateStatus(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "enable"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testSecurityReportImportState(rName),
				ImportStateVerifyIgnore: []string{
					"report_range",
					"language",
				},
			},
		},
	})
}

func testSecurityReport_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_layout" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s_layout"
  used_by      = "SECURITY_REPORT"
  boa_version  = "v3"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testSecurityReport_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_secmaster_security_report" "test" {
  workspace_id   = "%[2]s"
  report_name    = "%[3]s"
  report_period  = "daily"
  language       = "zh-cn"
  layout_id      = huaweicloud_secmaster_layout.test.id
  binding_wizard = jsonencode({
    binding_buttons = []
    boa_version     = "v4"
  })

  report_range {
    start = "1782262800000"
    end   = "1782349200000"
  }

  status = "disable"
}
`, testSecurityReport_base(name), acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testSecurityReport_updateStatus(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_secmaster_security_report" "test" {
  workspace_id   = "%[2]s"
  report_name    = "%[3]s"
  report_period  = "daily"
  language       = "zh-cn"
  layout_id      = huaweicloud_secmaster_layout.test.id
  binding_wizard = jsonencode({
    binding_buttons = []
    boa_version     = "v4"
  })

  report_range {
    start = "1782262800000"
    end   = "1782349200000"
  }

  status = "enable"
}
`, testSecurityReport_base(name), acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testSecurityReportImportState(name string) resource.ImportStateIdFunc {
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
