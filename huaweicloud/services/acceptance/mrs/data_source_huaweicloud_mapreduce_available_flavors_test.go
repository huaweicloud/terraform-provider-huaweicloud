package mrs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAvailableFlavors_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_mapreduce_available_flavors.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataAvailableFlavors_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "version_name"),
					resource.TestMatchResourceAttr(all, "available_flavors.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "available_flavors.0.az_code"),
					resource.TestCheckResourceAttrSet(all, "available_flavors.0.az_name"),
					resource.TestMatchResourceAttr(all, "available_flavors.0.master.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "available_flavors.0.master.0.flavor_name"),
					resource.TestMatchResourceAttr(all, "available_flavors.0.core.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "available_flavors.0.core.0.flavor_name"),
					resource.TestMatchResourceAttr(all, "available_flavors.0.task.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "available_flavors.0.task.0.flavor_name"),
				),
			},
		},
	})
}

const testDataAvailableFlavors_basic = `
data "huaweicloud_mapreduce_versions" "test" {}

data "huaweicloud_mapreduce_available_flavors" "test" {
  version_name = try(data.huaweicloud_mapreduce_versions.test.versions[0], null)
}
`
