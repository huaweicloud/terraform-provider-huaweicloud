---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_domain_associate_certificate"
description: |-
  Manages a resource to associate the certificate to the domain within HuaweiCloud.
---

# huaweicloud_waf_domain_associate_certificate

Manages a resource to associate the certificate to the domain within HuaweiCloud.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> 1. The current resource is a one-time resource, and destroying this resource will not change the current status.
<br/>2. If you use this resource associate a certificate to a domain, the opration may trigger changes to the existing
resource: `huaweicloud_waf_domain` or `huaweicloud_waf_dedicated_domain`.
<br/>3. The resource only can used when the domain `client_protocol` is set to **HTTPS**.

## Example Usage

```hcl
variable "certificate_id" {}
variable "cloud_host_ids" {
  type = list(string)
}

resource "huaweicloud_waf_domain_associate_certificate" "test" {
  certificate_id = var.certificate_id
  cloud_host_ids = var.cloud_host_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `certificate_id` - (Required, String, NonUpdatable) Specifies the ID of the certificate.

* `cloud_host_ids` - (Optional, List, NonUpdatable) Specifies the ID lsit of the domain in cloud mode.

* `premium_host_ids` - (Optional, List, NonUpdatable) Specifies the ID list of the domain in dedicated mode.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
