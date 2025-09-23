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

func getIndicatorResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getIndicator: Query the SecMaster indicator detail
	var (
		getIndicatorHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/indicators/{id}"
		getIndicatorProduct = "secmaster"
	)
	getIndicatorClient, err := cfg.NewServiceClient(getIndicatorProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	getIndicatorPath := getIndicatorClient.Endpoint + getIndicatorHttpUrl
	getIndicatorPath = strings.ReplaceAll(getIndicatorPath, "{project_id}", getIndicatorClient.ProjectID)
	getIndicatorPath = strings.ReplaceAll(getIndicatorPath, "{workspace_id}", state.Primary.Attributes["workspace_id"])
	getIndicatorPath = strings.ReplaceAll(getIndicatorPath, "{id}", state.Primary.ID)

	getIndicatorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getIndicatorResp, err := getIndicatorClient.Request("GET", getIndicatorPath, &getIndicatorOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getIndicatorResp)
}

func TestAccIndicator_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	nameUpdate := name + "_update"
	rName := "huaweicloud_secmaster_indicator.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIndicatorResourceFunc,
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
				Config: testIndicator_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "threat_degree", "Black"),
					resource.TestCheckResourceAttr(rName, "status", "Open"),
					resource.TestCheckResourceAttr(rName, "confidence", "80"),
					resource.TestCheckResourceAttr(rName, "first_occurrence_time", "2023-10-24T17:23:55.000+08:00"),
					resource.TestCheckResourceAttr(rName, "last_occurrence_time", "2023-10-25T11:15:30.000+08:00"),
					resource.TestCheckResourceAttr(rName, "granularity", "1"),
					resource.TestCheckResourceAttr(rName, "value", "test.terraform.com"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testIndicator_basic_update(nameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", nameUpdate),
					resource.TestCheckResourceAttr(rName, "threat_degree", "Gray"),
					resource.TestCheckResourceAttr(rName, "status", "Closed"),
					resource.TestCheckResourceAttr(rName, "confidence", "90"),
					resource.TestCheckResourceAttr(rName, "first_occurrence_time", "2023-10-26T09:33:55.000+08:00"),
					resource.TestCheckResourceAttr(rName, "last_occurrence_time", "2023-10-27T21:15:30.000+08:00"),
					resource.TestCheckResourceAttr(rName, "granularity", "1"),
					resource.TestCheckResourceAttr(rName, "value", "1.1.1.1"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testIndicatorImportState(rName),
			},
		},
	})
}

func testIndicator_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_indicator" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  
  type {
    category       = "Domain"
    indicator_type = "Domain"
    id             = "%[3]s"
  }

  data_source {
    source_type     = "1"
    product_feature = "hss"
    product_name    = "hss"
  }

  status                = "Open"
  confidence            = "80"
  first_occurrence_time = "2023-10-24T17:23:55.000+08:00"
  last_occurrence_time  = "2023-10-25T11:15:30.000+08:00"
  threat_degree         = "Black"
  granularity           = "1"
  value                 = "test.terraform.com"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name, acceptance.HW_SECMASTER_INDICATOR_TYPE_ID)
}

func testIndicator_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_indicator" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  
  type {
    category       = "IPv4"
    indicator_type = "IPv4"
    id             = "%[3]s"
  }

  data_source {
    source_type     = "1"
    product_feature = "hss"
    product_name    = "hss"
  }

  status                = "Closed"
  confidence            = "90"
  first_occurrence_time = "2023-10-26T09:33:55.000+08:00"
  last_occurrence_time  = "2023-10-27T21:15:30.000+08:00"
  threat_degree         = "Gray"
  granularity           = "1"
  value                 = "1.1.1.1"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name, acceptance.HW_SECMASTER_INDICATOR_TYPE_ID_UPDATE)
}

func testIndicatorImportState(name string) resource.ImportStateIdFunc {
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
