package secmaster

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPlaybookResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getPlaybook: Query the SecMaster playbook detail
	var (
		getPlaybookHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks"
		getPlaybookProduct = "secmaster"
	)
	getPlaybookClient, err := cfg.NewServiceClient(getPlaybookProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	getPlaybookPath := getPlaybookClient.Endpoint + getPlaybookHttpUrl
	getPlaybookPath = strings.ReplaceAll(getPlaybookPath, "{project_id}", getPlaybookClient.ProjectID)
	getPlaybookPath = strings.ReplaceAll(getPlaybookPath, "{workspace_id}", state.Primary.Attributes["workspace_id"])

	getPlaybookqueryParams := buildGetPlaybookQueryParams(state.Primary.Attributes["name"])
	getPlaybookPath += getPlaybookqueryParams

	getPlaybookResp, err := pagination.ListAllItems(
		getPlaybookClient,
		"offset",
		getPlaybookPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving Playbook: %s", err)
	}

	getPlaybookRespJson, err := json.Marshal(getPlaybookResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Playbook: %s", err)
	}
	var getPlaybookRespBody interface{}
	err = json.Unmarshal(getPlaybookRespJson, &getPlaybookRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Playbook: %s", err)
	}

	jsonPath := fmt.Sprintf("data[?id=='%s']|[0]", state.Primary.ID)
	getPlaybookRespBody = utils.PathSearch(jsonPath, getPlaybookRespBody, nil)
	if getPlaybookRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getPlaybookRespBody, nil
}

func TestAccPlaybook_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_secmaster_playbook.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPlaybookResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMaster(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPlaybook_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testPlaybook_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPlaybookImportState(rName),
			},
		},
	})
}

func testPlaybook_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_playbook" "test" {
  workspace_id = "%s"
  name         = "%s"
  description  = "created by terraform"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testPlaybook_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_playbook" "test" {
  workspace_id = "%s"
  name         = "%s_update"
  description  = ""
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func buildGetPlaybookQueryParams(name string) string {
	if name != "" {
		return fmt.Sprintf("?name=%v", name)
	}

	return ""
}

func testPlaybookImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.ID), nil
	}
}
