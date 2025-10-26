$ErrorActionPreference = 'Stop'

$packageName = 'trlo'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url = 'https://github.com/danbruder/trlo/releases/download/v1.0.0/trlo-windows-amd64.exe'
$url64 = 'https://github.com/danbruder/trlo/releases/download/v1.0.0/trlo-windows-arm64.exe'

$packageArgs = @{
  packageName   = $packageName
  unzipLocation = $toolsDir
  fileType      = 'exe'
  url           = $url
  url64bit      = $url64
  softwareName  = 'trlo*'
  checksum      = ''
  checksumType  = 'sha256'
  checksum64    = ''
  checksumType64= 'sha256'
}

Install-ChocolateyPackage @packageArgs
