FROM public.ecr.aws/lambda/provided:al2023.2024.02.07.17 AS runtime

RUN dnf install -y elfutils

ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib:/usr/local/lib64
ENV PATH=$PATH:/usr/local/bin
