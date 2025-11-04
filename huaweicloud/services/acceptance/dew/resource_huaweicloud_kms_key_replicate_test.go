package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKmsKeyReplicate_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKmsKeyID(t)
			// The region and project of the target key cannot be the same as the key to be copied.
			acceptance.TestAccPreCheckKmsKeyReplicateRegion(t)
			acceptance.TestAccPreCheckKmsKeyReplicateProjectId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyReplicate_basic(rName),
			},
		},
	})
}

func testAccKmsKeyReplicate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key_replicate" "test" {
  key_id             = "%[1]s"
  key_alias          = "%[2]s"
  replica_region     = "%[3]s"
  replica_project_id = "%[4]s"
  key_description    = "test description"

  tags = {
    environment = "production"
    owner       = "security"
  }
}
`, acceptance.HW_KMS_KEY_ID, name, acceptance.HW_KMS_KEY_REPLICATE_REGION, acceptance.HW_KMS_KEY_REPLICATE_PROJECT_ID)
}
