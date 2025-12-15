resource "huaweicloud_cdn_cache_refresh" "test" {
  count = length(var.refresh_file_urls) > 0 ? 1 : 0

  type                  = "file"
  urls                  = var.refresh_file_urls
  mode                  = "all"
  zh_url_encode         = var.zh_url_encode
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_cdn_cache_preheat" "test" {
  count = length(var.preheat_urls) > 0 ? 1 : 0

  urls                  = var.preheat_urls
  zh_url_encode         = var.zh_url_encode
  enterprise_project_id = var.enterprise_project_id

  depends_on = [
    huaweicloud_cdn_cache_refresh.test,
  ]
}
