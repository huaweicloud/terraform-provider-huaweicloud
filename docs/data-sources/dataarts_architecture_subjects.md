---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_subjects"
description: |-
  Use this data source to get the list of subjects within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_subjects

Use this data source to get the list of subjects within HuaweiCloud.

## Example Usage

### Query all subjects and without any filter

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_subjects" "test" {
  workspace_id = var.workspace_id
}
```

### Query the subjects and using status filter

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_subjects" "test" {
  workspace_id = var.workspace_id
  status       = "PUBLISHED"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the subjects are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the subjects belong.

* `name` - (Optional, String) Specifies the name or code of the subject (fuzzy query).

* `create_by` - (Optional, String) Specifies the creator of the subject.

* `owner` - (Optional, String) Specifies the owner of the subject.

* `status` - (Optional, String) Specifies the business status of the subject.  
  The valid values are as follows:
  + **DRAFT**: Draft.
  + **PUBLISH_DEVELOPING**: Publishing.
  + **PUBLISHED**: Published.
  + **OFFLINE_DEVELOPING**: Going offline.
  + **OFFLINE**: Offline.
  + **REJECT**: Rejected.

* `begin_time` - (Optional, String) Specifies the start time for filtering, in RFC3339 format.

* `end_time` - (Optional, String) Specifies the end time for filtering, in RFC3339 format.

* `parent_id` - (Optional, String) Specifies the parent directory ID of the subject.  
  The root node does not have this ID. An empty value means all nodes,
  and `-1` means nodes under the root node.

* `level` - (Optional, Int) Specifies the level of the subject.

* `with_relation` - (Optional, Bool) Specifies whether to include relation information.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `subjects` - The list of subjects that matched filter parameters.  
  The [subjects](#dataarts_architecture_subject) structure is documented below.

<a name="dataarts_architecture_subject"></a>
The `subjects` block supports:

* `id` - The ID of the subject.

* `name` - The Chinese name of the subject.

* `name_en` - The English name of the subject.

* `description` - The description of the subject.

* `qualified_name` - The qualified name of the subject.

* `guid` - The GUID of the subject, automatically generated.

* `code` - The code of the subject.

* `alias` - The alias of the subject.

* `status` - The business status of the subject.

* `new_biz` - The business version management information, in JSON format.

* `data_owner` - The data owner of the subject.

* `data_owner_list` - The data owner list of the subject.

* `data_department` - The data department of the subject.

* `path` - The path of the subject.

* `level` - The level of the subject.

* `ordinal` - The ordinal of the subject.

* `owner` - The owner of the subject.

* `parent_id` - The parent directory ID of the subject.

* `swap_order_id` - The swap order ID of the subject.

* `qualified_id` - The qualified ID of the subject.

* `from_public` - Whether the subject is from public layer.

* `created_by` - The creator of the subject.

* `updated_by` - The last editor of the subject.

* `created_at` - The creation time of the subject, in RFC3339 format.

* `updated_at` - The latest update time of the subject, in RFC3339 format.

* `children_num` - The number of child processes.

* `children` - The child directories, in JSON format.

* `self_defined_field` - The self-defined field in JSON format.

* `relations` - The list of relations.  
  The [relations](#dataarts_architecture_subject_relation) structure is documented below.

<a name="dataarts_architecture_subject_relation"></a>
The `relations` block supports:

* `id` - The ID of the relation.

* `source_table_id` - The source table ID of the relation.

* `target_table_id` - The target table ID of the relation.

* `name` - The name of the relation.

* `source_table_name` - The source table name of the relation.

* `target_table_name` - The target table name of the relation.

* `role` - The role of the relation.

* `source_type` - The source type of the relation.
  + **ONE**
  + **ZERO_OR_ONE**
  + **ZERO_OR_N**
  + **ONE_OR_N**

* `target_type` - The target type of the relation.
  + **ONE**
  + **ZERO_OR_ONE**
  + **ZERO_OR_N**
  + **ONE_OR_N**

* `created_by` - The creator of the relation.

* `updated_by` - The last editor of the relation.

* `created_at` - The creation time of the relation, in RFC3339 format.

* `updated_at` - The latest update time of the relation, in RFC3339 format.

* `mappings` - The list of mappings.  
  The [mappings](#dataarts_architecture_subject_relation_mapping) structure is documented below.

<a name="dataarts_architecture_subject_relation_mapping"></a>
The `mappings` block supports:

* `id` - The ID of the mapping.

* `relation_id` - The relation ID of the mapping.

* `source_field_id` - The source field ID of the mapping.

* `target_field_id` - The target field ID of the mapping.

* `source_field_name` - The source field name of the mapping.

* `target_field_name` - The target field name of the mapping.

* `created_by` - The creator of the mapping.

* `updated_by` - The last editor of the mapping.

* `created_at` - The creation time of the mapping, in RFC3339 format.

* `updated_at` - The latest update time of the mapping, in RFC3339 format.
