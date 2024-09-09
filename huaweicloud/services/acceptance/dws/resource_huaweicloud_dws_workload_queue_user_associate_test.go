package dws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getWorkloadQueueUserAssociateFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("dws", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}

	clusterId := state.Primary.Attributes["cluster_id"]
	queueName := state.Primary.Attributes["queue_name"]
	return dws.GetAssociatedUserNames(client, clusterId, queueName)
}

func getUserNames() (userName, updateUserName string) {
	agencyNames := strings.Split(acceptance.HW_DWS_ASSOCIATE_USER_NAMES, ",")
	if len(agencyNames) < 2 {
		return "", ""
	}
	return agencyNames[0], agencyNames[1]
}

func TestAccWorkloadQueueUserAssociate_basic(t *testing.T) {
	var (
		resp                     interface{}
		name                     = acceptance.RandomAccResourceName()
		rName                    = "huaweicloud_dws_workload_queue_user_associate.test"
		userName, updateUserName = getUserNames()

		rc = acceptance.InitResourceCheck(
			rName,
			&resp,
			getWorkloadQueueUserAssociateFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
			acceptance.TestAccPreCheckDwsClusterUserNames(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			// Bind `${userName}` to the queue.
			{
				Config: testAccWorkloadQueueUserAssociate_basic_step1(name, userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttrPair(rName, "queue_name",
						"huaweicloud_dws_workload_queue.test", "name"),
					resource.TestCheckResourceAttr(rName, "user_names.#", "1"),
					resource.TestCheckResourceAttr(rName, "user_names.0", userName),
				),
			},
			// Unbind `${userName}` from the queue and bind `${updateUserName}` to the queue.
			{
				Config: testAccWorkloadQueueUserAssociate_basic_step2(name, updateUserName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "user_names.#", "1"),
					resource.TestCheckResourceAttr(rName, "user_names.0", updateUserName),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccWorkloadQueueUserAssociateImportStateFunc(rName),
			},
		},
	})
}

func testAccWorkloadQueueUserAssociateImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["cluster_id"] == "" || rs.Primary.Attributes["queue_name"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<cluster_id>/<queue_name>', but got '%s/%s'",
				rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["queue_name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["queue_name"]), nil
	}
}

func testAccWorkloadQueueUserAssociate_basic_step1(name, userName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dws_workload_queue_user_associate" "test" {
  cluster_id = "%[2]s"
  queue_name = huaweicloud_dws_workload_queue.test.id
  user_names = ["%[3]s"]
}
`, testAccWorkloadQueue_basic(name), acceptance.HW_DWS_CLUSTER_ID, userName)
}

func testAccWorkloadQueueUserAssociate_basic_step2(name, updateUserName string) string {
	return testAccWorkloadQueueUserAssociate_basic_step1(name, updateUserName)
}
