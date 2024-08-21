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

func getDashboardResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	listHttpUrl := "v2/{project_id}/aom/dashboards"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeaders(state),
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dashboards: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening dashboards: %s", err)
	}

	jsonPath := fmt.Sprintf("dashboards[?dashboard_id=='%s']|[0]", state.Primary.ID)
	dashboard := utils.PathSearch(jsonPath, listRespBody, nil)
	if dashboard == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return dashboard, nil
}

func TestAccDashboard_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	newName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_dashboard.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDashboardResourceFunc,
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
				Config: testDashboard_basic(rName, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "dashboard_title", rName),
					resource.TestCheckResourceAttrPair(resourceName, "folder_title", "huaweicloud_aom_dashboards_folder.test", "folder_title"),
					resource.TestCheckResourceAttr(resourceName, "dashboard_type", rName),
					resource.TestCheckResourceAttr(resourceName, "is_favorite", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "enterprise_project_id",
						"huaweicloud_aom_dashboards_folder.test", "enterprise_project_id"),
					resource.TestCheckResourceAttr(resourceName, "dashboard_tags.0.key", rName),
				),
			},
			{
				Config: testDashboard_basic(newName, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "dashboard_title", newName),
					resource.TestCheckResourceAttrPair(resourceName, "folder_title", "huaweicloud_aom_dashboards_folder.test", "folder_title"),
					resource.TestCheckResourceAttr(resourceName, "dashboard_type", newName),
					resource.TestCheckResourceAttr(resourceName, "is_favorite", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "enterprise_project_id",
						"huaweicloud_aom_dashboards_folder.test", "enterprise_project_id"),
					resource.TestCheckResourceAttr(resourceName, "dashboard_tags.0.key", newName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testDashboard_basic(name string, isFavorite bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aom_dashboard" "test" {
  depends_on = [huaweicloud_aom_dashboards_folder.test]

  dashboard_title       = "%[2]s"
  folder_title          = huaweicloud_aom_dashboards_folder.test.folder_title
  dashboard_type        = "%[2]s"
  is_favorite           = %[3]t
  enterprise_project_id = huaweicloud_aom_dashboards_folder.test.enterprise_project_id
  dashboard_tags        = [
    {
      key = "%[2]s"
    }
  ]
}`, testDashboardsFolder_basic(name, false), name, isFavorite)
}
