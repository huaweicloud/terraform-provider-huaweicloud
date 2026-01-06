package mrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/mrs/v1/cluster"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

type GroupNodeNum struct {
	AnalysisCoreNum int
	StreamCoreNum   int
	AnalysisTaskNum int
	StreamTaskNum   int
}

func TestAccMrsMapReduceCluster_basic(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_update(rName, password, rName+"Update"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"Update"),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "update_value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_keypair(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_keypair(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_analysis(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_analysis(rName, password, buildGroupNodeNumbers(2, 0, 1, 0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "ANALYSIS"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "1"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_analysis(rName, password, buildGroupNodeNumbers(3, 0, 2, 0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "ANALYSIS"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "2"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_analysis(rName, password, buildGroupNodeNumbers(2, 0, 1, 0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "ANALYSIS"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_stream(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_stream(rName, password, buildGroupNodeNumbers(0, 2, 0, 1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "streaming_task_nodes.0.node_number", "1"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_stream(rName, password, buildGroupNodeNumbers(0, 3, 0, 2)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "streaming_task_nodes.0.node_number", "2"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_stream(rName, password, buildGroupNodeNumbers(0, 2, 0, 1)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "streaming_task_nodes.0.node_number", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_hybrid(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterFlavorID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_hybrid(rName, password, buildGroupNodeNumbers(3, 3, 1, 0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "MIXED"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.host_ips.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.host_ips.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.host_ips.#", "1"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_hybrid(rName, password, buildGroupNodeNumbers(4, 3, 2, 0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "MIXED"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.host_ips.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.host_ips.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.host_ips.#", "2"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_hybrid(rName, password, buildGroupNodeNumbers(3, 4, 1, 0)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "MIXED"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.node_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.node_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "analysis_core_nodes.0.host_ips.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "streaming_core_nodes.0.host_ips.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "analysis_task_nodes.0.host_ips.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_custom_compact(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsCustom(t)
			acceptance.TestAccPreCheckMrsClusterFlavorID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_custom_compact_step1(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "mrs_ecs_default_agency", "MRS_ECS_DEFAULT_AGENCY"),
					resource.TestCheckResourceAttr(resourceName, "master_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.host_ips.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.1.node_number", "2"),
				),
			},
			{
				Config: testAccCluster_custom_compact_step2(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "master_nodes.0.node_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.node_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.1.node_number", "3"),
				),
			},

			{
				Config: testAccCluster_custom_compact_step3(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "master_nodes.0.node_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.1.node_number", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"template_id",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_custom_separate(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsCustom(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_customSeparate(rName, password, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.host_ips.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"template_id",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_custom_fullsize(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsCustom(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_customFullsize(rName, password, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.host_ips.#", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"template_id",
				},
			},
		},
	})
}

func TestAccMrsMapReduceCluster_externalDataSources(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsCustom(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsClusterConfig_externalDataSources(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"external_datasources",
					"template_id",
				},
			},
		},
	})
}

func buildGroupNodeNumbers(analysisCoreNum, streamCoreNum, analysisTaskNum, streamTaskNum int) GroupNodeNum {
	return GroupNodeNum{
		AnalysisCoreNum: analysisCoreNum,
		StreamCoreNum:   streamCoreNum,
		AnalysisTaskNum: analysisTaskNum,
		StreamTaskNum:   streamTaskNum,
	}
}

func testAccCheckMRSV2ClusterDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := cfg.MrsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_mapreduce_cluster" {
			continue
		}

		clusterGet, err := cluster.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("the MRS cluster (%s) is still exists", rs.Primary.ID)
		}
		if clusterGet.Clusterstate == "terminated" {
			return nil
		}
	}

	return nil
}

func testAccCheckMRSV2ClusterExists(n string, clusterGet *cluster.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no MRS cluster ID")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		mrsClient, err := config.MrsV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating MRS client: %s ", err)
		}

		found, err := cluster.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*clusterGet = *found
		return nil
	}
}

func testAccMrsMapReduceClusterConfig_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/20"
  vpc_id     = huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}
`, rName, rName)
}

// The task node has not contain data disks.
func testAccMrsMapReduceClusterConfig_basic(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 1.9.2"
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  component_list     = ["Storm"]

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd)
}

func testAccMrsMapReduceClusterConfig_update(rName, pwd, newName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 1.9.2"
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  component_list     = ["Storm"]

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  tags = {
    foo1 = "bar"
    key  = "update_value"
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), newName, pwd, pwd)
}

func testAccMrsMapReduceClusterConfig_keypair(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_kps_keypair" "test" {
  name = "%s"

  lifecycle {
    ignore_changes = [
      public_key,
    ]
  }
}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 1.9.2"
  manager_admin_pass = "%s"
  node_key_pair      = huaweicloud_kps_keypair.test.name
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  component_list     = ["Storm"]

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, rName, pwd)
}

func testAccMrsMapReduceClusterConfig_analysis(rName, pwd string, nodeNums GroupNodeNum) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "ANALYSIS"
  version            = "MRS 3.5.0-LTS"
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  component_list     = ["Hadoop", "Hive", "Tez"]

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 600
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 600
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  analysis_task_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 600
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd,
		nodeNums.AnalysisCoreNum, nodeNums.AnalysisTaskNum)
}

func testAccMrsMapReduceClusterConfig_stream(rName, pwd string, nodeNums GroupNodeNum) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 1.9.2"
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  component_list     = ["Storm"]

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd,
		nodeNums.StreamCoreNum, nodeNums.StreamTaskNum)
}

func testAccMrsMapReduceClusterConfig_hybrid(rName, pwd string, nodeNums GroupNodeNum) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%[2]s"
  type               = "MIXED"
  version            = "MRS 3.5.0-LTS"
  safe_mode          = true
  manager_admin_pass = "%[3]s"
  node_admin_pass    = "%[3]s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  component_list     = ["Hadoop", "ZooKeeper", "Ranger","JobGateway", "Tez", "Kafka", "Flume", "DBService"]

  master_nodes {
    flavor            = "%[7]s"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "%[7]s"
    node_number       = %[4]d
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "%[7]s"
    node_number       = %[5]d
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  analysis_task_nodes {
    flavor            = "%[7]s"
    node_number       = %[6]d
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd,
		nodeNums.AnalysisCoreNum, nodeNums.StreamCoreNum, nodeNums.AnalysisTaskNum,
		acceptance.HW_MRS_CLUSTER_FLAVOR_ID)
}

func testAccCluster_custom_compact_step1(rName, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  name                   = "%[2]s"
  type                   = "CUSTOM"
  version                = "MRS 3.5.0-LTS"
  safe_mode              = true
  manager_admin_pass     = "%[3]s"
  node_admin_pass        = "%[3]s"
  subnet_id              = huaweicloud_vpc_subnet.test.id
  vpc_id                 = huaweicloud_vpc.test.id
  template_id            = "mgmt_control_combined_v4.1"
  component_list         = ["Hadoop", "ZooKeeper", "Ranger", "DBService"]
  mrs_ecs_default_agency = "MRS_ECS_DEFAULT_AGENCY"

  master_nodes {
    flavor            = "%[4]s"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "%[4]s"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_2"
    flavor            = "%[4]s"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, acceptance.HW_MRS_CLUSTER_FLAVOR_ID)
}

func testAccCluster_custom_compact_step2(rName, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  name                   = "%[2]s"
  type                   = "CUSTOM"
  version                = "MRS 3.5.0-LTS"
  safe_mode              = true
  manager_admin_pass     = "%[3]s"
  node_admin_pass        = "%[3]s"
  subnet_id              = huaweicloud_vpc_subnet.test.id
  vpc_id                 = huaweicloud_vpc.test.id
  template_id            = "mgmt_control_combined_v4.1"
  component_list         = ["Hadoop", "ZooKeeper", "Ranger", "DBService"]
  mrs_ecs_default_agency = "MRS_ECS_DEFAULT_AGENCY"

  master_nodes {
    flavor            = "%[4]s"
    node_number       = 4
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "%[4]s"
    node_number       = 4
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_2"
    flavor            = "%[4]s"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, acceptance.HW_MRS_CLUSTER_FLAVOR_ID)
}

func testAccCluster_custom_compact_step3(rName, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  name                   = "%[2]s"
  type                   = "CUSTOM"
  version                = "MRS 3.5.0-LTS"
  safe_mode              = true
  manager_admin_pass     = "%[3]s"
  node_admin_pass        = "%[3]s"
  subnet_id              = huaweicloud_vpc_subnet.test.id
  vpc_id                 = huaweicloud_vpc.test.id
  template_id            = "mgmt_control_combined_v4.1"
  component_list         = ["Hadoop", "ZooKeeper", "Ranger", "DBService"]
  mrs_ecs_default_agency = "MRS_ECS_DEFAULT_AGENCY"

  master_nodes {
    flavor            = "%[4]s"
    node_number       = 4
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "%[4]s"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_2"
    flavor            = "%[4]s"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, acceptance.HW_MRS_CLUSTER_FLAVOR_ID)
}

func testAccMrsClusterConfig_externalDataSources(rName, pwd string) string {
	return fmt.Sprintf(`
%[1]s


data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type           = "MySQL"
  group_type        = "general"
  db_version        = "5.7"
  instance_mode     = "single"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors.0.name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"
  fixed_ip          = "192.168.0.58"

  db {
    password = "%[3]s"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  name          = "%[2]s"
  character_set = "utf8"
}

resource "huaweicloud_rds_mysql_account" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s"
  password    = "%[3]s"
}

resource "huaweicloud_rds_mysql_database_privilege" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_mysql_database.test.name

  users {
    name = huaweicloud_rds_mysql_account.test.name
  }
}

resource "huaweicloud_mapreduce_data_connection" "test" {
  name        = "%[2]s"
  source_type = "RDS_MYSQL"
  source_info {
    db_instance_id = huaweicloud_rds_instance.test.id
    db_name        = huaweicloud_rds_mysql_database.test.name
    user_name      = huaweicloud_rds_mysql_account.test.name
    password       = "%[3]s"
  }
}


resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%[2]s"
  type               = "CUSTOM"
  version            = "MRS 3.1.5"
  safe_mode          = true
  manager_admin_pass = "%[3]s"
  node_admin_pass    = "%[3]s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  template_id        = "mgmt_control_combined_v4.1"
  component_list     = ["Hadoop", "ZooKeeper", "Ranger", "Hive"]

  master_nodes {
    flavor            = "ac7.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:2,3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "MetaStore:1,2",
      "WebHCat:3",
      "HiveServer:1,2",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "ac7.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_2"
    flavor            = "ac7.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  external_datasources {
    component_name     = "Hive"
    role_type          = "hive_metastore"
    source_type        = "RDS_MYSQL"
    data_connection_id = huaweicloud_mapreduce_data_connection.test.id
  }
}`, common.TestBaseNetwork(rName), rName, pwd)
}

func testAccMrsMapReduceClusterConfig_customSeparate(rName, pwd string, nodeNum1 int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "CUSTOM"
  version            = "MRS 3.1.0"
  safe_mode          = true
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  template_id        = "mgmt_control_separated_v4"
  component_list     = ["DBService", "Hadoop", "ZooKeeper", "Ranger"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 5
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:3,4",
      "KerberosServer:3,4",
      "KerberosAdmin:3,4",
      "quorumpeer:3,4,5",
      "NameNode:4,5",
      "Zkfc:4,5",
      "JournalNode:3,4,5",
      "ResourceManager:4,5",
      "JobHistoryServer:5",
      "DBServer:3,5",
      "HttpFS:3,5",
      "TimelineServer:5",
      "RangerAdmin:3,4",
      "UserSync:4",
      "TagSync:4",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd, nodeNum1)
}

func testAccMrsMapReduceClusterConfig_customFullsize(rName, pwd string, nodeNum1 int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "CUSTOM"
  version            = "MRS 3.1.5"
  safe_mode          = true
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  template_id        = "mgmt_control_data_separated_v4.1"
  component_list     = ["Hadoop", "Ranger", "ZooKeeper"]

  master_nodes {
    flavor            = "ac7.4xlarge.4.linux.bigdata"
    node_number       = 9
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:5,6",
      "KerberosServer:5,6",
      "KerberosAdmin:5,6",
      "quorumpeer:5,6,7,8,9",
      "NameNode:3,4",
      "Zkfc:3,4",
      "JournalNode:5,6,7",
      "ResourceManager:8,9",
      "JobHistoryServer:8,9",
      "DBServer:8,9",
      "HttpFS:8,9",
      "TimelineServer:5",
      "RangerAdmin:4,5",
      "UserSync:5",
      "TagSync:5",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "ac7.4xlarge.4.linux.bigdata"
    node_number       = %d
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd, nodeNum1)
}

func TestAccMrsMapReduceCluster_Eip_id(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	eipResourceName := "huaweicloud_vpc_eip.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))
	bName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_Eip_id(rName, password, bName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", eipResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", eipResourceName, "address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
				},
			},
		},
	})
}

func testAccMrsMapReduceClusterConfig_Eip_id(rName, pwd, bName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
	name        = "%s"
    share_type  = "PER"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 1.9.2"
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  eip_id             = huaweicloud_vpc_eip.test.id
  component_list     = ["Storm"]

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), bName, rName, pwd, pwd)
}

func TestAccMrsMapReduceCluster_Eip_publicIp(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	eipResourceName := "huaweicloud_vpc_eip.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))
	bName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_Eip_publicIp(rName, password, bName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "eip_id", eipResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "public_ip", eipResourceName, "address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
				},
			},
		},
	})
}

func testAccMrsMapReduceClusterConfig_Eip_publicIp(rName, pwd, bName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
	name        = "%s"
    share_type  = "PER"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 1.9.2"
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  public_ip          = huaweicloud_vpc_eip.test.address
  component_list     = ["Storm"]

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), bName, rName, pwd, pwd)
}

func TestAccMrsMapReduceCluster_bootstrap(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsBootstrapScript(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_bootstrap(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "false"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.name", "bootstrap_0"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.uri",
						acceptance.HW_MAPREDUCE_BOOTSTRAP_SCRIPT),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.parameters", "a"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.before_component_start", "false"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.execute_need_sudo_root", "true"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.fail_action", "continue"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.active_master", "false"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.nodes.0", "master_node_default_group"),
					resource.TestCheckResourceAttr(resourceName, "bootstrap_scripts.0.nodes.1", "node_group_1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"template_id",
				},
			},
		},
	})
}

func testAccMrsMapReduceClusterConfig_bootstrap(rName, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%[2]s"
  version            = "MRS 3.1.5"
  type               = "CUSTOM"
  safe_mode          = false
  manager_admin_pass = "%[3]s"
  node_admin_pass    = "%[3]s"
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  template_id        = "mgmt_control_combined_v4.1"
  component_list     = ["DBService", "Hadoop", "ZooKeeper", "Ranger"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:2,3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  bootstrap_scripts {
    name                   = "bootstrap_0"
    uri                    = "%[4]s"
    parameters             = "a"
    before_component_start = false
    execute_need_sudo_root = true
    fail_action            = "continue"
    active_master          = false
    nodes = [
      "master_node_default_group",
      "node_group_1"
    ]
  }
}
`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, acceptance.HW_MAPREDUCE_BOOTSTRAP_SCRIPT)
}

func TestAccMrsMapReduceCluster_alarm(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_alarm(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"template_id",
					"smn_notify",
				},
			},
		},
	})
}

func testAccMrsMapReduceClusterConfig_alarm(rName, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "topic_1" {
  name = "%[2]s"
}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%[2]s"
  version            = "MRS 3.1.5"
  type               = "CUSTOM"
  safe_mode          = false
  manager_admin_pass = "%[3]s"
  node_admin_pass    = "%[3]s"
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  template_id        = "mgmt_control_combined_v4.1"
  component_list     = ["DBService", "Hadoop", "ZooKeeper", "Ranger"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles    = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:2,3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 100
    data_volume_type  = "SAS"
    data_volume_size  = 200
    data_volume_count = 1
    assigned_roles    = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  smn_notify {
    topic_urn         = huaweicloud_smn_topic.topic_1.topic_urn
    subscription_name = "subscription-test"
  }
}
`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd)
}

func TestAccMrsMapReduceCluster_updateWithEpsId(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	password := acceptance.RandomPassword()
	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_withEpsId(rName, password, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_withEpsId(rName, password, destEPS),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func testAccMrsMapReduceClusterConfig_withEpsId(rName, pwd, epsId string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  name                   = "%s"
  type                   = "STREAMING"
  version                = "MRS 1.9.2"
  manager_admin_pass     = "%s"
  node_admin_pass        = "%s"
  enterprise_project_id  = "%s"
  subnet_id              = huaweicloud_vpc_subnet.test.id
  vpc_id                 = huaweicloud_vpc.test.id
  component_list         = ["Storm"]

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd, epsId)
}

func TestAccMrsMapReduceCluster_prepaid_basic(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mapreduce_cluster.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	password := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceCluster_prepaid_basic(rName, password, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
				),
			},
			{
				Config: testAccMrsMapReduceCluster_prepaid_basic(rName, password, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"charging_mode",
					"period",
					"period_unit",
					"auto_renew",
					"master_nodes.0.charging_mode",
					"master_nodes.0.period",
					"master_nodes.0.period_unit",
					"master_nodes.0.auto_renew",
				},
			},
		},
	})
}

func testAccMrsMapReduceCluster_prepaid_basic(rName, pwd, autoRenew string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%[2]s"
  type               = "STREAMING"
  version            = "MRS 1.9.2"
  manager_admin_pass = "%[3]s"
  node_admin_pass    = "%[3]s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  component_list     = ["Storm"]

  charging_mode = "prePaid"
  period        = 1
  period_unit   = "month"
  auto_renew    = "%[4]s"

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1

    charging_mode = "prePaid"
    period        = 1
    period_unit   = "month"
    auto_renew    = "%[4]s"
  }
  streaming_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, autoRenew)
}

func TestAccMrsMapReduceCluster_prepaid_custom(t *testing.T) {
	var (
		clusterGet   cluster.Cluster
		resourceName = "huaweicloud_mapreduce_cluster.test"
		rName        = acceptance.RandomAccResourceName()
		password     = acceptance.RandomPassword()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterFlavorID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCluster_prepaid_custom_step1(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "master_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.node_number", "3"),
				),
			},
			{
				Config: testAccCluster_prepaid_custom_step2(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "master_nodes.0.node_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "custom_nodes.0.node_number", "4"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
					"template_id",
					"charging_mode",
					"period",
					"period_unit",
					"auto_renew",
				},
			},
		},
	})
}

func testAccCluster_prepaid_custom_step1(rName, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%[2]s"
  type               = "CUSTOM"
  version            = "MRS 3.5.0-LTS"
  safe_mode          = true
  manager_admin_pass = "%[3]s"
  node_admin_pass    = "%[3]s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  template_id        = "mgmt_control_combined_v4.1"
  component_list     = ["Hadoop", "Ranger", "ZooKeeper", "DBService"]

  charging_mode = "prePaid"
  period        = 1
  period_unit   = "month"
  auto_renew    = "true"

  master_nodes {
    flavor            = "%[4]s"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "%[4]s"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, acceptance.HW_MRS_CLUSTER_FLAVOR_ID)
}

func testAccCluster_prepaid_custom_step2(rName, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%[2]s"
  type               = "CUSTOM"
  version            = "MRS 3.5.0-LTS"
  safe_mode          = true
  manager_admin_pass = "%[3]s"
  node_admin_pass    = "%[3]s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  template_id        = "mgmt_control_combined_v4.1"
  component_list     = ["Hadoop", "Ranger", "ZooKeeper", "DBService"]

  charging_mode = "prePaid"
  period        = 1
  period_unit   = "month"
  auto_renew    = "true"

  master_nodes {
    flavor            = "%[4]s"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "OMSServer:1,2",
      "SlapdServer:1,2",
      "KerberosServer:1,2",
      "KerberosAdmin:1,2",
      "quorumpeer:1,2,3",
      "NameNode:2,3",
      "Zkfc:2,3",
      "JournalNode:1,2,3",
      "ResourceManager:2,3",
      "JobHistoryServer:3",
      "DBServer:1,3",
      "HttpFS:1,3",
      "TimelineServer:3",
      "RangerAdmin:1,2",
      "UserSync:2",
      "TagSync:2",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }

  custom_nodes {
    group_name        = "node_group_1"
    flavor            = "%[4]s"
    node_number       = 4
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1

    assigned_roles = [
      "DataNode",
      "NodeManager",
      "KerberosClient",
      "SlapdClient",
      "meta"
    ]
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, acceptance.HW_MRS_CLUSTER_FLAVOR_ID)
}
