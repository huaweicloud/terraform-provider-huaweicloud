package taurusdb

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceQuota(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/quotas"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listMysqlDatabasesResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, err
	}

	listRespJson, err := json.Marshal(listMysqlDatabasesResp)
	if err != nil {
		return nil, err
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return nil, err
	}
	expression := fmt.Sprintf("quota_list[?enterprise_project_id=='%s']|[0]", state.Primary.Attributes["enterprise_project_id"])
	quota := utils.PathSearch(expression, listRespBody, nil)
	if quota == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return quota, nil
}

func TestAccGaussDBMysqlQuota_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_gaussdb_mysql_quota.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceQuota,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBMysqlQuota_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "instance_quota", "10"),
					resource.TestCheckResourceAttr(resourceName, "vcpus_quota", "20"),
					resource.TestCheckResourceAttr(resourceName, "ram_quota", "30"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_name"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_instance_quota"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_vcpus_quota"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_ram_quota"),
				),
			},
			{
				Config: testAccGaussDBMysqlQuota_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "instance_quota", "0"),
					resource.TestCheckResourceAttr(resourceName, "vcpus_quota", "50"),
					resource.TestCheckResourceAttr(resourceName, "ram_quota", "-1"),
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

func testAccGaussDBMysqlQuota_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_mysql_quota" "test" {
  enterprise_project_id = "%s"
  instance_quota        = 10
  vcpus_quota           = 20
  ram_quota             = 30
}`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccGaussDBMysqlQuota_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_mysql_quota" "test" {
  enterprise_project_id = "%s"
  instance_quota        = 0
  vcpus_quota           = 50
}`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
