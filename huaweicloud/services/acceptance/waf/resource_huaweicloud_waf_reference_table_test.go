package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/valuelists"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWafReferenceTableV1_basic(t *testing.T) {
	var referencTable valuelists.WafValueList
	resourceName := "huaweicloud_waf_reference_table.ref_table"
	name := acceptance.RandomAccResourceName()
	updateName := name + "_update"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafReferenceTableV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafReferenceTableV1_conf(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafReferenceTableV1Exists(resourceName, &referencTable),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "tf acc"),
					resource.TestCheckResourceAttr(resourceName, "type", "url"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
				),
			},
			{
				Config: testAccWafReferenceTableV1_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafReferenceTableV1Exists(resourceName, &referencTable),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "type", "url"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),

					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWafReferenceTableV1_withEpsID(t *testing.T) {
	var referencTable valuelists.WafValueList
	resourceName := "huaweicloud_waf_reference_table.ref_table"
	name := acceptance.RandomAccResourceName()
	updateName := name + "_update"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafReferenceTableV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafReferenceTableV1_conf_withEpsID(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafReferenceTableV1Exists(resourceName, &referencTable),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "tf acc"),
					resource.TestCheckResourceAttr(resourceName, "type", "url"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
				),
			},
			{
				Config: testAccWafReferenceTableV1_update_withEpsID(updateName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafReferenceTableV1Exists(resourceName, &referencTable),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "type", "url"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFResourceImportState(resourceName),
			},
		},
	})
}

func testAccCheckWafReferenceTableV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	wafClient, err := config.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_waf_reference_table" {
			continue
		}

		_, err := valuelists.GetWithEpsID(wafClient, rs.Primary.ID, rs.Primary.Attributes["enterprise_project_id"])
		if err == nil {
			return fmt.Errorf("WAF reference table still exists")
		}
	}

	return nil
}

func testAccCheckWafReferenceTableV1Exists(n string, valueList *valuelists.WafValueList) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		wafClient, err := config.WafV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating WAF client: %s", err)
		}

		found, err := valuelists.GetWithEpsID(wafClient, rs.Primary.ID, rs.Primary.Attributes["enterprise_project_id"])
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("WAF reference table not found")
		}

		*valueList = *found

		return nil
	}
}

func testAccWafReferenceTableV1_conf(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_reference_table" "ref_table" {
  name        = "%s"
  type        = "url"
  description = "tf acc"

  conditions = [
    "/admin",
    "/manage"
  ]

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstanceV1_conf(name), name)
}

func testAccWafReferenceTableV1_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_reference_table" "ref_table" {
  name        = "%s"
  type        = "url"
  description = ""

  conditions = [
    "/bill",
    "/sql"
  ]

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstanceV1_conf(name), name)
}

func testAccWafReferenceTableV1_conf_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_reference_table" "ref_table" {
  name                  = "%s"
  type                  = "url"
  description           = "tf acc"
  enterprise_project_id = "%s"

  conditions = [
    "/admin",
    "/manage"
  ]

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstance_epsId(name, epsID), name, epsID)
}

func testAccWafReferenceTableV1_update_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_reference_table" "ref_table" {
  name                  = "%s"
  type                  = "url"
  description           = ""
  enterprise_project_id = "%s"

  conditions = [
    "/bill",
    "/sql"
  ]

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}
`, testAccWafDedicatedInstance_epsId(name, epsID), name, epsID)
}
