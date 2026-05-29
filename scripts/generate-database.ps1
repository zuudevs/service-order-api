param(
	[string]$schema_path,
	[string]$output_path
)

function GetAppName {
	return Split-Path $PSCommandPath -Leaf
}

function PrintHelp {
	Write-Host "Usage:" -ForegroundColor Cyan
	Write-Host "$(GetAppName) [schema_path] [output_path]"
}

if (-not $schema_path -or -not $output_path) {
	Write-Host "invalid syntax" -ForegroundColor Red
	PrintHelp
	exit 1
}

sqlite3 $output_path ".read $schema_path"