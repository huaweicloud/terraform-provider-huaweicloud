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

func getDcsHotKeyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getHotKeyAnalysis: query DCS hot key analysis
	var (
		getHotKeyAnalysisHttpUrl = "v2/{project_id}/instances/{instance_id}/hotkey-task/{hotkey_id}"
		getHotKeyAnalysisProduct = "dcs"
	)
	getHotKeyAnalysisClient, err := cfg.NewServiceClient(getHotKeyAnalysisProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	getHotKeyAnalysisPath := getHotKeyAnalysisClient.Endpoint + getHotKeyAnalysisHttpUrl
	getHotKeyAnalysisPath = strings.ReplaceAll(getHotKeyAnalysisPath, "{project_id}", getHotKeyAnalysisClient.ProjectID)
	getHotKeyAnalysisPath = strings.ReplaceAll(getHotKeyAnalysisPath, "{instance_id}", instanceId)
	getHotKeyAnalysisPath = strings.ReplaceAll(getHotKeyAnalysisPath, "{hotkey_id}", state.Primary.ID)

	getHotKeyAnalysisOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getHotKeyAnalysisResp, err := getHotKeyAnalysisClient.Request("GET", getHotKeyAnalysisPath,
		&getHotKeyAnalysisOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving DCS hot key analysis: %s", err)
	}

	getHotKeyAnalysisRespBody, err := utils.FlattenResponse(getHotKeyAnalysisResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DCS hot key analysis: %s", err)
	}
	return getHotKeyAnalysisRespBody, nil
}

func TestAccHotKeyAnalysis_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_hotkey_analysis.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDcsHotKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHotKeyAnalysis_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_dcs_instance.instance_1", "id"),
					resource.TestCheckResourceAttr(rName, "scan_type", "manual"),
					resource.TestCheckResourceAttr(rName, "status", "success"),
					resource.TestCheckResourceAttrSet(rName, "num"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "started_at"),
					resource.TestCheckResourceAttrSet(rName, "finished_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testHotKeyAnalysisResourceImportState(rName),
			},
		},
	})
}

func testAccDcsInstance_base(instanceName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode     = "ha"
  capacity       = 0.125
  engine_version = "5.0"
}

resource "huaweicloud_dcs_instance" "instance_1" {
  name               = "%s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  port               = 6388
  capacity           = 0.125
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name

  parameters {
    id    = "2"
    name  = "maxmemory-policy"
    value = "volatile-lfu"
  }
}`, instanceName)
}

func testHotKeyAnalysis_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_hotkey_analysis" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
}
`, testAccDcsInstance_base(name))
}

// testHotKeyAnalysisResourceImportState is used to return an import id with format <instance_id>/<id>
func testHotKeyAnalysisResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceID, rs.Primary.ID), nil
	}
}
