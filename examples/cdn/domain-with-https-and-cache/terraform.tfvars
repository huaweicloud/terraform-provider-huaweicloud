# Domain Configuration
domain_name                   = "example.com"
domain_type                   = "web"
service_area                  = "outside_mainland_china"
origin_protocol               = "https"
origin_type                   = "domain"
origin_server                 = "hostaddress"
http_port                     = 80
https_port                    = 443
ipv6_enable                   = false
range_based_retrieval_enabled = false
domain_description            = "CDN domain for example.com"
# HTTPS Configuration
https_enabled                 = true
certificate_name              = "terraform_test_cert"
certificate_source            = "0"
certificate_body_path         = "/path/to/your/certificate.crt"
private_key_path              = "/path/to/your/private.key"
http2_enabled                 = true
ocsp_stapling_status          = "on"
# Cache Rules Configuration
cache_rules                   = [
  {
    rule_type          = "all"
    content            = ""
    ttl                = 2592000
    ttl_type           = "s"
    priority           = 1
    url_parameter_type = "full_url"
  },
  {
    rule_type          = "file_extension"
    content            = ".php;.jsp;.asp;.aspx"
    ttl                = 2592000
    ttl_type           = "s"
    priority           = 2
    url_parameter_type = "full_url"
  }
]
domain_tags                   = {
  Environment = "production"
  Project     = "cdn-example"
}
