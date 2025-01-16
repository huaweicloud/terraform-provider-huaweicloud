package gaussdb

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

func getResourceQuota(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/enterprise-projects/quotas"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)

	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	expression := fmt.Sprintf("eps_quotas[?enterprise_project_id=='%s']|[0]", state.Primary.ID)
	quota := utils.PathSearch(expression, getRespBody, nil)
	if quota == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return quota, nil
}

func TestAccOpenGaussQuota_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_gaussdb_opengauss_quota.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceQuota,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussQuota_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "instance_quota", "10"),
					resource.TestCheckResourceAttr(resourceName, "vcpus_quota", "20"),
					resource.TestCheckResourceAttr(resourceName, "ram_quota", "30"),
					resource.TestCheckResourceAttr(resourceName, "volume_quota", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_name"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_used"),
					resource.TestCheckResourceAttrSet(resourceName, "vcpus_used"),
					resource.TestCheckResourceAttrSet(resourceName, "ram_used"),
					resource.TestCheckResourceAttrSet(resourceName, "volume_used"),
				),
			},
			{
				Config: testAccOpenGaussQuota_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "instance_quota", "20"),
					resource.TestCheckResourceAttr(resourceName, "vcpus_quota", "30"),
					resource.TestCheckResourceAttr(resourceName, "ram_quota", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume_quota", "-1"),
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

func testAccOpenGaussQuota_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_opengauss_quota" "test" {
  enterprise_project_id = "%s"
  instance_quota        = 10
  vcpus_quota           = 20
  ram_quota             = 30
  volume_quota          = 40
}`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccOpenGaussQuota_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_opengauss_quota" "test" {
  enterprise_project_id = "%s"
  instance_quota        = 20
  vcpus_quota           = 30
  ram_quota             = 0
  volume_quota          = -1
}`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
