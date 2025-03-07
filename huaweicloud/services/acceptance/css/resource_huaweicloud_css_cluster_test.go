package css

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getCssClusterFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CSS v1 client: %s", err)
	}

	return cluster.Get(client, state.Primary.ID)
}

func TestAccCssCluster_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_basic(rName, "Test@passw0rd", 7, "bar"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "7.10.2"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "security_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "https_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.flavor", "ess.spec-4u8g"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "7"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "is_period", "false"),
					resource.TestCheckResourceAttr(resourceName, "backup_available", "true"),
					resource.TestCheckResourceAttr(resourceName, "disk_encrypted", "false"),
				),
			},
			{
				Config: testAccCssCluster_basic(rName, "Test@passw0rd.", 8, "bar_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
				),
			},
			{
				Config: testAccCssCluster_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "security_mode", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test_update", "id"),
				),
			},
			{
				Config: testAccCssCluster_basic(rName, "Test@passw0rd", 8, "bar_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "security_mode", "true"),
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

func TestAccCssCluster_opensearchBasic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_opensearchBasic(rName, "Test@passw0rd", 7, "bar"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "1.3.6"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "opensearch"),
					resource.TestCheckResourceAttr(resourceName, "security_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "https_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.flavor", "ess.spec-4u8g"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "7"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttr(resourceName, "is_period", "false"),
					resource.TestCheckResourceAttr(resourceName, "backup_available", "true"),
					resource.TestCheckResourceAttr(resourceName, "disk_encrypted", "false"),
				),
			},
			{
				Config: testAccCssCluster_opensearchBasic(rName, "Test@passw0rd.", 8, "bar_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
				),
			},
			{
				Config: testAccCssCluster_opensearchBasicUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "security_mode", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test_update", "id"),
				),
			},
			{
				Config: testAccCssCluster_opensearchBasic(rName, "Test@passw0rd", 8, "bar_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "security_mode", "true"),
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

func TestAccCssCluster_access(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{ // open public access on creating cluster
				Config: testAccCssCluster_access(rName, "Test@passw0rd", "116.204.111.47", 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "public_access.0.whitelist", "116.204.111.47"),
					resource.TestCheckResourceAttr(resourceName, "kibana_public_access.0.whitelist", "116.204.111.47"),
					resource.TestCheckResourceAttrSet(resourceName, "public_access.0.public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "kibana_public_access.0.public_ip"),
				),
			},
			{ // close public access
				Config: testAccCssCluster_accessClose(rName, "Test@passw0rd"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{ // reopen public access
				Config: testAccCssCluster_access(rName, "Test@passw0rd", "116.204.111.47", 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "public_access.0.whitelist", "116.204.111.47"),
					resource.TestCheckResourceAttr(resourceName, "kibana_public_access.0.whitelist", "116.204.111.47"),
					resource.TestCheckResourceAttrSet(resourceName, "public_access.0.public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "kibana_public_access.0.public_ip"),
				),
			},
			{ // update whitelist and bandwidth
				Config: testAccCssCluster_access(rName, "Test@passw0rd.", "116.204.111.47,121.37.117.211", 5),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "public_access.0.whitelist",
						"116.204.111.47,121.37.117.211"),
					resource.TestCheckResourceAttr(resourceName, "kibana_public_access.0.whitelist",
						"116.204.111.47,121.37.117.211"),
				),
			},
		},
	})
}

func TestAccCssCluster_localDisk(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSSLocalDiskFlavor(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_localDisk(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "security_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "https_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.flavor",
						acceptance.HW_CSS_LOCAL_DISK_FLAVOR),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.flavor",
						acceptance.HW_CSS_LOCAL_DISK_FLAVOR),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.instance_number", "1"),
				),
			},
			{
				Config: testAccCssCluster_localDisk(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.instance_number", "2"),
				),
			},
		},
	})
}

func TestAccCssCluster_prePaid(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_prePaid(rName, "Test@passw0rd", 1, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "security_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "https_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.flavor", "ess.spec-4u8g"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "vpcep_endpoint.0.endpoint_with_dns_name", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcep_endpoint_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcep_ip"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccCssCluster_prePaid(rName, "Test@passw0rd", 1, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func TestAccCssCluster_updateWithEpsId(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)
	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_withEpsId(rName, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccCssCluster_withEpsId(rName, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func TestAccCssCluster_flavor(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"
	flavor := "ess.spec-4u8g"
	updateFlavor := "ess.spec-4u16g"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_favor(rName, flavor),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.flavor", flavor),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.flavor", flavor),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.flavor", flavor),
					resource.TestCheckResourceAttr(resourceName, "client_node_config.0.flavor", flavor),
				),
			},
			{
				Config: testAccCssCluster_favor(rName, updateFlavor),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.flavor", updateFlavor),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.flavor", updateFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.flavor", updateFlavor),
					resource.TestCheckResourceAttr(resourceName, "client_node_config.0.flavor", updateFlavor),
				),
			},
		},
	})
}

func TestAccCssCluster_addMasterAndClientNode(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_extend_basic(rName, 3, 40),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccCssCluster_nodeChangeBasic(rName, 3, 3, 3, 1, 40),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.instance_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "client_node_config.0.instance_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_node_config.0.volume.0.size", "40"),
				),
			},
		},
	})
}

func TestAccCssCluster_extend(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_nodeChangeBasic(rName, 3, 3, 3, 1, 40),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.instance_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.instance_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "client_node_config.0.instance_number", "1"),
				),
			},
			{
				Config: testAccCssCluster_nodeChangeBasic(rName, 4, 4, 5, 3, 60),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.volume.0.size", "60"),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.instance_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.volume.0.size", "60"),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.instance_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "client_node_config.0.instance_number", "3"),
				),
			},
		},
	})
}

func TestAccCssCluster_shrink(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_nodeChangeBasic(rName, 4, 4, 5, 4, 40),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.instance_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.instance_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "client_node_config.0.instance_number", "4"),
				),
			},
			{
				Config: testAccCssCluster_nodeChangeBasic(rName, 3, 3, 3, 3, 60),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "cold_node_config.0.instance_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.instance_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "client_node_config.0.instance_number", "3"),
				),
			},
		},
	})
}

func TestAccCssCluster_prePaid_extend(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"
	flavor := "ess.spec-4u8g"
	updateFlavor := "ess.spec-4u16g"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCssClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_extend_prePaid(rName, flavor, 3, 3, 40),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.flavor", flavor),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.flavor", flavor),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.instance_number", "3"),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.volume.0.size", "40"),
				),
			},
			{
				Config: testAccCssCluster_extend_prePaid(rName, updateFlavor, 4, 5, 60),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.flavor", updateFlavor),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.instance_number", "4"),
					resource.TestCheckResourceAttr(resourceName, "ess_node_config.0.volume.0.size", "60"),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.flavor", updateFlavor),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.instance_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "master_node_config.0.volume.0.size", "40"),
				),
			},
		},
	})
}

func TestAccCssCluster_changeToPeriod(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_cluster.test"

	var obj cluster.ClusterDetailResponse
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getLogstashClusterFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_changeToPrepaidBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "is_period", "false"),
				),
			},
			{
				Config: testAccCssCluster_toPrePaid(rName, "false"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_period", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccCssCluster_toPrePaid(rName, "true"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "is_period", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
			{
				Config:      testAccCssCluster_toPostPaid(rName),
				ExpectError: regexp.MustCompile(`only support changing the CSS cluster form post-paid to pre-paid`),
			},
		},
	})
}

func testAccCssBase(rName string) string {
	bucketName := acceptance.RandomAccResourceNameWithDash()
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_obs_bucket" "cssObs" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true
}
`, common.TestBaseNetwork(rName), bucketName)
}

func testAccSecGroupUpdate(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test_update" {
  name                 = "%s_update"
  delete_default_rules = true
}
`, name)
}

func testAccCssCluster_basic(rName, pwd string, keepDays int, tag string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[3]s"
  engine_version = "7.10.2"
  security_mode  = true
  password       = "%[4]s"

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

  backup_strategy {
    keep_days   = %[5]d
    start_time  = "00:00 GMT+08:00"
    prefix      = "snapshot"
    bucket      = huaweicloud_obs_bucket.cssObs.bucket
    agency      = "css_obs_agency"
    backup_path = "css_repository/acctest"
  }

  tags = {
    foo = "%[6]s"
    key = "value"
  }

  lifecycle {
    ignore_changes = [
      ess_node_config.0.shrink_node_ids,
    ]
  }
}
`, testAccCssBase(rName), testAccSecGroupUpdate(rName), rName, pwd, keepDays, tag)
}

func testAccCssCluster_access(rName, pwd string, whiteList string, bwSize int) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[3]s"
  engine_version = "7.10.2"
  security_mode  = true
  https_enabled  = true
  password       = "%[4]s"

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

  backup_strategy {
    keep_days   = 7
    start_time  = "00:00 GMT+08:00"
    prefix      = "snapshot"
    bucket      = huaweicloud_obs_bucket.cssObs.bucket
    agency      = "css_obs_agency"
    backup_path = "css_repository/acctest"
  }

  public_access {
    bandwidth         = %[6]d
    whitelist_enabled = true
    whitelist         = "%[5]s"
  }

  kibana_public_access {
    bandwidth         = %[6]d
    whitelist_enabled = true
    whitelist         = "%[5]s"
  }
}
`, testAccCssBase(rName), testAccSecGroupUpdate(rName), rName, pwd, whiteList, bwSize)
}

func testAccCssCluster_accessClose(rName, pwd string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[3]s"
  engine_version = "7.10.2"
  security_mode  = true
  https_enabled  = true
  password       = "%[4]s"

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

  backup_strategy {
    keep_days   = 7
    start_time  = "00:00 GMT+08:00"
    prefix      = "snapshot"
    bucket      = huaweicloud_obs_bucket.cssObs.bucket
    agency      = "css_obs_agency"
    backup_path = "css_repository/acctest"
  }
}
`, testAccCssBase(rName), testAccSecGroupUpdate(rName), rName, pwd)
}

func testAccCssCluster_basic_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[3]s"
  engine_version = "7.10.2"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test_update.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  backup_strategy {
    keep_days   = 8
    start_time  = "00:00 GMT+08:00"
    prefix      = "snapshot"
    bucket      = huaweicloud_obs_bucket.cssObs.bucket
    agency      = "css_obs_agency"
    backup_path = "css_repository/acctest"
  }

  tags = {
    foo = "bar_update"
    key = "value"
  }
}
`, testAccCssBase(rName), testAccSecGroupUpdate(rName), rName)
}

func testAccCssCluster_localDisk(rName string, nodeNum int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"
  security_mode  = true
  password       = "Test@passw0rd"

  ess_node_config {
    flavor          = "%[4]s"
    instance_number = %[3]d
  }

  cold_node_config {
    flavor          = "%[4]s"
    instance_number = %[3]d
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
}
`, testAccCssBase(rName), rName, nodeNum, acceptance.HW_CSS_LOCAL_DISK_FLAVOR)
}

func testAccCssCluster_prePaid(rName, pwd string, period int, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"
  security_mode  = true
  password       = "%[3]s"

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = %[4]d
  auto_renew    = "%[5]v"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  vpcep_endpoint {
    endpoint_with_dns_name = true
  }
}
`, testAccCssBase(rName), rName, pwd, period, isAutoRenew)
}

func testAccCssCluster_withEpsId(rName string, epsId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"
  security_mode  = true
  password       = "Test@passw0rd"

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  enterprise_project_id = "%[3]s"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }
}
`, testAccCssBase(rName), rName, epsId)
}

func testAccCssCluster_extend_basic(rName string, essNodeNum, size int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = %[3]d
    volume {
      volume_type = "HIGH"
      size        = %[4]d
    }
  }
  cold_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = %[3]d
    volume {
      volume_type = "HIGH"
      size        = %[4]d
    }
  }
}
`, testAccCssBase(rName), rName, essNodeNum, size)
}

func testAccCssCluster_favor(rName, flavor string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  ess_node_config {
    flavor          = "%[3]s"
    instance_number = 3
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  cold_node_config {
    flavor          = "%[3]s"
    instance_number = 3
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  master_node_config {
    flavor          = "%[3]s"
    instance_number = 3
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  client_node_config {
    flavor          = "%[3]s"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }
}
`, testAccCssBase(rName), rName, flavor)
}

func testAccCssCluster_nodeChangeBasic(rName string, essNum, coldNum, masterNum, clientNum, size int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = %[3]d
    volume {
      volume_type = "HIGH"
      size        = %[7]d
    }
  }

  cold_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = %[4]d
    volume {
      volume_type = "HIGH"
      size        = %[7]d
    }
  }

  master_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = %[5]d
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  client_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = %[6]d
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }
}
`, testAccCssBase(rName), rName, essNum, coldNum, masterNum, clientNum, size)
}

func testAccCssCluster_extend_prePaid(rName, flavorNmae string, essNodeNum, masterNodeNum, size int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

  ess_node_config {
    flavor          = "%[3]s"
    instance_number = %[4]d
    volume {
      volume_type = "HIGH"
      size        = %[6]d
    }
  }

  master_node_config {
    flavor          = "%[3]s"
    instance_number = %[5]d
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }
}
`, testAccCssBase(rName), rName, flavorNmae, essNodeNum, masterNodeNum, size)
}

func testAccCssCluster_changeToPrepaidBasic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"

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
`, testAccCssBase(rName), rName)
}

func testAccCssCluster_toPrePaid(rName, autoRenew string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[3]s"

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
}
`, testAccCssBase(rName), rName, autoRenew)
}

func testAccCssCluster_toPostPaid(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[2]s"
  engine_version = "7.10.2"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  charging_mode = "postPaid"

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
}
`, testAccCssBase(rName), rName)
}

func testAccCssCluster_opensearchBasic(rName, pwd string, keepDays int, tag string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[3]s"
  engine_version = "1.3.6"
  engine_type    = "opensearch"
  security_mode  = true
  password       = "%[4]s"

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

  backup_strategy {
    keep_days   = %[5]d
    start_time  = "00:00 GMT+08:00"
    prefix      = "snapshot"
    bucket      = huaweicloud_obs_bucket.cssObs.bucket
    agency      = "css_obs_agency"
    backup_path = "css_repository/acctest"
  }

  tags = {
    foo = "%[6]s"
    key = "value"
  }

  lifecycle {
    ignore_changes = [
      ess_node_config.0.shrink_node_ids,
    ]
  }
}
`, testAccCssBase(rName), testAccSecGroupUpdate(rName), rName, pwd, keepDays, tag)
}

func testAccCssCluster_opensearchBasicUpdate(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_css_cluster" "test" {
  name           = "%[3]s"
  engine_version = "1.3.6"
  engine_type    = "opensearch"

  ess_node_config {
    flavor          = "ess.spec-4u8g"
    instance_number = 1
    volume {
      volume_type = "HIGH"
      size        = 40
    }
  }

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  security_group_id = huaweicloud_networking_secgroup.test_update.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  backup_strategy {
    keep_days   = 8
    start_time  = "00:00 GMT+08:00"
    prefix      = "snapshot"
    bucket      = huaweicloud_obs_bucket.cssObs.bucket
    agency      = "css_obs_agency"
    backup_path = "css_repository/acctest"
  }

  tags = {
    foo = "bar_update"
    key = "value"
  }
}
`, testAccCssBase(rName), testAccSecGroupUpdate(rName), rName)
}
