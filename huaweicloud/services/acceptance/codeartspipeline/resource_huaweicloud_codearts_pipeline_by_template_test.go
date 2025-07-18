package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPipelineByTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_by_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelineByTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test"),
					resource.TestCheckResourceAttr(rName, "is_publish", "false"),
					resource.TestCheckResourceAttr(rName, "banned", "true"),
					resource.TestCheckResourceAttrSet(rName, "definition"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "sources.#"),
					resource.TestCheckResourceAttrSet(rName, "concurrency_control.0.concurrency_number"),
					resource.TestCheckResourceAttrSet(rName, "variables.#"),
				),
			},
			{
				Config: testPipelineByTemplate_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "test-update"),
					resource.TestCheckResourceAttr(rName, "is_publish", "false"),
					resource.TestCheckResourceAttr(rName, "banned", "false"),
					resource.TestCheckResourceAttrSet(rName, "definition"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "sources.#"),
					resource.TestCheckResourceAttr(rName, "concurrency_control.0.concurrency_number", "1"),
					resource.TestCheckResourceAttrSet(rName, "variables.#"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testPipelineImportState(rName),
				ImportStateVerifyIgnore: []string{"template_id"},
			},
		},
	})
}

func testPipelineByTemplate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

%[4]s

resource "huaweicloud_codearts_pipeline_tag" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[5]s"
  color      = "#0b81f6"
}
  
resource "huaweicloud_codearts_pipeline_group" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[5]s"
}

resource "huaweicloud_codearts_pipeline_by_template" "test" {
  project_id       = huaweicloud_codearts_project.test.id
  template_id      = huaweicloud_codearts_pipeline_template.test.id
  name             = "%[5]s"
  is_publish       = false
  banned           = true
  description      = "test"
  parameter_groups = [huaweicloud_codearts_pipeline_parameter_group.test.id]
  group_id         = huaweicloud_codearts_pipeline_group.test.id
  tags             = [huaweicloud_codearts_pipeline_tag.test.id]

  sources {
    type = "code"

    params {
      codehub_id     = huaweicloud_codearts_repository.test.id
      git_type       = "codehub"
      git_url        = huaweicloud_codearts_repository.test.https_url
      ssh_git_url    = huaweicloud_codearts_repository.test.ssh_url
      repo_name      = huaweicloud_codearts_repository.test.name
      default_branch = "master"
    }
  }
}
`, testProject_basic(name), testRepository_basic(name), testPipelineTemplate_basic(name), testPipeline_parameterGroup(name), name)
}

func testPipelineByTemplate_update(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

resource "huaweicloud_codearts_pipeline_by_template" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  template_id = huaweicloud_codearts_pipeline_template.test.id
  name        = "%[4]s"
  is_publish  = false
  banned      = false
  description = "test-update"

  sources {
    type = "code"

    params {
      codehub_id     = huaweicloud_codearts_repository.test.id
      git_type       = "codehub"
      git_url        = huaweicloud_codearts_repository.test.https_url
      ssh_git_url    = huaweicloud_codearts_repository.test.ssh_url
      repo_name      = huaweicloud_codearts_repository.test.name
      default_branch = "master"
    }
  }

  triggers {
    git_url        = huaweicloud_codearts_repository.test.https_url
    git_type       = "codehub"
    is_auto_commit = false
    repo_id        = huaweicloud_codearts_repository.test.id

    events {
      type   = "push"
      enable = false
    }

    events {
      type   = "merge_request"
      enable = true
    }

    events {
      type   = "tag_push"
      enable = true
    }
  }

  schedules {
    type          = "periodic"
    name          = "test_update"
    enable        = true
    days_of_week  = [2,3]
    time_zone     = "China Standard Time"
    start_time    = "18:51"
    end_time      = "19:27"
    interval_time = "3600"
    interval_unit = "s"
  }

  concurrency_control {
    concurrency_number = 1
    exceed_action      = "QUEUE"
    enable             = true
  }
}
`, testProject_basic(name), testRepository_basic(name), testPipelineTemplate_basic(name), name)
}

func TestAccPipelineByTemplate_createWithUpdate(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_by_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelineByTemplate_createWithUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test"),
					resource.TestCheckResourceAttr(rName, "is_publish", "false"),
					resource.TestCheckResourceAttr(rName, "banned", "true"),
					resource.TestCheckResourceAttrSet(rName, "definition"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "sources.#"),
					resource.TestCheckResourceAttr(rName, "concurrency_control.0.concurrency_number", "1"),
					resource.TestCheckResourceAttrSet(rName, "variables.#"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testPipelineImportState(rName),
				ImportStateVerifyIgnore: []string{"template_id"},
			},
		},
	})
}

func testPipelineByTemplate_createWithUpdate(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

resource "huaweicloud_codearts_pipeline_by_template" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  template_id = huaweicloud_codearts_pipeline_template.test.id
  name        = "%[4]s"
  is_publish  = false
  banned      = true
  description = "test"

  sources {
    type = "code"

    params {
      codehub_id     = huaweicloud_codearts_repository.test.id
      git_type       = "codehub"
      git_url        = huaweicloud_codearts_repository.test.https_url
      ssh_git_url    = huaweicloud_codearts_repository.test.ssh_url
      repo_name      = huaweicloud_codearts_repository.test.name
      default_branch = "master"
    }
  }

  concurrency_control {
    concurrency_number = 1
    exceed_action      = "QUEUE"
    enable             = true
  }
}
`, testProject_basic(name), testRepository_basic(name), testPipelineTemplate_basic(name), name)
}
