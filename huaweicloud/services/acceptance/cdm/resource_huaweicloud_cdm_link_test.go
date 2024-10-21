package cdm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cdm/v1/link"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdm"
)

func getCdmLinkResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CdmV11Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CDM v1 client, err=%s", err)
	}
	clusterId, linkName, err := cdm.ParseLinkInfoFromId(state.Primary.ID)
	if err != nil {
		return nil, err
	}
	return link.Get(client, clusterId, linkName)
}

// Link to OBS
func TestAccResourceCdmLink_basic(t *testing.T) {
	var obj link.LinkCreateOpts
	resourceName := "huaweicloud_cdm_link.test"
	name := acceptance.RandomAccResourceName()
	bucketName := acceptance.RandomAccResourceNameWithDash()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdmLinkResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdmLinkResource_basic(name, bucketName, acceptance.HW_ACCESS_KEY,
					acceptance.HW_SECRET_KEY),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "connector", "obs-connector"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "config.port", "443"),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id", "huaweicloud_cdm_cluster.test", "id"),
				),
			},
			{
				Config: testAccCdmLinkResource_update(name, bucketName, acceptance.HW_ACCESS_KEY,
					acceptance.HW_SECRET_KEY),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "connector", "obs-connector"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "config.port", "80"),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id", "huaweicloud_cdm_cluster.test", "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

func testAccCdmLinkResource_basic(name, bucketName, ak, sk string) string {
	clusterConfig := testAccCdmCluster_basic(name)

	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_cdm_link" "test" {
  name       = "%s"
  connector  = "obs-connector"
  cluster_id = huaweicloud_cdm_cluster.test.id
  enabled    = true

  config = {
    "storageType" = "OBS"
    "server"      = trimprefix(huaweicloud_obs_bucket.test.bucket_domain_name, "${huaweicloud_obs_bucket.test.bucket}.")
    "port"        = "443"
    "properties"  = jsonencode(
      {
        connectionTimeout = "10000",
        socketTimeout     = "20000"
      }
    )
  }

  access_key = "%s"
  secret_key = "%s"
}
`, clusterConfig, bucketName, name, ak, sk)
}

func testAccCdmLinkResource_update(name, bucketName, ak, sk string) string {
	clusterConfig := testAccCdmCluster_basic(name)

	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_cdm_link" "test" {
  name       = "%s"
  connector  = "obs-connector"
  cluster_id = huaweicloud_cdm_cluster.test.id
  enabled    = true

  config = {
    "storageType" = "OBS"
    "server"      = trimprefix(huaweicloud_obs_bucket.test.bucket_domain_name, "${huaweicloud_obs_bucket.test.bucket}.")
    "port"        = "80"
    "properties"  = jsonencode(
      {
        connectionTimeout = "5000",
        socketTimeout     = "2000"
      }
    )
  }

  access_key = "%s"
  secret_key = "%s"
}
`, clusterConfig, bucketName, name, ak, sk)
}

// Link to DLI
func TestAccResourceCdmLink_DLI(t *testing.T) {
	var obj link.LinkCreateOpts
	resourceName := "huaweicloud_cdm_link.testDLI"
	name := acceptance.RandomAccResourceName()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdmLinkResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAKAndSK(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdmLinkResource_DLI(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "connector", "dli-connector"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id", "huaweicloud_cdm_cluster.test", "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

func testAccCdmLinkResource_DLI(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cdm_link" "testDLI" {
  name       = "%[2]s"
  connector  = "dli-connector"
  cluster_id = huaweicloud_cdm_cluster.test.id
  enabled    = true

  config = {
    "region"    = "%[3]s"
    "projectId" = "%[4]s"
  }

  access_key = "%[5]s"
  secret_key = "%[6]s"
}
`, testAccCdmCluster_basic(name), name, acceptance.HW_REGION_NAME, acceptance.HW_PROJECT_ID,
		acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}

// Link to CSS
func TestAccResourceCdmLink_CSS(t *testing.T) {
	var obj link.LinkCreateOpts
	resourceName := "huaweicloud_cdm_link.testCSS"
	name := acceptance.RandomAccResourceName()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdmLinkResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdmLinkResource_CSS(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "connector", "elasticsearch-connector"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_id", "huaweicloud_cdm_cluster.test", "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCdmLinkResource_CSS(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "rule" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  ports             = "9200"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "out_v4_all" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "egress"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"
  security_mode  = true
  https_enabled  = true
  password       = "Test@passw0rd"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
}

resource "huaweicloud_cdm_link" "testCSS" {
  name       = "%[2]s"
  connector  = "elasticsearch-connector"
  cluster_id = huaweicloud_cdm_cluster.test.id
  enabled    = true

  config = {
    "linkType" = "CSS"
    "safemode" = "true"
    "host"     = huaweicloud_css_cluster.test.endpoint
    "user"     = "admin"
  }

  password = "Test@passw0rd"

  lifecycle {
    ignore_changes = [
      config,
    ]
  }
}
`, testAccCdmCluster_basic(name), name)
}
