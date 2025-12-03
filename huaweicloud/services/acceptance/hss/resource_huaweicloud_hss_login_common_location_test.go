package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
)

func getLoginCommonLocationResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	var (
		region       = acceptance.HW_REGION_NAME
		epsId        = "all_granted_eps"
		testAreaCode = 110109
		product      = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	return hss.QueryLoginCommonLocation(client, testAreaCode, epsId)
}

func TestAccLoginCommonLocation_basic(t *testing.T) {
	var (
		policy interface{}
		rName  = "huaweicloud_hss_login_common_location.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&policy,
		getLoginCommonLocationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLoginCommonLocation_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "host_id_list.0", acceptance.HW_HSS_HOST_ID1),
				),
			},
			{
				Config: testAccLoginCommonLocation_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "host_id_list.0", acceptance.HW_HSS_HOST_ID2),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateIdFunc: testAccLoginCommonLocationImportStateIDFunc(rName),
			},
		},
	})
}

func testAccLoginCommonLocation_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_login_common_location" "test" {
  area_code             = 110109
  host_id_list          = ["%s"]
  enterprise_project_id = "all_granted_eps"
}
`, acceptance.HW_HSS_HOST_ID1)
}

func testAccLoginCommonLocation_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_login_common_location" "test" {
  area_code             = 110109
  host_id_list          = ["%s"]
  enterprise_project_id = "all_granted_eps"
}
`, acceptance.HW_HSS_HOST_ID2)
}

func testAccLoginCommonLocationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", resourceName)
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		areaCode := rs.Primary.Attributes["area_code"]

		if epsId == "" || areaCode == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<enterprise_project_id>/<area_code>', but got '%s/%s'", epsId, areaCode)
		}
		return fmt.Sprintf("%s/%s", epsId, areaCode), nil
	}
}
