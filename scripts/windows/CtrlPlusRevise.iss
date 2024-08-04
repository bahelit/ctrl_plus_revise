; -- CtrlPlusRevise.iss --

[Setup]
AppName=Ctrl Plus Revise
AppVersion=0.7
WizardStyle=modern
DefaultDirName={autopf}\Ctrl Plus Revise
DefaultGroupName=Ctrl Plus Revise
UninstallDisplayIcon={app}\CtrlPlusRevise.exe
Compression=lzma2
SolidCompression=yes
OutputDir=userdocs:Ctrl Plus Revise Docs
AppPublisher=Ctrl+Revise
AppPublisherURL=https://www.ctrlplusrevise.com/
SetupIconFile=icon.ico

[Files]
Source: "CtrlPlusRevise.exe"; DestDir: "{app}"
Source: "README.md"; DestDir: "{app}"; Flags: isreadme

[Icons]
Name: "{group}\Ctrl Plus Revise"; Filename: "{app}\CtrlPlusRevise.exe"

[Tasks]
Name: StartAfterInstall; Description: Start Ctrl+Revise after install

[Run]
Filename: {app}\CtrlPlusRevise.exe; Flags: shellexec skipifsilent nowait; Tasks: StartAfterInstall