---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_ou"
description: |-
  Manage an organizational unit (OU) resource of the Workspace within HuaweiCloud.
---

# huaweicloud_workspace_ou

Manage an organizational unit (OU) resource of the Workspace within HuaweiCloud.

-> Before using this resource, please ensure that the OU already exists on the AD server and the Workspace service
   has been registered to the AD domain.

## Example Usage

```hcl
variable "ou_name" {}
variable "ad_domain_name" {}

resource "huaweicloud_workspace_ou" "test" {
  ou_name = var.ou_name
  domain  = var.ad_domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the OU is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `ou_name` - (Required, String) Specifies the name of the OU.  
  The name only allows Chinese and English characters, digits, spaces and special characters (-_$!@*?.).
  Multiple levels of OUs are separated by slashes (/) and no space is allowed before or after the slashes (/),
  a maximum of five layers of OUs are supported. e.g. `ab/cd/ef`.

* `domain` - (Required, String, NonUpdatable) Specifies the AD domain name to which the OU belongs.  

* `description` - (Optional, String) Specifies the description of the OU.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also OU ID.

* `ou_dn` - The distinguished name (DN) of the OU.

* `domain_id` - The ID of the AD domain.

## Import

The OU resource can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_ou.test <id>
```

If the OU ID is unknown, the OU name can be used as an alternative to ID.

```bash
$ terraform import huaweicloud_workspace_ou.test <ou_name>
```
