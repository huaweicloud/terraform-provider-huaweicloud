package apig

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAppQuotaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/app-quotas/{app_quota_id}"
		product = "apig"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getPath = strings.ReplaceAll(getPath, "{app_quota_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving APIG application quota: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccResourceAppQuota_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_apig_application_quota.test"
		baseConfig   = testAccAppQuota_basic_base()
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getAppQuotaResourceFunc)
		name         = acceptance.RandomAccResourceName()
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
				Config: testAccAppQuota_basic(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "call_limits", "5"),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "SECOND"),
					resource.TestCheckResourceAttr(resourceName, "time_interval", "3"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccAppQuota_updateBasic(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttr(resourceName, "name", "update_"+name),
					resource.TestCheckResourceAttr(resourceName, "call_limits", "6"),
					resource.TestCheckResourceAttr(resourceName, "time_unit", "DAY"),
					resource.TestCheckResourceAttr(resourceName, "time_interval", "8"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceAppQuotaImportFunc(resourceName),
			},
		},
	})
}

func testAccAppQuota_basic_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
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

func testAccAppQuota_basic(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application_quota" "test" {
  instance_id   = local.instance_id
  name          = "%[2]s"
  time_unit     = "SECOND"
  call_limits   = 5
  time_interval = 3
  description   = "Created by terraform script"
}
`, baseConfig, name)
}

func testAccAppQuota_updateBasic(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application_quota" "test" {
  instance_id   = local.instance_id
  name          = "update_%[2]s"
  time_unit     = "DAY"
  call_limits   = 6
  time_interval = 8
  description   = "Updated by terraform script"
}
`, baseConfig, name)
}

func testAccResourceAppQuotaImportFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		quotaId := rs.Primary.ID
		if instanceId == "" || quotaId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want to'<instance_id>/<id>', but got '%s/%s'",
				instanceId, quotaId)
		}
		return fmt.Sprintf("%s/%s", instanceId, quotaId), nil
	}
}
