package secmaster

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAlertResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAlert: Query the SecMaster alert detail
	var (
		getAlertHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/alerts/{id}"
		getAlertProduct = "secmaster"
	)
	getAlertClient, err := cfg.NewServiceClient(getAlertProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	getAlertPath := getAlertClient.Endpoint + getAlertHttpUrl
	getAlertPath = strings.ReplaceAll(getAlertPath, "{project_id}", getAlertClient.ProjectID)
	getAlertPath = strings.ReplaceAll(getAlertPath, "{workspace_id}", state.Primary.Attributes["workspace_id"])
	getAlertPath = strings.ReplaceAll(getAlertPath, "{id}", state.Primary.ID)

	getAlertOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getAlertResp, err := getAlertClient.Request("GET", getAlertPath, &getAlertOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getAlertResp)
}

func TestAccAlert_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_secmaster_alert.test"

	now := time.Now()
	firstTime := strings.ReplaceAll(utils.GetBeforeOrAfterDate(now, -3), "+08:00", ".000+08:00")
	lastTime := strings.ReplaceAll(utils.GetBeforeOrAfterDate(now, -2), "+08:00", ".000+08:00")
	firstTimeUpdate := strings.ReplaceAll(utils.GetBeforeOrAfterDate(now, -1), "+08:00", ".000+08:00")
	lastTimeUpdate := strings.ReplaceAll(utils.GetBeforeOrAfterDate(now, 0), "+08:00", ".000+08:00")

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAlertResourceFunc,
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
				Config: testAlert_basic(name, firstTime, lastTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "severity", "Tips"),
					resource.TestCheckResourceAttr(rName, "status", "Open"),
					resource.TestCheckResourceAttr(rName, "stage", "Preparation"),
					resource.TestCheckResourceAttr(rName, "verification_status", "Unknown"),
					resource.TestCheckResourceAttr(rName, "data_source.0.product_feature", "hss"),
					resource.TestCheckResourceAttr(rName, "data_source.0.product_name", "hss"),
					resource.TestCheckResourceAttr(rName, "data_source.0.source_type", "1"),
					resource.TestCheckResourceAttr(rName, "first_occurrence_time", firstTime),
					resource.TestCheckResourceAttr(rName, "last_occurrence_time", lastTime),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAlert_basic_update(name, firstTimeUpdate, lastTimeUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "severity", "Medium"),
					resource.TestCheckResourceAttr(rName, "first_occurrence_time", firstTimeUpdate),
					resource.TestCheckResourceAttr(rName, "last_occurrence_time", lastTimeUpdate),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAlertImportState(rName),
			},
		},
	})
}

func testAlert_basic(name, firstTime, lastTime string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_alert" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  description  = "created by terraform"

  type {
    category   = "Abnormal network behavior"
    alert_type = "Abnormal access frequency of IP address"
  }

  data_source {
    source_type     = "1"
    product_feature = "hss"
    product_name    = "hss"
  }

  first_occurrence_time = "%[3]s"
  last_occurrence_time  = "%[4]s"

  severity            = "Tips"
  status              = "Open"
  verification_status = "Unknown"
  stage               = "Preparation"

  lifecycle {
    ignore_changes = [
      name, status,
    ]
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name, firstTime, lastTime)
}

func testAlert_basic_update(name, firstTimeUpdate, lastTimeUpdate string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_alert" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  description  = ""

  type {
    category   = "Abnormal network behavior"
    alert_type = "Abnormal access frequency of IP address"
  }

  data_source {
    source_type     = "1"
    product_feature = "hss"
    product_name    = "hss"
  }

  first_occurrence_time = "%[3]s"
  last_occurrence_time  = "%[4]s"

  severity            = "Medium"
  status              = "Open"
  verification_status = "Unknown"
  stage               = "Preparation"

  lifecycle {
    ignore_changes = [
      name, status,
    ]
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name, firstTimeUpdate, lastTimeUpdate)
}

func testAlertImportState(name string) resource.ImportStateIdFunc {
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
