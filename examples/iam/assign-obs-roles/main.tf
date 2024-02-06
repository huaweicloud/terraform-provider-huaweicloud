terraform {
  required_providers {
    huaweicloud = {
      source = "huaweicloud/huaweicloud"
      version = ">=1.48.0"
    }
  }
}

resource "huaweicloud_identity_role" "obs_role" {
  name        = "obs_bucket_create_role"
  type        = "XA"
  description = "test for obs role assignment"
  policy      = <<EOT
{
  "Version": "1.1",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "obs:*:*"
      ]
    },
    {
      "Effect": "Deny",
      "Action": [
        "obs:object:DeleteObjectVersion",
        "obs:object:DeleteAccessLabel",
        "obs:bucket:DeleteDirectColdAccessConfiguration",
        "obs:object:AbortMultipartUpload",
        "obs:bucket:DeleteBucketWebsite",
        "obs:object:DeleteObject",
        "obs:bucket:DeleteBucketPolicy",
        "obs:bucket:DeleteBucketCustomDomainConfiguration",
        "obs:object:RestoreObject",
        "obs:bucket:DeleteBucket",
        "obs:object:ModifyObjectMetaData",
        "obs:bucket:DeleteBucketInventoryConfiguration",
        "obs:bucket:DeleteReplicationConfiguration",
        "obs:bucket:DeleteBucketTagging"
      ]
    }
  ]
}
EOT
}

resource "huaweicloud_identity_group" "test" {
  name = var.group_name
}

# OBS policy must grant global permission.
resource "huaweicloud_identity_group_role_assignment" "test" {
  group_id  = huaweicloud_identity_group.test.id
  role_id   = huaweicloud_identity_role.obs_role.id
  domain_id = var.account_domain_id
}

data "huaweicloud_identity_projects" "test" {
  # The special project for managing OBS service billing.
  name = "MOS"
}

resource "huaweicloud_identity_group_role_assignment" "mos" {
  role_id    = huaweicloud_identity_role.obs_role.id
  group_id   = huaweicloud_identity_group.test.id
  project_id = try(data.huaweicloud_identity_projects.test.projects[0].id, "")
}

data "huaweicloud_identity_users" "test" {
  name = var.user_name
}

resource "huaweicloud_identity_group_membership" "test" {
  depends_on = [
    huaweicloud_identity_group_role_assignment.test,
    huaweicloud_identity_group_role_assignment.mos,
  ]

  group = huaweicloud_identity_group.test.id
  users = [
    try(data.huaweicloud_identity_users.test.users[0].id, "")
  ]
}
