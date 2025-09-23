package cbh

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

func getAssetAgencyAuthorizationResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "cbh"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CBH client: %s", err)
	}

	basePath := client.Endpoint + "v2/{project_id}/cbs/agency/authorization"
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	baseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", basePath, &baseOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CBH asset agency authorization: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccAssetAgencyAuthorization_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_cbh_asset_agency_authorization.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAssetAgencyAuthorizationResourceFunc,
	)

	// Avoid CheckDestroy, because there is nothing in the resource destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssetAgencyAuthorization_config(true, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "csms", "true"),
					resource.TestCheckResourceAttr(rName, "kms", "true"),
				),
			},
			{
				Config: testAccAssetAgencyAuthorization_config(false, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "csms", "false"),
					resource.TestCheckResourceAttr(rName, "kms", "true")),
			},
			{
				Config: testAccAssetAgencyAuthorization_config(true, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "csms", "true"),
					resource.TestCheckResourceAttr(rName, "kms", "false")),
			},
		},
	})
}

func testAccAssetAgencyAuthorization_config(isCsmsEnabled, isKmsEnabled bool) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbh_asset_agency_authorization" "test" {
  csms = "%[1]v"
  kms  = "%[2]v"
}
`, isCsmsEnabled, isKmsEnabled)
}
