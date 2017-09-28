# -*- mode: ruby -*-
# vi: set ft=ruby :
#
# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "ubuntu/yakkety64"
  config.vm.hostname = "cern-lmt"
  config.vm.network :private_network, type: "dhcp"
  config.vm.network "forwarded_port", guest: 8080, host: 8080
  config.vm.synced_folder ".", "/home/ubuntu/go/src/gitlab.cern.ch/fts/lmt", :mount_options => ["dmode=775", "fmode=666"]
  config.vm.provision :shell, path: "bootstrap-cern-lmt.sh"
end
