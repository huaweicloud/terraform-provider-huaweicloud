package aom

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

func getDashboardsFolderResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	listHttpUrl := "v2/{project_id}/aom/dashboards-folder"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeaders(state),
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dashboards folder: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening dashboards folder: %s", err)
	}

	jsonPath := fmt.Sprintf("[?folder_id=='%s']|[0]", state.Primary.ID)
	folder := utils.PathSearch(jsonPath, listRespBody, nil)
	if folder == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return folder, nil
}

func TestAccDashboardsFolder_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	newName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_dashboards_folder.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDashboardsFolderResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDashboardsFolder_basic(rName, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "folder_title", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "delete_all", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "is_template"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
			{
				Config: testDashboardsFolder_basic(newName, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "folder_title", newName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "delete_all", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "is_template"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_all"},
			},
		},
	})
}

func testDashboardsFolder_basic(name string, deleteAll bool) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_dashboards_folder" "test" {
  folder_title          = "%[1]s"
  enterprise_project_id = "%[2]s"
  delete_all            = %[3]t
}`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, deleteAll)
}
