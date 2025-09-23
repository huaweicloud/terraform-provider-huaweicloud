---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_migrate_domain"
description: |-
  Manages a resource to migrate the domain from an enterprise project to another within HuaweiCloud.
---

# huaweicloud_waf_migrate_domain

Manages a resource to migrate the domain from an enterprise project to another within HuaweiCloud.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> 1. The current resource is a one-time resource, and destroying this resource will not change the current status.
<br/>2. If you use this resource to migrate a domain, the opration will trigger changes to the existing
resource: `huaweicloud_waf_domain` or `huaweicloud_waf_dedicated_domain`.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "target_enterprise_project_id" {}
variable "policy_id" {}
variable "host_ids" {
  type = list(string)
}

resource "huaweicloud_waf_migrate_domain" "test" {
  enterprise_project_id        = var.enterprise_project_id
  target_enterprise_project_id = var.target_enterprise_project_id
  policy_id                    = var.policy_id
  host_ids                     = var.host_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the ID of the source enterprise project to which
  the domain belongs.

* `target_enterprise_project_id` - (Required, String, NonUpdatable) Specifies the ID of the destination enterprise
  project.

* `policy_id` - (Required, String, NonUpdatable) Specifies the ID of the policy under destination enterprise project.

* `host_ids` - (Required, List, NonUpdatable) Specifies the ID lsit of the domains.

* `certificate_id` - (Optional, String, NonUpdatable) Specifies the ID of the certificate under destination enterprise
  project.

  -> This parameter is valid and mandatory only when the domain `client_protocol` is set to **HTTPS**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
