package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnvironments_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_apig_environments.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		rName          = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnvironments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "environments.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

func testAccDataEnvironments_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_environment" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}

data "huaweicloud_apig_environments" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = huaweicloud_apig_environment.test.name
}
`, testAccEnvironment_base(rName), rName)
}
