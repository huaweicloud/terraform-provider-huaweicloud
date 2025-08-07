key_alias             = "tf-test-obs-key"
bucket_name           = "tf-test-obs-bucket"
object_name           = "tf-test-obs-object"
object_upload_content = <<EOT
def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
EOT
