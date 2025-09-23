package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccessControlAttributeConfigurations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identitycenter_access_control_attribute_configurations.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAccessControlAttributeConfigurations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "access_control_attributes.#"),
					resource.TestCheckResourceAttr(dataSource, "access_control_attributes.0.key", rName+"_1"),
					resource.TestCheckResourceAttr(dataSource, "access_control_attributes.0.value.0", "${user:email}"),
				),
			},
		},
	})
}

func testDataSourceAccessControlAttributeConfigurations_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identitycenter_instance" "test" {}

data "huaweicloud_identitycenter_access_control_attribute_configurations" "test" {
  instance_id = data.huaweicloud_identitycenter_instance.test.id
}
`, testAccessControlAttributeConfiguration_basic(name))
}
