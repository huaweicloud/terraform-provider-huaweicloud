package eg

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/eg/v1/channel/custom"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getCustomEventChannelFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.EgV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating EG v1 client: %s", err)
	}

	return custom.Get(client, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"])
}

// Without enterprise project association (Notes: not default enterprise project)
func TestAccCustomEventChannel_basic(t *testing.T) {
	var (
		obj custom.Channel

		rName = "huaweicloud_eg_custom_event_channel.test"
		name  = acceptance.RandomAccResourceName()
		rc    = acceptance.InitResourceCheck(rName, &obj, getCustomEventChannelFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomEventChannel_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by acceptance test"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testAccCustomEventChannel_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "cross_account_ids.#", "2"),
					resource.TestCheckResourceAttr(rName, "cross_account_ids.0", "account1"),
					resource.TestCheckResourceAttr(rName, "cross_account_ids.1", "account2"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				// Restored test.
				Config: testAccCustomEventChannel_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by acceptance test"),
					resource.TestCheckResourceAttr(rName, "cross_account_ids.#", "0"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCustomEventChannel_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name        = "%[1]s"
  description = "Created by acceptance test"
}
`, name)
}

func testAccCustomEventChannel_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name              = "%[1]s"
  cross_account_ids = ["account1", "account2"]
}
`, name)
}

func testAccCustomEventChannel_basic_step3(name string) string {
	// Test whether the relevant parameter configuration of the resource can be restored.
	return testAccCustomEventChannel_basic_step1(name)
}

func TestAccCustomEventChannel_withEpsId(t *testing.T) {
	var (
		obj custom.Channel

		rName = "huaweicloud_eg_custom_event_channel.test"
		name  = acceptance.RandomAccResourceName()
		rc    = acceptance.InitResourceCheck(rName, &obj, getCustomEventChannelFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomEventChannel_withEpsId_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "description", "Created by acceptance test"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testAccCustomEventChannel_withEpsId_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "cross_account_ids.#", "2"),
					resource.TestCheckResourceAttr(rName, "cross_account_ids.0", "account1"),
					resource.TestCheckResourceAttr(rName, "cross_account_ids.1", "account2"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccCustomEventChannel_withEpsId_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "description", "Created by acceptance test"),
					resource.TestCheckResourceAttr(rName, "cross_account_ids.#", "0"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCustomEventChannelImportIdWithEpsIdFunc(rName),
			},
		},
	})
}

func testAccCustomEventChannelImportIdWithEpsIdFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		channelId := rs.Primary.ID
		epsId := rs.Primary.Attributes["enterprise_project_id"]
		if channelId == "" || epsId == "" {
			return "", fmt.Errorf("invalid format specified for import ID (custom event channel with enterprise "+
				"project associated), want '<id>/<enterprise_project_id>', but got '%s/%s'",
				channelId, epsId)
		}
		return fmt.Sprintf("%s/%s", channelId, epsId), nil
	}
}

func testAccCustomEventChannel_withEpsId_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name                  = "%[1]s"
  description           = "Created by acceptance test"
  enterprise_project_id = "%[2]s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCustomEventChannel_withEpsId_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name                  = "%[1]s"
  cross_account_ids     = ["account1", "account2"]
  enterprise_project_id = "%[2]s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCustomEventChannel_withEpsId_step3(name string) string {
	// Test whether the relevant parameter configuration of the resource can be restored.
	return testAccCustomEventChannel_withEpsId_step1(name)
}
