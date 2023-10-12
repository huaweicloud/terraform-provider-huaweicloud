package cmdb

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
	client, err := conf.NewServiceClient("cmdb", acceptance.HW_REGION_NAME)
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
		return nil, fmt.Errorf("error retrieving Application: %s", err)
	}

	getApplicationRespBody, err := utils.FlattenResponse(getApplicationResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Application: %s", err)
	}

	return getApplicationRespBody, nil
}

func TestAccApplication_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_application.test"

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
				Config: tesAomApplication_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "display_name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acceptance"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "register_type", "API"),
					resource.TestCheckResourceAttrSet(resourceName, "creator"),
					resource.TestCheckResourceAttrSet(resourceName, "modifier"),
				),
			},
			{
				Config: tesAomApplication_updated(rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("%s-display", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by acceptance"),
					resource.TestCheckResourceAttrSet(resourceName, "creator"),
					resource.TestCheckResourceAttrSet(resourceName, "modifier"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// enterprise_project_id is required for the testing account
func tesAomApplication_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_application" "test" {
  name                  = "%s"
  description           = "created by acceptance"
  enterprise_project_id = "0"
}`, name)
}

// updating `enterprise_project_id` is optional
// if HW_ENTERPRISE_PROJECT_ID_TEST is not specified, using the default value "0"
func tesAomApplication_updated(name, epsID string) string {
	if epsID == "" {
		epsID = "0"
	}

	return fmt.Sprintf(`
resource "huaweicloud_aom_application" "test" {
  name                  = "%[1]s-update"
  display_name          = "%[1]s-display"
  description           = "updated by acceptance"
  enterprise_project_id = "%[2]s"
}`, name, epsID)
}
