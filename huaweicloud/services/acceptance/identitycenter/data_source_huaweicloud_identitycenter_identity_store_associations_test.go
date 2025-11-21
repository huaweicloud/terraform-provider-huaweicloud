package identitycenter

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIdentityCenterIdentityStoreAssociations_basic(t *testing.T) {
	rName := "data.huaweicloud_identitycenter_identity_store_associations.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenterIdentityStoreAssociations_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "identity_store_id"),
					resource.TestCheckResourceAttrSet(rName, "identity_store_type"),
					resource.TestCheckResourceAttrSet(rName, "authentication_type"),
					resource.TestCheckResourceAttrSet(rName, "provisioning_type.#"),
					resource.TestCheckResourceAttrSet(rName, "status")),
			},
		},
	})
}

const testAccDatasourceIdentityCenterIdentityStoreAssociations_basic = `
data "huaweicloud_identitycenter_instance" "test" {}
 
data "huaweicloud_identitycenter_identity_store_associations" "test"{
  depends_on  = [data.huaweicloud_identitycenter_instance.test]
  instance_id = data.huaweicloud_identitycenter_instance.test.id
}
`
