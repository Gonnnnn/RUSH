# About the project

해당 프로젝트는 러쉬의 출석 자동화를 위해 개발된 프로그램입니다.

## 1. Principles

- 해당 프로그램은 금전적 이익을 추구하는 것이 아닌 러쉬 관리자의 편의를 위해 개발되었습니다. 따라서 owner를 포함한 모든 contributor는 기여와 리뷰에 강제성을 가지지 않습니다.
- 코드 소스파일 내부의 모든 주석 및 코드는 모두 영어를 기본으로 작성합니다.
- 해당 문서를 제외한 모든 문서를 최소화합니다. 문서 또한 유지보수돼야하는 리소스이기 때문입니다.

## 2. How to run

```sh
ENV_FILE=.env.local go run rush
```

local ENV file은 관리자에게 문의하세요.

## 3. Structure

### 3.1. Packages

메인 패키지들은 다음과 같습니다.

```sh
.
├── attendance
├── auth
├── golang
├── http
├── job
├── main.go
├── oauth
├── permission
├── server
├── session
├── ui
└── user
```

`main.go`

- Entry point.

`http`

- 라우터 및 HTTP request/resopnse 핸들러
- 요청에 대한 권한은 모두 http에서 먼저 검사합니다.

`permission`

- 권한 관리를 위한 패키지. respository 전체에서 접근 가능한 visibility를 가집니다.

`server`

- 출석 관리를 위한 메인 로직
- 기본적으로 모든 로직은 server에 구현됩니다. 로직이 성장함에 따라 readability가 떨어질 때, 혹은 사용처가 많아지는 중복 로직의 경우 세부 패키지로 분리합니다.

`attendance`, `session`, `user`

- Attendance, session, user 관련 데이터베이스 로직 및 기타 로직

`job`

- 배치로 실행되는 작업

`oauth`

- 외부 인증 서비스 처리 로직

`auth`

- 러쉬 자체 인증 로직

`golang`

- helpers

`ui`

- 러쉬 웹 클라이언트

### 3.2. Deployment

Cloud type

### 3.3. Branch strategy

Trunk based development. 절대 기능 단위로 브랜치를 생성해 유지하지 않습니다.

## 4. Contributing

### 4.1. Issue

필수적으로, Issue 생성을 통해 해당 문제나 안건의 우선순위 및 필요성에 대해 admin과 논의합니다.
하나의 문제 상황에 대한 해결책은 여러가지입니다. 따라서 Issue의 제목은 TODO 형식이 아닌, 문제 상황을 간결하게 설명합니다.
E.g.,

- Wrong: "Fix the session list page design"
- Correct: "Users are not engaging with session items"

Issue의 body에서는 해당 문제에 대한 설명을 간결하지만 명료하게 작성합니다. 주로 배경, 실제 문제, 추측되는 원인, 가능한 해결 방안을 제시합니다. 필요한 경우 스크린샷등을 첨부해 이해를 돕습니다. 요점은, 누구라도 해당 issue를 읽고 어떤 문제가 있는지 이해할 수 있어야합니다.

꼭 어떤 문제를 해결하지 않더라도, 러쉬 및 운영진에게 가치있는 제안을 하는 것은 기여의 일부입니다.

### 4.2. PR

각 문제의 해결 방안은 PR Author에게 달려있습니다. Contributor는 Issue에 대한 문제를 해결하는 PR을 생성합니다.
각 PR은 동작하는 기능 단위여야합니다. 각 PR이 main 브랜치에 병합된 후, 배포됐을 때, 시스템은 항상 동작해야합니다. 이를 위해서 해결 방안을 작게 쪼개는 것을 권장합니다.

각 PR의 크기는 모든 변경사항을 포함해 200자 안팎을 권장하며, 400자가 넘지 않도록 합니다. Issue와 PR은 1:N 관계입니다. 해결 방안을 작게 쪼개는 것을 권장합니다.

PR은 생성시 제공되는 template에 따라 작성합니다.

### 4.3. Code review

[Google code review guideline](https://google.github.io/eng-practices/review/)을 따릅니다.
Author는 자신의 해결 방안 및 코드를 제시할 권리가 있으며, reviewer는 해당 코드를 승인하거나 수정을 요구할 권리가 있습니다. Reviewer의 승인 없이는 해당 PR을 병합할 수 없습니다.

## 5. Owner's note

- 개발환경 세팅, CI, 문서화 등, 시스템의 안정성을 확보하면서 여러 contributor가 쉬운 참여를 하기에는 난이도가 높은 상황입니다. 이점에 유의하여 참여하시길 바랍니다.
