## Ironite
#### Server List For Cubyz

**WARNING**: Ironite is currently a WIP, code is not production-ready. Please report any bugs you find.

### Test Running
You can also use the script for production

```bash
./run.sh
```

### Cross-Compile/Build

```bash
GOOS=linux GOARCH=amd64 go build -o ironite.exe ./src/ # for linux
GOOS=windows GOARCH=amd64 go build -o ironite.exe ./src/ # for windows
GOOS=darwin GOARCH=amd64 go build -o ironite.exe ./src/ # for macOS
```