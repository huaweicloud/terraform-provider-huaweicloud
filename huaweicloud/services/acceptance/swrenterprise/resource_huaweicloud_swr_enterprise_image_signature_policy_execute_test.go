package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSwrEnterpriseImageSignaturePolicyExecute_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseImageSignaturePolicyExecute_basic(rName),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccSwrEnterpriseImageSignaturePolicyExecute_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_enterprise_image_signature_policy_execute" "test" {
  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  namespace_name = "library"
  policy_id      = huaweicloud_swr_enterprise_image_signature_policy.test.id
}
`, testAccSwrEnterpriseImageSignaturePolicy_basic(rName))
}
