# minijs

[![go report][go_report_img]][go_report_url]
[![go doc][go_doc_img]][go_doc_url]
[![release][repo_releases_img]][repo_releases_url]
[![ci][repo_ci_img]][repo_ci_url]
[![code coverage][go_code_coverage_img]][go_code_coverage_url]

**minijs** is a JavaScript bytecode virtual machine implemented in Go. It compiles JavaScript code into bytecode and executes it in a virtual machine to optimize performance. With **high compatibility with Go**, it can be seamlessly embedded into a wide range of Go-based applications.

## Key Features

- **Bytecode Execution**: Compiles JavaScript code into bytecode for execution in the virtual machine, optimizing performance.
- **High Compatibility with Go**: Implemented in Go, it offers high compatibility with Go environments and can be easily embedded into various Go-based applications.

## Installation

**minijs** runs in a Go environment. Use the following commands to install and build the project locally:

```bash
git clone https://github.com/siyul-park/minijs.git
cd minijs
make build
```

## Usage

You can run JavaScript code interactively using the REPL (Read-Eval-Print Loop).

### Running REPL

```bash
minijs
```

### Example

```bash
> 'b'+'a'+ +'a'+'a'
"baNaNa"
```

<!-- Go -->

[go_download_url]: https://golang.org/dl/
[go_version_img]: https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go
[go_code_coverage_img]: https://codecov.io/gh/siyul-park/minijs/graph/badge.svg?token=quEl9AbBcW
[go_code_coverage_url]: https://codecov.io/gh/siyul-park/minijs
[go_report_img]: https://goreportcard.com/badge/github.com/siyul-park/minijs
[go_report_url]: https://goreportcard.com/report/github.com/siyul-park/minijs
[go_doc_img]: https://godoc.org/github.com/siyul-park/minijs?status.svg
[go_doc_url]: https://godoc.org/github.com/siyul-park/minijs

<!-- Repository -->

[repo_url]: https://github.com/siyul-park/minijs
[repo_issues_url]: https://github.com/siyul-park/minijs/issues
[repo_pull_request_url]: https://github.com/siyul-park/minijs/pulls
[repo_discussions_url]: https://github.com/siyul-park/minijs/discussions
[repo_releases_img]: https://img.shields.io/github/release/siyul-park/minijs.svg
[repo_releases_url]: https://github.com/siyul-park/minijs/releases
[repo_wiki_url]: https://github.com/siyul-park/minijs/wiki
[repo_wiki_img]: https://img.shields.io/badge/docs-wiki_page-blue?style=for-the-badge&logo=none
[repo_wiki_faq_url]: https://github.com/siyul-park/minijs/wiki/FAQ
[repo_ci_img]: https://github.com/siyul-park/minijs/actions/workflows/ci.yml/badge.svg
[repo_ci_url]: https://github.com/siyul-park/minijs/actions/workflows/ci.yml