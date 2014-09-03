# Vagrantfile API/syntax version.
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  config.vm.box = "misheska/ubuntu1404-i386"
  config.vm.boot_timeout = 600

  # fix "stdin: is not a tty" error
  config.ssh.shell = "bash -c 'BASH_ENV=/etc/profile exec bash'"

  config.vm.define "docker_one" do |docker_one|
    docker_one.vm.hostname = "dockerone"
    docker_one.vm.network "private_network", ip: "192.168.50.5"
    docker_one.vm.provision :shell, :path => "install_docker.sh"
    docker_one.vm.provider :virtualbox do |vb|
      vb.memory = 2048
      vb.cpus = 3
    end
  end


  config.vm.define "docker_two" do |docker_two|
    docker_two.vm.hostname = "dockertwo"
    docker_two.vm.network "private_network", ip: "192.168.50.7"
    docker_two.vm.provision :shell, :path => "install_docker.sh"
    docker_two.vm.provider :virtualbox do |vb|
      vb.memory = 2048
      vb.cpus = 3
    end
  end

end
