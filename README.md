# minijs

[![go report][go_report_img]][go_report_url]
[![go doc][go_doc_img]][go_doc_url]
[![release][repo_releases_img]][repo_releases_url]
[![ci][repo_ci_img]][repo_ci_url]
[![code coverage][go_code_coverage_img]][go_code_coverage_url]

**minijs** is a JavaScript bytecode virtual machine implemented in Go. It converts JavaScript code into bytecode for execution and, with **high compatibility with Go**, can be seamlessly embedded into various Go-based applications.

## **Key Features**

- **Bytecode Execution**: Optimizes performance by converting JavaScript code into bytecode and executing it in a virtual machine.
- **High Compatibility with Go**: Developed in Go, making it easy to integrate into Go-based applications.

## **Installation**

**minijs** runs in a Go environment. You can clone the repository and build it using the following commands:

```bash
git clone https://github.com/siyul-park/minijs.git  
cd minijs  
make build  
```

## **Usage**

### **Running the REPL**

You can execute JavaScript code interactively in the REPL (Read-Eval-Print Loop).

```bash
minijs  
```

```bash
> 'b'+'a'+ +'a'+'a'  
"baNaNa"  
```

### **Bytecode Output**

To print the corresponding bytecode, use the `-print-bytecode` flag.

```bash
minijs --print-bytecode  
```

```bash
> 'b'+'a'+ +'a'+'a'  
section .text:
        str.load 0x00000000 0x00000001
        str.load 0x00000002 0x00000001
        str.add
        str.load 0x00000002 0x00000001
        str.to_f64
        f64.to_str
        str.add
        str.load 0x00000002 0x00000001
        str.add
        pop

.section .data:
        b
        a

"baNaNa"  
```

### **Executing a JavaScript File**

When executing a file, **minijs** applies optimization processes to make the bytecode more efficient. To run a JavaScript file, use the following command:

```bash
minijs banana.js  
```

### **Printing Bytecode from a File**

To print the bytecode while executing a file, use the `-print-bytecode` flag.

```bash
minijs -print-bytecode banana.js  
```

```text
section .text:
        str.load 0x00000000 0x00000006
        pop

.section .data:
        baNaNa
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