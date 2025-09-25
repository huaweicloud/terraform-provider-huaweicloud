package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseImageSignaturePolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_image_signature_policies.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseImageSignaturePolicies_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.signature_method"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.signature_key"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.namespace_id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.namespace_name"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseImageSignaturePolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_image_signature_policies" "test" {
  depends_on = [huaweicloud_swr_enterprise_image_signature_policy.test]

  instance_id = huaweicloud_swr_enterprise_instance.test.id
}
`, testAccSwrEnterpriseImageSignaturePolicy_basic(name))
}
