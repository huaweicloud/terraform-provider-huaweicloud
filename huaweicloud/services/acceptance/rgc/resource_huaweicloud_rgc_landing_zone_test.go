package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getLandingZoneResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	// getLandingZone: Query RGC landing zone via rgc API
	var (
		region                = acceptance.HW_REGION_NAME
		getLandingZoneHttpUrl = "v1/landing-zone/status"
		getLandingZoneProduct = "rgc"
	)
	getLandingZoneClient, err := cfg.NewServiceClient(getLandingZoneProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RGC client: %s", err)
	}

	getLandingZonePath := getLandingZoneClient.Endpoint + getLandingZoneHttpUrl

	getLandingZoneOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getLandingZoneResp, err := getLandingZoneClient.Request("GET", getLandingZonePath, &getLandingZoneOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving landing zone: %s", err)
	}

	getLandingZoneStatusRespBody, err := utils.FlattenResponse(getLandingZoneResp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("landing_zone_status", getLandingZoneStatusRespBody, "").(string)
	actionType := utils.PathSearch("landing_zone_action_type", getLandingZoneStatusRespBody, "").(string)

	if status != "succeeded" {
		if actionType == "CREATE" {
			return getLandingZoneStatusRespBody, nil
		}

		message := utils.PathSearch("message", getLandingZoneStatusRespBody, "")
		return nil, fmt.Errorf("status: %s; message: %s", status, message)
	}

	if actionType == "DELETE" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getLandingZoneStatusRespBody, nil
}

func TestAccLandingZone_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_rgc_landing_zone.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLandingZoneResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAKAndSK(t)
			acceptance.TestAccPreCheckRGCLandingZone(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLandingZone_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "landing_zone_status"),
					resource.TestCheckResourceAttrSet(rName, "deployed_version"),
				),
			},
		},
	})
}

func testLandingZone_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rgc_landing_zone" "test" {
  home_region            = "%[1]s"
  identity_center_status = "ENABLE"
  identity_store_email   = "%[6]s"
  cloud_trail_type       = true
  organization_structure_type = "NON_STANDARD"
  region_configuration_list {
    region                      = "%[1]s"
    region_configuration_status = "ENABLED"
  }
  organization_structure {
    organizational_unit_type = "CORE"
    accounts {
      account_name = "%[2]s"
      account_type = "LOGGING"
      account_id   = "%[3]s"
    }
    accounts {
      account_name  = "%[4]s"
      account_type  = "SECURITY"
      account_id    = "%[5]s"
      account_email = "%[4]s@huawei.com"
    }
  }
  logging_configuration {
    logging_bucket {
      retention_days = 365
    }
    access_logging_bucket {
      retention_days = 3650
    }
  }
}
`, acceptance.HW_REGION_NAME, acceptance.HW_RGC_LOGGING_ACCOUNT_NAME, acceptance.HW_RGC_LOGGING_ACCOUNT_ID,
		acceptance.HW_RGC_AUDIT_ACCOUNT_NAME, acceptance.HW_RGC_AUDIT_ACCOUNT_ID, acceptance.HW_RGC_MANAGE_ACCOUNT_EMAIL)
}
