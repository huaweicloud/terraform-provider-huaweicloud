---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_publishers"
description: |-
  Use this data source to get a list of CodeArts pipeline publishers.
---

# huaweicloud_codearts_pipeline_publishers

Use this data source to get a list of CodeArts pipeline publishers.

## Example Usage

```hcl
data "huaweicloud_codearts_pipeline_publishers" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the publisher name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `publishers` - Indicates the publisher list.
  The [publishers](#attrblock--publishers) structure is documented below.

<a name="attrblock--publishers"></a>
The `publishers` block supports:

* `id` - Indicates the publisher ID.

* `name` - Indicates the publisher name.

* `en_name` - Indicates the publisher English name.

* `logo_url` - Indicates the logo URL.

* `source_url` - Indicates the source URL.

* `support_url` - Indicates the support URL.

* `website` - Indicates the website URL.

* `description` - Indicates the description.

* `auth_status` - Indicates the authorization status.

* `last_update_time` - Indicates the update time.

* `last_update_user_id` - Indicates the updater ID.

* `last_update_user_name` - Indicates the updater name.

* `user_id` - Indicates the user ID.
