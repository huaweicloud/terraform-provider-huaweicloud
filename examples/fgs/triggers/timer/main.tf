resource "huaweicloud_fgs_function" "test" {
  name        = var.function_name
  app         = "default"
  handler     = "index.handler"
  memory_size = var.function_memory_size
  timeout     = var.function_timeout
  runtime     = var.function_runtime
  code_type   = "inline"
  func_code   = base64encode(var.function_code)
  description = var.function_description
}

resource "huaweicloud_fgs_function_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = var.trigger_status
  event_data   = jsonencode({
    "name"           = var.trigger_name
    "schedule_type"  = var.trigger_schedule_type
    "sync_execution" = var.trigger_sync_execution
    "user_event"     = var.trigger_user_event
    "schedule"       = var.trigger_schedule
  })
}
