package dataarts

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

func getScriptResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getScriptHttpUrl = "v1/{project_id}/scripts/{script_name}"
		getScriptProduct = "dataarts-dlf"
	)
	client, err := cfg.NewServiceClient(getScriptProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts client: %s", err)
	}

	split := strings.Split(state.Primary.ID, "/")
	if len(split) != 2 {
		return nil, fmt.Errorf("error resolving the id: %s", state.Primary.ID)
	}
	scriptName := split[1]

	getPath := client.Endpoint + getScriptHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{script_name}", scriptName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": state.Primary.Attributes["workspace_id"]},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func TestAccResourceScript_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dataarts_factory_script.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getScriptResourceFunc,
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccScript_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "DLISQL"),
					resource.TestCheckResourceAttr(resourceName, "workspace_id",
						acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(resourceName, "content", "#test"),
					resource.TestCheckResourceAttr(resourceName, "connection_name",
						acceptance.HW_DATAARTS_CONNECTION_NAME),
					resource.TestCheckResourceAttr(resourceName, "directory", "/basic"),
					resource.TestCheckResourceAttr(resourceName, "queue_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "configuration.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
			{
				Config: testAccScript_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "DLISQL"),
					resource.TestCheckResourceAttr(resourceName, "workspace_id",
						acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(resourceName, "content", "#test_update"),
					resource.TestCheckResourceAttr(resourceName, "connection_name",
						acceptance.HW_DATAARTS_CONNECTION_NAME),
					resource.TestCheckResourceAttr(resourceName, "directory", "/update"),
					resource.TestCheckResourceAttr(resourceName, "queue_name", "update"),
					resource.TestCheckResourceAttr(resourceName, "description", "test_update"),
					resource.TestCheckResourceAttr(resourceName, "configuration.%", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccResourceScriptImportStateIDFunc(resourceName),
				ImportStateVerifyIgnore: []string{"approvers", "target_status"},
			},
		},
	})
}

func testAccScript_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_factory_script" "test" {
  workspace_id    = "%s"
  name            = "%s"
  type            = "DLISQL"
  content         = "#test"
  connection_name = "%s"
  directory       = "/basic"
  queue_name      = "default"
  description     = "test"
  configuration   = {
    "spark.sql.files.maxRecordsPerFile" = "1"
  }
}`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_CONNECTION_NAME)
}

func testAccScript_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_factory_script" "test" {
  workspace_id    = "%s"
  name            = "%s"
  type            = "DLISQL"
  content         = "#test_update"
  connection_name = "%s"
  directory       = "/update"
  queue_name      = "update"
  description     = "test_update"
  configuration   = {
    "spark.sql.files.maxRecordsPerFile" = "1"
    "dli.sql.job.timeout"               = "1"
  }
}`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_CONNECTION_NAME)
}

func testAccResourceScriptImportStateIDFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		workspaceID := rs.Primary.Attributes["workspace_id"]
		scriptName := rs.Primary.Attributes["name"]
		if workspaceID == "" || scriptName == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<workspace_id>/<name>', but got '%s/%s'",
				workspaceID, scriptName)
		}
		return fmt.Sprintf("%s/%s", workspaceID, scriptName), nil
	}
}
