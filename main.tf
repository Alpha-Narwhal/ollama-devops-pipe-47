provider "aws" {
  region = "us-east-1"
}

resource "aws_instance" "web" {
  ami           = "ami-0c02fb55956c7d316"
  instance_type = "t2.micro"

  user_data = <<-EOF
              #!/bin/bash
              yum update -y
              yum install -y nginx
              systemctl start nginx
              systemctl enable nginx
              echo "Hello from Terraform" > /usr/share/nginx/html/index.html
              EOF

  tags = {
    Name = "Ollama"
  }
}