Echo "Setup For next Run"
erase main
erase main.zip

set GOOS=linux

echo "Build the binary - main"
go build -o main main.go

echo "Create a zip file to upload"
%USERPROFILE%\Go\bin\build-lambda-zip.exe -o main.zip main

echo "Cleaning up"
erase main
