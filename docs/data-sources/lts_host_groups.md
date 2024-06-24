---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_host_groups"
description: |-
  Use this data source to get the list of LTS host groups.
---

# huaweicloud_lts_host_groups

Use this data source to get the list of LTS host groups.

## Example Usage

```hcl
data "huaweicloud_lts_host_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `host_group_id` - (Optional, String) Speicifies the ID of the host group.

* `name` - (Optional, String) Speicifies the name of the host group.

* `type` - (Optional, String) Speicifies the type of the host group.  
  The valid values are as follows:
  + **windows**
  + **linux**

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the host group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - All host groups that match the filter parameters.

  The [groups](#groups_struct) structure is documented below.

<a name="groups_struct"></a>
The `groups` block supports:

* `id` - The ID of the host group.

* `name` - The name of the host group.

* `type` - The type of the host group.

* `host_ids` - The ID list of hosts to associate with the host group.

* `tags` - The key/value pairs to associate with the host group.

* `created_at` - The creation time of the host group, in RFC3339 format.

* `updated_at` - The latest update time of the host group, in RFC3339 format.
