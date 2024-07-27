---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_multi_account_aggregation_rules"
description: |-
  Use this data source to get the list of AOM multi account aggregation rules.
---

# huaweicloud_aom_multi_account_aggregation_rules

Use this data source to get the list of AOM multi account aggregation rules.

## Example Usage

```hcl
data "huaweicloud_aom_multi_account_aggregation_rules" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the prometheus instance belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - Indicates the multi account aggregation rules list.
  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `instance_id` - Indicates the prometheus instance ID.

* `accounts` - Indicates the accounts list.
  The [accounts](#accounts_struct) structure is documented below.

* `services` - Indicates the services list.
  The [services](#services_struct) structure is documented below.

* `send_to_source_account` - Indicates whether the member accounts retain metric data after they are connected to the
  prometheus instance for aggregation.

<a name="accounts_struct"></a>
The `accounts` block supports:

* `id` - Indicates the account ID.

* `name` - Indicates the account name.

* `urn` - Indicates the uniform resource name of the account.

* `join_method` - Indicates the method how the account joined in the organization.

* `joined_at` - Indicates the time when the account joined in the organization.

<a name="services_struct"></a>
The `services` block supports:

* `service` - Indicates the service name.

* `metrics` - Indicates the metrics List.
