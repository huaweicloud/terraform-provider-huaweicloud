package apig

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getApplicationAclFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG Client: %s", err)
	}

	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-acl"
	applicationId := state.Primary.Attributes["application_id"]
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getPath = strings.ReplaceAll(getPath, "{app_id}", applicationId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the ACL rules from the application (%s): %s", applicationId, err)
	}
	return utils.FlattenResponse(requestResp)
}

func TestAccApplicationAcl_basic(t *testing.T) {
	var (
		obj interface{}

		baseConfig = testAccApplicationAcl_basic_base()
		rName      = "huaweicloud_apig_application_acl.test"
		rc         = acceptance.InitResourceCheck(rName, &obj, getApplicationAclFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationAcl_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(rName, "application_id", "huaweicloud_apig_application.test", "id"),
					resource.TestCheckResourceAttr(rName, "type", "PERMIT"),
					resource.TestCheckResourceAttr(rName, "values.#", "2"),
					resource.TestCheckResourceAttr(rName, "values.0", "127.0.0.1"),
					resource.TestCheckResourceAttr(rName, "values.1", "192.145.0.0/16"),
				),
			},
			{
				Config: testAccApplicationAcl_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "PERMIT"),
					resource.TestCheckResourceAttr(rName, "values.#", "2"),
					resource.TestCheckResourceAttr(rName, "values.0", "192.145.0.0/16"),
					resource.TestCheckResourceAttr(rName, "values.1", "127.0.0.2-192.144.0.1"),
				),
			},
			{
				Config: testAccApplicationAcl_basic_step3(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "DENY"),
					resource.TestCheckResourceAttr(rName, "values.#", "2"),
					resource.TestCheckResourceAttr(rName, "values.0", "192.145.0.0/16"),
					resource.TestCheckResourceAttr(rName, "values.1", "127.0.0.2-192.144.0.1"),
				),
			},
			{
				Config: testAccApplicationAcl_basic_step4(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "DENY"),
					resource.TestCheckResourceAttr(rName, "values.#", "2"),
					resource.TestCheckResourceAttr(rName, "values.0", "127.0.0.1"),
					resource.TestCheckResourceAttr(rName, "values.1", "192.145.0.0/16"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationAclImportIdFunc(rName),
			},
		},
	})
}

func testAccApplicationAclImportIdFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, applicationId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rsName, rs)
		}

		instanceId = rs.Primary.Attributes["instance_id"]
		applicationId = rs.Primary.ID
		if instanceId == "" || applicationId == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, applicationId)
		}
		return fmt.Sprintf("%s/%s", instanceId, applicationId), nil
	}
}

func testAccApplicationAcl_basic_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
variable "cidrs" {
  type    = list(string)
  default = ["127.0.0.1", "192.145.0.0/16", "127.0.0.2-192.144.0.1"]
}

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_application" "test" {
  instance_id = local.instance_id
  name        = "%[2]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccApplicationAcl_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application_acl" "test" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test.id
  type           = "PERMIT"
  values         = slice(var.cidrs, 0, 2)
}
`, baseConfig)
}

func testAccApplicationAcl_basic_step2(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application_acl" "test" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test.id
  type           = "PERMIT"
  values         = slice(var.cidrs, 1, 3)
}
`, baseConfig)
}

func testAccApplicationAcl_basic_step3(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application_acl" "test" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test.id
  type           = "DENY"
  values         = slice(var.cidrs, 1, 3)
}
`, baseConfig)
}

func testAccApplicationAcl_basic_step4(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application_acl" "test" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test.id
  type           = "DENY"
  values         = slice(var.cidrs, 0, 2)
}
`, baseConfig)
}
