$ErrorActionPreference = 'Stop'

$packageName = 'trello-cli'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url = 'https://github.com/danbruder/trello-cli/releases/download/v1.0.0/trello-cli-windows-amd64.exe'
$url64 = 'https://github.com/danbruder/trello-cli/releases/download/v1.0.0/trello-cli-windows-arm64.exe'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  fileType      = 'exe'
  url           = $url
  url64bit      = $url64
  softwareName  = 'trello-cli*'
  checksum      = ''
  checksumType  = 'sha256'
  checksum64    = ''
  checksumType64= 'sha256'
}

Install-ChocolateyPackage @packageArgs
