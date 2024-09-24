---
subcategory: "Database Security Service (DBSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dbss_audit_rule_scopes"
description: |-
  Use this data source to get a list of audit rule scopes.
---

# huaweicloud_dbss_audit_rule_scopes

Use this data source to get a list of audit rule scopes.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dbss_audit_rule_scopes" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the audit instance ID to which the audit scopes belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `scopes` - The list of the audit scopes.

  The [scopes](#scopes_struct) structure is documented below.

<a name="scopes_struct"></a>
The `scopes` block supports:

* `id` - The ID of the audit scope.

* `name` - The name of the audit scope.

* `status` - The status of the audit scope.

* `action` - The action of the audit scope.

* `exception_ips` - The exception IP addresses of the audit scope.

* `source_ips` - The source IP addresses of the audit scope.

* `source_ports` - The source ports of the audit scope.

* `db_ids` - The database IDs associated with the audit scope.

* `db_names` - The database names associated with the audit scope.

* `db_users` - The database accounts associated with the audit scope.

* `all_audit` - Whether is full audit.
