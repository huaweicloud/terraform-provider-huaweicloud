---
subcategory: "Advanced Anti-DDoS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aad_domain"
description: |-
  Manages an AAD domain resource within HuaweiCloud.
---

# huaweicloud_aad_domain

Manages an AAD domain resource within HuaweiCloud.

-> If the WAF cname of the origin server you want to add shares an IP address and port with the WAF cname of another
  protected domain, **passby** the WAF at this time, it will affect the protection of all related domains.

-> One user can only create one protected domain within 3 seconds. If multiple protected domains need to be created,
  please refer to the example **Create multiple domain**.

## Example Usage

### Create one domain

```hcl
variable "domain_name" {}
variable "enterprise_project_id" {}
variable "real_server_type" {}
variable "real_server" {}

variable "instance_ips" {
  type = list(string)
}
variable "port_http" {
  type = list(int)
}
variable "port_https" {
  type = list(int)
}

resource "huaweicloud_aad_domain" "test" {
  domain_name           = var.domain_name
  enterprise_project_id = var.enterprise_project_id
  real_server_type      = var.real_server_type
  real_server           = var.real_server
  vips                  = var.instance_ips
  port_http             = var.port_http
  port_https            = var.port_https
}
```

### Create multiple domain

```hcl
variable "domain_name_0" {}
variable "domain_name_1" {}
variable "enterprise_project_id" {}
variable "real_server_type" {}
variable "real_server" {}

variable "instance_ips" {
  type = list(string)
}
variable "port_http" {
  type = list(int)
}
variable "port_https" {
  type = list(int)
}

resource "huaweicloud_aad_domain" "test0" {
  domain_name           = var.domain_name_0
  enterprise_project_id = var.enterprise_project_id
  real_server_type      = var.real_server_type
  real_server           = var.real_server
  vips                  = var.instance_ips
  port_http             = var.port_http
  port_https            = var.port_https
}

resource "null_resource" "wait_3_seconds"{
  provisinoer "local-exec" {
    command = "sleep 3"
  }

  depends_on = [huaweicloud_aad_domain.test0]
}

resource "huaweicloud_aad_domain" "test1" {
  domain_name           = var.domain_name_1
  enterprise_project_id = var.enterprise_project_id
  real_server_type      = var.real_server_type
  real_server           = var.real_server
  vips                  = var.instance_ips
  port_http             = var.port_http
  port_https            = var.port_https

  depends_on = [null_resource.wait_3_seconds]
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String, NonUpdatable) Specifies the domain name to be protected by AAD instance.
  The domain name must be put on record otherwise it can not be used.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the enterprise project ID. The enterprise project
  ID should be consistent with the enterprise project to which the AAD instance belongs. The enterprise project ID can be
  viewed in Huawei Cloud EPS service, and the default enterprise project ID is **0**.

* `real_server_type` - (Required, Int) Specifies the origin server type.
  The valid values are as follows:  
  `0`: Indicates the type is IP address.  
  `1`: Indicates the type is domain.

* `real_server` - (Required, String) Specifies the value of the origin server.
  When the `real_server_type` is set to `0`, there can be maximum of `20` IP addresses, using commas(,) to separate multiple
  IP addresses. Each IP address is unique and invalid IP addresses are as follows:
  `127.0.0.1`, `172.16.*.*`, `192.168.*.*`, `10.0~255.*.*`.  
  When the `real_server_type` is set to `1`, you can enter a domain such as `www.domain.com`.
  For multiple second-level domains, enter `*.domain.com.`.

* `vips` - (Optional, List, NonUpdatable) Specifies the list of AAD instance IP addresses. Defense instance IP address must
  belong to the same enterprise project.

* `instance_ids` - (Optional, List, NonUpdatable) Specifies the list of AAD instance IDs.

-> Exactly one of `vips` or `instance_ids` must be set.

* `port_http` - (Optional, List, NonUpdatable) Specifies the port when forwarding protocol is HTTP.
  The valid values are as follows:  
  `80`, `81`, `82`, `83`, `84`, `85`, `88`, `133`, `134`, `140`, `141`, `144`, `151`, `881`,<br>
  `1,135`, `1,139`, `7,000`, `7,001`, `8,001`, `8,006`, `8,078`, `8,080`, `8,087`, `8,088`,<br>
   `8,089`, `8,090`, `8,093`, `8,097`, `8,100`, `8,182`, `8,200`, `8,813`, `8,814`, `8,888`,<br>
   `9,000`, `9,001`, `9,002`, `9,003`, `18,080`, `19,101`, `19,501`, `21,028`, `40,010`.

* `port_https` - (Optional, List, NonUpdatable) Specifies the port when forwarding protocol is HTTPS.
  The valid values are as follows:  
  `443`, `882`, `1,818`, `4,006`, `4,430`, `4,443`, `5,443`, `6,443`, `7,443`, `8,033`, `8,081`,<br>
  `8,082`, `8,083`, `8,443`, `8,445`, `8,553`, `8,663`, `8,750`, `8,804`, `8,805`, `7,443`,<br>
  `9,999`, `13,080`, `14,443`, `18,000`, `18,443`, `18,980`, `20,000`, `28,443`, `30,001`,<br>
  `30,003`, `30,004`, `30,005`.

-> Exactly one of `port_http` or `port_https` must be set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is also the domain ID.

* `cname` - The cname of domain.

* `protocol` - The protocol of the domain.

* `waf_status` - The protect status of WAF server.

## Import

The AAD domain can be imported using the `id`, e.g.

```bash
terraform import huaweicloud_aad_domain.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `vips`, `port_http` and `port_https`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_aad_domain" "test" {
  ...

  lifecycle {
    ignore_changes = [
      vips, port_http, port_https
    ]
  }
}
```
