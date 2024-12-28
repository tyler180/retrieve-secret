# resource "aws_lambda_layer_version" "retrieve_secret_lambda_layer" {
#   filename   = "retrieve_secret.zip"
#   layer_name = "retrieve_secret"

#   compatible_runtimes = ["go1.x"]
# }

resource "aws_lambda_layer_version" "retrieve_secret_lambda_layer" {
  layer_name          = "secrets-go-layer"
  description         = "Go source for retrieving AWS Secrets"
  compatible_runtimes = ["go1.x"]
  
  # If you have the zip locally:
  filename = "${path.module}/retrieve_secret_layer.zip"

  # Alternatively, if the zip is in S3, use these instead:
  # s3_bucket = "your-bucket"
  # s3_key    = "layers/mysecrets_layer.zip"
  
  # Optionally set license info, environment variables, etc.
}