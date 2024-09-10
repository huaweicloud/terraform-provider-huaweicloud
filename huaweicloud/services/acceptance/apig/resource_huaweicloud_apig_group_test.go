package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apigroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getGroupFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return apigroups.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccGroup_basic(t *testing.T) {
	var (
		group apigroups.Group

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		rNameBasic = "huaweicloud_apig_group.basic"
		rcBasic    = acceptance.InitResourceCheck(rNameBasic, &group, getGroupFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcBasic.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Check whether illegal group name ​​can be intercepted normally (create phase).
				Config:      testAccGroup_basic_step1(),
				ExpectError: regexp.MustCompile("Invalid parameter value"),
			},
			{
				Config: testAccGroup_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcBasic.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameBasic, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(rNameBasic, "name", name),
					resource.TestCheckResourceAttr(rNameBasic, "description", "Created by script"),
					resource.TestCheckResourceAttrSet(rNameBasic, "created_at"),
				),
			},
			{
				Config: testAccGroup_basic_step3(updateName),
				Check: resource.ComposeTestCheckFunc(
					rcBasic.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameBasic, "name", updateName),
					resource.TestCheckResourceAttr(rNameBasic, "description", ""),
					resource.TestCheckResourceAttrSet(rNameBasic, "updated_at"),
				),
			},
			{
				ResourceName:      rNameBasic,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupImportStateFunc(rNameBasic),
			},
		},
	})
}

func testAccGroupImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rsName)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.ID == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<id>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccGroup_basic_general(name, desc string) string {
	return fmt.Sprintf(`
variable "instance_description" {
  type    = string
  default = "%[3]s"
}

resource "huaweicloud_apig_group" "basic" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = var.instance_description != "" ? var.instance_description : null
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, desc)
}

func testAccGroup_basic_step1() string {
	return testAccGroup_basic_general("INVALID_GROUP_NAME_WITH_SPECIAL_CHAR!", "")
}

func testAccGroup_basic_step2(name string) string {
	return testAccGroup_basic_general(name, "Created by script")
}

func testAccGroup_basic_step3(name string) string {
	return testAccGroup_basic_general(name, "")
}

func TestAccGroup_withVariables(t *testing.T) {
	var (
		group apigroups.Group

		name = acceptance.RandomAccResourceName()

		rNameWithVariables = "huaweicloud_apig_group.with_variables"
		rcWithVariables    = acceptance.InitResourceCheck(rNameWithVariables, &group, getGroupFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcWithVariables.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGroup_withVariables_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithVariables.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithVariables, "environment.#", "2"),
				),
			},
			{
				Config: testAccGroup_withVariables_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithVariables.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithVariables, "environment.#", "2"),
				),
			},
			{
				ResourceName:      rNameWithVariables,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupImportStateFunc(rNameWithVariables),
			},
		},
	})
}

// Create two environments for the group, and add a total of three variables to the two environments.
// Each of the two environments has a variable with the same name and different value.
func testAccGroup_withVariables_general(name string, offset int) string {
	return fmt.Sprintf(`
variable "variables_configuration" {
  type = list(object({
    name  = string
    value = string
  }))
  default = [
    {name="TEST_VAR_1", value="TEST_VALUE_1"},
    {name="TEST_VAR_2", value="TEST_VALUE_2"},
    {name="TEST_VAR_3", value="TEST_VALUE_3"},
    {name="TEST_VAR_2", value="TEST_VALUE_4"}, // same variable name, but value is different.
  ]
}

resource "huaweicloud_apig_environment" "test" {
  count = 2

  instance_id = "%[1]s"
  name        = format("%[2]s_%%d", count.index)
}

resource "huaweicloud_apig_group" "with_variables" {
  instance_id = "%[1]s"
  name        = "%[2]s"

  environment {
    environment_id = huaweicloud_apig_environment.test[0].id

    dynamic "variable" {
      for_each = slice(var.variables_configuration, 0+%[3]d, 2+%[3]d)

      content {
        name  = variable.value.name
        value = variable.value.value
      }
    }
  }
  environment {
    environment_id = huaweicloud_apig_environment.test[1].id

    dynamic "variable" {
      for_each = slice(var.variables_configuration, 1+%[3]d, 3+%[3]d)

      content {
        name  = variable.value.name
        value = variable.value.value
      }
    }
  }
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, offset)
}

func testAccGroup_withVariables_step1(name string) string {
	return testAccGroup_withVariables_general(name, 0)
}

func testAccGroup_withVariables_step2(name string) string {
	return testAccGroup_withVariables_general(name, 1)
}

func TestAccGroup_withUrlDomain(t *testing.T) {
	var (
		group apigroups.Group

		name = acceptance.RandomAccResourceName()

		rNameWithUrlDomain = "huaweicloud_apig_group.with_url_domain"
		rcWithUrlDomain    = acceptance.InitResourceCheck(rNameWithUrlDomain, &group, getGroupFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcWithUrlDomain.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGroup_withUrlDomain_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithUrlDomain.CheckResourceExists(),
					// since the order in the schema is inconsistent with the order of data obtained by the API, other parameters are not verified.
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.#", "2"),
					resource.TestCheckResourceAttrSet(rNameWithUrlDomain, "url_domains.0.min_ssl_version"),
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.0.is_http_redirect_to_https", "false"),
				),
			},
			{
				ResourceName:      rNameWithUrlDomain,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupImportStateFunc(rNameWithUrlDomain),
			},
			{
				Config: testAccGroup_withUrlDomain_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithUrlDomain.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.#", "1"),
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.0.name", "www.terraform.test3.com"),
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.0.min_ssl_version", "TLSv1.1"),
					resource.TestCheckResourceAttr(rNameWithUrlDomain, "url_domains.0.is_http_redirect_to_https", "true"),
				),
			},
			{
				// Check whether illegal URL domain ​​can be intercepted normally (update phase).
				Config:      testAccGroup_withUrlDomain_step3(name),
				ExpectError: regexp.MustCompile("error binding domain name to the API group"),
			},
		},
	})
}

func testAccGroup_withUrlDomain_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_group" "with_url_domain" {
  instance_id = "%[1]s"
  name        = "%[2]s"

  url_domains {
    name = "www.terraform.test1.com"
  }
  url_domains {
    name = "www.terraform.test2.com"
  }
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccGroup_withUrlDomain_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_group" "with_url_domain" {
  instance_id = "%[1]s"
  name        = "%[2]s"

  url_domains {
    name                      = "www.terraform.test3.com"
    min_ssl_version           = "TLSv1.1"
    is_http_redirect_to_https = true
  }
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccGroup_withUrlDomain_step3(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_group" "with_url_domain" {
  instance_id = "%[1]s"
  name        = "%[2]s"

  url_domains {
    name                      = "INVALID_URL_DOMAIN"
    min_ssl_version           = "TLSv1.1"
    is_http_redirect_to_https = true
  }
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func TestAccGroup_withDomainAccess(t *testing.T) {
	var (
		group apigroups.Group

		name = acceptance.RandomAccResourceName()

		rNameWithDomainAccess = "huaweicloud_apig_group.with_domain_access"
		rcWithDomainAccess    = acceptance.InitResourceCheck(rNameWithDomainAccess, &group, getGroupFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcWithDomainAccess.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGroup_withDomainAccess_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithDomainAccess.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithDomainAccess, "domain_access_enabled", "false"),
				),
			},
			{
				Config: testAccGroup_withDomainAccess_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithDomainAccess.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithDomainAccess, "domain_access_enabled", "true"),
				),
			},
			{
				ResourceName:      rNameWithDomainAccess,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGroupImportStateFunc(rNameWithDomainAccess),
			},
		},
	})
}

func testAccGroup_withDomainAccess_general(name string, accessEnabled bool) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_group" "with_domain_access" {
  instance_id           = "%[1]s"
  name                  = "%[2]s"
  domain_access_enabled = %v
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name, accessEnabled)
}

func testAccGroup_withDomainAccess_step1(name string) string {
	return testAccGroup_withDomainAccess_general(name, false)
}

func testAccGroup_withDomainAccess_step2(name string) string {
	return testAccGroup_withDomainAccess_general(name, true)
}
