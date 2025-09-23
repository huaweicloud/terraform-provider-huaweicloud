package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourcePrivateCertificateRevoke_basic(t *testing.T) {
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
				Config: testAccResourcePrivateCertificateRevoke_basic(rName),
			},
		},
	})
}

func testAccResourcePrivateCertificateRevoke_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_private_ca" "test" {
  type                = "ROOT"
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  charging_mode       = "postPaid"
  auto_renew          = false

  distinguished_name {
    common_name         = "%[1]s"
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

resource "huaweicloud_ccm_private_certificate" "test" {
  issuer_id           = huaweicloud_ccm_private_ca.test.id
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA256"
  distinguished_name {
    common_name         = "%[1]s"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }
  validity {
    type  = "DAY"
    value = "1"
  }
}
`, name)
}

func testAccResourcePrivateCertificateRevoke_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ccm_private_certificate_revoke" "test" {
  certificate_id = huaweicloud_ccm_private_certificate.test.id
  reason         = "SUPERSEDED"
}
`, testAccResourcePrivateCertificateRevoke_base(name))
}
