package fgs

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getApplicationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/fgs/applications/{id}"
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting FunctionGraph application: %s", err)
	}
	return utils.FlattenResponse(requestResp)
}

func TestAccApplication_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_fgs_application.test"
		name         = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getApplicationFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please read the instructions carefully before use to ensure sufficient permissions.
			acceptance.TestAccPreCheckFgsAgency(t)
			acceptance.TestAccPreCheckFgsTemplateId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApplication_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestMatchResourceAttr(resourceName, "stack_resources.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttr(resourceName, "repository.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"template_id",
					"agency_name",
					"params",
				},
			},
		},
	})
}

func testAccApplication_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_application" "test" {
  name        = "%[1]s"
  agency_name = "%[2]s"
  template_id = "%[3]s"
  description = "Created by terraform script"
}
`, name, acceptance.HW_FGS_AGENCY_NAME, acceptance.HW_FGS_TEMPLATE_ID)
}
