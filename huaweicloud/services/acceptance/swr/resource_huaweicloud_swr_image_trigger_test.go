package swr

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

func getSwrImageTriggerResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getSwrImageTrigger: Query SWR image trigger
	var (
		getSwrImageTriggerHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/triggers/{trigger}"
		getSwrImageTriggerProduct = "swr"
	)
	getSwrImageTriggerClient, err := cfg.NewServiceClient(getSwrImageTriggerProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	organization := state.Primary.Attributes["organization"]
	repository := strings.ReplaceAll(state.Primary.Attributes["repository"], "/", "$")
	trigger := state.Primary.Attributes["name"]

	getSwrImageTriggerPath := getSwrImageTriggerClient.Endpoint + getSwrImageTriggerHttpUrl
	getSwrImageTriggerPath = strings.ReplaceAll(getSwrImageTriggerPath, "{namespace}", organization)
	getSwrImageTriggerPath = strings.ReplaceAll(getSwrImageTriggerPath, "{repository}", repository)
	getSwrImageTriggerPath = strings.ReplaceAll(getSwrImageTriggerPath, "{trigger}", trigger)

	getSwrImageTriggerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSwrImageTriggerResp, err := getSwrImageTriggerClient.Request("GET",
		getSwrImageTriggerPath, &getSwrImageTriggerOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SWR image trigger: %s", err)
	}
	return utils.FlattenResponse(getSwrImageTriggerResp)
}

func TestAccSwrImageTrigger_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_swr_image_trigger.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSwrImageTriggerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkloadType(t)
			acceptance.TestAccPreCheckWorkloadName(t)
			acceptance.TestAccPreCheckCceClusterId(t)
			acceptance.TestAccPreCheckWorkloadNameSpace(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSwrImageTrigger_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"huaweicloud_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"huaweicloud_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "workload_type", acceptance.HW_WORKLOAD_TYPE),
					resource.TestCheckResourceAttr(rName, "workload_name", acceptance.HW_WORKLOAD_NAME),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "namespace", acceptance.HW_WORKLOAD_NAMESPACE),
					resource.TestCheckResourceAttr(rName, "condition_value", ".*"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "cce"),
					resource.TestCheckResourceAttr(rName, "condition_type", "all"),
				),
			},
			{
				Config: testSwrImageTrigger_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"huaweicloud_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"huaweicloud_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "workload_type", acceptance.HW_WORKLOAD_TYPE),
					resource.TestCheckResourceAttr(rName, "workload_name", acceptance.HW_WORKLOAD_NAME),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "namespace", acceptance.HW_WORKLOAD_NAMESPACE),
					resource.TestCheckResourceAttr(rName, "condition_value", ".*"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "cce"),
					resource.TestCheckResourceAttr(rName, "condition_type", "all"),
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

func TestAccSwrImageTrigger_cci(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_swr_image_trigger.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSwrImageTriggerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkloadType(t)
			acceptance.TestAccPreCheckWorkloadName(t)
			acceptance.TestAccPreCheckWorkloadNameSpace(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSwrImageTrigger_cci(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "organization",
						"huaweicloud_swr_organization.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "repository",
						"huaweicloud_swr_repository.test", "name"),
					resource.TestCheckResourceAttr(rName, "workload_type", acceptance.HW_WORKLOAD_TYPE),
					resource.TestCheckResourceAttr(rName, "workload_name", acceptance.HW_WORKLOAD_NAME),
					resource.TestCheckResourceAttr(rName, "namespace", acceptance.HW_WORKLOAD_NAMESPACE),
					resource.TestCheckResourceAttr(rName, "condition_value", ".*"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "cci"),
					resource.TestCheckResourceAttr(rName, "condition_type", "all"),
				),
			},
		},
	})
}

func testSwrImageTrigger_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_image_trigger" "test" {
  organization    = huaweicloud_swr_organization.test.name
  repository      = huaweicloud_swr_repository.test.name
  workload_type   = "%[2]s"
  workload_name   = "%[3]s"
  cluster_id      = "%[4]s"
  namespace       = "%[5]s"
  condition_value = ".*"
  name            = "%[6]s"
  type            = "cce"
  condition_type  = "all"
}
`, testAccSWRRepository_basic(name), acceptance.HW_WORKLOAD_TYPE, acceptance.HW_WORKLOAD_NAME,
		acceptance.HW_CCE_CLUSTER_ID, acceptance.HW_WORKLOAD_NAMESPACE, name)
}

func testSwrImageTrigger_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_image_trigger" "test" {
  organization    = huaweicloud_swr_organization.test.name
  repository      = huaweicloud_swr_repository.test.name
  workload_type   = "%[2]s"
  workload_name   = "%[3]s"
  cluster_id      = "%[4]s"
  namespace       = "%[5]s"
  condition_value = ".*"
  enabled         = "false"
  name            = "%[6]s"
  type            = "cce"
  condition_type  = "all"
}
`, testAccSWRRepository_basic(name), acceptance.HW_WORKLOAD_TYPE, acceptance.HW_WORKLOAD_NAME,
		acceptance.HW_CCE_CLUSTER_ID, acceptance.HW_WORKLOAD_NAMESPACE, name)
}

func testSwrImageTrigger_cci(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_image_trigger" "test" {
  organization    = huaweicloud_swr_organization.test.name
  repository      = huaweicloud_swr_repository.test.name
  workload_type   = "%[2]s"
  workload_name   = "%[3]s"
  namespace       = "%[4]s"
  condition_value = ".*"
  name            = "%[5]s"
  type            = "cci"
  condition_type  = "all"
}
`, testAccSWRRepository_basic(name), acceptance.HW_WORKLOAD_TYPE, acceptance.HW_WORKLOAD_NAME,
		acceptance.HW_WORKLOAD_NAMESPACE, name)
}
