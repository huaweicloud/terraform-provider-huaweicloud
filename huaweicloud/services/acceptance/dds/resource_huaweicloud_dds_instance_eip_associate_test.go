package dds

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

func getDDSInstanceEipAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.DdsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS client: %s ", err)
	}

	instID := state.Primary.Attributes["instance_id"]
	getInstanceInfoHttpUrl := "v3/{project_id}/instances?id={instance_id}"
	getInstanceInfoPath := client.Endpoint + getInstanceInfoHttpUrl
	getInstanceInfoPath = strings.ReplaceAll(getInstanceInfoPath, "{project_id}", client.ProjectID)
	getInstanceInfoPath = strings.ReplaceAll(getInstanceInfoPath, "{instance_id}", instID)
	getInstanceInfoOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getInstanceInfoResp, err := client.Request("GET", getInstanceInfoPath, &getInstanceInfoOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting instance(%s) info: %s", instID, err)
	}

	getInstanceInfoRespBody, err := utils.FlattenResponse(getInstanceInfoResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten response: %s", err)
	}

	jsonPaths := fmt.Sprintf("instances|[0].groups[*].nodes[?id=='%s'][]|[0].public_ip", state.Primary.Attributes["node_id"])
	publicIP := utils.PathSearch(jsonPaths, getInstanceInfoRespBody, "")
	if publicIP.(string) == "" {
		return nil, fmt.Errorf("error retrieving public IP")
	}

	return getInstanceInfoRespBody, nil
}

func TestAccDDSV3InstanceBindEIP_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dds_instance_eip_associate.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDDSInstanceEipAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceBindEIP_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", "huaweicloud_vpc_eip.test", "address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceDDSInstanceNodeImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccDDSInstanceBindEIP_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    share_type  = "PER"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "huaweicloud_dds_instance_eip_associate" "test" { 
  instance_id = huaweicloud_dds_instance.instance.id
  node_id     = huaweicloud_dds_instance.instance.nodes.1.id
  public_ip   = huaweicloud_vpc_eip.test.address
}`, testAccDDSInstanceReplicaSetBasic(rName), rName)
}
