package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dsc"
)

func getResourceMeasureInfoFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	lifeCycle := state.Primary.Attributes["life_cycle"]
	return dsc.GetMeasureInfoById(client, state.Primary.ID, lifeCycle)
}

func TestAccResourceMeasureInfo_basic(t *testing.T) {
	var (
		rName   = "huaweicloud_dsc_measure_info.test"
		measure = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceMeasureInfoFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMeasureInfo_basic(measure),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", measure),
					resource.TestCheckResourceAttr(rName, "life_cycle", "COLLECTION"),
					resource.TestCheckResourceAttr(rName, "description", "test_description"),
				),
			},
			{
				Config: testAccResourceMeasureInfo_update(measure),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", measure)),
					resource.TestCheckResourceAttr(rName, "description", "test_description_update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateIdFunc: testAccResourceMeasureInfoImportStateIdFunc(rName),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceMeasureInfo_basic(measure string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_measure_info" "test" {
  name        = "%s"
  life_cycle  = "COLLECTION"
  description = "test_description"
  measure_type = "CUSTOM"
}
`, measure)
}

func testAccResourceMeasureInfo_update(measure string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_measure_info" "test" {
  name        = "%s_update"
  life_cycle  = "COLLECTION"
  description = "test_description_update"
  measure_type = "CUSTOM"
}
`, measure)
}

func testAccResourceMeasureInfoImportStateIdFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var measureId, lifeCycle string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		measureId = rs.Primary.ID
		lifeCycle = rs.Primary.Attributes["life_cycle"]

		if measureId == "" || lifeCycle == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<id>/<life_cycle>', but got '%s/%s'",
				measureId, lifeCycle)
		}

		return fmt.Sprintf("%s/%s", measureId, lifeCycle), nil
	}
}
