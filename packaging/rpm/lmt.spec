%global import_path    gitlab.cern.ch/fts/lmt

%if %{?fedora}%{!?fedora:0} >= 17 || %{?rhel}%{!?rhel:0} >= 7
%global systemd 1
%else
%global systemd 0
%endif

Name: lmt
Version: 0.0.1
Release: 1
License: Apache 2.0
Url: https://gitlab.cern.ch/fts/lmt
Summary: FT Last mile proxy service written in Proxy in GO
Source0: %{name}-%{version}.tar.gz

BuildRequires:  %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}
%description
LMT is a proxy service that extends the File Transfer Service in order to enable local data transfers on the WLCG infrastructure.

%prep
%setup -q -n %{name}-%{version}

%build
%define debug_package %{nil}
mkdir -p src/gitlab.cern.ch/fts/
ln -s ../../../ src/gitlab.cern.ch/fts/lmt

export GOPATH=$(pwd):%{gopath}
go build -o bin/lmt %import_path

%install
mkdir -p %{buildroot}/%{_sysconfdir}/lmt
mkdir -p %{buildroot}/%{_sbindir}
%if %systemd
mkdir -p %{buildroot}/%{_unitdir}
cp etc/lmt.service %{buildroot}/%{_unitdir}
%endif
cp config.yml %{buildroot}/%{_sysconfdir}/lmt
cp bin/lmt %{buildroot}/%{_sbindir}/lmt

%files
%{_sbindir}/lmt
%{_sysconfdir}/lmt/config.yml
%if %systemd
%attr(0644,root,root) %{_unitdir}/lmt.service
%endif

%clean

%changelog
* Mon Jan 29 2018  Andrea Manzi <amanzi@cern.ch> - 0.0.1
- first version
