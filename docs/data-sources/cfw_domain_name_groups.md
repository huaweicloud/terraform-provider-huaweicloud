---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_domain_name_groups"
description: |-
  Use this data source to get the list of CFW domain name groups.
---

# huaweicloud_cfw_domain_name_groups

Use this data source to get the list of CFW domain name groups.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "object_id" {}

data "huaweicloud_cfw_domain_name_groups" "test" {
  fw_instance_id = var.fw_instance_id
  object_id      = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall instance ID.

* `object_id` - (Required, String) Specifies the protected object ID.

* `key_word` - (Optional, String) Specifies the key word.

* `type` - (Optional, String) Specifies the domain name group type.
  The value can be:
  + **0**: means application type;
  + **1**: means network type.

* `config_status` - (Optional, String) Specifies the configuration status.
  The valid values are as follows:
  + **-1**: not configured.
  + **0**: configuration failed.
  + **1**: configuration succeeded.
  + **2**: configuration in progress.
  + **3**: normal.
  + **4**: configuration exception - domain group usage.

* `ref_count` - (Optional, String) Specifies the domain name group reference count.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `group_id` - (Optional, String) Specifies the domain name group ID.

* `name` - (Optional, String) Specifies the name of a domain name group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The domain name group list.

  The [records](#data_records_struct) structure is documented below.

<a name="data_records_struct"></a>
The `records` block supports:

* `group_id` - The domain name group ID.

* `name` - The name of the domain name group.

* `description` - The domain name group description.

* `ref_count` - The domain name group reference count.

* `type` - The domain name group type.

* `config_status` - The configuration status.

* `message` - The configuration message.

* `rules` - The used rule list.

  The [rules](#records_rules_struct) structure is documented below.

* `domain_names` - The list of domain names.

  The [domain_names](#records_domain_names_struct) structure is documented below.

<a name="records_rules_struct"></a>
The `rules` block supports:

* `id` - The rule ID.

* `name` - The rule name.

<a name="records_domain_names_struct"></a>
The `domain_names` block supports:

* `domain_name` - The domain name.

* `description` - The description.

* `domain_address_id` - The domain address ID.
