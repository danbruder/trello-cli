Name:           trlo
Version:        1.0.0
Release:        1%{?dist}
Summary:        A comprehensive Trello CLI tool optimized for LLM integration

License:        MIT
URL:            https://github.com/danbruder/trello-cli
Source0:        https://github.com/danbruder/trello-cli/archive/v%{version}.tar.gz

BuildArch:      noarch
Requires:       glibc

%description
A comprehensive Trello CLI tool built in Go that provides full access to Trello's API
with features optimized for LLM integration including context optimization, batch operations,
and flexible output formats.

Features:
- Full CRUD operations on boards, lists, cards, labels, checklists, members, and attachments
- LLM-optimized output formats (Markdown and JSON)
- Batch operations support
- Context optimization with token limits and field filtering
- Flexible authentication (environment variables, config file, command-line flags)
- Scripting support with quiet mode
- Cross-platform support

%prep
%setup -q

%build
# Binary will be downloaded from GitHub releases
# No build step needed

%install
mkdir -p %{buildroot}%{_bindir}
# Download and install the binary
curl -L -o %{buildroot}%{_bindir}/trlo https://github.com/danbruder/trello-cli/releases/download/v%{version}/trlo-linux-amd64
chmod +x %{buildroot}%{_bindir}/trlo

%files
%{_bindir}/trlo

%changelog
* $(date +'%a %b %d %Y') Dan Bruder <dan@example.com> - 1.0.0-1
- Initial release
