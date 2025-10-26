$ErrorActionPreference = 'Stop'

$packageName = 'trello-cli'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$filePath = Join-Path $toolsDir 'trello-cli.exe'

# Determine architecture; Chocolatey treats both AMD64 and ARM64 as 64-bit.
$isArm64 = $false
try {
  $procArch = [System.Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture
  if ($procArch -eq 'Arm64') { $isArm64 = $true }
} catch {
  if ($env:PROCESSOR_IDENTIFIER -match 'ARM') { $isArm64 = $true }
}

if ($isArm64) {
  $downloadUrl = 'https://github.com/danbruder/trello-cli/releases/download/v1.0.4/trello-cli-windows-arm64.exe'
  $checksum    = '08637a4bf3255932ff9bfad976d0d99fcc850b36245cb74e07bedd8d6da71e57'
} else {
  $downloadUrl = 'https://github.com/danbruder/trello-cli/releases/download/v1.0.4/trello-cli-windows-amd64.exe'
  $checksum    = '63b79da4499780961b668227003c6e88a1841e07c9d73160179efcc2af8761ac'
}

Get-ChocolateyWebFile -PackageName $packageName `
  -FileFullPath $filePath `
  -Url64bit $downloadUrl `
  -Checksum64 $checksum `
  -ChecksumType64 'sha256'

# Chocolatey will automatically shim trello-cli.exe placed in the tools folder.
