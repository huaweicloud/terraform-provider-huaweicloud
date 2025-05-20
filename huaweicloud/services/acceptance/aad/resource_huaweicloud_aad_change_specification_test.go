package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccChangeSpecification_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a valid AAD instance ID and config it to the environment variable.
			acceptance.TestAccPreCheckAadInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testChangeSpecification_basic(),
			},
		},
	})
}

func testChangeSpecification_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_change_specification" "test0" {
  instance_id = "%[1]s"
  upgrade_data {}
}

resource "huaweicloud_aad_change_specification" "test1" {
  depends_on = [huaweicloud_aad_change_specification.test0]

  instance_id = "%[1]s"
  upgrade_data {
    basic_bandwidth   = "10"
    elastic_bandwidth = "10"
  }
}

resource "huaweicloud_aad_change_specification" "test2" {
  depends_on = [huaweicloud_aad_change_specification.test1]

  instance_id = "%[1]s"
  upgrade_data {
    basic_bandwidth   = "20"
    elastic_bandwidth = "20"
  }
}

resource "huaweicloud_aad_change_specification" "test3" {
  depends_on = [huaweicloud_aad_change_specification.test2]

  instance_id = "%[1]s"
  upgrade_data {
    basic_bandwidth   = "10"
    elastic_bandwidth = "10"
  }
}
`, acceptance.HW_AAD_INSTANCE_ID)
}
