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

func getMaskAlgorithmResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	return dsc.GetMaskAlgorithmById(client, state.Primary.ID)
}

// Before this test, please ensure that the DSC instance has been created.
func TestAccResourceMaskAlgorithm_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dsc_mask_algorithm.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getMaskAlgorithmResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDscInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMaskAlgorithm_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "algorithm_name", name),
					resource.TestCheckResourceAttr(rName, "algorithm", "PRESNM"),
					resource.TestCheckResourceAttr(rName, "algorithm_type", "MASK_BY_OVERWRITE"),
					resource.TestCheckResourceAttr(rName, "category", "BUILT_SELF"),
					resource.TestCheckResourceAttrSet(rName, "parameter"),
					resource.TestCheckResourceAttr(rName, "data", "110108012345671"),
					resource.TestCheckResourceAttr(rName, "processed_data", "110108*****5671"),
				),
			},
			{
				Config: testAccResourceMaskAlgorithm_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "algorithm_name", updateName),
					resource.TestCheckResourceAttr(rName, "algorithm", "KEYWORD"),
					resource.TestCheckResourceAttr(rName, "algorithm_type", "MASK_BY_KEYWORDS_EXCHANGE"),
					resource.TestCheckResourceAttrSet(rName, "parameter"),
					resource.TestCheckResourceAttr(rName, "data", ""),
					resource.TestCheckResourceAttr(rName, "processed_data", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"algorithm_type",
				},
			},
		},
	})
}

func testAccResourceMaskAlgorithm_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_mask_algorithm" "test" {
  algorithm_name = "%[1]s"
  algorithm      = "PRESNM"
  algorithm_type = "MASK_BY_OVERWRITE"
  category       = "BUILT_SELF"

  parameter = jsonencode({
    type   = "CHAR"
    first  = 6
    second = 4
    method = "*"
  })

  data           = "110108012345671"
  processed_data = "110108*****5671"
}
`, name)
}

func testAccResourceMaskAlgorithm_basic_step2(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_mask_algorithm" "test" {
  algorithm_name = "%[1]s"
  algorithm      = "KEYWORD"
  algorithm_type = "MASK_BY_KEYWORDS_EXCHANGE"
  category       = "BUILT_SELF"

  parameter = jsonencode({
    key    = "keyword"
    target = "TF"
  })
}
`, updateName)
}
