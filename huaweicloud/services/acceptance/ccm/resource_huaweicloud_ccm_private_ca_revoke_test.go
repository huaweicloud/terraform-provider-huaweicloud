package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourcePrivateCaRevoke_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePrivateCaRevoke_basic(rName),
			},
		},
	})
}

func testAccResourcePrivateCaRevoke_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_private_ca" "root" {
  type                = "ROOT"
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  charging_mode       = "postPaid"
  auto_renew          = false

  distinguished_name {
    common_name         = "%[1]s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "DAY"
    value = 5
  }
}

resource "huaweicloud_ccm_private_ca" "test" {
  type                = "SUBORDINATE"
  issuer_id           = huaweicloud_ccm_private_ca.root.id
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  charging_mode       = "postPaid"
  auto_renew          = false

  distinguished_name {
    common_name         = "%[1]s-sub"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "DAY"
    value = 1
  }
}
`, name)
}

func testAccResourcePrivateCaRevoke_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ccm_private_ca_revoke" "test" {
  ca_id  = huaweicloud_ccm_private_ca.test.id
  reason = "SUPERSEDED"
}
`, testAccResourcePrivateCaRevoke_base(name))
}
