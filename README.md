# minijs

[![go report][go_report_img]][go_report_url]  
[![go doc][go_doc_img]][go_doc_url]  
[![release][repo_releases_img]][repo_releases_url]  
[![ci][repo_ci_img]][repo_ci_url]  
[![code coverage][go_code_coverage_img]][go_code_coverage_url]

**minijs** is a JavaScript bytecode virtual machine implemented in Go. It compiles JavaScript code into bytecode for execution, optimizing performance while ensuring **high compatibility with Go**, making it easy to embed in various Go-based applications.

## Key Features

- **Bytecode Execution**: Transforms JavaScript code into bytecode and executes it in a virtual machine, optimizing performance.
- **High Compatibility with Go**: Developed in Go, allowing seamless integration with Go-based applications.

## Installation

**minijs** runs in a Go environment. Use the following commands to clone and build the project:

```bash
git clone https://github.com/siyul-park/minijs.git
cd minijs
make build
```

## Usage

### Running REPL

You can execute JavaScript code interactively using the REPL (Read-Eval-Print Loop):

```bash
minijs
```

```bash
> 'b'+'a'+ +'a'+'a'
"baNaNa"
```

### Executing a Script File

To run a JavaScript script from a file, use:

```bash
minijs script.js
```

### Printing Bytecode

To output the compiled bytecode, use the `-print-bytecode` flag:

```bash
minijs -print-bytecode script.js
```

```text
section .text:
 global _main

_main:
 0000   sload 0x0 0x0
 0009   sload 0x2 0x2
 0018   sadd
 0019   sload 0x2 0x2
 0028   s2f64
 0029   f642s
 0030   sadd
 0031   sload 0x2 0x2
 0040   sadd
 0041   pop

.section .data:
 0000   b
 0002   a
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