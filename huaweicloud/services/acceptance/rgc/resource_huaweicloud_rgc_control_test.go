package rgc

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

func getControlResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getControl: Query RGC control via rgc API
	var (
		region            = acceptance.HW_REGION_NAME
		getControlHttpUrl = "v1/governance/managed-organizational-units/{managed_organizational_unit_id}/controls/{control_id}"
		getControlProduct = "rgc"
	)

	getControlClient, err := cfg.NewServiceClient(getControlProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RGC client: %s", err)
	}

	getControlPath := getControlClient.Endpoint + getControlHttpUrl
	getControlPath = strings.ReplaceAll(getControlPath, "{managed_organizational_unit_id}", acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID)
	getControlPath = strings.ReplaceAll(getControlPath, "{control_id}", acceptance.HW_RGC_CONTROL_IDENTIFIER)

	getControlOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getControlResp, err := getControlClient.Request("GET", getControlPath, &getControlOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Control: %s", err)
	}

	return utils.FlattenResponse(getControlResp)
}

func TestAccControl_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_rgc_control.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getControlResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCControl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testControl_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "state"),
					resource.TestCheckResourceAttrSet(rName, "version"),
				),
			},
		},
	})
}

func testControl_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rgc_control" "test" {
  identifier        = "%[1]s"
  target_identifier = "%[2]s"
}
`, acceptance.HW_RGC_CONTROL_IDENTIFIER, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID)
}
