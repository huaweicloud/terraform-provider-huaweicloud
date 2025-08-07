# Upload an OBS object with content

This example provides best practice code for using Terraform to create an Object Storage Service (OBS) bucket in
HuaweiCloud and upload an object with custom content.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the OBS bucket is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `bucket_name` - The name of the OBS bucket
* `object_name` - The name of the OBS object to be uploaded
* `object_upload_content` - The content of the OBS object to be uploaded

#### Optional Variables

* `bucket_encryption` - Whether to enable encryption for the OBS bucket (default: true)
* `bucket_encryption_key_id` - The encryption key ID of the OBS bucket (default: "")
* `key_alias` - The alias of the KMS key (default: "")
  The alias of the KMS key (required when `bucket_encryption` is true and `bucket_encryption_key_id` is empty)
* `key_usage` - The usage of the KMS key (default: "ENCRYPT_DECRYPT")
* `bucket_storage_class` - The storage class of the OBS bucket (default: "STANDARD")
* `bucket_acl` - The ACL of the OBS bucket (default: "private")
* `bucket_sse_algorithm` - The SSE algorithm of the OBS bucket (default: "kms")
* `bucket_force_destroy` - Whether to force destroy the OBS bucket (default: true)
* `bucket_tags` - The tags of the OBS bucket (default: {})
* `object_extension_name` - The extension name of the OBS object to be uploaded (default: ".txt")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  bucket_name           = "your_obs_bucket_name"
  object_name           = "your_object_name"
  object_upload_content = <<EOT
  def main():
      print("Hello, World!")

  if __name__ == "__main__":
      main()
  EOT
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Import the existing resources (optional):

  ```bash
  $ terraform import huaweicloud_kms_key.test[0] xxxxxxxx-xxx-xxx-xxx-xxxxxxxxxxxx
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Features

This example demonstrates the following features:

1. **OBS Bucket Creation**: Creates a complete OBS bucket with all necessary components
2. **KMS Encryption**: Enables KMS encryption for enhanced data security
3. **Flexible KMS Key Configuration**: Supports both creating new KMS key and using existing KMS key
4. **Object Upload**: Uploads custom content as an object to the bucket
5. **Storage Class Configuration**: Configurable storage class for cost optimization
6. **Access Control**: Configurable ACL for bucket access management
7. **Tagging Support**: Supports custom tags for resource management

## Encryption Options

### Option 1: Create New KMS Key

If you don't provide an existing KMS key ID, the example will create a new KMS key with the specified alias:

```hcl
bucket_encryption = true
key_alias         = "your_kms_key_alias"
```

### Option 2: Use Existing KMS Key

If you have an existing KMS key, you can use it directly:

```hcl
bucket_encryption        = true
bucket_encryption_key_id = "your_existing_kms_key_id"
```

### Option 3: Disable Encryption

If you don't need encryption, you can disable it:

```hcl
bucket_encryption = false
```

## Content Upload Options

### Option 1: Simple Text Content

```hcl
object_upload_content = "This is a simple text file content."
```

### Option 2: Multi-line Content (Heredoc)

```hcl
object_upload_content = <<EOT
def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
EOT
```

### Option 3: JSON Content

```hcl
object_upload_content = <<EOT
{
  "name": "example",
  "version": "1.0.0",
  "description": "Sample JSON content"
}
EOT
```

## Storage Classes

The example supports different storage classes for cost optimization:

* `STANDARD` - Standard storage for frequently accessed data (default)
* `WARM` - Infrequent access storage for data accessed less than once per month
* `COLD` - Archive storage for data accessed less than once per year

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the OBS bucket and object upload is usually instantaneous
* This example creates the OBS bucket, optionally a KMS key for encryption, and uploads an object with custom content
* KMS encryption provides server-side encryption for enhanced data security
* All resources will be created in the specified region
* Bucket names must be globally unique across all HuaweiCloud accounts
* When `bucket_force_destroy` is set to true, the bucket can be destroyed even if it contains objects
* The uploaded object will have content type "application/xml" by default
* The object name will be constructed as `object_name + object_extension_name` if extension is provided

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.64.3 |
| random | >= 3.0.0 |
