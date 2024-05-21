---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_parameter_templates"
description: |-
  Use this data source to get the list of DDS parameter templates.
---

# huaweicloud_dds_parameter_templates

Use this data source to get the list of DDS parameter templates.

## Example Usage

```hcl
data "huaweicloud_dds_parameter_templates" "test1" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the parameter template name.

* `node_type` - (Optional, String) The the node type of parameter template.
  Valid value:
  + **mongos**: the mongos node type.
  + **shard**: the shard node type.
  + **config**: the config node type.
  + **replica**: the replica node type.
  + **single**: the single node type.

* `datastore_version` - (Optional, String) Specifies the database (DB Engine) version.
  The value can be `4.4`, `4.2`, `4.0`, `3.4` or `3.2`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - DDS parameter templates list.

  The [configurations](#configurations_struct) structure is documented below.

<a name="configurations_struct"></a>
The `configurations` block supports:

* `id` - The parameter template ID.

* `name` - The parameter template name.

* `node_type` - The the node type of parameter template.
  Valid value:
  + **mongos**: the mongos node type.
  + **shard**: the shard node type.
  + **config**: the config node type.
  + **replica**: the replica node type.
  + **single**: the single node type.

* `datastore_version` - Database (DB Engine) version.

* `datastore_name` - Database (DB Engine) type.

* `user_defined` - Whether the parameter template is a custom template.
  + **false**: default parameter template.
  + **true**: custom template.

* `description` - The parameter template description.

* `updated_at` - The creation time of the parameter template.

* `created_at` - The update time of the parameter template.
