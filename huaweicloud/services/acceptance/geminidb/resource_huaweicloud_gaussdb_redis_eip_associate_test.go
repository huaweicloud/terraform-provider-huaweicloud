package geminidb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGaussRedisEipAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getGaussRedisEipAssociate: Query GaussDB Redis node EIP associate
	var (
		getGaussRedisEipAssociateHttpUrl = "v3/{project_id}/instances"
		getGaussRedisEipAssociateProduct = "geminidb"
	)
	getGaussRedisEipAssociateClient, err := cfg.NewServiceClient(getGaussRedisEipAssociateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB for Redis Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<node_id>")
	}
	instanceID := parts[0]
	nodeID := parts[1]

	getGaussRedisEipAssociatePath := getGaussRedisEipAssociateClient.Endpoint + getGaussRedisEipAssociateHttpUrl
	getGaussRedisEipAssociatePath = strings.ReplaceAll(getGaussRedisEipAssociatePath, "{project_id}",
		getGaussRedisEipAssociateClient.ProjectID)

	getGaussRedisEipAssociateQueryParams := buildGetGaussRedisEipAssociateQueryParams(instanceID)
	getGaussRedisEipAssociatePath += getGaussRedisEipAssociateQueryParams

	getGaussRedisEipAssociateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getGaussRedisEipAssociateResp, err := getGaussRedisEipAssociateClient.Request("GET",
		getGaussRedisEipAssociatePath, &getGaussRedisEipAssociateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving EipAssociate: %s", err)
	}

	getGaussRedisEipAssociateRespBody, err := utils.FlattenResponse(getGaussRedisEipAssociateResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GaussDB Redis EIP associate")
	}

	publicIP := utils.PathSearch(fmt.Sprintf("instances[?id=='%s']|[0].groups[0].nodes[?id=='%s']|[0].public_ip",
		instanceID, nodeID), getGaussRedisEipAssociateRespBody, "")
	if publicIP == "" {
		return nil, fmt.Errorf("error retrieving GaussDB Redis EIP associate")
	}
	return getGaussRedisEipAssociateRespBody, nil
}

func buildGetGaussRedisEipAssociateQueryParams(instanceID string) string {
	return fmt.Sprintf("?id=%s", instanceID)
}

func TestAccGaussRedisEipAssociate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	password := fmt.Sprintf("Acc%s@123", acctest.RandString(5))
	rName := "huaweicloud_gaussdb_redis_eip_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussRedisEipAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGaussRedisEipAssociate_basic(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_gaussdb_redis_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "node_id",
						"huaweicloud_gaussdb_redis_instance.test", "nodes.0.id"),
					resource.TestCheckResourceAttrPair(rName, "public_ip",
						"huaweicloud_vpc_eip.test", "address"),
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

func testGaussRedisEipAssociate_basic(name, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%[2]s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_gaussdb_redis_eip_associate" "test" {
  instance_id = huaweicloud_gaussdb_redis_instance.test.id
  node_id     = huaweicloud_gaussdb_redis_instance.test.nodes[0].id
  public_ip   = huaweicloud_vpc_eip.test.address
}
`, testAccGaussRedisInstanceConfig_basic(name, password), name)
}
