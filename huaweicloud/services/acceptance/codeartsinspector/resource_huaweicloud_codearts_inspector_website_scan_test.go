package codeartsinspector

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getInspectorWebsiteScanResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v3/{project_id}/webscan/tasks"
		product = "vss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts inspector client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += fmt.Sprintf("?task_id=%s", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CodeArts inspector website scan: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	taskStatus := utils.PathSearch("task_status", getRespBody, "")
	if taskStatus == "" {
		return nil, fmt.Errorf("error retrieving CodeArts inspector website scan: field `task_status` is not found" +
			" in detail API response")
	}

	if taskStatus == "canceled" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccInspectorWebsiteScan_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_inspector_website_scan.test"
	// The normal trigger time needs to be after the current time.
	timer := utils.FormatTimeStampUTC(time.Now().Add(48 * time.Hour).Unix())

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInspectorWebsiteScanResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInspectorWebsiteScan_basic(name, timer),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", name),
					resource.TestCheckResourceAttr(rName, "task_type", "normal"),
					resource.TestCheckResourceAttrPair(rName, "url",
						"huaweicloud_codearts_inspector_website.test", "website_address"),
					resource.TestCheckResourceAttr(rName, "timer", timer),
					resource.TestCheckResourceAttr(rName, "scan_mode", "deep"),
					resource.TestCheckResourceAttr(rName, "port_scan", "true"),
					resource.TestCheckResourceAttr(rName, "weak_pwd_scan", "false"),
					resource.TestCheckResourceAttr(rName, "cve_check", "false"),
					resource.TestCheckResourceAttr(rName, "text_check", "false"),
					resource.TestCheckResourceAttr(rName, "picture_check", "false"),
					resource.TestCheckResourceAttr(rName, "malicious_code", "false"),
					resource.TestCheckResourceAttr(rName, "malicious_link", "false"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "task_status"),
					resource.TestCheckResourceAttrSet(rName, "progress"),
					resource.TestCheckResourceAttrSet(rName, "reason"),
					resource.TestCheckResourceAttrSet(rName, "pack_num"),
					resource.TestCheckResourceAttrSet(rName, "score"),
					resource.TestCheckResourceAttrSet(rName, "safe_level"),
					resource.TestCheckResourceAttrSet(rName, "high"),
					resource.TestCheckResourceAttrSet(rName, "middle"),
					resource.TestCheckResourceAttrSet(rName, "low"),
					resource.TestCheckResourceAttrSet(rName, "hint"),
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

func TestAccInspectorWebsiteScan_monitor(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_inspector_website_scan.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInspectorWebsiteScanResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInspectorWebsiteScan_monitor(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", name),
					resource.TestCheckResourceAttr(rName, "task_type", "monitor"),
					resource.TestCheckResourceAttrPair(rName, "url",
						"huaweicloud_codearts_inspector_website.test", "website_address"),
					resource.TestCheckResourceAttr(rName, "trigger_time", "2023-11-26 16:10:53"),
					resource.TestCheckResourceAttr(rName, "task_period", "everyweek"),
					resource.TestCheckResourceAttr(rName, "scan_mode", "normal"),
					resource.TestCheckResourceAttr(rName, "port_scan", "true"),
					resource.TestCheckResourceAttr(rName, "weak_pwd_scan", "false"),
					resource.TestCheckResourceAttr(rName, "cve_check", "false"),
					resource.TestCheckResourceAttr(rName, "text_check", "false"),
					resource.TestCheckResourceAttr(rName, "picture_check", "false"),
					resource.TestCheckResourceAttr(rName, "malicious_code", "false"),
					resource.TestCheckResourceAttr(rName, "malicious_link", "false"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "task_status"),
					resource.TestCheckResourceAttrSet(rName, "schedule_status"),
					resource.TestCheckResourceAttrSet(rName, "progress"),
					resource.TestCheckResourceAttrSet(rName, "reason"),
					resource.TestCheckResourceAttrSet(rName, "pack_num"),
					resource.TestCheckResourceAttrSet(rName, "score"),
					resource.TestCheckResourceAttrSet(rName, "safe_level"),
					resource.TestCheckResourceAttrSet(rName, "high"),
					resource.TestCheckResourceAttrSet(rName, "middle"),
					resource.TestCheckResourceAttrSet(rName, "low"),
					resource.TestCheckResourceAttrSet(rName, "hint"),
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

// Before testing this test, upgrade the vulnerability management service version to the professional version.
func TestAccInspectorWebsiteScan_professional(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_inspector_website_scan.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getInspectorWebsiteScanResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCodeArtsEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testInspectorWebsiteScan_professional(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", name),
					resource.TestCheckResourceAttr(rName, "task_type", "monitor"),
					resource.TestCheckResourceAttrPair(rName, "url",
						"huaweicloud_codearts_inspector_website.test", "website_address"),
					resource.TestCheckResourceAttr(rName, "trigger_time", "2023-11-26 16:10:53"),
					resource.TestCheckResourceAttr(rName, "task_period", "everyweek"),
					resource.TestCheckResourceAttr(rName, "scan_mode", "normal"),
					resource.TestCheckResourceAttr(rName, "port_scan", "true"),
					resource.TestCheckResourceAttr(rName, "weak_pwd_scan", "true"),
					resource.TestCheckResourceAttr(rName, "cve_check", "true"),
					resource.TestCheckResourceAttr(rName, "text_check", "true"),
					resource.TestCheckResourceAttr(rName, "picture_check", "false"),
					resource.TestCheckResourceAttr(rName, "malicious_code", "false"),
					resource.TestCheckResourceAttr(rName, "malicious_link", "false"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "task_status"),
					resource.TestCheckResourceAttrSet(rName, "schedule_status"),
					resource.TestCheckResourceAttrSet(rName, "progress"),
					resource.TestCheckResourceAttrSet(rName, "reason"),
					resource.TestCheckResourceAttrSet(rName, "pack_num"),
					resource.TestCheckResourceAttrSet(rName, "score"),
					resource.TestCheckResourceAttrSet(rName, "safe_level"),
					resource.TestCheckResourceAttrSet(rName, "high"),
					resource.TestCheckResourceAttrSet(rName, "middle"),
					resource.TestCheckResourceAttrSet(rName, "low"),
					resource.TestCheckResourceAttrSet(rName, "hint"),
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

func testInspectorWebsiteScan_basic(name, timer string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_inspector_website_scan" "test" {
  task_name = "%s"
  task_type = "normal"
  url       = huaweicloud_codearts_inspector_website.test.website_address
  timer     = "%s"
  scan_mode = "deep"
  port_scan = true
}
`, testInspectorWebsite_basic(name), name, timer)
}

func testInspectorWebsiteScan_monitor(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_inspector_website_scan" "test" {
  task_name    = "%s"
  task_type    = "monitor"
  url          = huaweicloud_codearts_inspector_website.test.website_address
  trigger_time = "2023-11-26 16:10:53"
  task_period  = "everyweek"
  scan_mode    = "normal"
  port_scan    = true
}
`, testInspectorWebsite_basic(name), name)
}

func testInspectorWebsiteScan_professional(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_inspector_website_scan" "test" {
  task_name     = "%s"
  task_type     = "monitor"
  url           = huaweicloud_codearts_inspector_website.test.website_address
  trigger_time  = "2023-11-26 16:10:53"
  task_period   = "everyweek"
  scan_mode     = "normal"
  port_scan     = true
  weak_pwd_scan = true
  cve_check     = true
  text_check    = true
}
`, testInspectorWebsite_basic(name), name)
}
