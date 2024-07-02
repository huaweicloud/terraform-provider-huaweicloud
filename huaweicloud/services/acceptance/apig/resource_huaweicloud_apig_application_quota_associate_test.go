package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getApplicationQuotaAssociateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	quotaId := state.Primary.ID
	associatedApps, err := apig.QueryQuotaAssociatedApplications(client, state.Primary.Attributes["instance_id"], quotaId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the associated application(s) from the application quota (%s)", quotaId)
	}
	if len(associatedApps) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return associatedApps, nil
}

func TestAccApplicationQuotaAssociate_basic(t *testing.T) {
	var (
		obj interface{}

		baseConfig = testAccApplicationQuotaAssociate_basic_base()

		rNamePart1 = "huaweicloud_apig_application_quota_associate.part1"
		rcPart1    = acceptance.InitResourceCheck(rNamePart1, &obj, getApplicationQuotaAssociateFunc)
		rNamePart2 = "huaweicloud_apig_application_quota_associate.part2"
		rcPart2    = acceptance.InitResourceCheck(rNamePart2, &obj, getApplicationQuotaAssociateFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcPart1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationQuotaAssociate_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rNamePart1, "instance_id",
						"data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(rNamePart1, "quota_id",
						"huaweicloud_apig_application_quota.test.0", "id"),
					resource.TestCheckResourceAttr(rNamePart1, "applications.#", "2"),
					rcPart2.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rNamePart2, "instance_id",
						"data.huaweicloud_apig_instances.test", "instances.0.id"),
					resource.TestCheckResourceAttrPair(rNamePart2, "quota_id",
						"huaweicloud_apig_application_quota.test.1", "id"),
					resource.TestCheckResourceAttr(rNamePart2, "applications.#", "1"),
				),
			},
			{
				Config: testAccApplicationQuotaAssociate_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNamePart1, "applications.#", "1"),
					rcPart2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNamePart2, "applications.#", "2"),
				),
			},
			{
				ResourceName:      rNamePart1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationQuotaAssociateImportIdFunc(rNamePart1),
			},
			{
				ResourceName:      rNamePart1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApplicationQuotaAssociateImportIdFunc(rNamePart1),
			},
		},
	})
}

func testAccApplicationQuotaAssociateImportIdFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, quotaId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rsName, rs)
		}

		instanceId = rs.Primary.Attributes["instance_id"]
		quotaId = rs.Primary.ID
		if instanceId == "" || quotaId == "" {
			return "", fmt.Errorf("missing some attributes, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, quotaId)
		}
		return fmt.Sprintf("%s/%s", instanceId, quotaId), nil
	}
}

func testAccApplicationQuotaAssociate_basic_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_application" "test" {
  count = 3

  instance_id = local.instance_id
  name        = format("%[3]s_%%d", count.index)
}

resource "huaweicloud_apig_application_quota" "test" {
  count = 2

  instance_id   = local.instance_id
  name          = format("%[3]s_%%d", count.index)
  time_unit     = "MINUTE"
  call_limits   = 200
  time_interval = 3
}
`, common.TestBaseNetwork(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccApplicationQuotaAssociate_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application_quota_associate" "part1" {
  instance_id = local.instance_id
  quota_id    = huaweicloud_apig_application_quota.test[0].id

  dynamic "applications" {
    for_each = slice(huaweicloud_apig_application.test[*].id, 0, 2)

    content {
      id = applications.value
    }
  }
}

resource "huaweicloud_apig_application_quota_associate" "part2" {
  instance_id = local.instance_id
  quota_id    = huaweicloud_apig_application_quota.test[1].id

  dynamic "applications" {
    for_each = slice(huaweicloud_apig_application.test[*].id, 2, 3)

    content {
      id = applications.value
    }
  }
}
`, baseConfig)
}

func testAccApplicationQuotaAssociate_basic_step2(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application_quota_associate" "part1" {
  instance_id = local.instance_id
  quota_id    = huaweicloud_apig_application_quota.test[0].id

  dynamic "applications" {
    for_each = slice(huaweicloud_apig_application.test[*].id, 0, 1)

    content {
      id = applications.value
    }
  }
}

resource "huaweicloud_apig_application_quota_associate" "part2" {
  instance_id = local.instance_id
  quota_id    = huaweicloud_apig_application_quota.test[1].id

  dynamic "applications" {
    for_each = slice(huaweicloud_apig_application.test[*].id, 1, 3)

    content {
      id = applications.value
    }
  }
}
`, baseConfig)
}
