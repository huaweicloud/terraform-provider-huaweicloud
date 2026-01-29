package lakeformation

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/lakeformation"
)

func getInstanceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("lakeformation", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LakeFormation client: %s", err)
	}

	return lakeformation.GetInstanceById(client, state.Primary.ID)
}

func TestAccInstance_basic(t *testing.T) {
	var (
		instance interface{}

		resourceName = "huaweicloud_lakeformation_instance.test"
		rc           = acceptance.InitResourceCheck(resourceName, &instance, getInstanceFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
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
				Config: testAccInstance_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "shared", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					// After the shared instance is created, some additional specifications will be returned, and the
					// value of `stride_num` will most likely not be 0.
					resource.TestMatchResourceAttr(resourceName, "specs.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "to_recycle_bin", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccInstance_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "shared", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestMatchResourceAttr(resourceName, "specs.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "baar"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "to_recycle_bin", "true"),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// Used to compare the differences between the script configuration and the remote configuration,
					// only prompts changes in the acceptance test, and changes will be hidden by DiffSuppress in actual use.
					"specs_origin",
					"to_recycle_bin",
				},
			},
		},
	})
}

func testAccInstance_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lakeformation_instance" "test" {
  name                  = "%[1]s"
  shared                = true
  description           = "Created by terraform script"
  enterprise_project_id = "%[2]s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccInstance_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lakeformation_instance" "test" {
  name                  = "%[1]s"
  shared                = true
  description           = ""
  enterprise_project_id = "%[2]s"

  tags = {
    foo   = "baar"
    owner = "terraform"
  }

  to_recycle_bin = true
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// Please delete the dedicated instance 15 minutes after the test is completed, otherwise additional charges will apply.
// - After deletion, the instance will be placed in the recycle bin and billed until it is removed from the recycle bin.
// - Instances will be automatically deleted after being stored in the recycle bin for more than one day and cannot be recovered.
// - To prevent your business from being affected, please wait 15 minutes after placing the instance in the recycle bin before forcibly deleting it.
// - Once an instance is deleted, the task cannot be recovered.
func TestAccInstance_dedicated(t *testing.T) {
	var (
		instance interface{}

		resourceName = "huaweicloud_lakeformation_instance.test"
		rc           = acceptance.InitResourceCheck(resourceName, &instance, getInstanceFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
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
				Config: testAccInstance_dedicatedPrePaid_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "shared", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					// After the dedicated instance is created, some additional specifications will be returned (not
					// just the QPS specification), and the value of `stride_num` will most likely not be 0.
					resource.TestMatchResourceAttr(resourceName, "specs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(resourceName, "specs.0.spec_code", "lakeformation.unit.basic.qps"),
					resource.TestCheckResourceAttr(resourceName, "specs.0.stride_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "to_recycle_bin", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccInstance_dedicatedPrePaid_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "shared", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					// Although the number of specification list objects returned after update is likely to remain unchanged,
					// the value of `metadata.stride_num` may change.
					resource.TestMatchResourceAttr(resourceName, "specs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(resourceName, "specs.0.spec_code", "lakeformation.unit.basic.qps"),
					resource.TestCheckResourceAttr(resourceName, "specs.0.stride_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "to_recycle_bin", "true"),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// Used to compare the differences between the script configuration and the remote configuration,
					// only prompts changes in the acceptance test, and changes will be hidden by DiffSuppress in actual use.
					"specs_origin",
					"to_recycle_bin",
				},
			},
		},
	})
}

func testAccInstance_dedicatedPrePaid_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lakeformation_instance" "test" {
  name                  = "%[1]s"
  shared                = false
  description           = "Created by terraform script"
  enterprise_project_id = "%[2]s"

  tags = {
    foo = "bar"
    key = "value"
  }

  specs {
    spec_code  = "lakeformation.unit.basic.qps"
    stride_num = 1 # 2000 QPS (1 * 2000)
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccInstance_dedicatedPrePaid_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lakeformation_instance" "test" {
  name                  = "%[1]s"
  shared                = false
  description           = "Updated by terraform script"
  enterprise_project_id = "%[2]s"

  tags = {
    foo   = "baar"
    owner = "terraform"
  }

  specs {
    spec_code  = "lakeformation.unit.basic.qps"
    stride_num = 2 # 4000 QPS (2 * 2000)
  }

  to_recycle_bin = true
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
