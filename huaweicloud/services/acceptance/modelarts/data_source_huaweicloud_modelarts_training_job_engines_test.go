package modelarts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTrainingJobEngines_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelarts_training_job_engines.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTrainingJobEngines_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "engines.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(all, "engines.0.image_info.#", "1"),
					resource.TestCheckResourceAttrSet(all, "engines.0.engine_id"),
					resource.TestCheckResourceAttrSet(all, "engines.0.engine_name"),
					resource.TestCheckResourceAttrSet(all, "engines.0.engine_version"),
					resource.TestCheckResourceAttrSet(all, "engines.0.v1_compatible"),
					resource.TestCheckResourceAttrSet(all, "engines.0.image_info.0.image_version"),
					resource.TestCheckOutput("is_image_url_set_and_valid", "true"),
				),
			},
		},
	})
}

const testAccDataTrainingJobEngines_basic = `
# Without any filter parameters.
data "huaweicloud_modelarts_training_job_engines" "test" {}

locals {
  engine = try(data.huaweicloud_modelarts_training_job_engines.test.engines[0], {})
}

# At least one of gpu_image_url and cpu_image_url is not empty.
output "is_image_url_set_and_valid" {
  value = try(local.engine.image_info[0].cpu_image_url != "" || local.engine.image_info[0].gpu_image_url != "", false)
}
`
