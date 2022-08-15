package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataKpsKeypairs_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "data.huaweicloud_kps_keypairs.test"
	publicKey, _, _ := acctest.RandSSHKeyPair("Generated-by-AccTest")

	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKpsKeypairs_basic(rName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "keypairs.0.name", rName),
					resource.TestCheckResourceAttr(resourceName, "keypairs.0.public_key", publicKey),
					resource.TestCheckResourceAttrPair(resourceName, "keypairs.0.scope",
						"huaweicloud_kps_keypair.test", "scope"),
					resource.TestCheckResourceAttrPair(resourceName, "keypairs.0.fingerprint",
						"huaweicloud_kps_keypair.test", "fingerprint"),
					resource.TestCheckResourceAttrPair(resourceName, "keypairs.0.is_managed",
						"huaweicloud_kps_keypair.test", "is_managed"),
				),
			},
		},
	})
}

func testAccDataKpsKeypairs_basic(rName, key string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_kps_keypairs" "test" {
  name = huaweicloud_kps_keypair.test.name

  depends_on = [huaweicloud_kps_keypair.test]
}
`, testKeypair_publicKey(rName, key))
}
