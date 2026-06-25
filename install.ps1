param(
  [string]$Agent = "codex",
  [string]$InstallDir = "$HOME\AppData\Local\Programs\loomloom",
  [string]$SkillDir = "",
  [string]$Version = "latest",
  [ValidateSet("stable", "beta", "rc", "internal")]
  [string]$Channel = "stable",
  [ValidateSet("github", "gitee")]
  [string]$Source = "github"
)

$ErrorActionPreference = "Stop"

$GithubRepo = "Cogfoundry-ai/loomloom"
$GiteeRepo = if ($env:GITEE_REPO) { $env:GITEE_REPO } else { "shengsuanyun/loomloom" }
$Repo = if ($Source -eq "gitee") { $GiteeRepo } else { $GithubRepo }
$ApiBase = if ($Source -eq "gitee") { "https://gitee.com/api/v5/repos/$Repo" } else { "https://api.github.com/repos/$Repo" }

function Get-ReleaseHeaders {
  if ($Source -eq "gitee") {
    return @{ Accept = "application/json"; "User-Agent" = "loomloom-installer" }
  }
  return @{ Accept = "application/vnd.github+json"; "User-Agent" = "loomloom-installer" }
}

function Resolve-SkillDir {
  param([string]$AgentName, [string]$Override)
  if ($Override) { return $Override }
  switch ($AgentName) {
    "codex" { return "$HOME\.codex\skills\loomloom" }
    "claude" { return "$HOME\.claude\skills\loomloom" }
    "openclaw" { return "$HOME\.openclaw\workspace\skills\loomloom" }
    default { throw "unsupported agent: $AgentName" }
  }
}

function Resolve-Tag {
  param([string]$Requested, [string]$ChannelName)
  if ($Requested -ne "latest") { return $Requested }
  if ($ChannelName -ne "stable") {
    $releases = Invoke-RestMethod -Uri "$ApiBase/releases?per_page=100" -Headers (Get-ReleaseHeaders)
    $pattern = "^v[0-9]+\.[0-9]+\.[0-9]+-$ChannelName\.[0-9]+$"
    $release = @($releases | Where-Object { $_.prerelease -and $_.tag_name -match $pattern } | Select-Object -First 1)
    if (-not $release -or -not $release[0].tag_name) { throw "failed to resolve latest $ChannelName release tag" }
    return [string]$release[0].tag_name
  }
  $resp = Invoke-RestMethod -Uri "$ApiBase/releases/latest" -Headers (Get-ReleaseHeaders)
  if (-not $resp.tag_name) { throw "failed to resolve latest release tag" }
  return [string]$resp.tag_name
}

function Get-ChecksumMap {
  param([string]$ChecksumsPath)
  $map = @{}
  Get-Content $ChecksumsPath | ForEach-Object {
    if ($_ -match '^([0-9a-fA-F]+)\s+(.+)$') {
      $map[$matches[2]] = $matches[1].ToLowerInvariant()
    }
  }
  return $map
}

function Assert-Checksum {
  param(
    [string]$AssetName,
    [string]$FilePath,
    [hashtable]$ChecksumMap
  )
  if (-not $ChecksumMap.ContainsKey($AssetName)) { return }
  $actual = (Get-FileHash -Path $FilePath -Algorithm SHA256).Hash.ToLowerInvariant()
  $expected = $ChecksumMap[$AssetName]
  if ($actual -ne $expected) {
    throw "checksum mismatch for $AssetName"
  }
}

$arch = switch ($env:PROCESSOR_ARCHITECTURE.ToLowerInvariant()) {
  "amd64" { "amd64" }
  "arm64" { "arm64" }
  default { throw "unsupported architecture: $env:PROCESSOR_ARCHITECTURE" }
}

$tag = Resolve-Tag -Requested $Version -ChannelName $Channel
$cliAsset = "loomloom-windows-$arch.zip"
$skillsAsset = "loomloom-skills.zip"
$checksumsAsset = "checksums.txt"
$baseUrl = if ($Source -eq "gitee") { "https://gitee.com/$Repo/releases/download/$tag" } else { "https://github.com/$Repo/releases/download/$tag" }

$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ("LoomLoom-" + [System.Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tmpDir | Out-Null
try {
  $cliZip = Join-Path $tmpDir $cliAsset
  $skillsZip = Join-Path $tmpDir $skillsAsset
  $checksumsPath = Join-Path $tmpDir $checksumsAsset

  Write-Host "LoomLoom installer"
  Write-Host "repo: $Repo"
  Write-Host "source: $Source"
  Write-Host "version: $tag"
  Write-Host "channel: $Channel"
  Write-Host "agent: $Agent"
  Write-Host "install dir: $InstallDir"
  Write-Host "skill dir: $(Resolve-SkillDir -AgentName $Agent -Override $SkillDir)"
  Write-Host ""

  Invoke-WebRequest -Uri "$baseUrl/$cliAsset" -OutFile $cliZip
  Invoke-WebRequest -Uri "$baseUrl/$checksumsAsset" -OutFile $checksumsPath
  $checksumMap = Get-ChecksumMap -ChecksumsPath $checksumsPath
  Assert-Checksum -AssetName $cliAsset -FilePath $cliZip -ChecksumMap $checksumMap

  $cliExtract = Join-Path $tmpDir "cli"
  Expand-Archive -LiteralPath $cliZip -DestinationPath $cliExtract -Force
  New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
  Copy-Item -LiteralPath (Join-Path $cliExtract "loomloom.exe") -Destination (Join-Path $InstallDir "loomloom.exe") -Force

  Invoke-WebRequest -Uri "$baseUrl/$skillsAsset" -OutFile $skillsZip
  Assert-Checksum -AssetName $skillsAsset -FilePath $skillsZip -ChecksumMap $checksumMap

  $skillsExtract = Join-Path $tmpDir "skills"
  Expand-Archive -LiteralPath $skillsZip -DestinationPath $skillsExtract -Force
  $finalSkillDir = Resolve-SkillDir -AgentName $Agent -Override $SkillDir
  New-Item -ItemType Directory -Force -Path $finalSkillDir | Out-Null
  Copy-Item -Path (Join-Path $skillsExtract "skills\$Agent\loomloom\*") -Destination $finalSkillDir -Recurse -Force

  Write-Host "installed:"
  Write-Host "  $(Join-Path $InstallDir 'loomloom.exe')"
  Write-Host "  $(Join-Path (Resolve-SkillDir -AgentName $Agent -Override $SkillDir) 'SKILL.md')"
  Write-Host ""
  Write-Host "next:"
  Write-Host "  Add $InstallDir to PATH if needed"
  Write-Host "  `$env:LOOMLOOM_SERVER='<your LoomLoom server URL>'"
  Write-Host "  `$env:LOOMLOOM_TOKEN='your-token'"
  Write-Host "  loomloom doctor"
}
finally {
  if (Test-Path $tmpDir) {
    Remove-Item -LiteralPath $tmpDir -Recurse -Force
  }
}
