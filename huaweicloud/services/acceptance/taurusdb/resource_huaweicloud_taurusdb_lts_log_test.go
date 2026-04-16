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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getTaurusDBLtsLogResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/logs/lts-configs"
		product = "gaussdb"
	)

	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += fmt.Sprintf("?instance_id=%s", state.Primary.Attributes["instance_id"])

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, err
	}
	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return nil, err
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return nil, err
	}
	searchPath := fmt.Sprintf("instance_lts_configs|[0].lts_configs|[?log_type=='%s' && enabled]|[0]",
		state.Primary.Attributes["log_type"])
	ltsConfig := utils.PathSearch(searchPath, listRespBody, nil)
	if ltsConfig == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return ltsConfig, nil
}

func TestAccTaurusDBLtsLog_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_taurusdb_lts_log.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getTaurusDBLtsLogResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBLtsLog_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_taurusdb_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id",
						"huaweicloud_lts_group.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id",
						"huaweicloud_lts_stream.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "log_type", "error_log"),
				),
			},
			{
				Config: testAccTaurusDBLtsLog_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_taurusdb_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_group_id",
						"huaweicloud_lts_group.test.1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "lts_stream_id",
						"huaweicloud_lts_stream.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "log_type", "error_log"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccTaurusDBLtsLogImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccTaurusDBLtsLog_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_taurusdb_flavors" "test" {
  engine  = "gaussdb-mysql"
  version = "8.0"
}

resource "huaweicloud_taurusdb_instance" "test" {
  name                     = "%[2]s"
  password                 = "Test@12345678"
  flavor                   = data.huaweicloud_taurusdb_flavors.test.flavors[0].name
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  availability_zone_mode   = "multi"
  master_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  read_replicas            = 2
  enterprise_project_id    = "0"
}

resource "huaweicloud_lts_group" "test" {
  count = 2

  group_name  = "%[2]s_${count.index}"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  count = 2

  group_id    = huaweicloud_lts_group.test[count.index].id
  stream_name = "%[2]s_${count.index}"
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccTaurusDBLtsLog_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_lts_log" "test" {
  instance_id   = huaweicloud_taurusdb_instance.test.id
  log_type      = "error_log"
  lts_group_id  = huaweicloud_lts_group.test[0].id
  lts_stream_id = huaweicloud_lts_stream.test[0].id
}`, testAccTaurusDBLtsLog_base(rName), rName)
}

func testAccTaurusDBLtsLog_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_lts_log" "test" {
  instance_id   = huaweicloud_taurusdb_instance.test.id
  log_type      = "error_log"
  lts_group_id  = huaweicloud_lts_group.test[1].id
  lts_stream_id = huaweicloud_lts_stream.test[1].id
}`, testAccTaurusDBLtsLog_base(rName), rName)
}

func testAccTaurusDBLtsLogImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("the resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("attribute (instance_id) of resource(%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["log_type"] == "" {
			return "", fmt.Errorf("attribute (log_type) of Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["log_type"]), nil
	}
}
