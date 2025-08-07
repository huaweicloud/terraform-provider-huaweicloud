# Create an OBS bucket with website hosting

This example provides best practice code for using Terraform to create an Object Storage Service (OBS) bucket in
HuaweiCloud with static website hosting capabilities.

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
* `website_configurations` - The configurations of the OBS bucket website

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

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  bucket_name            = "your_obs_bucket_name"
  website_configurations = {
    index = {
      file_name = "index.html"
      content   = <<EOT
  <html>
    <head>
      <title>Hello OBS!</title>
      <meta charset="utf-8">
    </head>
    <body>
      <p>Welcome to use OBS static website hosting.</p>
      <p>This is the homepage.</p>
    </body>
  </html>
  EOT
    }
    error = {
      file_name = "error.html"
      content   = <<EOT
  <html>
    <head>
      <title>Hello OBS!</title>
      <meta charset="utf-8">
    </head>
    <body>
      <p>Welcome to use OBS static website hosting.</p>
      <p>This is the 404 error page.</p>
    </body>
  </html>
  EOT
    }
  }
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
3. **Static Website Hosting**: Configures the bucket for static website hosting
4. **Flexible KMS Key Configuration**: Supports both creating new KMS key and using existing KMS key
5. **Website Configuration**: Supports custom index and error pages
6. **Storage Class Configuration**: Configurable storage class for cost optimization
7. **Access Control**: Configurable ACL for bucket access management
8. **Tagging Support**: Supports custom tags for resource management
9. **Bucket Policy**: Automatically configures bucket policy for public read access

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

## Website Configuration Options

### Option 1: Simple HTML Pages

```hcl
website_configurations = {
  index = {
    file_name = "index.html"
    content   = "<html><body><h1>Welcome!</h1></body></html>"
  }
  error = {
    file_name = "error.html"
    content   = "<html><body><h1>Page Not Found</h1></body></html>"
  }
}
```

### Option 2: Complex HTML with Heredoc

```hcl
website_configurations = {
  index = {
    file_name = "index.html"
    content   = <<EOT
<html>
  <head>
    <title>My Website</title>
    <meta charset="utf-8">
    <style>
      body { font-family: Arial, sans-serif; margin: 40px; }
    </style>
  </head>
  <body>
    <h1>Welcome to My Website</h1>
    <p>This is a static website hosted on OBS.</p>
  </body>
</html>
EOT
  }
  error = {
    file_name = "error.html"
    content   = <<EOT
<html>
  <head>
    <title>Page Not Found</title>
    <meta charset="utf-8">
  </head>
  <body>
    <h1>404 - Page Not Found</h1>
    <p>The page you are looking for does not exist.</p>
    <a href="/">Return to Homepage</a>
  </body>
</html>
EOT
  }
}
```

### Option 3: Custom File Names

```hcl
website_configurations = {
  index = {
    file_name = "home.html"
    content   = "<html><body><h1>Home Page</h1></body></html>"
  }
  error = {
    file_name = "404.html"
    content   = "<html><body><h1>404 Error</h1></body></html>"
  }
}
```

## Storage Classes

The example supports different storage classes for cost optimization:

* `STANDARD` - Standard storage for frequently accessed data (default)
* `WARM` - Infrequent access storage for data accessed less than once per month
* `COLD` - Archive storage for data accessed less than once per year

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the OBS bucket and website configuration is usually instantaneous
* This example creates the OBS bucket, optionally a KMS key for encryption, and configures static website hosting
* KMS encryption provides server-side encryption for enhanced data security
* All resources will be created in the specified region
* Bucket names must be globally unique across all HuaweiCloud accounts
* When `bucket_force_destroy` is set to true, the bucket can be destroyed even if it contains objects
* The website_configurations map must contain both 'index' and 'error' keys
* The bucket policy allows public read access to all objects in the bucket
* Website files are automatically created and uploaded to the bucket

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.64.3 |
| random | >= 3.0.0 |
