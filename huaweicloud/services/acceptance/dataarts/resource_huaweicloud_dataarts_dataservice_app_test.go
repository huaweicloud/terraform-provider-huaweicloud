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

func getDataServiceAppResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getApp: Query the app
	var (
		getAppHttpUrl = "v1/{project_id}/service/apps/{id}"
		getAppProduct = "dataarts"
	)
	getAppClient, err := cfg.NewServiceClient(getAppProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	getAppPath := getAppClient.Endpoint + getAppHttpUrl
	getAppPath = strings.ReplaceAll(getAppPath, "{project_id}", getAppClient.ProjectID)
	getAppPath = strings.ReplaceAll(getAppPath, "{id}", state.Primary.ID)

	getAppOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    state.Primary.Attributes["workspace_id"],
			"dlm-type":     state.Primary.Attributes["dlm_type"],
		},
	}

	getAppResp, err := getAppClient.Request("GET", getAppPath, &getAppOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving app: %s", err)
	}

	getAppRespBody, err := utils.FlattenResponse(getAppResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving app: %s", err)
	}

	return getAppRespBody, nil
}

func TestAccDataServiceApp_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dataarts_dataservice_app.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataServiceAppResourceFunc,
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
				Config: testDataServiceApp_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "dlm_type", "SHARED"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "app_type", "APP"),
					resource.TestCheckResourceAttrSet(rName, "description"),
					resource.TestCheckResourceAttrSet(rName, "app_key"),
					resource.TestCheckResourceAttrSet(rName, "app_secret"),
				),
			},
			{
				Config: testDataServiceApp_basic_update(name + "update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "dlm_type", "SHARED"),
					resource.TestCheckResourceAttr(rName, "name", name+"update"),
					resource.TestCheckResourceAttr(rName, "app_type", "APP"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrSet(rName, "app_key"),
					resource.TestCheckResourceAttrSet(rName, "app_secret"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataServiceAppImportState(rName),
			},
		},
	})
}

func TestAccDataServiceApp_iam(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_dataarts_dataservice_app.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataServiceAppResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainName(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataServiceApp_iam(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "dlm_type", "EXCLUSIVE"),
					resource.TestCheckResourceAttr(rName, "app_type", "IAM"),
					resource.TestCheckResourceAttr(rName, "app_key", "NO DATA"),
					resource.TestCheckResourceAttr(rName, "app_secret", "NO DATA"),
					resource.TestCheckResourceAttrSet(rName, "description"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataServiceAppImportState(rName),
			},
		},
	})
}

func testDataServiceApp_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_dataservice_app" "test" {
  workspace_id = "%[1]s"
  dlm_type     = "SHARED"
  name         = "%[2]s"
  description  = "created by acceptance"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testDataServiceApp_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_dataservice_app" "test" {
  workspace_id = "%[1]s"
  dlm_type     = "SHARED"
  name         = "%[2]s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testDataServiceApp_iam() string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_dataservice_app" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  dlm_type     = "EXCLUSIVE"
  app_type     = "IAM"
  description  = "IAM authentication with EXCLUSIVE DLM engine"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DOMAIN_NAME)
}

func testDataServiceAppImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["dlm_type"] == "" {
			return "", fmt.Errorf("attribute (dlm_type) of resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["workspace_id"] + "/" +
			rs.Primary.Attributes["dlm_type"] + "/" +
			rs.Primary.ID, nil
	}
}
