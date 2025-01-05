# libディレクトリが存在しない場合は作成
if (!(Test-Path "lib")) {
    New-Item -ItemType Directory -Path "lib"
} 