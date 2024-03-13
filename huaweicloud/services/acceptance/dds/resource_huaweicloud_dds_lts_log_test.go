package dds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dds/v3/instances"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDdsLtsLogResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		getDdsLtsLogHttpUrl = "v3/{project_id}/instances/logs/lts-configs"
		getDdsLtsLogProduct = "dds"
	)

	getDdsLtsLogClient, err := cfg.NewServiceClient(getDdsLtsLogProduct, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS client: %s", err)
	}

	getDdsLtsLogPath := getDdsLtsLogClient.Endpoint + getDdsLtsLogHttpUrl
	getDdsLtsLogPath = strings.ReplaceAll(getDdsLtsLogPath, "{project_id}", getDdsLtsLogClient.ProjectID)

	getDdsLtsLogResp, err := pagination.ListAllItems(
		getDdsLtsLogClient,
		"offset",
		getDdsLtsLogPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving DDS LTS configs: %s", err)
	}

	getDdsLtsLogRespJson, err := json.Marshal(getDdsLtsLogResp)
	if err != nil {
		return nil, fmt.Errorf("error marshaling DDS LTS configs: %s", err)
	}

	var getDdsLtsLogRespBody interface{}
	err = json.Unmarshal(getDdsLtsLogRespJson, &getDdsLtsLogRespBody)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling DDS LTS configs: %s", err)
	}

	jsonPath := fmt.Sprintf("instance_lts_configs[?instance.id=='%s']|[0].lts_configs|[0]", state.Primary.ID)
	ltsConfig := utils.PathSearch(jsonPath, getDdsLtsLogRespBody, nil)
	if !utils.PathSearch("enabled", ltsConfig, false).(bool) {
		return nil, fmt.Errorf("error retrieving DDS LTS config: the LTS config can not be found")
	}

	return ltsConfig, nil
}

func TestAccDdsLtsLog_basic(t *testing.T) {
	var instance instances.InstanceResponse
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dds_lts_log.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDdsLtsLogResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDdsLtsLog_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_dds_instance.instance", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "log_type", "audit_log"),
				),
			},
			{
				Config: testAccDdsLtsLog_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_dds_instance.instance", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "log_type", "audit_log"),
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

func testAccDdsLtsLog_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_dds_lts_log" "test" {
  instance_id   = huaweicloud_dds_instance.instance.id
  log_type      = "audit_log"
  lts_group_id  = huaweicloud_lts_group.test.id
  lts_stream_id = huaweicloud_lts_stream.test.id
}`, testAccDDSInstanceV3Config_basic(rName, 8800), rName)
}

func testAccDdsLtsLog_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s-update"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s-update"
}

resource "huaweicloud_dds_lts_log" "test" {
  instance_id   = huaweicloud_dds_instance.instance.id
  log_type      = "audit_log"
  lts_group_id  = huaweicloud_lts_group.test.id
  lts_stream_id = huaweicloud_lts_stream.test.id
}`, testAccDDSInstanceV3Config_basic(rName, 8800), rName)
}
