package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSwrEnterpriseTemporaryCredential_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_swr_enterprise_temporary_credential.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseTemporaryCredential_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "expire_date"),
					resource.TestCheckResourceAttrSet(resourceName, "auth_token"),
				),
			},
		},
	})
}

func testAccSwrEnterpriseTemporaryCredential_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_enterprise_temporary_credential" "test" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
}
`, testAccSwrEnterpriseInstance_update(rName))
}
