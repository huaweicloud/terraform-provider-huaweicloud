package codeartsbuild

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBuildTaskAction_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testBuildTaskAction_basic(name),
			},
		},
	})
}

func testBuildTaskAction_base(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_codearts_build_task" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[3]s"
  arch       = "x86-64"

  scms {
    url      = huaweicloud_codearts_repository.test.ssh_url
    scm_type = "codehub"
    web_url  = huaweicloud_codearts_repository.test.web_url
    repo_id  = huaweicloud_codearts_repository.test.repository_id
    branch   = "master"
  }

  parameters {
    name = "hudson.model.StringParameterDefinition"

    params {
      name  = "name"
      value = "test"
    }
    params {
      name  = "type"
      value = "normalparam"
    }
    params {
      name  = "defaultValue"
      value = "demo"
    }
    params {
      name  = "staticVar"
      value = "false"
    }
    params {
      name  = "sensitiveVar"
      value = "false"
    }
    params {
      name  = "deletion"
      value = "false"
    }
    params {
      name  = "defaults"
      value = "false"
    }
  }

  steps {
    enable     = true
    module_id  = "devcloud2018.codeci_action_20017.action"
    name       = "Execute Shell Command"
    properties = {
      command = jsonencode("sleep 600")
      image   = jsonencode("shell4.2.46-git1.8.3-zip6.00")
    }
  }
}
`, testProject_basic(name), testRepository_basic(name), name)
}

func testBuildTaskAction_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_build_task_action" "execute" {
  job_id = huaweicloud_codearts_build_task.test.id
  action = "execute"

  parameter {
    name  = "test"
    value = "set"
  }

  scm {
    build_tag = "test"
  }
}

resource "huaweicloud_codearts_build_task_action" "stop" {
  job_id   = huaweicloud_codearts_build_task.test.id
  action   = "stop"
  build_no = huaweicloud_codearts_build_task_action.execute.build_no
}
`, testBuildTaskAction_base(name))
}
