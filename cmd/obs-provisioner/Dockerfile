# Based on centos
FROM centos:7.4.1708
LABEL maintainers="Kubernetes Authors"
LABEL description="OBS Provisioner"

# Copy from build directory
COPY obs-provisioner /obs-provisioner

# Update
RUN yum -y update

# Define default command
ENTRYPOINT ["/obs-provisioner"]
