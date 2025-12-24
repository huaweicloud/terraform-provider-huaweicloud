# ST.001 Disable
resource "huaweicloud_obs_bucket" "source" {
  bucket        = var.source_bucket_name
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket" "target" {
  bucket        = var.target_bucket_name
  acl           = "private"
  force_destroy = true
}
# ST.001 Enable

data "huaweicloud_fgs_dependencies" "test" {
  type = "public"
  name = "pillow-7.1.2"
}

data "huaweicloud_fgs_dependency_versions" "test" {
  dependency_id = try(data.huaweicloud_fgs_dependencies.test.packages[0].id, "NOT_FOUND")
  version       = 1
}

resource "huaweicloud_fgs_function" "test" {
  name        = var.function_name
  agency      = var.function_agency_name
  app         = "default"
  handler     = "index.handler"
  memory_size = var.function_memory_size
  timeout     = var.function_timeout
  runtime     = var.function_runtime
  code_type   = "inline"
  func_code   = base64encode(var.function_code)
  description = var.function_description
  depend_list = data.huaweicloud_fgs_dependency_versions.test.versions[*].id

  user_data = jsonencode({
    "output_bucket" = var.target_bucket_name
    "obs_endpoint"  = format("obs.%s.myhuaweicloud.com", var.region_name)
  })
}

data "huaweicloud_eg_event_channels" "test" {
  provider_type = "OFFICIAL"
  name          = "default"
}

resource "huaweicloud_fgs_function_trigger" "test" {
  depends_on = [huaweicloud_obs_bucket.source]

  function_urn                   = huaweicloud_fgs_function.test.urn
  type                           = "EVENTGRID"
  cascade_delete_eg_subscription = true
  status                         = var.trigger_status
  event_data                     = jsonencode({
    "channel_id"   = try(data.huaweicloud_eg_event_channels.test.channels[0].id, "")
    "channel_name" = try(data.huaweicloud_eg_event_channels.test.channels[0].name, "")
    "source_name"  = "HC.OBS.DWR"
    "trigger_name" = var.trigger_name_suffix
    "agency"       = var.trigger_agency_name
    "bucket"       = var.source_bucket_name
    "event_types"  = ["OBS:DWR:ObjectCreated:PUT", "OBS:DWR:ObjectCreated:POST"]
    "Key_encode"   = true
  })

  lifecycle {
    ignore_changes = [
      event_data,
    ]
  }
}
