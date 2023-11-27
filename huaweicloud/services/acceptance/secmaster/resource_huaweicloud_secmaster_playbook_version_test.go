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

func getPlaybookVersionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getPlaybookVersion: Query the SecMaster playbook version detail
	var (
		getPlaybookVersionHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}/versions"
		getPlaybookVersionProduct = "secmaster"
	)
	getPlaybookVersionClient, err := cfg.NewServiceClient(getPlaybookVersionProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	getPlaybookVersionPath := getPlaybookVersionClient.Endpoint + getPlaybookVersionHttpUrl
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{project_id}", getPlaybookVersionClient.ProjectID)
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{workspace_id}", state.Primary.Attributes["workspace_id"])
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{playbook_id}", state.Primary.Attributes["playbook_id"])
	getPlaybookVersionPath += "?limit=1000"

	getPlaybookVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getPlaybookVersionResp, err := getPlaybookVersionClient.Request("GET", getPlaybookVersionPath, &getPlaybookVersionOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving PlaybookVersion: %s", err)
	}

	getPlaybookVersionRespBody, err := utils.FlattenResponse(getPlaybookVersionResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving PlaybookVersion: %s", err)
	}

	jsonPath := fmt.Sprintf("data[?id=='%s']|[0]", state.Primary.ID)
	getPlaybookVersionRespBody = utils.PathSearch(jsonPath, getPlaybookVersionRespBody, nil)
	if getPlaybookVersionRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getPlaybookVersionRespBody, nil
}

func TestAccPlaybookVersion_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_secmaster_playbook_version.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPlaybookVersionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMaster(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			// due to service problem, dataclass_id is fixed to "909494e3-558e-46b6-a9eb-07a8e18ca62f"
			// and updating test is skipped
			{
				Config: testPlaybookVersion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "playbook_id",
						"huaweicloud_secmaster_playbook.test", "id"),
					resource.TestCheckResourceAttr(rName, "dataclass_id", "909494e3-558e-46b6-a9eb-07a8e18ca62f"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "creator_id"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPlaybookVersionImportState(rName),
			},
		},
	})
}

func testPlaybookVersion_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_secmaster_playbook_version" "test" {
  workspace_id = "%s"
  playbook_id  = huaweicloud_secmaster_playbook.test.id
  dataclass_id = "909494e3-558e-46b6-a9eb-07a8e18ca62f"
  description  = "created by terraform"
}
`, testPlaybook_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}

func testPlaybookVersionImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.Attributes["playbook_id"], rs.Primary.ID), nil
	}
}
