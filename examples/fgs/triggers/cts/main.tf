resource "huaweicloud_fgs_function" "test" {
  name        = var.function_name
  app         = "default"
  handler     = "index.handler"
  agency      = var.function_agency_name
  memory_size = var.function_memory_size
  timeout     = var.function_timeout
  runtime     = var.function_runtime
  code_type   = "inline"
  func_code   = base64encode(var.function_code)
  description = var.function_description
}

resource "huaweicloud_fgs_function_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "CTS"
  status       = var.trigger_status
  event_data   = jsonencode({
    "name"       = var.trigger_name
    "operations" = var.trigger_operations
  })
}
