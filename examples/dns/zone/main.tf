# Create a DNS public zone
resource "huaweicloud_dns_zone" "test" {
  name        = var.name
  description = var.description
  ttl         = var.ttl
  dnssec      = var.dnssec
}
