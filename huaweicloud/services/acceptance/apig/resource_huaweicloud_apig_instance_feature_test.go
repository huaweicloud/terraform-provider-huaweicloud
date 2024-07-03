package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getInstanceFeatureFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	return apig.GetInstanceFeature(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccInstanceFeature_basic(t *testing.T) {
	var (
		feature interface{}
		rName   = "huaweicloud_apig_instance_feature.test"
		name    = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&feature,
		getInstanceFeatureFunc,
	)

	// Avoid CheckDestroy because this resource already exists and does not need to be deleted.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceFeature_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "ratelimit"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "config", "{\"api_limits\":200}"),
				),
			},
			{
				Config: testAccInstanceFeature_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "ratelimit"),
					resource.TestCheckResourceAttr(rName, "config", "{\"api_limits\":300}"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccInstanceFeatureResourceImportStateFunc(rName),
			},
		},
	})
}

func testAccInstanceFeatureResourceImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		featureName := rs.Primary.ID
		if instanceId == "" || featureName == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<name>', but '%s/%s'",
				instanceId, featureName)
		}
		return fmt.Sprintf("%s/%s", instanceId, featureName), nil
	}
}

func testAccInstanceFeature_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_instance_feature" "test" {
  instance_id = "%[1]s"
  name        = "ratelimit"
  enabled     = true

  config = jsonencode({
    api_limits = 200
  })
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccInstanceFeature_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_instance_feature" "test" {
  instance_id = "%[1]s"
  name        = "ratelimit"
  enabled     = true

  config = jsonencode({
    api_limits = 300
  })
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}
