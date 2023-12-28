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

func getArchitectureReviewerResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/design/approvals/users?approver_name={user_name}"
		product = "dataarts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{user_name}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": state.Primary.Attributes["workspace_id"]},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DataArts Architecture reviewer: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	reviewer := utils.PathSearch("data.value.records|[0]", getRespBody, nil)
	if reviewer == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return reviewer, nil
}

func TestAccResourceArchitectureReviewer_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dataarts_architecture_reviewer.test"
	)
	workspaceId := acceptance.HW_DATAARTS_WORKSPACE_ID
	userId := acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID
	userName := acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getArchitectureReviewerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsArchitectureReviewer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccArchitectureReviewer_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "workspace_id", workspaceId),
					resource.TestCheckResourceAttr(resourceName, "user_id", userId),
					resource.TestCheckResourceAttr(resourceName, "user_name", userName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"email", "phone_number",
				},
				ImportStateIdFunc: testAccResourceArchitectureReviewerImportFunc(resourceName),
			},
		},
	})
}

func testAccArchitectureReviewer_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_reviewer" "test" {
  workspace_id = "%s"
  user_id      = "%s"
  user_name    = "%s"
  email        = "test@example.com"
  phone_number = "12345678901"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

func testAccResourceArchitectureReviewerImportFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceID := rs.Primary.Attributes["workspace_id"]
		userName := rs.Primary.ID
		if workspaceID == "" || userName == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceID, userName)
		}
		return fmt.Sprintf("%s/%s", workspaceID, userName), nil
	}
}
