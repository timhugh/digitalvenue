FROM public.ecr.aws/lambda/provided:al2023.2024.02.07.17 AS runtime

RUN dnf install -y elfutils
