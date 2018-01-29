%global import_path    gitlab.cern.ch/fts/lmt

Name: lmt
Version: 0.0.1
Release: 1
License: LICENSE
Url: https://gitlab.cern.ch/fts/lmt
Summary: Proxy service in GO
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
go build -o bin/lmt gitlab.cern.ch/fts/lmt

%install
mkdir -p %{buildroot}/%{_bindir}/root/go/src/gitlab.cern.ch/fts/lmt
cp bin/lmt %{buildroot}/%{_bindir}/root/go/src/gitlab.cern.ch/fts/lmt


%files
%{_bindir}/root/go/src/gitlab.cern.ch/fts/lmt/lmt

%clean

%changelog
* Mon Jan 29 2018  Andrea Manzi <amanzi@cern.ch> - 0.0.1
- first version
