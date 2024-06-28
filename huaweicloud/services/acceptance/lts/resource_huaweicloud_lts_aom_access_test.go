package lts

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

func getAOMAccessResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/lts/aom-mapping/{rule_id}"
		product = "lts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{rule_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		// there is no special error code here
		return nil, fmt.Errorf("error retrieving AOM to LTS log mapping rule: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	arrayBody, ok := getRespBody.([]interface{})
	if !ok {
		return nil, fmt.Errorf("error retrieving AOM to LTS log mapping rule: the API response is not array")
	}
	if len(arrayBody) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}
	return arrayBody[0], nil
}

func TestAccAOMAccess_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_aom_access.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAOMAccessResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLtsAomAccess(t)
			acceptance.TestAccPreCheckLtsAomAccessUpdate(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAOMAccess_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_LTS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "cluster_name", acceptance.HW_LTS_CLUSTER_NAME),
					resource.TestCheckResourceAttr(rName, "namespace", "default"),
					resource.TestCheckResourceAttr(rName, "workloads.0", "__ALL_DEPLOYMENTS__"),
					resource.TestCheckResourceAttr(rName, "container_name", "test_container"),
					resource.TestCheckResourceAttr(rName, "access_rules.0.file_name", "/test/*"),
					resource.TestCheckResourceAttr(rName, "access_rules.1.file_name", "/demo/demo.log"),
					resource.TestCheckResourceAttrPair(rName, "access_rules.0.log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "access_rules.0.log_group_name",
						"huaweicloud_lts_group.test", "group_name"),
					resource.TestCheckResourceAttrPair(rName, "access_rules.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "access_rules.0.log_stream_name",
						"huaweicloud_lts_stream.test", "stream_name"),
				),
			},
			{
				Config: testAOMAccess_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_LTS_CLUSTER_ID_ANOTHER),
					resource.TestCheckResourceAttr(rName, "cluster_name", acceptance.HW_LTS_CLUSTER_NAME_ANOTHER),
					resource.TestCheckResourceAttr(rName, "namespace", "test_namespace"),
					resource.TestCheckResourceAttr(rName, "workloads.0", "WORKLOAD_1"),
					resource.TestCheckResourceAttr(rName, "workloads.1", "WORKLOAD_2"),
					resource.TestCheckResourceAttr(rName, "container_name", "test_container_update"),
					resource.TestCheckResourceAttr(rName, "access_rules.0.file_name", "/test/update/*"),
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

func TestAccAOMAccess_withCCICluster(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_aom_access.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAOMAccessResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAOMAccess_withCCICluster(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "cluster_id", "CCI-ClusterID"),
					resource.TestCheckResourceAttr(rName, "cluster_name", "CCI-Cluster"),
					resource.TestCheckResourceAttr(rName, "namespace", "default"),
					resource.TestCheckResourceAttr(rName, "workloads.0", "__ALL_DEPLOYMENTS__"),
					resource.TestCheckResourceAttr(rName, "container_name", "test_container"),
					resource.TestCheckResourceAttr(rName, "access_rules.0.file_name", "__ALL_FILES__"),
					resource.TestCheckResourceAttrPair(rName, "access_rules.0.log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "access_rules.0.log_group_name",
						"huaweicloud_lts_group.test", "group_name"),
					resource.TestCheckResourceAttrPair(rName, "access_rules.0.log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "access_rules.0.log_stream_name",
						"huaweicloud_lts_stream.test", "stream_name"),
				),
			},
			{
				Config: testAOMAccess_withCCICluster_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "namespace", "test_namespace"),
					resource.TestCheckResourceAttr(rName, "workloads.0", "WORKLOAD_3"),
					resource.TestCheckResourceAttr(rName, "container_name", "test_container_update"),
					resource.TestCheckResourceAttr(rName, "access_rules.0.file_name", "/test/*"),
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

func testAOMAccess_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_aom_access" "test" {
  name           = "%[2]s"
  cluster_id     = "%[3]s"
  cluster_name   = "%[4]s"
  namespace      = "default"
  workloads      = ["__ALL_DEPLOYMENTS__"]
  container_name = "test_container"

  access_rules {
    file_name       = "/test/*"
    log_group_id    = huaweicloud_lts_group.test.id
    log_group_name  = huaweicloud_lts_group.test.group_name
    log_stream_id   = huaweicloud_lts_stream.test.id
    log_stream_name = huaweicloud_lts_stream.test.stream_name
  }

  access_rules {
    file_name       = "/demo/demo.log"
    log_group_id    = huaweicloud_lts_group.test.id
    log_group_name  = huaweicloud_lts_group.test.group_name
    log_stream_id   = huaweicloud_lts_stream.test.id
    log_stream_name = huaweicloud_lts_stream.test.stream_name
  }
}
`, testAccLtsStream_basic(name), name, acceptance.HW_LTS_CLUSTER_ID, acceptance.HW_LTS_CLUSTER_NAME)
}

func testAOMAccess_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_aom_access" "test" {
  name           = "%[2]s_update"
  cluster_id     = "%[3]s"
  cluster_name   = "%[4]s"
  namespace      = "test_namespace"
  workloads      = ["WORKLOAD_1", "WORKLOAD_2"]
  container_name = "test_container_update"

  access_rules {
    file_name       = "/test/update/*"
    log_group_id    = huaweicloud_lts_group.test.id
    log_group_name  = huaweicloud_lts_group.test.group_name
    log_stream_id   = huaweicloud_lts_stream.test.id
    log_stream_name = huaweicloud_lts_stream.test.stream_name
  }
}
`, testAccLtsStream_basic(name), name, acceptance.HW_LTS_CLUSTER_ID_ANOTHER, acceptance.HW_LTS_CLUSTER_NAME_ANOTHER)
}

func testAOMAccess_withCCICluster(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_aom_access" "test" {
  name           = "%[2]s"
  cluster_id     = "CCI-ClusterID"
  cluster_name   = "CCI-Cluster"
  namespace      = "default"
  workloads      = ["__ALL_DEPLOYMENTS__"]
  container_name = "test_container"

  access_rules {
    file_name       = "__ALL_FILES__"
    log_group_id    = huaweicloud_lts_group.test.id
    log_group_name  = huaweicloud_lts_group.test.group_name
    log_stream_id   = huaweicloud_lts_stream.test.id
    log_stream_name = huaweicloud_lts_stream.test.stream_name
  }
}
`, testAccLtsStream_basic(name), name)
}

func testAOMAccess_withCCICluster_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_aom_access" "test" {
  name           = "%[2]s_update"
  cluster_id     = "CCI-ClusterID"
  cluster_name   = "CCI-Cluster"
  namespace      = "test_namespace"
  workloads      = ["WORKLOAD_3"]
  container_name = "test_container_update"

  access_rules {
    file_name       = "/test/*"
    log_group_id    = huaweicloud_lts_group.test.id
    log_group_name  = huaweicloud_lts_group.test.group_name
    log_stream_id   = huaweicloud_lts_stream.test.id
    log_stream_name = huaweicloud_lts_stream.test.stream_name
  }
}
`, testAccLtsStream_basic(name), name)
}
