package das

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataInstanceGroups_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_instance_groups.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataInstanceGroups_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "groups.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "groups.0.id"),
					resource.TestCheckResourceAttrSet(all, "groups.0.name"),
					resource.TestCheckResourceAttrSet(all, "groups.0.description"),
				),
			},
		},
	})
}

const testAccDataInstanceGroups_basic = `
data "huaweicloud_das_instance_groups" "all" {
  datastore_type = "MySQL"
}
`
