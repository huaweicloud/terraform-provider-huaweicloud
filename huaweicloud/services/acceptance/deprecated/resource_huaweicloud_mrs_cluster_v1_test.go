package deprecated

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccMRSV1Cluster_basic(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "huaweicloud_mrs_cluster.cluster1"
	rName := fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
	password := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV1ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMRSV1ClusterConfig_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV1ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(
						resourceName, "cluster_state", "running"),
				),
			},
		},
	})
}

func testAccCheckMRSV1ClusterDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	mrsClient, err := config.MrsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating huaweicloud mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_mrs_cluster_v1" {
			continue
		}

		clusterGet, err := cluster.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmtp.Errorf("cluster still exists. err : %s", err)
		}
		if clusterGet.Clusterstate == "terminated" {
			return nil
		}
	}

	return nil
}

func testAccCheckMRSV1ClusterExists(n string, clusterGet *cluster.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set. ")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		mrsClient, err := config.MrsV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating huaweicloud mrs client: %s ", err)
		}

		found, err := cluster.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Clusterid != rs.Primary.ID {
			return fmtp.Errorf("Cluster not found. ")
		}

		*clusterGet = *found

		return nil
	}
}

func testAccMrsClusterConfig_base(rName string) string {
	return fmt.Sprintf(`
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

func testAccMRSV1ClusterConfig_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_mrs_cluster" "cluster1" {
  cluster_name          = "%s"
  billing_type          = 12
  master_node_num       = 2
  core_node_num         = 3
  master_node_size      = "c6.4xlarge.4.linux.bigdata"
  core_node_size        = "c6.4xlarge.4.linux.bigdata"
  available_zone_id     = "effdcbc7d4d64a02aa1fa26b42f56533"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  cluster_version       = "MRS 3.1.0"
  volume_type           = "SAS"
  volume_size           = 600
  safe_mode             = 0
  cluster_type          = 0
  node_password         = "%s"
  cluster_admin_secret  = "%s"

  component_list {
    component_name = "Hadoop"
  }
  component_list {
    component_name = "ZooKeeper"
  }
  component_list {
    component_name = "Ranger"
  }
  component_list {
    component_name = "Spark2x"
  }
  component_list {
    component_name = "Hive"
  }
}`, testAccMrsClusterConfig_base(rName), rName, password, password)
}
