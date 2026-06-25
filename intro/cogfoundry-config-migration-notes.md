# CogFoundry 配置迁移核对备忘

这份文档是本地迁移决策备忘，记录已经核对属实、后续还需要处理或确认的配置差异。文件已加入 `.gitignore`，不会提交到仓库。

## 运行时配置

### `LOOMLOOM_SERVER`

状态：属实。

含义：CogFoundry Batch API 的基础地址。凡是需要请求服务端的 CLI 命令，都依赖这个地址。

出现位置：

- `README.md:38` 提醒用户替换 `<your LoomLoom server URL>`。
- `README.md:42` 示例中写了 `My server URL is <your LoomLoom server URL>`。
- `README.md:96` 示例中写了 `export LOOMLOOM_SERVER="<your LoomLoom server URL>"`。
- `install.sh:346` 安装完成后提示 `export LOOMLOOM_SERVER=<your LoomLoom server URL>`。
- `install.ps1:131` Windows 安装完成后提示 `$env:LOOMLOOM_SERVER='<your LoomLoom server URL>'`。
- `cli/internal/cmd/root.go:21` 实际读取 `LOOMLOOM_SERVER`，并兼容 `BATCHJOB_SERVER`。
- `cli/internal/client/http.go:51` 没有 server URL 时会报错。
- `cli/internal/cmd/doctor.go:67` 提示用户设置 `LOOMLOOM_SERVER` 和 `LOOMLOOM_TOKEN`。
- `docs/template-spec/03-conversational-authoring.md:57` 提到需要提供 `LOOMLOOM_SERVER` 和 `LOOMLOOM_TOKEN`。
- `cli/internal/template_spec_docs/template-spec/03-conversational-authoring.md:57` 是内嵌文档副本。
- `skills/codex/loomloom/docs/template-spec/03-conversational-authoring.md:57` 是 skill 文档副本。
- `skills/claude/loomloom/docs/template-spec/03-conversational-authoring.md:57` 是 skill 文档副本。
- `skills/openclaw/loomloom/docs/template-spec/03-conversational-authoring.md:57` 是 skill 文档副本。

影响场景：

- 这是 CLI 正常访问 CogFoundry 服务的必要配置。
- 需要替换为真实 CogFoundry Batch API URL，或者在文档里明确说明用户应填写自己的 workspace API URL。
- 这不只是 README 占位问题，而是核心运行时入口。

### `LOOMLOOM_TOKEN`

状态：属实。

含义：CLI 访问 CogFoundry 服务时使用的认证 token。

出现位置：

- `README.md:38` 提醒用户替换 `your-token`。
- `README.md:97` 示例中写了 `export LOOMLOOM_TOKEN="your-token"`。
- `README.md:102` 写着从 CogFoundry workspace 获取 token。
- `README.md:300` 写着从 CogFoundry workspace 获取 token。
- `install.sh:347` 安装完成后提示 `export LOOMLOOM_TOKEN=your-token`。
- `install.ps1:132` Windows 安装完成后提示 `$env:LOOMLOOM_TOKEN='your-token'`。
- `cli/internal/cmd/root.go:22` 实际读取 `LOOMLOOM_TOKEN`，并兼容 `BATCHJOB_TOKEN`。
- `cli/internal/cmd/doctor.go:67` 提示用户设置 `LOOMLOOM_SERVER` 和 `LOOMLOOM_TOKEN`。
- `skills/codex/loomloom/SKILL.md:22` 提到缺少 `LOOMLOOM_TOKEN` 的情况。
- `skills/claude/loomloom/SKILL.md:22` 提到缺少 `LOOMLOOM_TOKEN` 的情况。
- `skills/openclaw/loomloom/SKILL.md:22` 提到缺少 `LOOMLOOM_TOKEN` 的情况。

影响场景：

- 实际访问 CogFoundry 服务时通常必须配置。
- token 获取页面或操作步骤属于文档/用户引导内容，不是 CLI 代码里写死的逻辑。

### `BATCHJOB_SERVER`

状态：属实。

含义：旧版服务地址环境变量。

出现位置：

- `README.md:100` 说明兼容 `BATCHJOB_SERVER`。
- `cli/internal/cmd/root.go:21` 在 `LOOMLOOM_SERVER` 为空时 fallback 到 `BATCHJOB_SERVER`。

影响场景：

- 新的 CogFoundry 用户不需要主动使用。
- 保留它可以兼容旧脚本。
- 直接删除会造成破坏性变更。

### `BATCHJOB_TOKEN`

状态：属实。

含义：旧版 token 环境变量。

出现位置：

- `README.md:100` 说明兼容 `BATCHJOB_TOKEN`。
- `cli/internal/cmd/root.go:22` 在 `LOOMLOOM_TOKEN` 为空时 fallback 到 `BATCHJOB_TOKEN`。

影响场景：

- 新的 CogFoundry 用户不需要主动使用。
- 保留它可以兼容旧脚本。
- 直接删除会造成破坏性变更。

### `LOOMLOOM_CLI_RELEASE_API`

状态：属实。

含义：CLI 自更新或版本检查时使用的 release API 覆盖地址。

出现位置：

- `cli/internal/version/version.go:90` 读取 `LOOMLOOM_CLI_RELEASE_API`。
- `cli/internal/version/version_test.go:61` 在测试中清理该变量。
- `cli/internal/version/version_test.go:71` 在测试中设置该变量。

影响场景：

- 不影响普通 CLI 命令运行。
- 主要用于测试或覆盖默认 GitHub release API。

### `BATCHJOB_CLI_RELEASE_API`

状态：属实。

含义：旧版 release API 覆盖地址。

出现位置：

- `cli/internal/version/version.go:93` 在 `LOOMLOOM_CLI_RELEASE_API` 为空时读取 `BATCHJOB_CLI_RELEASE_API`。
- `cli/internal/version/version_test.go:62` 在测试中清理该变量。

影响场景：

- 新的 CogFoundry 用户不需要主动使用。
- 保留它可以兼容旧脚本。

## 文档和 workspace URL

### CogFoundry Console URL

状态：属实。

含义：用户查看 workspace、运行进度或运行结果的 CogFoundry 控制台地址。

出现位置：

- `README.md:197` 写着使用 CogFoundry workspace console URL 查看运行进度。
- `README.md:315` troubleshooting 中列了 Console 地址。
- `skills/codex/loomloom/SKILL.md:302` 说明 console 链接必须来自用户提供的 CogFoundry workspace 信息。
- `skills/codex/loomloom/SKILL.md:306` 提到如果用户提供了 CogFoundry console URL，可以让用户去那里看运行状态。
- `skills/codex/loomloom/SKILL.md:307` 提到没有 console URL 时，可以询问用户。
- `skills/claude/loomloom/SKILL.md:185` 有相同的 console 引导。
- `skills/openclaw/loomloom/SKILL.md:185` 有相同的 console 引导。

影响场景：

- 不影响 CLI 命令执行。
- 影响 README、skill、助手引导里给用户的跳转说明。
- 需要等真实 CogFoundry workspace console URL 确定后再补。

### CogFoundry Workflow Runs URL

状态：属实，当前没有固定配置。

含义：如果 CogFoundry 有单独的 workflow run 详情页，这里指那个直接查看运行记录的 URL。

出现位置：

- 当前没有固定的 Workflow Runs URL。
- 现有文档只泛指 CogFoundry workspace console。

影响场景：

- 不影响 CLI 命令执行。
- 只影响文档和运行结果交接体验。

### CogFoundry Token 获取 URL 或操作步骤

状态：属实，目前还是占位级说明。

含义：用户生成或复制 CogFoundry token 的页面地址，或者准确操作步骤。

出现位置：

- `README.md:102` 写着 `Get the token from your CogFoundry workspace.`
- `README.md:300` 写着从 CogFoundry workspace 获取 token。

影响场景：

- 如果用户已经有 token，不影响 CLI 二进制运行。
- 但 README 上手体验需要补真实步骤，否则新用户不知道 token 从哪里拿。

### CogFoundry 官网 URL

状态：属实，当前缺失。

含义：CogFoundry 对外官网地址。

出现位置：

- README 中当前没有正式官网 URL。
- `README.md:4` 现在链接的是 `github.com/Cogfoundry-ai/loomloom`，这是 GitHub 仓库，不是官网。

影响场景：

- 不影响 CLI 命令执行。
- 影响文档完整性和产品对外展示。

## Gitee 分发

### `GITEE_REPO`

状态：属实。

含义：安装脚本在选择 `--source gitee` 时使用的 Gitee 仓库。

出现位置：

- `install.sh:5` 默认值是 `shengsuanyun/loomloom`。
- `install.ps1:15` 默认值是 `shengsuanyun/loomloom`。

影响场景：

- 不影响默认 GitHub 安装。
- 不影响已经安装好的 CLI 运行。
- 会影响使用 `--source gitee` 安装的用户。
- 如果不处理，Gitee 安装路径仍会指向旧的胜算云仓库。

### `LOOMLOOM_RELEASE_SOURCE`

状态：属实。

含义：选择安装时从 GitHub 还是 Gitee 下载 release。

出现位置：

- `install.sh:6` 默认值是 `github`。
- `install-gitee.sh:5` 默认值设为 `gitee`。

影响场景：

- 不影响已安装 CLI 的运行。
- 影响安装来源和 release asset 下载来源。

### `GITEE_INSTALL_URL`

状态：属实。

含义：`install-gitee.sh` 用来下载真正安装脚本的地址。

出现位置：

- `install-gitee.sh:4` 默认值是 `https://gitee.com/shengsuanyun/loomloom/raw/main/install.sh`。

影响场景：

- 不影响普通 CLI 运行。
- 影响 `install-gitee.sh`。
- 如果公开暴露，用户仍会被导向旧的胜算云 Gitee 仓库。

### `LOOMLOOM_GITEE_TOKEN`

状态：属实。

含义：Gitee release 发布时使用的 token。

出现位置：

- `.workflow/loomloom-cli-release.yml:13` 声明了 `LOOMLOOM_GITEE_TOKEN`。
- `scripts/ci-release-gitee.sh:18` 要求必须存在。
- `scripts/ci-release-gitee.sh:23` 会输出 token 已存在。
- `scripts/publish-gitee-release.sh:19` 说明它是必需 token。
- `scripts/publish-gitee-release.sh:59` 读取它，并兼容 `GITEE_TOKEN`。
- `scripts/publish-gitee-release.sh:61` 缺失时会报错。

影响场景：

- 不影响 CLI 运行。
- 只有 CogFoundry 继续保留 Gitee release 发布时才需要。

### `GITEE_TOKEN`

状态：属实。

含义：`LOOMLOOM_GITEE_TOKEN` 的本地别名。

出现位置：

- `scripts/publish-gitee-release.sh:20` 文档中说明了它。
- `scripts/publish-gitee-release.sh:59` 实际读取它。

影响场景：

- 不影响 CLI 运行。
- 只影响手动或 CI 的 Gitee 发布。

### `GITEE_OWNER`

状态：属实。

含义：Gitee 发布脚本使用的 owner。

出现位置：

- `scripts/publish-gitee-release.sh:4` 默认值是 `shengsuanyun`。
- `scripts/publish-gitee-release.sh:21` 文档中说明默认值是 `shengsuanyun`。

影响场景：

- 不影响 CLI 运行。
- 如果保留 Gitee 发布，需要替换为 CogFoundry 的 Gitee owner，或者由 CI 显式传入。

### `GITEE_REPO_NAME`

状态：属实。

含义：Gitee 发布脚本使用的仓库名。

出现位置：

- `scripts/publish-gitee-release.sh:5` 默认值是 `loomloom`。
- `scripts/publish-gitee-release.sh:22` 文档中说明默认值是 `loomloom`。

影响场景：

- 不影响 CLI 运行。
- 影响 Gitee release 发布目标。

### `GITEE_API_BASE`

状态：属实。

含义：Gitee release 发布脚本使用的 API 地址。

出现位置：

- `scripts/publish-gitee-release.sh:10` 默认值是 `https://gitee.com/api/v5`。
- `scripts/publish-gitee-release.sh:25` 文档中说明默认值。

影响场景：

- 不影响 CLI 运行。
- 只影响 Gitee release 发布。

### `.workflow/loomloom-cli-release.yml`

状态：属实。

含义：Gitee 或 Codeup 风格的 release 流水线配置。

出现位置：

- `.workflow/loomloom-cli-release.yml:10` 允许 `gitee-release-test.*` 标签触发。
- `.workflow/loomloom-cli-release.yml:21` 定义了 `build_and_publish_gitee_release`。
- `.workflow/loomloom-cli-release.yml:27` 执行 `scripts/ci-release-gitee.sh`。

影响场景：

- 不影响 CLI 运行。
- 影响 Gitee 侧 CI release 发布。

## Homebrew 分发

### `LOOMLOOM_HOMEBREW_TAP`

状态：属实。

含义：安装脚本可选使用的 Homebrew tap。

出现位置：

- `install.sh:13` 读取 `LOOMLOOM_HOMEBREW_TAP`。
- `install.sh:198` 判断是否可以使用 Homebrew。
- `install.sh:203` 要求 tap 非空。
- `install.sh:276` 在条件满足时选择 Homebrew 安装。
- `install.sh:303` 执行 `brew install "$HOMEBREW_TAP/loomloom"`。

影响场景：

- 不影响 CLI 运行。
- 默认安装不受影响，除非显式设置了 `LOOMLOOM_HOMEBREW_TAP`。
- 真实 CogFoundry tap 没准备好之前，应该保持禁用。

### `HOMEBREW_TAP_GITHUB_TOKEN`

状态：属实。

含义：GitHub release CI 更新 Homebrew tap 时使用的 token。

出现位置：

- `.github/workflows/release.yml:61` 从 GitHub secrets 中读取它。
- `.github/workflows/release.yml:64` 缺失时会报错。
- `.github/workflows/release.yml:70` 用它 clone `github.com/Cogfoundry-ai/homebrew-tap.git`。

影响场景：

- 不影响 CLI 运行。
- 影响稳定版 GitHub release 发布。
- 如果 secret 和 tap 仓库不存在，稳定版发布会在 Homebrew 更新步骤失败。

### `Cogfoundry-ai/homebrew-tap`

状态：属实。

含义：当前 release workflow 预期存在的 CogFoundry Homebrew tap 仓库。

出现位置：

- `.github/workflows/release.yml:70` clone `github.com/Cogfoundry-ai/homebrew-tap.git`。
- `.github/workflows/release.yml:71` 写入 `Formula/loomloom.rb`。

影响场景：

- 不影响 CLI 运行。
- 只有启用 Homebrew 发布时才需要。

### `scripts/render-homebrew-formula.sh`

状态：属实。

含义：生成 `Formula/loomloom.rb` 内容的脚本。

出现位置：

- `scripts/render-homebrew-formula.sh:6` 说明需要 `--tag` 和 `--checksums`。
- `.github/workflows/release.yml:71` 在稳定版发布时调用它。

影响场景：

- 不影响 CLI 运行。
- 只影响 Homebrew release 发布。

## 旧服务域名单测样例

### `loomloom.shengsuanyun.com/batch`

状态：当前目标项目中已经不存在。

出现位置：

- 当前在 `loomloom-cogfoundry` 中没有匹配。
- 之前文档里说 `cli/internal/client/http_test.go` 还有这个域名，这个说法已经过期。

影响场景：

- 这项当前不需要处理。
- 目前剩余的 `shengsuanyun` 主要集中在 Gitee 分发默认值，不在 CLI 核心运行逻辑中。
