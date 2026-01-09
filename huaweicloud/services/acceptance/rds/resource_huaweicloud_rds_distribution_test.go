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

func getDistributionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/distribution"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	instanceID := state.Primary.Attributes["instance_id"]
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	distributorInstanceId := utils.PathSearch("distributor_instance_id", getRespBody, nil)
	if distributorInstanceId == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccDistribution_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_distribution.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDistributionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDistribution_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "distributor_instance_name"),
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

func testDistribution_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = data.huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2022_SE"
  instance_mode = "single"
}

resource "huaweicloud_rds_instance" "test" {
  depends_on        = [huaweicloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  tde_enabled       = true

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2022_SE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testDistribution_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_distribution" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}
`, testDistribution_base(name))
}
