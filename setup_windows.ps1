# 管理者権限チェック
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Error "このスクリプトは管理者権限で実行してください。"
    exit 1
}

# カレントディレクトリの取得
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$configPath = Join-Path $scriptPath "config.json"

# 設定ファイルの読み込み
if (Test-Path $configPath) {
    $config = Get-Content $configPath | ConvertFrom-Json
} else {
    Write-Error "config.jsonが見つかりません。"
    exit 1
}

# PlantUMLコマンドのパスを取得
$plantumlCmd = Join-Path $scriptPath $config.plantuml.command
$plantumlAlias = $config.plantuml.alias

# PowerShellプロファイルのパスを取得
$profilePath = $PROFILE.CurrentUserAllHosts
$profileDir = Split-Path -Parent $profilePath

# プロファイルディレクトリが存在しない場合は作成
if (-not (Test-Path $profileDir)) {
    New-Item -ItemType Directory -Path $profileDir -Force
}

# プロファイルファイルが存在しない場合は作成
if (-not (Test-Path $profilePath)) {
    New-Item -ItemType File -Path $profilePath -Force
}

# 既存のエイリアスを削除（更新のため）
$content = Get-Content $profilePath | Where-Object { $_ -notmatch "function.*$plantumlAlias.*" -and $_ -notmatch "Set-Alias.*$plantumlAlias.*" }
Set-Content $profilePath $content

# 新しいエイリアスを追加
$newFunction = @"
function global:Invoke-PlantUML {
    param([Parameter(ValueFromRemainingArguments=`$true)][string[]]`$args)
    & "$plantumlCmd" `$args
}
Set-Alias -Name $plantumlAlias -Value Invoke-PlantUML -Scope Global
"@

Add-Content $profilePath $newFunction

Write-Host "セットアップが完了しました。"
Write-Host "新しいPowerShellウィンドウを開くか、以下のコマンドを実行してエイリアスを有効にしてください："
Write-Host ". `$PROFILE"

# 環境変数PATHにlibディレクトリを追加
$libPath = Join-Path $scriptPath "lib"
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$libPath*") {
    [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$libPath", "User")
    Write-Host "環境変数PATHにlibディレクトリを追加しました。"
} 