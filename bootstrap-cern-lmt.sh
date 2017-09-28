#!/usr/bin/env bash

# DMC Ubuntu repository.
DMC_REPO="deb [allow-insecure=yes trusted=true] http://grid-deployment.web.cern.ch/grid-deployment/dms/dmc/repos/apt/ubuntu/yakkety dmc/"
SOURCES="/etc/apt/sources.list"
# Add repo to sources.list.
grep -qF "${DMC_REPO}" "${SOURCES}" || echo "${DMC_REPO}" >> "${SOURCES}"
# Update sources.
apt-get update
# Install python, voms-clients and gfal2.
apt-get install -fy python voms-clients gfal2-util gfal2-plugin-gridftp gfal2-plugin-file
# Set up golang.
apt-get install -fy git golang-go
mkdir -p /home/ubuntu/go
chown -R ubuntu:ubuntu /home/ubuntu/go
echo "export GOPATH=/home/ubuntu/go" >> /home/ubuntu/.bashrc

# Default SSH user for the yakkety vbox is 'ubuntu'.
SSH_USER="ubuntu"

# Synced folder mount point.
SYNCED_DIR="/vagrant"
GRID_CERTS="${SYNCED_DIR}/grid-security"
# Copy grid-security certs.
cp -nr "${GRID_CERTS}" /etc/
# Copy user's key and certificate to $HOME.
USER_CERT="${SYNCED_DIR}/voms-config/.globus"
cp -nr "${USER_CERT}" "/home/${SSH_USER}"
# Copy VOMS config files.
VOMS_CONF="${SYNCED_DIR}/voms-config/.glite"
cp -nr "${VOMS_CONF}" "/home/${SSH_USER}"
