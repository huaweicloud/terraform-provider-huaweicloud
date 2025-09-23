---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_logs"
description: |-
  Use this data source to get the list of LTS logs under specified DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_logs

Use this data source to get the list of LTS logs under specified DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_cluster_logs" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID to which the logs belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - Whether the LTS log function is enabled.
  + **OPEN**
  + **CLOSE**

* `logs` - All LTS logs that match the filter parameters.

  The [logs](#logs_struct) structure is documented below.

<a name="logs_struct"></a>
The `logs` block supports:

* `id` - The ID of the log.

* `type` - The type of the log.

* `description` - The description of the log.

* `access_url` - The URL to access the LTS log.
