resource "huaweicloud_obs_bucket" "mywebsite" {
  bucket = "mywebsite"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}

# granting the Read-Only permission to anonymous users
resource "huaweicloud_obs_bucket_policy" "policy" {
  bucket = huaweicloud_obs_bucket.mywebsite.bucket
  policy = <<POLICY
{
  "Statement": [
    {
      "Sid": "AddPerm",
      "Effect": "Allow",
      "Principal": {"ID": "*"},
      "Action": ["GetObject"],
      "Resource": "mywebsite/*"
    } 
  ]
}
POLICY
}

# put index.html
resource "huaweicloud_obs_bucket_object" "index" {
  bucket = huaweicloud_obs_bucket.mywebsite.bucket
  key    = "index.html"
  source = "index.html"
}

# put error.html
resource "huaweicloud_obs_bucket_object" "error" {
  bucket = huaweicloud_obs_bucket.mywebsite.bucket
  key    = "error.html"
  source = "error.html"
}
