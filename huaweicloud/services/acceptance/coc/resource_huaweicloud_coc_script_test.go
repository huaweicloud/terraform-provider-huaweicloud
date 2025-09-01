package coc

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

func getScriptResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	getScriptHttpUrl := "v1/job/scripts/{id}"
	getScriptPath := client.Endpoint + getScriptHttpUrl
	getScriptPath = strings.ReplaceAll(getScriptPath, "{id}", state.Primary.ID)

	getScriptOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getScriptResp, err := client.Request("GET", getScriptPath, &getScriptOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving COC script: %s", err)
	}

	getScriptRespBody, err := utils.FlattenResponse(getScriptResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving COC script: %s", err)
	}

	return getScriptRespBody, nil
}

func TestAccScript_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_script.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getScriptResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesScript_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "a demo script"),
					resource.TestCheckResourceAttr(resourceName, "risk_level", "LOW"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "parameters.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: tesScript_updated(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "a new demo script"),
					resource.TestCheckResourceAttr(resourceName, "risk_level", "MEDIUM"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttr(resourceName, "parameters.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
		},
	})
}

func TestAccScript_reviewers(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_script.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getScriptResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesScript_reviewers(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "a demo script"),
					resource.TestCheckResourceAttr(resourceName, "risk_level", "LOW"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "parameters.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: tesScript_reviewers_updated(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "a new demo script"),
					resource.TestCheckResourceAttr(resourceName, "risk_level", "MEDIUM"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.1"),
					resource.TestCheckResourceAttr(resourceName, "parameters.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
		},
	})
}

func tesScript_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_script" "test" {
  name                  = "%s"
  description           = "a demo script"
  risk_level            = "LOW"
  version               = "1.0.0"
  type                  = "SHELL"
  enterprise_project_id = "0"

  content = <<EOF
#! /bin/bash
echo "hello $${name}!"
EOF

  parameters {
    name        = "name"
    value       = "world"
    description = "the first parameter"
  }
}`, name)
}

func tesScript_updated(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_script" "test" {
  name                  = "%s"
  description           = "a new demo script"
  risk_level            = "MEDIUM"
  version               = "1.0.1"
  type                  = "SHELL"
  enterprise_project_id = "0"

  content = <<EOF
#! /bin/bash
echo "hello $${name}@$${company}!"
EOF

  parameters {
    name        = "name"
    value       = "world"
    description = "the first parameter"
  }
  parameters {
    name        = "company"
    value       = "Huawei"
    description = "the second parameter"
    sensitive   = true
  }
}`, name)
}

func tesScript_reviewers(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_script" "test" {
  name                  = "%s"
  description           = "a demo script"
  risk_level            = "LOW"
  version               = "1.0.0"
  type                  = "SHELL"
  enterprise_project_id = "0"

  content = <<EOF
#! /bin/bash
echo "hello $${name}!"
EOF

  parameters {
    name        = "name"
    value       = "world"
    description = "the first parameter"
  }
}`, name)
}

func tesScript_reviewers_updated(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_script" "test" {
  name                  = "%[1]s"
  description           = "a new demo script"
  risk_level            = "MEDIUM"
  version               = "1.0.1"
  type                  = "SHELL"
  enterprise_project_id = "0"

  content = <<EOF
#! /bin/bash
echo "hello $${name}@$${company}!"
EOF

  protocol = "DEFAULT"
  reviewers {
    reviewer_id   = "%[2]s"
    reviewer_name = "%[3]s"
  }

  parameters {
    name        = "name"
    value       = "world"
    description = "the first parameter"
  }
  parameters {
    name        = "company"
    value       = "Huawei"
    description = "the second parameter"
    sensitive   = true
  }
}`, name, acceptance.HW_USER_ID, acceptance.HW_USER_NAME)
}
