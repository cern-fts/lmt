%global import_path    gitlab.cern.ch/fts/lmt

Name: lmt
Version: 0.0.1
Release: 1
License: LICENSE
Url: https://gitlab.cern.ch/fts/lmt
Summary: Proxy service in GO
BuildRequires:  %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}
%description
LMT is a proxy service that extends the File Transfer Service in order to enable local data transfers on the WLCG infrastructure.

%prep

%build
%gobuild -o bin/NAME %{import_path}

%install
mkdir -p %{buildroot}/%{_bindir}/root/go/src/gitlab.cern.ch/fts/lmt
cp bin/lmt %{buildroot}/%{_bindir}/root/go/src/gitlab.cern.ch/fts/lmt


%files
%{_bindir}/root/go/src/gitlab.cern.ch/fts/lmt/lmt


%clean

%changelog
* Mon Jan 29 2018  Andrea Manzi <amanzi@cern.ch> - 0.0.1
- first version
