echo 'Cleaning Up...'
rm -rf bin/
mkdir bin


echo 'Building Linux amd64...'
env GOOS=linux GOARCH=amd64 go build -o bin/elara-linux-amd64

echo 'Building Linux x86...'
env GOOS=linux GOARCH=amd64 go build -o bin/elara-linux-x86-64

echo 'Building Windows amd64...'
env GOOS=windows GOARCH=amd64 go build -o bin/elara-windows-amd64.exe

echo 'Done!'
