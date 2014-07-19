See: https://gist.github.com/mattes/2d0ffd027cb16571895c#file-dockerfile-tmpl-L21

,,,bash
# generate Dockerfile from Dockerfile.tmpl
chmod +x build_docker.sh
./build_docker.sh

# build the actual boot2docker.iso with virtual box guest additions
docker build -t mattes/boot2docker-vbga .

# the following line is proposed in many tutorials, but does not work for me
# (it outputs an iso that won't work)
docker run -i -t --rm mattes/boot2docker-vbga > boot2docker.iso

# so I do:
docker run -i -t --rm mattes/boot2docker-vbga /bin/bash
# then in a second shell:
docker cp <Container-ID>:boot2docker.iso boot2docker.iso

# use the new boot2docker.iso
boot2docker stop
mv ~/.boot2docker/boot2docker.iso ~/.boot2docker/boot2docker.iso.backup
mv boot2docker.iso ~/.boot2docker/boot2docker.iso
VBoxManage sharedfolder add boot2docker-vm -name home -hostpath /Users
boot2docker up
boot2docker ssh "ls /Users" # to verify if it worked
,,,
