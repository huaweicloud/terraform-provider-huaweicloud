package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceSupportedFeatures_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_apig_instance_supported_features.test"
		dc    = acceptance.InitDataSourceCheck(rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceSupportedFeatures_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "features.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

func testAccDataSourceInstanceSupportedFeatures_basic() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_instance_supported_features" "test" {
  instance_id = huaweicloud_apig_instance.test.id
}
`, testAccInstanceFeature_base(name))
}
