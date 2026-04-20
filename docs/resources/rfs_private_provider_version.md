---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_private_provider_version"
description: |-
  Manages a RFS private provider version resource within HuaweiCloud.
---

# huaweicloud_rfs_private_provider_version

Manages a RFS private provider version resource within HuaweiCloud.

## Example Usage

```hcl
variable "provider_name" {}
variable "provider_version" {}
variable "function_graph_urn" {}

resource "huaweicloud_rfs_private_provider_version" "test" {
  provider_name      = var.provider_name
  provider_version   = var.provider_version
  function_graph_urn = var.function_graph_urn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource is located.  
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `provider_name` - (Required, String, NonUpdatable) Specifies the name of a private provider.

* `provider_version` - (Required, String, NonUpdatable) Specifies the provider version number. The version number must
  follow semantic versioning (Semantic Version).

* `function_graph_urn` - (Required, String, NonUpdatable) Specifies the URN of the FunctionGraph function used to run
  the private provider. Only a function in the same region as RFS is supported.  
  For a detailed explanation of this parameter,
  please refer to [documentation](https://support.huaweicloud.com/api-functiongraph/functiongraph_06_0102.html).

* `provider_id` - (Optional, String, NonUpdatable) Specifies the ID of a private provider.

* `version_description` - (Optional, String, NonUpdatable) Specifies the description of a private provider version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The format is `<provider_name>/<provider_version>`.

* `provider_source` - The source parameter that users need to specify when defining the required providers information
  in Terraform templates using private providers.

  -> This parameter is concatenated in the form of "huawei.com/private-provider/{provider_name}". For detailed
  information on the use of the `provider_name` and `provider_source` fields in the template, please refer to the API
  description for creating a private provider.

* `create_time` - The creation time of a private provider version. It is represented in UTC format
  (YYYY-MM-DDTHH:mm:ss.SSSZ), such as **1970-01-01T00:00:00.000Z**.

## Import

Private provider version can be imported using `provider_name` and `provider_version` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rfs_private_provider_version.test <provider_name>/<provider_version>
```
