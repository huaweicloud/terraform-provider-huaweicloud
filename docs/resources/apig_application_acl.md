---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_application_acl"
description: |-
  Manages multiple ACL rules of the same type under an APIG application within HuaweiCloud.
---

# huaweicloud_apig_application_acl

Manages multiple ACL rules of the same type under an APIG application within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "application_id" {}

resource "huaweicloud_apig_application_acl" "test" {
  instance_id    = var.instance_id
  application_id = var.application_id
  type           = "PERMIT"
  values         = ["192.145.0.0/16"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application and ACL rules are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the application
  belongs.  
  Changing this will create a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the ID of the application to which the ACL rules belong.  
  Changing this will create a new resource.

* `type` - (Required, String) Specifies the ACL type.  
  The valid values are as follows:
  + **PERMIT**: Whitelist.
  + **DENY**: Blacklist.

* `values` - (Required, List) Specifies the ACL values.  
  The valid formats are as follows:
  + Common IP address, e.g. `127.0.0.1` or `::1`.
  + IP address with mask, e.g. `192.145.0.0/16` or `2407:c080:17ef:ffff::3104:703a/64`.
  + IP address range, e.g. `127.0.0.1-192.145.0.1` or `2407:c080:17ef:ffff::3104:703a-2407:c080:17ef:ffff::3104:704a`.

  -> Before entering an IPv6 address, ensure that the dedicated instance supports the **IPv6** protocol.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID(, also the application ID).

## Import

ACL rules can be imported using the related dedicate instance ID (`instance_id`) and the resource ID (`id`, also the
`application_id`) and the ID of the , separated by a slash, e.g.

```bash
$ terraform import huaweicloud_apig_application_acl.test <instance_id>/<id>
```
