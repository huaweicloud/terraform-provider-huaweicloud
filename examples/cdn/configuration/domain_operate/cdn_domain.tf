# 华为云CDN加速域名创建示例
resource "huaweicloud_cdn_domain" "example" {
  name = "terraform-test.nanguapi.com"
  type = "web"
  sources {
    origin = "third-bucket-addr"
    origin_type = "third_bucket"
    http_port = 80
    https_port = 443
    bucket_access_key = "test-ak-12345"
    bucket_secret_key = "test-sk-12345"
    bucket_region = "北京四"
    bucket_name = "third-bucket"
    retrieval_host = "terraform-test-2.nanguapi.com"
  }
  service_area = "mainland_china"
}

data "huaweicloud_cdn_domains" "example_domain" {
  name = "terraform-test.nanguapi.com"
}

output "example_domain" {
  value = data.huaweicloud_cdn_domains.example_domain
}