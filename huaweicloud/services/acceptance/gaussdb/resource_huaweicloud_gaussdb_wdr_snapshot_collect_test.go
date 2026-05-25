package gaussdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGaussDbWdrSnapshotCollectResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/wdr-snapshots"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])

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
	var getRespBody any
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, err
	}

	snapshot := utils.PathSearch(fmt.Sprintf("wdr_snapshots[?job_id=='%s']|[0]", state.Primary.ID), getRespBody, nil)
	if snapshot == nil {
		return nil, errors.New("error retrieving GaussDB WDR snapshot collect")
	}

	return snapshot, nil
}

func TestAccGaussDbWdrSnapshotCollect_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_gaussdb_wdr_snapshot_collect.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDbWdrSnapshotCollectResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
			acceptance.TestAccPreCheckGaussDBTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBWdrSnapshotCollect_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_GAUSSDB_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "start_time", acceptance.HW_GAUSSDB_START_TIME),
					resource.TestCheckResourceAttr(rName, "end_time", acceptance.HW_GAUSSDB_END_TIME),
					resource.TestCheckResourceAttr(rName, "wdr_type_attr", "cluster"),
					resource.TestCheckResourceAttrSet(rName, "file_size"),
					resource.TestCheckResourceAttrSet(rName, "job_create_time"),
					resource.TestCheckResourceAttrSet(rName, "download_url"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "notes"),
					resource.TestCheckResourceAttrSet(rName, "file_name"),
					resource.TestCheckResourceAttrSet(rName, "file_path"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.#"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.name"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.type"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.url"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.port"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.domain_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccGaussDbWdrSnapshotCollectImportStateFunc(rName),
				ImportStateVerifyIgnore: []string{
					"wdr_type",
					"download_url",
				},
			},
		},
	})
}

func TestAccGaussDbWdrSnapshotCollect_nodes(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_gaussdb_wdr_snapshot_collect.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDbWdrSnapshotCollectResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBWdrSnapshotCollect_nodes(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_GAUSSDB_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "start_time", acceptance.HW_GAUSSDB_START_TIME),
					resource.TestCheckResourceAttr(rName, "end_time", acceptance.HW_GAUSSDB_END_TIME),
					resource.TestCheckResourceAttr(rName, "wdr_type_attr", "component"),
					resource.TestCheckResourceAttrSet(rName, "file_size"),
					resource.TestCheckResourceAttrSet(rName, "wdr_type_attr"),
					resource.TestCheckResourceAttrSet(rName, "job_create_time"),
					resource.TestCheckResourceAttrSet(rName, "download_url"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "notes"),
					resource.TestCheckResourceAttrSet(rName, "file_name"),
					resource.TestCheckResourceAttrSet(rName, "file_path"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.#"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.name"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.type"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.url"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.port"),
					resource.TestCheckResourceAttrSet(rName, "obs_bucket.0.domain_id"),
				),
			},
		},
	})
}

func testAccGaussDBWdrSnapshotCollect_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_wdr_snapshot_collect" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  wdr_type    = ["cluster"]
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, acceptance.HW_GAUSSDB_START_TIME, acceptance.HW_GAUSSDB_END_TIME)
}

func testAccGaussDBWdrSnapshotCollect_nodes() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_nodes" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_gaussdb_wdr_snapshot_collect" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  wdr_type    = [
    data.huaweicloud_gaussdb_instance_nodes.test.nodes[0].components[0].id
  ]
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, acceptance.HW_GAUSSDB_START_TIME, acceptance.HW_GAUSSDB_END_TIME)
}

func testAccGaussDbWdrSnapshotCollectImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("attribute (instance_id) of Resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}
