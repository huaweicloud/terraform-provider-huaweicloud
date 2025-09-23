package dws

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

func getDwsExtDataSourceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDwsExtDataSource: Query the DWS external data source.
	var (
		getDwsExtDataSourceHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ext-data-sources"
		getDwsExtDataSourceProduct = "dws"
	)
	getDwsExtDataSourceClient, err := cfg.NewServiceClient(getDwsExtDataSourceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	getDwsExtDataSourcePath := getDwsExtDataSourceClient.Endpoint + getDwsExtDataSourceHttpUrl
	getDwsExtDataSourcePath = strings.ReplaceAll(getDwsExtDataSourcePath, "{project_id}", getDwsExtDataSourceClient.ProjectID)
	getDwsExtDataSourcePath = strings.ReplaceAll(getDwsExtDataSourcePath, "{cluster_id}",
		fmt.Sprintf("%v", state.Primary.Attributes["cluster_id"]))

	getDwsExtDataSourcePath += fmt.Sprintf("?type=%s", state.Primary.Attributes["type"])

	getDwsExtDataSourceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
		OkCodes: []int{
			200,
		},
	}
	getDwsExtDataSourceResp, err := getDwsExtDataSourceClient.Request("GET", getDwsExtDataSourcePath, &getDwsExtDataSourceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DWS external data source: %s", err)
	}

	getDwsExtDataSourceRespBody, err := utils.FlattenResponse(getDwsExtDataSourceResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DWS external data source: %s", err)
	}

	jsonPath := fmt.Sprintf("data_sources[?id=='%s']|[0]", state.Primary.ID)
	rawData := utils.PathSearch(jsonPath, getDwsExtDataSourceRespBody, nil)
	if rawData == nil {
		return nil, fmt.Errorf("error retrieving DWS external data source: %s", err)
	}

	return rawData, nil
}

func TestAccDwsExtDataSource_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dws_ext_data_source.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDwsExtDataSourceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDwsExtDataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttrPair(rName, "name", "huaweicloud_mapreduce_cluster.test", "name"),
					resource.TestCheckResourceAttr(rName, "type", "MRS"),
					resource.TestCheckResourceAttr(rName, "user_name", "admin"),
					resource.TestCheckResourceAttrPair(rName, "data_source_id", "huaweicloud_mapreduce_cluster.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttrSet(rName, "configure_status"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testDwsExtDataSourceImportState(rName),
				ImportStateVerifyIgnore: []string{"user_pwd", "reboot"},
			},
		},
	})
}

func testDwsExtDataSource_basic(name string) string {
	pwd := acceptance.RandomPassword()
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dws_ext_data_source" "test" {
  cluster_id     = "%[2]s"
  name           = huaweicloud_mapreduce_cluster.test.name
  type           = "MRS"
  data_source_id = huaweicloud_mapreduce_cluster.test.id
  user_name      = "admin"
  user_pwd       = "%[3]s"
  description    = "This is a demo"
}
`, testDwsExtDataSourceMrs(name, pwd), acceptance.HW_DWS_CLUSTER_ID, pwd)
}

func testDwsExtDataSourceMrs(rName, pwd string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

# For MRS data source, the DWS cluster and MRS cluster must be in the same VPC.
data "huaweicloud_dws_clusters" "test" {}

locals {
  dws_cluster_info = [for v in data.huaweicloud_dws_clusters.test.clusters : v if v.id == "%[1]s"]
}

resource "huaweicloud_mapreduce_cluster" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  name               = "%[2]s"
  type               = "ANALYSIS"
  version            = "MRS 1.9.2"
  manager_admin_pass = "%[3]s"
  node_admin_pass    = "%[3]s"
  subnet_id          = try(local.dws_cluster_info[0].subnet_id, "")
  vpc_id             = try(local.dws_cluster_info[0].vpc_id, "")
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
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 600
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  analysis_task_nodes {
    flavor            = "c6.2xlarge.4.linux.bigdata"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 600
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
}`, acceptance.HW_DWS_CLUSTER_ID, rName, pwd)
}

func getAgencyNames() (agencyName, updateAgencyName string) {
	agencyNames := strings.Split(acceptance.HW_DWS_OBS_AGENCY_NAMES, ",")
	if len(agencyNames) < 2 {
		return "", ""
	}
	return agencyNames[0], agencyNames[1]
}

func TestAccDwsExtDataSource_obs(t *testing.T) {
	var (
		obj                          interface{}
		name                         = acceptance.RandomAccResourceName()
		rName                        = "huaweicloud_dws_ext_data_source.test"
		agencyName, updateAgencyName = getAgencyNames()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDwsExtDataSourceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
			acceptance.TestAccPreCheckDwsExtDataSourceAgencyNames(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDwsExtDataSource_obs(name, agencyName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "OBS"),
					resource.TestCheckResourceAttr(rName, "user_name", agencyName),
					resource.TestCheckResourceAttr(rName, "connect_info", "gaussdb"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
				),
			},
			{
				Config: testDwsExtDataSource_obs(name, updateAgencyName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "user_name", updateAgencyName),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDwsExtDataSourceImportState(rName),
			},
		},
	})
}

func testDwsExtDataSource_obs(name, agencyName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_ext_data_source" "test" {
  cluster_id   = "%[1]s"
  name         = "%[2]s"
  type         = "OBS"
  user_name    = "%[3]s"
  connect_info = "gaussdb"
  description  = "Created by terraform script"
}
`, acceptance.HW_DWS_CLUSTER_ID, name, agencyName)
}

func testDwsExtDataSourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["cluster_id"] == "" {
			return "", fmt.Errorf("attribute (cluster_id) of resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["cluster_id"] + "/" +
			rs.Primary.ID, nil
	}
}
