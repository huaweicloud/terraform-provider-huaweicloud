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

func getApplicationResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	getApplicationHttpUrl := "v1/applications/{id}"
	getApplicationPath := client.Endpoint + getApplicationHttpUrl
	getApplicationPath = strings.ReplaceAll(getApplicationPath, "{id}", state.Primary.ID)

	getApplicationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getApplicationResp, err := client.Request("GET", getApplicationPath, &getApplicationOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CMDB application: %s", err)
	}

	getApplicationRespBody, err := utils.FlattenResponse(getApplicationResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CMDB application: %s", err)
	}

	return getApplicationRespBody, nil
}

func TestAccCmdbApplication_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_cmdb_application.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getApplicationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesCmdbApplication_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acceptance"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "register_type", "API"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: tesCmdbApplication_updated(rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("%s-display", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
		},
	})
}

// enterprise_project_id is required for the testing account
func tesCmdbApplication_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_cmdb_application" "test" {
  name                  = "%s"
  description           = "created by acceptance"
  enterprise_project_id = "0"
}`, name)
}

// updating `enterprise_project_id` is optional
// if HW_ENTERPRISE_PROJECT_ID_TEST is not specified, using the default value "0"
func tesCmdbApplication_updated(name, epsID string) string {
	if epsID == "" {
		epsID = "0"
	}

	// clear the description field
	return fmt.Sprintf(`
resource "huaweicloud_aom_cmdb_application" "test" {
  name                  = "%[1]s-update"
  display_name          = "%[1]s-display"
  enterprise_project_id = "%[2]s"
}`, name, epsID)
}
