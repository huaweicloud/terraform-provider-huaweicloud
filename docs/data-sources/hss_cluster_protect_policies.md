---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_cluster_protect_policies"
description: |-
  Use this data source to get the list of HSS cluster protect policies within HuaweiCloud.
---

# huaweicloud_hss_cluster_protect_policies

Use this data source to get the list of HSS cluster protect policies within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_cluster_protect_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `general_policy_num` - The general policies count.

* `malicious_image_policy_num` - The malicious image policies count.

* `security_policy_num` - The security policies count.

* `data_list` - The list of cluster protect policies information.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `content` - The policy content.

* `deploy_content` - The deployment content.

* `parameters` - The parameters.

* `policy_name` - The policy name.

* `policy_id` - The policy ID.

* `resources` - The resources.

  The [resources](#resources_struct) structure is documented below.

* `template_id` - The template ID.

* `template_name` - The template name.

* `template_type` - The template type.

* `description` - The description.

* `update_time` - The update time.

* `create_time` - The creation time.

* `image_num` - The number of protective images.

* `labels_num` - The number of protective labels.

* `status` - The status.

* `white_images` - The whitelist image.

  The [white_images](#white_images_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `images` - The images.

* `labels` - The labels.

* `namespace` - The namespace.

<a name="white_images_struct"></a>
The `white_images` block supports:

* `cluster_id` - The cluster ID.

* `image_name` - The image name.

* `image_version` - The image version.
