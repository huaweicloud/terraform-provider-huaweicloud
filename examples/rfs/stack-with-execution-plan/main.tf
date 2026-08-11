resource "huaweicloud_rfs_stack" "test" {
  name        = var.stack_name
  description = var.stack_description
}

resource "huaweicloud_rfs_execution_plan_v2" "test" {
  stack_name          = huaweicloud_rfs_stack.test.name
  stack_id            = huaweicloud_rfs_stack.test.id
  execution_plan_name = var.execution_plan_name
  description         = var.execution_plan_description
  template_body       = var.execution_plan_template_body
}
