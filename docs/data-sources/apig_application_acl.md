---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_application_acl"
description: |-
  Use this data source to query the multiple ACL rules from an APIG application within HuaweiCloud.
---

# huaweicloud_apig_application_acl

Use this data source to query the multiple ACL rules from an APIG application within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "application_id" {}

data "huaweicloud_apig_application_acl" "test" {
  instance_id    = var.instance_id
  application_id = var.application_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the application and ACL rules are located.  
  If omitted, the provider-level region will be used

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the application
  belongs.

* `application_id` - (Required, String) Specifies the ID of the application to which the ACL rules belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `type` - The ACL type.
  + **PERMIT**: Whitelist.
  + **DENY**: Blacklist.

* `values` - The ACL values.  
  The elements of the rule list contain the following:
  + Common IP address, e.g. `127.0.0.1` or `::1`.
  + IP address with mask, e.g. `192.145.0.0/16` or `2407:c080:17ef:ffff::3104:703a/64`.
  + IP address range, e.g. `127.0.0.1-192.145.0.1` or `2407:c080:17ef:ffff::3104:703a-2407:c080:17ef:ffff::3104:704a`.

-> If the values of the attribute `type` and `values` are both empty, means no ACL rule is configured.
