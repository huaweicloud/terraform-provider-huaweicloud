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

func getDcsBigKeyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getBigKeyAnalysis: query DCS big key analysis
	var (
		getBigKeyAnalysisHttpUrl = "v2/{project_id}/instances/{instance_id}/bigkey-task/{bigkey_id}"
		getBigKeyAnalysisProduct = "dcs"
	)
	getBigKeyAnalysisClient, err := cfg.NewServiceClient(getBigKeyAnalysisProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	getBigKeyAnalysisPath := getBigKeyAnalysisClient.Endpoint + getBigKeyAnalysisHttpUrl
	getBigKeyAnalysisPath = strings.ReplaceAll(getBigKeyAnalysisPath, "{project_id}", getBigKeyAnalysisClient.ProjectID)
	getBigKeyAnalysisPath = strings.ReplaceAll(getBigKeyAnalysisPath, "{instance_id}", instanceId)
	getBigKeyAnalysisPath = strings.ReplaceAll(getBigKeyAnalysisPath, "{bigkey_id}", state.Primary.ID)

	getBigKeyAnalysisOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getBigKeyAnalysisResp, err := getBigKeyAnalysisClient.Request("GET", getBigKeyAnalysisPath,
		&getBigKeyAnalysisOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving DCS big key analysis: %s", err)
	}

	getBigKeyAnalysisRespBody, err := utils.FlattenResponse(getBigKeyAnalysisResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DCS big key analysis: %s", err)
	}
	return getBigKeyAnalysisRespBody, nil
}

func TestAccBigKeyAnalysis_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_bigkey_analysis.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDcsBigKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBigKeyAnalysis_basic(name),
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
				ImportStateIdFunc: testBigKeyAnalysisResourceImportState(rName),
			},
		},
	})
}

func testBigKeyAnalysis_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_bigkey_analysis" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
}
`, testAccDcsV1Instance_basic(name))
}

// testBigKeyAnalysisResourceImportState is used to return an import id with format <instance_id>/<id>
func testBigKeyAnalysisResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceID, rs.Primary.ID), nil
	}
}
