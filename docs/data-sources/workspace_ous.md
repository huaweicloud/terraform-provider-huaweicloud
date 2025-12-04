---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_ous"
description: |-
  Use this data source to get the list of organizational units (OUs) within HuaweiCloud.
---

# huaweicloud_workspace_ous

Use this data source to get the list of organizational units (OUs) within HuaweiCloud.

## Example Usage

### Query all OUs

```hcl
data "huaweicloud_workspace_ous" "test" {}
```

### Filter OUs by name

```hcl
variable "ou_name" {}

data "huaweicloud_workspace_ous" "test" {
  ou_name = var.ou_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the OUs are located.  
  If omitted, the provider-level region will be used.

* `ou_name` - (Optional, String) Specifies the name of the OU.  
  Fuzzy match is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ous` - The list of OUs that match the filter parameters.  
  The [ous](#workspace_ous) structure is documented below.

<a name="workspace_ous"></a>
The `ous` block supports:

* `id` - The ID of the OU.

* `name` - The name of the OU.

* `domain_id` - The ID of the AD domain.

* `domain` - The AD domain name to which the OU belongs.

* `ou_dn` - The distinguished name (DN) of the OU.

* `description` - The description of the OU.
