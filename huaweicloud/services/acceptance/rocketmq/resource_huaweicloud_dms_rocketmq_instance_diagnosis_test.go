package rocketmq

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

func getDiagnosisResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		getDiagnosisHttpUrl = "v2/{engine}/{project_id}/diagnosis/{report_id}"
		getDiagnosisProduct = "dmsv2"
	)
	getDiagnosisClient, err := cfg.NewServiceClient(getDiagnosisProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	getDiagnosisPath := getDiagnosisClient.Endpoint + getDiagnosisHttpUrl
	getDiagnosisPath = strings.ReplaceAll(getDiagnosisPath, "{engine}", "rocketmq")
	getDiagnosisPath = strings.ReplaceAll(getDiagnosisPath, "{project_id}", getDiagnosisClient.ProjectID)
	getDiagnosisPath = strings.ReplaceAll(getDiagnosisPath, "{report_id}", state.Primary.ID)

	getDiagnosisOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getDiagnosisResp, err := getDiagnosisClient.Request("GET", getDiagnosisPath, &getDiagnosisOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving DMS RocketMQ instance Diagnosis: %s", err)
	}

	return utils.FlattenResponse(getDiagnosisResp)
}

func TestAccDiagnosis_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dms_rocketmq_instance_diagnosis.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDiagnosisResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDiagnosis_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "report_id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "abnormal_item_sum"),
					resource.TestCheckResourceAttrSet(resourceName, "faulted_node_sum"),
					resource.TestCheckResourceAttrSet(resourceName, "consumer_nums"),
					resource.TestCheckResourceAttrSet(resourceName, "online"),
					resource.TestCheckResourceAttrSet(resourceName, "message_accumulation"),
					resource.TestCheckResourceAttrSet(resourceName, "subscription_consistency"),
					resource.TestCheckResourceAttrSet(resourceName, "subscriptions.#"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnosis_node_reports.#"),
				),
			},
		},
	})
}

func testAccDiagnosis_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_instance_nodes" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_dms_rocketmq_consumer_group" "test" {
  instance_id     = "%[1]s"
  name            = "%[2]s"
  retry_max_times = "3"
  description     = "terraform test"
}

resource "huaweicloud_dms_rocketmq_instance_diagnosis" "test" {
  instance_id = "%[1]s"
  group_name  = "%[2]s"
  node_ids    = data.huaweicloud_dms_rocketmq_instance_nodes.test.nodes[*].id
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, name)
}
