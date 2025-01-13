package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataKpsKeypairs_basic(t *testing.T) {
	var (
		rName           = acceptance.RandomAccResourceName()
		publicKey, _, _ = acctest.RandSSHKeyPair("Generated-by-AccTest")

		resourceName = "data.huaweicloud_kps_keypairs.test"
		dc           = acceptance.InitDataSourceCheck(resourceName)

		byName   = "data.huaweicloud_kps_keypairs.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byPublicKey   = "data.huaweicloud_kps_keypairs.filter_by_public_key"
		dcByPublicKey = acceptance.InitDataSourceCheck(byPublicKey)

		byFingerprint   = "data.huaweicloud_kps_keypairs.filter_by_fingerprint"
		dcByFingerprint = acceptance.InitDataSourceCheck(byFingerprint)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKpsKeypairs_basic(rName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "keypairs.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "keypairs.0.public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "keypairs.0.scope"),
					resource.TestCheckResourceAttrSet(resourceName, "keypairs.0.fingerprint"),
					resource.TestCheckResourceAttrSet(resourceName, "keypairs.0.is_managed"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					dcByPublicKey.CheckResourceExists(),
					resource.TestCheckOutput("is_public_key_filter_useful", "true"),

					dcByFingerprint.CheckResourceExists(),
					resource.TestCheckOutput("is_fingerprint_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataKpsKeypairs_base(rName, key string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kps_keypair" "test" {
  name       = "%s"
  public_key = "%s"
}
`, rName, key)
}

func testAccDataKpsKeypairs_basic(rName, key string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_kps_keypairs" "test" {
  depends_on = [huaweicloud_kps_keypair.test]
}

# Filter by name
locals {
  name = data.huaweicloud_kps_keypairs.test.keypairs.0.name
}

data "huaweicloud_kps_keypairs" "filter_by_name" {
  depends_on = [huaweicloud_kps_keypair.test]

  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_kps_keypairs.filter_by_name.keypairs[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by public_key
locals {
  public_key = data.huaweicloud_kps_keypairs.test.keypairs.0.public_key
}

data "huaweicloud_kps_keypairs" "filter_by_public_key" {
  depends_on = [huaweicloud_kps_keypair.test]

  public_key = local.public_key
}

locals {
  public_key_filter_result = [
    for v in data.huaweicloud_kps_keypairs.filter_by_public_key.keypairs[*].public_key : v == local.public_key
  ]
}

output "is_public_key_filter_useful" {
  value = length(local.public_key_filter_result) > 0 && alltrue(local.public_key_filter_result)
}

# Filter by fingerprint
locals {
  fingerprint = data.huaweicloud_kps_keypairs.test.keypairs.0.fingerprint
}

data "huaweicloud_kps_keypairs" "filter_by_fingerprint" {
  depends_on = [huaweicloud_kps_keypair.test]

  fingerprint = local.fingerprint
}

locals {
  fingerprint_filter_result = [
    for v in data.huaweicloud_kps_keypairs.filter_by_fingerprint.keypairs[*].fingerprint : v == local.fingerprint
  ]
}

output "is_fingerprint_filter_useful" {
  value = length(local.fingerprint_filter_result) > 0 && alltrue(local.fingerprint_filter_result)
}
`, testAccDataKpsKeypairs_base(rName, key))
}
