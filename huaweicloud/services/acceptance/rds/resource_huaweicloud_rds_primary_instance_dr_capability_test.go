package rds

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

func getPrimaryInstanceDrCapabilityResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/disaster-recovery-infos"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOpt.JSONBody = utils.RemoveNil(buildGetPrimaryInstanceDrCapabilityQueryParams(state.Primary.ID))

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS primary instance DR capability: %s", err)
	}

	instanceDrInfo := utils.PathSearch("instance_dr_infos|[0]", getRespBody, nil)
	if instanceDrInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func buildGetPrimaryInstanceDrCapabilityQueryParams(id string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id": id,
	}
	return bodyParams
}

func TestAccPrimaryInstanceDrCapability_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_rds_primary_instance_dr_capability.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrimaryInstanceDrCapabilityResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckRdsTargetInstanceId(t)
			acceptance.TestAccPreCheckRdsTargetProjectId(t)
			acceptance.TestAccPreCheckRdsTargetRegion(t)
			acceptance.TestAccPreCheckRdsTargetIp(t)
			acceptance.TestAccPreCheckRdsTargetSubnet(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPrimaryInstanceDrCapability_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_RDS_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "target_instance_id", acceptance.HW_RDS_TARGET_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "target_project_id", acceptance.HW_RDS_TARGET_PROJECT_ID),
					resource.TestCheckResourceAttr(rName, "target_region", acceptance.HW_RDS_TARGET_REGION),
					resource.TestCheckResourceAttr(rName, "target_ip", acceptance.HW_RDS_TARGET_IP),
					resource.TestCheckResourceAttr(rName, "target_subnet", acceptance.HW_RDS_TARGET_SUBNET),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "time"),
					resource.TestCheckResourceAttrSet(rName, "build_process"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testPrimaryInstanceDrCapability_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_primary_instance_dr_capability" "test" {
  instance_id        = "%[1]s"
  target_instance_id = "%[2]s"
  target_project_id  = "%[3]s"
  target_region      = "%[4]s"
  target_ip          = "%[5]s"
  target_subnet      = "%[6]s"
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_TARGET_INSTANCE_ID, acceptance.HW_RDS_TARGET_PROJECT_ID,
		acceptance.HW_RDS_TARGET_REGION, acceptance.HW_RDS_TARGET_IP, acceptance.HW_RDS_TARGET_SUBNET)
}
