package dcs

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

func getDcsOfflineKeyAnalysisResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/offline/key-analysis/{task_id}"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getPath = strings.ReplaceAll(getPath, "{task_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DCS offline key analysis: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccOfflineKeyAnalysis_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_offline_key_analysis.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDcsOfflineKeyAnalysisResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOfflineKeyAnalysis_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_dcs_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "node_id",
						"data.huaweicloud_dcs_instance_nodes.test", "nodes.0.node_id"),
					resource.TestCheckResourceAttrSet(rName, "group_name"),
					resource.TestCheckResourceAttrSet(rName, "node_ip"),
					resource.TestCheckResourceAttrSet(rName, "node_type"),
					resource.TestCheckResourceAttrSet(rName, "analysis_type"),
					resource.TestCheckResourceAttrSet(rName, "started_at"),
					resource.TestCheckResourceAttrSet(rName, "finished_at"),
					resource.TestCheckResourceAttrSet(rName, "total_bytes"),
					resource.TestCheckResourceAttrSet(rName, "total_num"),
					resource.TestCheckResourceAttrSet(rName, "largest_key_prefixes.#"),
					resource.TestCheckResourceAttrSet(rName, "largest_keys.#"),
					resource.TestCheckResourceAttrSet(rName, "type_bytes.#"),
					resource.TestCheckResourceAttrSet(rName, "type_num.#"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testOfflineKeyAnalysisResourceImportState(rName),
			},
		},
	})
}

func TestAccOfflineKeyAnalysis_with_backup_id(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_offline_key_analysis.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDcsOfflineKeyAnalysisResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOfflineKeyAnalysis_with_backup_id(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_dcs_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "node_id",
						"data.huaweicloud_dcs_instance_nodes.test", "nodes.0.node_id"),
					resource.TestCheckResourceAttrPair(rName, "backup_id",
						"huaweicloud_dcs_backup.test", "backup_id"),
					resource.TestCheckResourceAttrSet(rName, "group_name"),
					resource.TestCheckResourceAttrSet(rName, "node_ip"),
					resource.TestCheckResourceAttrSet(rName, "node_type"),
					resource.TestCheckResourceAttrSet(rName, "analysis_type"),
					resource.TestCheckResourceAttrSet(rName, "started_at"),
					resource.TestCheckResourceAttrSet(rName, "finished_at"),
					resource.TestCheckResourceAttrSet(rName, "total_bytes"),
					resource.TestCheckResourceAttrSet(rName, "total_num"),
					resource.TestCheckResourceAttrSet(rName, "largest_key_prefixes.#"),
					resource.TestCheckResourceAttrSet(rName, "largest_keys.#"),
					resource.TestCheckResourceAttrSet(rName, "type_bytes.#"),
					resource.TestCheckResourceAttrSet(rName, "type_num.#"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testOfflineKeyAnalysisResourceImportState(rName),
			},
		},
	})
}

func testOfflineKeyAnalysis_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode       = "ha"
  capacity         = 1
  cpu_architecture = "x86_64"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}

data "huaweicloud_dcs_instance_nodes" "test"{
  instance_id = huaweicloud_dcs_instance.test.id
}
`, name)
}

func testOfflineKeyAnalysis_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_offline_key_analysis" "test" {
  depends_on = [data.huaweicloud_dcs_instance_nodes.test]

  instance_id = huaweicloud_dcs_instance.test.id
  node_id     = data.huaweicloud_dcs_instance_nodes.test.nodes[0].node_id
}
`, testOfflineKeyAnalysis_base(name))
}

func testOfflineKeyAnalysis_with_backup_id(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_backup" "test" {
  instance_id   = huaweicloud_dcs_instance.test.id
  backup_format = "rdb"
}

resource "huaweicloud_dcs_offline_key_analysis" "test" {
  depends_on = [data.huaweicloud_dcs_instance_nodes.test]

  instance_id = huaweicloud_dcs_instance.test.id
  node_id     = data.huaweicloud_dcs_instance_nodes.test.nodes[0].node_id
  backup_id   = huaweicloud_dcs_backup.test.backup_id
}
`, testOfflineKeyAnalysis_base(name))
}

func testOfflineKeyAnalysisResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceID, rs.Primary.ID), nil
	}
}
