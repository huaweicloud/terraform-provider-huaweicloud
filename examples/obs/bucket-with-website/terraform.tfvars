key_alias              = "tf-test-obs-key"
bucket_name            = "tf-test-obs-bucket"
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
    <p> This is the 404 error page.</p>
  </body>
</html>
EOT
  }
}
