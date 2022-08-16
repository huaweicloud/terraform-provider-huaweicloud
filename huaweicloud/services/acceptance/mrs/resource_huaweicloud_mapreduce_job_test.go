package mrs

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/mrs/v2/jobs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	mrsRes "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/mrs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccMrsMapReduceJob_basic(t *testing.T) {
	var job jobs.Job
	resourceName := "huaweicloud_mapreduce_job.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	pwd := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2JobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceJobConfig_basic(rName, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2JobExists(resourceName, &job),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", mrsRes.JobSparkSubmit),
					resource.TestCheckResourceAttr(resourceName, "program_path",
						"obs://obs-demo-analysis-tf/program/driver_behavior.jar"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMRSClusterSubResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCheckMRSV2JobDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.MrsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_mapreduce_job" {
			continue
		}

		_, err := jobs.Get(client, rs.Primary.Attributes["cluster_id"], rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("MRS cluster (%s) is still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMRSV2JobExists(n string, job *jobs.Job) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No MRS cluster ID")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.MrsV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloud MRS client: %s ", err)
		}

		found, err := jobs.Get(client, rs.Primary.Attributes["cluster_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*job = *found
		return nil
	}
}

func testAccMRSClusterSubResourceImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["cluster_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID), nil
	}
}

func testAccMrsMapReduceJobConfig_base(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "ANALYSIS"
  version            = "MRS 1.9.2"
  safe_mode          = false
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id
  component_list     = ["Hadoop", "Spark", "Hive", "Tez"]

  master_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SSD"
    root_volume_size  = 300
    data_volume_type  = "SSD"
    data_volume_size  = 480
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SSD"
    root_volume_size  = 300
    data_volume_type  = "SSD"
    data_volume_size  = 480
    data_volume_count = 1
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd)
}

func testAccMrsMapReduceJobConfig_basic(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mapreduce_job" "test" {
  cluster_id   = huaweicloud_mapreduce_cluster.test.id
  name         = "%s"
  type         = "SparkSubmit"
  program_path = "obs://obs-demo-analysis-tf/program/driver_behavior.jar"
  parameters   = "%s %s 1 obs://obs-demo-analysis-tf/input obs://obs-demo-analysis-tf/output"

  program_parameters = {
    "--class" = "com.huawei.bigdata.spark.examples.DriverBehavior"
  }
}`, testAccMrsMapReduceJobConfig_base(rName, pwd), rName, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}
