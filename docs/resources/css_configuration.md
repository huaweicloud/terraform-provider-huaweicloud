---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_configuration"
description: ""
---

# huaweicloud_css_configuration

Manages a CSS configuration resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

resource "huaweicloud_css_configuration" "test" {
  cluster_id                   = var.cluster_id
  thread_pool_force_merge_size = "3"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, ForceNew) The CSS cluster ID.

  Changing this parameter will create a new resource.

* `http_cors_allow_credetials` - (Optional, String) Whether to return the Access-Control-Allow-Credentials of
  the header during cross-domain access.  
  The value can be **true** or **false**. Default value: **false**.

* `http_cors_allow_origin` - (Optional, String) Origin IP address allowed for cross-domain access, for example, **122.122.122.122:9200**.

* `http_cors_max_age` - (Optional, String) Cache duration of the browser. The cache is automatically cleared
  after the time range you specify.  
  Unit: s, Default value: **1,728,000**.

* `http_cors_allow_headers` - (Optional, String) Headers allowed for cross-domain access.  
  Including **X-Requested-With**, **Content-Type**, and **Content-Length**.
  Use commas (,) and spaces to separate headers.

* `http_cors_enabled` - (Optional, String) Whether to allow cross-domain access.  
  The value can be **true** or **false**. Default value: **false**.

* `http_cors_allow_methods` - (Optional, String) Methods allowed for cross-domain access.  
  Including **OPTIONS**, **HEAD**, **GET**, **POST**, **PUT**, and **DELETE**.
  Use commas (,) and spaces to separate methods.

* `reindex_remote_whitelist` - (Optional, String) Configured for migrating data from the current cluster to
  the target cluster through the reindex API.
  The example value is **122.122.122.122:9200**.

* `indices_queries_cache_size` - (Optional, String) Cache size in the query phase. Value range: **1%** to **100%**.  
  Unit: %, Default value: **10%**.

* `thread_pool_force_merge_size` - (Optional, String) Queue size in the force merge thread pool.  
  Default value: **1**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `cluster_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

The CSS configuration can be imported using the `id` which equals the `cluster_id`, e.g.

```bash
$ terraform import huaweicloud_css_configuration.test <id>
```
