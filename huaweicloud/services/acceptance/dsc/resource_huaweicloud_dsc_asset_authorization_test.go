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

func getResourceDscAssetAuthorizationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + "v1/{project_id}/sdg/asset/authorization"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?type=%s", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DSC asset authorization: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func TestAccResourceDscAssetAuthorization_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dsc_asset_authorization.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceDscAssetAuthorizationFunc,
	)

	// Avoid CheckDestroy
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceDscAssetAuthorization_basic(true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "LTS"),
					resource.TestCheckResourceAttr(resourceName, "authorization_status", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "lts_authorization"),
				),
			},
			{
				Config: testResourceDscAssetAuthorization_basic(false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "authorization_status", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authorization_status"},
			},
		},
	})
}

func testResourceDscAssetAuthorization_basic(status bool) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_asset_authorization" "test" {
  type                 = "LTS"
  authorization_status = %[1]t
}
`, status)
}
