FROM centos

COPY s2i /usr/local/bin

RUN yum install -y git

WORKDIR click2cloud-dockm/

COPY dist app/

VOLUME /click2cloud-dockm/data

EXPOSE 9000

ENTRYPOINT ["app/dockm-linux-amd64", "-d", "data", "-a", "app", "--no-analytics"]