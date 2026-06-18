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

func getGaussDBInstanceLtsLogAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/logs/lts-config"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	logType := state.Primary.Attributes["log_type"]

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGaussdbInstanceLtsLogAssociateBodyParams(instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	listPath := fmt.Sprintf("instance_lts_configs[?instance.id=='%s']|[0].lts_configs", instanceId)
	ltsList := utils.PathSearch(listPath, getRespBody, []interface{}{}).([]interface{})

	filterPath := fmt.Sprintf("[?log_type=='%s']|[0]", logType)
	matchedConfig := utils.PathSearch(filterPath, ltsList, nil)
	if matchedConfig == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	enabled := utils.PathSearch("enabled", matchedConfig, false).(bool)
	if !enabled {
		return nil, golangsdk.ErrDefault404{}
	}
	return matchedConfig, nil
}

func buildGaussdbInstanceLtsLogAssociateBodyParams(instanceId string) string {
	res := ""
	res = fmt.Sprintf("%s&instance_id=%v", res, instanceId)
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
func TestAccResourceGaussdbInstanceLtsLogAssociate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	nameUpdate := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_instance_lts_log_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDBInstanceLtsLogAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGaussdbInstanceLtsLogAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttr(rName, "log_type", "audit_log"),
					resource.TestCheckResourceAttrSet(rName, "lts_group_id"),
					resource.TestCheckResourceAttrSet(rName, "lts_stream_id"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
				),
			},
			{
				Config: testAccResourceGaussdbInstanceLtsLogAssociate_update(name, nameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "log_type", "audit_log"),
					resource.TestCheckResourceAttrPair(rName, "lts_stream_id", "huaweicloud_lts_stream.test_new", "id"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceGaussdbInstanceLtsLogAssociateImportState(rName),
			},
		},
	})
}

func testAccResourceGaussdbInstanceLtsLogAssociate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_gaussdb_instance_lts_log_associate" "test" {
  instance_id    = huaweicloud_gaussdb_instance.test.id
  log_type       = "audit_log"
  lts_group_id   = huaweicloud_lts_group.test.id
  lts_stream_id  = huaweicloud_lts_stream.test.id
}
`, testDataSourceGaussdbInstanceMetrics_base(name), name)
}

func testAccResourceGaussdbInstanceLtsLogAssociate_update(name, newName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test_new" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_gaussdb_instance_lts_log_associate" "test" {
  instance_id    = huaweicloud_gaussdb_instance.test.id
  log_type       = "audit_log"
  lts_group_id   = huaweicloud_lts_group.test.id
  lts_stream_id  = huaweicloud_lts_stream.test_new.id
}
`, testDataSourceGaussdbInstanceMetrics_base(name), newName)
}

func testAccResourceGaussdbInstanceLtsLogAssociateImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", rName)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		logType := rs.Primary.Attributes["log_type"]
		return fmt.Sprintf("%s/%s", instanceId, logType), nil
	}
}
