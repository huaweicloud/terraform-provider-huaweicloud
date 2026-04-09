package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getClusterAssociatedEipFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	return dws.GetClusterAssociatedEipById(client, state.Primary.ID)
}

// Provide enterprise project ID to test if the DWS cluster is in the enterprise project.
func TestAccClusterEipAssociate_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_dws_cluster_eip_associate.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getClusterAssociatedEipFunc)

		name          = acceptance.RandomAccResourceName()
		randomUUID, _ = uuid.GenerateUUID()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testClusterEipAssociate_nonExistentEip(randomUUID),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`error associating EIP \(%s\) to the DWS cluster \(%s\)`,
					randomUUID, acceptance.HW_DWS_CLUSTER_ID)),
			},
			{
				Config: testClusterEipAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttrPair(rName, "eip_id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestMatchResourceAttr(rName, "public_ip", regexp.MustCompile(`^\d{1,3}(\.\d{1,3}){3}$`)),
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

func testClusterEipAssociate_nonExistentEip(randomUUID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_cluster_eip_associate" "test" {
  cluster_id = "%[1]s"
  eip_id     = "%[2]s"
}
`, acceptance.HW_DWS_CLUSTER_ID, randomUUID)
}

func testClusterEipAssociate_basic(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  default = "%[1]s"
}

resource "huaweicloud_vpc_eip" "test" {
  name                  = "%[2]s"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[2]s"
    size        = 1
    charge_mode = "traffic"
  }
}

resource "huaweicloud_dws_cluster_eip_associate" "test" {
  cluster_id = "%[3]s"
  eip_id     = huaweicloud_vpc_eip.test.id
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name, acceptance.HW_DWS_CLUSTER_ID)
}
