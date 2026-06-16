package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
)

func getResourceConfigurationDictionaryFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetConfigurationDictionaryInfo(client, state.Primary.ID)
}

func TestAccResourceConfigurationDictionary_basic(t *testing.T) {
	var (
		rName          = "huaweicloud_secmaster_configuration_dictionary.test"
		dictCode       = acceptance.RandomAccResourceName()
		updatedictCode = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceConfigurationDictionaryFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfigurationDictionary_basic(dictCode),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "dict_key", "alert_comments"),
					resource.TestCheckResourceAttr(rName, "dict_code", dictCode),
					resource.TestCheckResourceAttr(rName, "dict_val", "Open"),
					resource.TestCheckResourceAttr(rName, "language", "zh"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				Config: testAccConfigurationDictionary_update(updatedictCode),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "dict_code", updatedictCode),
					resource.TestCheckResourceAttr(rName, "dict_val", "Closed"),
					resource.TestCheckResourceAttr(rName, "dict_pkey", "alert_comments_pkey_update"),
					resource.TestCheckResourceAttr(rName, "description", "alert_comments_action_update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"is_built_in",
				},
			},
		},
	})
}

func testAccConfigurationDictionary_basic(dictCode string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_configuration_dictionary" "test" {
  dict_id      = "3027"
  dict_key     = "alert_comments"
  dict_code    = "%[1]s"
  dict_val     = "Open"
  language     = "zh"
  version      = "1.0.0"
  dict_pkey    = "alert_comments_pkey"
  scope        = "ALERT"
  description  = "alert_comments_action"
  extend_field = {"test": "extend"}
  is_built_in  = false
}
`, dictCode)
}

func testAccConfigurationDictionary_update(dictCode string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_configuration_dictionary" "test" {
  dict_id      = "3027"
  dict_key     = "alert_comments"
  dict_code    = "%[1]s"
  dict_val     = "Closed"
  language     = "zh"
  version      = "1.0.0"
  dict_pkey    = "alert_comments_pkey_update"
  scope        = "ALERT"
  description  = "alert_comments_action_update"
  extend_field = {"tests": "extends"}
  is_built_in  = false
}
`, dictCode)
}
