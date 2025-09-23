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
	var (
		obj interface{}

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		sharedForAPP   = "huaweicloud_dataarts_dataservice_app.shared_type_app"
		rcSharedForAPP = acceptance.InitResourceCheck(sharedForAPP, &obj, getDataServiceAppResourceFunc)

		exclusiveForAPP   = "huaweicloud_dataarts_dataservice_app.exclusive_type_app"
		rcExclusiveForAPP = acceptance.InitResourceCheck(exclusiveForAPP, &obj, getDataServiceAppResourceFunc)

		exclusiveForDLM   = "huaweicloud_dataarts_dataservice_app.exclusive_type_dlm"
		rcExclusiveForDLM = acceptance.InitResourceCheck(exclusiveForDLM, &obj, getDataServiceAppResourceFunc)

		exclusiveForIAM   = "huaweicloud_dataarts_dataservice_app.exclusive_type_iam"
		rcExclusiveForIAM = acceptance.InitResourceCheck(exclusiveForIAM, &obj, getDataServiceAppResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainName(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcSharedForAPP.CheckResourceDestroy(),
			rcExclusiveForAPP.CheckResourceDestroy(),
			rcExclusiveForDLM.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testDataServiceApp_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcSharedForAPP.CheckResourceExists(),
					resource.TestCheckResourceAttr(sharedForAPP, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(sharedForAPP, "dlm_type", "SHARED"),
					resource.TestCheckResourceAttr(sharedForAPP, "name", name+"_shared_type_app"),
					resource.TestCheckResourceAttrSet(sharedForAPP, "description"),
					resource.TestCheckResourceAttrSet(sharedForAPP, "app_key"),
					resource.TestCheckResourceAttrSet(sharedForAPP, "app_secret"),
					rcExclusiveForAPP.CheckResourceExists(),
					resource.TestCheckResourceAttr(exclusiveForAPP, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(exclusiveForAPP, "dlm_type", "EXCLUSIVE"),
					resource.TestCheckResourceAttr(exclusiveForAPP, "name", name+"_exclusive_type_app"),
					resource.TestCheckResourceAttr(exclusiveForAPP, "app_type", "APP"),
					resource.TestCheckResourceAttrSet(exclusiveForAPP, "description"),
					resource.TestCheckResourceAttrSet(exclusiveForAPP, "app_key"),
					resource.TestCheckResourceAttrSet(exclusiveForAPP, "app_secret"),
					rcExclusiveForDLM.CheckResourceExists(),
					resource.TestCheckResourceAttr(exclusiveForDLM, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(exclusiveForDLM, "dlm_type", "EXCLUSIVE"),
					resource.TestCheckResourceAttr(exclusiveForDLM, "name", name+"_exclusive_type_dlm"),
					resource.TestCheckResourceAttr(exclusiveForDLM, "app_type", "DLM"),
					resource.TestCheckResourceAttrSet(exclusiveForDLM, "description"),
					resource.TestCheckResourceAttrSet(exclusiveForDLM, "app_key"),
					resource.TestCheckResourceAttrSet(exclusiveForDLM, "app_secret"),
					rcExclusiveForIAM.CheckResourceExists(),
					resource.TestCheckResourceAttr(exclusiveForIAM, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(exclusiveForIAM, "dlm_type", "EXCLUSIVE"),
					resource.TestCheckResourceAttr(exclusiveForIAM, "name", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(exclusiveForIAM, "app_type", "IAM"),
					resource.TestCheckResourceAttrSet(exclusiveForIAM, "description"),
					resource.TestCheckResourceAttr(exclusiveForIAM, "app_key", "NO DATA"),
					resource.TestCheckResourceAttr(exclusiveForIAM, "app_secret", "NO DATA"),
				),
			},
			{
				ResourceName:      sharedForAPP,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataServiceAppImportState(sharedForAPP),
				ImportStateVerifyIgnore: []string{
					"app_type",
				},
			},
			{
				ResourceName:      exclusiveForAPP,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataServiceAppImportState(exclusiveForAPP),
				ImportStateVerifyIgnore: []string{
					"app_type",
				},
			},
			{
				ResourceName:      exclusiveForDLM,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataServiceAppImportState(exclusiveForDLM),
				ImportStateVerifyIgnore: []string{
					"app_type",
				},
			},
			{
				ResourceName:      exclusiveForIAM,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataServiceAppImportState(exclusiveForIAM),
				ImportStateVerifyIgnore: []string{
					"app_type",
				},
			},
			{
				Config: testDataServiceApp_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rcSharedForAPP.CheckResourceExists(),
					resource.TestCheckResourceAttr(sharedForAPP, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(sharedForAPP, "dlm_type", "SHARED"),
					resource.TestCheckResourceAttr(sharedForAPP, "name", updateName+"_shared_type_app"),
					resource.TestCheckResourceAttr(sharedForAPP, "description", ""),
					resource.TestCheckResourceAttrSet(sharedForAPP, "app_key"),
					resource.TestCheckResourceAttrSet(sharedForAPP, "app_secret"),
					rcExclusiveForAPP.CheckResourceExists(),
					resource.TestCheckResourceAttr(exclusiveForAPP, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(exclusiveForAPP, "dlm_type", "EXCLUSIVE"),
					resource.TestCheckResourceAttr(exclusiveForAPP, "name", updateName+"_exclusive_type_app"),
					resource.TestCheckResourceAttr(exclusiveForAPP, "app_type", "APP"),
					resource.TestCheckResourceAttr(exclusiveForAPP, "description", ""),
					resource.TestCheckResourceAttrSet(exclusiveForAPP, "app_key"),
					resource.TestCheckResourceAttrSet(exclusiveForAPP, "app_secret"),
					rcExclusiveForDLM.CheckResourceExists(),
					resource.TestCheckResourceAttr(exclusiveForDLM, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(exclusiveForDLM, "dlm_type", "EXCLUSIVE"),
					resource.TestCheckResourceAttr(exclusiveForDLM, "name", updateName+"_exclusive_type_dlm"),
					resource.TestCheckResourceAttr(exclusiveForDLM, "app_type", "DLM"),
					resource.TestCheckResourceAttr(exclusiveForDLM, "description", ""),
					resource.TestCheckResourceAttrSet(exclusiveForDLM, "app_key"),
					resource.TestCheckResourceAttrSet(exclusiveForDLM, "app_secret"),
					rcExclusiveForIAM.CheckResourceExists(),
					resource.TestCheckResourceAttr(exclusiveForIAM, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(exclusiveForIAM, "dlm_type", "EXCLUSIVE"),
					resource.TestCheckResourceAttr(exclusiveForIAM, "name", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(exclusiveForIAM, "app_type", "IAM"),
					resource.TestCheckResourceAttr(exclusiveForIAM, "description", ""),
					resource.TestCheckResourceAttr(exclusiveForIAM, "app_key", "NO DATA"),
					resource.TestCheckResourceAttr(exclusiveForIAM, "app_secret", "NO DATA"),
				),
			},
		},
	})
}

// The enum 'APP' includes these type:
// + APIG
// + APIGW
// + DLM (Exclusive application)
// + ROMA_APIC
func testDataServiceApp_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_dataservice_app" "shared_type_app" {
  workspace_id = "%[1]s"
  dlm_type     = "SHARED"
  name         = "%[2]s_shared_type_app"
  description  = "created by acceptance"
}

resource "huaweicloud_dataarts_dataservice_app" "exclusive_type_app" {
  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  app_type     = "APP"
  name         = "%[2]s_exclusive_type_app"
  description  = "created by acceptance"
}

resource "huaweicloud_dataarts_dataservice_app" "exclusive_type_dlm" {
  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  app_type     = "DLM"
  name         = "%[2]s_exclusive_type_dlm"
  description  = "created by acceptance"
}

resource "huaweicloud_dataarts_dataservice_app" "exclusive_type_iam" {
  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  app_type     = "IAM"
  name         = "%[3]s"
  description  = "created by acceptance"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DOMAIN_NAME)
}

func testDataServiceApp_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_dataservice_app" "shared_type_app" {
  workspace_id = "%[1]s"
  dlm_type     = "SHARED"
  name         = "%[2]s_shared_type_app"
}

resource "huaweicloud_dataarts_dataservice_app" "exclusive_type_app" {
  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  app_type     = "APP" # Required if the value of the DLM type is exclusive.
  name         = "%[2]s_exclusive_type_app"
}

resource "huaweicloud_dataarts_dataservice_app" "exclusive_type_dlm" {
  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  app_type     = "DLM"
  name         = "%[2]s_exclusive_type_dlm"
}

resource "huaweicloud_dataarts_dataservice_app" "exclusive_type_iam" {
  workspace_id = "%[1]s"
  dlm_type     = "EXCLUSIVE"
  app_type     = "IAM"
  name         = "%[3]s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DOMAIN_NAME)
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
