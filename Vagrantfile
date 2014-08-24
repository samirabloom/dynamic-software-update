# Vagrantfile API/syntax version.
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  config.vm.box = "misheska/ubuntu1404-i386"
  config.vm.boot_timeout = 600

  # fix "stdin: is not a tty" error
  config.ssh.shell = "bash -c 'BASH_ENV=/etc/profile exec bash'"

  config.vm.define "docker" do |docker|
    docker.vm.hostname = "docker"
    docker.vm.network "private_network", ip: "192.168.50.5"
    docker.vm.network "forwarded_port", guest: 1234, host: 1234
    docker.vm.network "forwarded_port", guest: 1025, host: 1025
    docker.vm.network "forwarded_port", guest: 1026, host: 1026
    docker.vm.network "forwarded_port", guest: 1027, host: 1027
    docker.vm.provision :shell, :path => "install_docker.sh"
    docker.vm.provider :virtualbox do |vb|
      vb.memory = 2048
      vb.cpus = 3
    end
  end

  config.vm.define "docker_and_go" do |docker_and_go|
    docker_and_go.vm.hostname = "docker-and-go"
    docker_and_go.vm.network "private_network", ip: "192.168.50.7"
    docker_and_go.vm.provision :shell, :path => "install_docker_and_go.sh"
    docker_and_go.vm.provider :virtualbox do |vb|
      vb.memory = 2048
      vb.cpus = 3
    end
  end

  config.vm.define "nginx" do |nginx|
    nginx.vm.hostname = "nginx"
    nginx.vm.network "private_network", ip: "192.168.50.20"
    nginx.vm.provision :shell, :path => "install_nginx.sh"
    nginx.vm.provider :virtualbox do |vb|
      vb.memory = 2048
      vb.cpus = 3
    end
  end

  config.vm.define "wordpress_3_9_1" do |wordpress_3_9_1|
    wordpress_3_9_1.vm.hostname = "wordpress-3.9.1"
    wordpress_3_9_1.vm.network "private_network", ip: "192.168.50.30"
    wordpress_3_9_1.vm.provision :shell, :path => "install_wordpress_3_9_1.sh"
    wordpress_3_9_1.vm.provider :virtualbox do |vb|
      vb.memory = 2048
      vb.cpus = 3
    end
  end


  config.vm.define "wordpress_3_9_2" do |wordpress_3_9_2|
    wordpress_3_9_2.vm.hostname = "wordpress-3.9.2"
    wordpress_3_9_2.vm.network "private_network", ip: "192.168.50.40"
    wordpress_3_9_2.vm.provision :shell, :path => "install_wordpress_3_9_2.sh"
    wordpress_3_9_2.vm.provider :virtualbox do |vb|
      vb.memory = 2048
      vb.cpus = 3
    end
  end

  config.vm.define "couchbase_2_5_1" do |couchbase_2_5_1|
    couchbase_2_5_1.vm.hostname = "couchbase-2.5.1"
    couchbase_2_5_1.vm.network "private_network", ip: "192.168.50.50"
    couchbase_2_5_1.vm.provision :shell, :path => "install_couchbase_2_5_1.sh"
    couchbase_2_5_1.vm.provider :virtualbox do |vb|
      vb.memory = 2048
      vb.cpus = 3
    end
  end

  config.vm.define "couchbase_3_0_0" do |couchbase_3_0_0|
    couchbase_3_0_0.vm.hostname = "couchbase-3.0.0"
    couchbase_3_0_0.vm.network "private_network", ip: "192.168.50.60"
    couchbase_3_0_0.vm.provision :shell, :path => "install_couchbase_3_0_0.sh"
    couchbase_3_0_0.vm.provider :virtualbox do |vb|
      vb.memory = 2048
      vb.cpus = 3
    end
  end

end
