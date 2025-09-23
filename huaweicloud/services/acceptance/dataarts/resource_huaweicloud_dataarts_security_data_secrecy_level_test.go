package dataarts

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getSecurityDataSecrecyLevelResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPath := client.Endpoint + "v1/{project_id}/security/data-classification/secrecy-level/{id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", state.Primary.ID)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": state.Primary.Attributes["workspace_id"]},
	}
	resp, err := client.Request("GET", getPath, &opts)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", dataarts.DataSecrecyLevelResourceNotFoundCodes...)
	}

	return utils.FlattenResponse(resp)
}

func TestAccResourceSecurityDataSecrecyLevel_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dataarts_security_data_secrecy_level.test"
	rName := acceptance.RandomAccResourceName()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getSecurityDataSecrecyLevelResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityDataSecrecyLevel_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "level_number"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccSecurityDataSecrecyLevel_basic_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceSecurityDataSecrecyLevelImportStateFunc(resourceName),
			},
		},
	})
}

func testAccResourceSecurityDataSecrecyLevelImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceId := rs.Primary.Attributes["workspace_id"]
		secrecyLevelId := rs.Primary.ID
		if workspaceId == "" || secrecyLevelId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceId, secrecyLevelId)
		}
		return fmt.Sprintf("%s/%s", workspaceId, secrecyLevelId), nil
	}
}

func testAccSecurityDataSecrecyLevel_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_data_secrecy_level" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  description  = "Created by terraform"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccSecurityDataSecrecyLevel_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_data_secrecy_level" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
