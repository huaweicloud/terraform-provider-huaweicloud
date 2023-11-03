package dsc

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

func getDscInstanceResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDscInstance: Query the DSC instance
	var (
		getDscInstanceHttpUrl = "v1/{project_id}/period/product/specification"
		getDscInstanceProduct = "dsc"
	)
	getDscInstanceClient, err := cfg.NewServiceClient(getDscInstanceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DscInstance Client: %s", err)
	}

	getDscInstancePath := getDscInstanceClient.Endpoint + getDscInstanceHttpUrl
	getDscInstancePath = strings.ReplaceAll(getDscInstancePath, "{project_id}", getDscInstanceClient.ProjectID)

	getDscInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDscInstanceResp, err := getDscInstanceClient.Request("GET", getDscInstancePath, &getDscInstanceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DscInstance: %s", err)
	}

	getDscInstanceRespBody, err := utils.FlattenResponse(getDscInstanceResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DscInstance: %s", err)
	}

	orderInfo := utils.PathSearch("orderInfo", getDscInstanceRespBody, []interface{}{})
	orders := orderInfo.([]interface{})
	if len(orders) == 0 {
		return nil, fmt.Errorf("error retrieving DscInstance: %s", err)
	}
	return orderInfo, nil
}

func TestAccDscInstance_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_dsc_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDscInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDscInstance_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "edition", "base_standard"),
					resource.TestCheckResourceAttr(rName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(rName, "period_unit", "month"),
					resource.TestCheckResourceAttr(rName, "period", "1"),
					resource.TestCheckResourceAttr(rName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(rName, "obs_expansion_package", "1"),
					resource.TestCheckResourceAttr(rName, "database_expansion_package", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"edition", "charging_mode", "auto_renew", "obs_expansion_package",
					"database_expansion_package"},
			},
		},
	})
}

func testDscInstance_basic() string {
	return `
resource "huaweicloud_dsc_instance" "test" {
  edition                    = "base_standard"
  charging_mode              = "prePaid"
  period_unit                = "month"
  period                     = 1
  auto_renew                 = false
  obs_expansion_package      = 1
  database_expansion_package = 1
}
`
}
