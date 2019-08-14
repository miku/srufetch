Summary:    Simple SRU fetch tool.
Name:       srufetch
Version:    0.1.2
Release:    0
License:    GPL
ExclusiveArch:  x86_64
BuildRoot:  %{_tmppath}/%{name}-build
Group:      System/Base
Vendor:     Leipzig University Library, https://www.ub.uni-leipzig.de
URL:        https://github.com/miku/srufetch

%description

Simple SRU fetch tool.

%prep

%build

%pre

%install

mkdir -p $RPM_BUILD_ROOT/usr/local/bin
install -m 755 srufetch $RPM_BUILD_ROOT/usr/local/bin

%post

%clean
rm -rf $RPM_BUILD_ROOT
rm -rf %{_tmppath}/%{name}
rm -rf %{_topdir}/BUILD/%{name}

%files
%defattr(-,root,root)

/usr/local/bin/srufetch

%changelog

* Mon Aug 12 2019 Martin Czygan
- 0.1.0 release

