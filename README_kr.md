# minijs

[![go report][go_report_img]][go_report_url]
[![go doc][go_doc_img]][go_doc_url]
[![release][repo_releases_img]][repo_releases_url]
[![ci][repo_ci_img]][repo_ci_url]
[![code coverage][go_code_coverage_img]][go_code_coverage_url]

**minijs**는 Go 언어로 구현된 자바스크립트 바이트코드 가상 머신입니다. 자바스크립트 코드를 바이트코드로 변환하고 이를 가상 머신에서 실행하여 성능을 최적화하며, **Go와의 높은 호환성**을 바탕으로 다양한 Go 기반 애플리케이션에 유연하게 내장될 수 있습니다.

## 주요 특징

- **바이트코드 실행**: 자바스크립트 코드를 바이트코드로 컴파일하여 가상 머신에서 실행, 성능을 최적화합니다.
- **Go와의 높은 호환성**: Go 언어로 구현되어 Go 환경과의 높은 호환성을 보장하며, 다양한 Go 기반 애플리케이션에 손쉽게 내장 가능합니다.

## 설치

**minijs**는 Go 환경에서 실행됩니다. 아래 명령어를 통해 프로젝트를 로컬에 설치하고 빌드할 수 있습니다:

```bash
git clone https://github.com/siyul-park/minijs.git
cd minijs
make build
```

## 사용법

대화형 셸(REPL)을 이용해 자바스크립트 코드를 실시간으로 실행할 수 있습니다.

### REPL 실행

```bash
minijs
```

### 예시

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