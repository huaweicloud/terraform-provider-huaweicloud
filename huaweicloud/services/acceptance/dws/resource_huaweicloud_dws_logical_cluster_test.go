package dws

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getLogicalClusterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/logical-clusters"
		product = "dws"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", state.Primary.Attributes["cluster_id"])
	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, err
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return nil, err
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, err
	}
	expression := fmt.Sprintf("logical_clusters[?logical_cluster_id=='%s']|[0]", state.Primary.ID)
	cluster := utils.PathSearch(expression, getRespBody, nil)
	if cluster == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return cluster, nil
}

// Two logical clusters are created to test concurrent creation and deletion scenarios.
func TestAccLogicalCluster_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dws_logical_cluster.test"
	rName2 := "huaweicloud_dws_logical_cluster.test2"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLogicalClusterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLogicalCluster_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "cluster_id",
						"huaweicloud_dws_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "logical_cluster_name", name),
					resource.TestCheckResourceAttr(rName, "cluster_rings.#", "2"),
					resource.TestCheckResourceAttr(rName2, "cluster_rings.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.host_name"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.back_ip"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.cpu_cores"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.memory"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.disk_size"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.1.host_name"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.1.back_ip"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.1.cpu_cores"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.1.memory"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.1.disk_size"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.2.host_name"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.2.back_ip"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.2.cpu_cores"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.2.memory"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.2.disk_size"),

					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "first_logical_cluster"),
					resource.TestCheckResourceAttrSet(rName, "edit_enable"),
					resource.TestCheckResourceAttrSet(rName, "restart_enable"),
					resource.TestCheckResourceAttrSet(rName, "delete_enable"),
				),
			},
			{
				Config: testLogicalCluster_basic_step2(name),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testLogicalClusterImportState(rName),
			},
		},
	})
}

func testLogicalCluster_base(name string) string {
	clusterBasic := testAccDwsCluster_basic_step1(name, 10, dws.PublicBindTypeAuto, "cluster123@!")
	return fmt.Sprintf(`
%s

data "huaweicloud_dws_logical_cluster_rings" "test" {
  cluster_id = huaweicloud_dws_cluster.test.id
}
`, clusterBasic)
}

func testLogicalCluster_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dws_logical_cluster" "test" {
  cluster_id           = huaweicloud_dws_cluster.test.id
  logical_cluster_name = "%s"

  cluster_rings {
    dynamic "ring_hosts" {
      for_each = data.huaweicloud_dws_logical_cluster_rings.test.cluster_rings.0.ring_hosts[*]
      content {
        host_name = ring_hosts.value.host_name
        back_ip   = ring_hosts.value.back_ip
        cpu_cores = ring_hosts.value.cpu_cores
        memory    = ring_hosts.value.memory
        disk_size = ring_hosts.value.disk_size
      }
    }
  }

  cluster_rings {
    dynamic "ring_hosts" {
      for_each = data.huaweicloud_dws_logical_cluster_rings.test.cluster_rings.1.ring_hosts[*]
      content {
        host_name = ring_hosts.value.host_name
        back_ip   = ring_hosts.value.back_ip
        cpu_cores = ring_hosts.value.cpu_cores
        memory    = ring_hosts.value.memory
        disk_size = ring_hosts.value.disk_size
      }
    }
  }
}

resource "huaweicloud_dws_logical_cluster" "test2" {
  cluster_id           = huaweicloud_dws_cluster.test.id
  logical_cluster_name = "%s_test2"

  cluster_rings {
    dynamic "ring_hosts" {
      for_each = data.huaweicloud_dws_logical_cluster_rings.test.cluster_rings.2.ring_hosts[*]
      content {
        host_name = ring_hosts.value.host_name
        back_ip   = ring_hosts.value.back_ip
        cpu_cores = ring_hosts.value.cpu_cores
        memory    = ring_hosts.value.memory
        disk_size = ring_hosts.value.disk_size
      }
    }
  }
}
`, testLogicalCluster_base(name), name, name)
}

func testLogicalCluster_basic_step1(name string) string {
	return testLogicalCluster_basic(name)
}

func testLogicalCluster_basic_step2(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dws_logical_cluster_restart" "test" {
  cluster_id         = huaweicloud_dws_cluster.test.id
  logical_cluster_id = huaweicloud_dws_logical_cluster.test.id
}

resource "huaweicloud_dws_logical_cluster_restart" "test2" {
  cluster_id         = huaweicloud_dws_cluster.test.id
  logical_cluster_id = huaweicloud_dws_logical_cluster.test2.id
}
`, testLogicalCluster_basic(name))
}

// testLogicalClusterImportState use to return an ID with format <cluster_id>/<id>
func testLogicalClusterImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		clusterID := rs.Primary.Attributes["cluster_id"]
		if clusterID == "" {
			return "", fmt.Errorf("attribute (cluster_id) of resource (%s) not found", name)
		}

		return clusterID + "/" + rs.Primary.ID, nil
	}
}
