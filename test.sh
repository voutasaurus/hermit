set -ex
curl http://localhost:8080/cmd -d '{"cmd":"mkfile 1k somefile"}'
curl http://localhost:8080/cmd -d '{"cmd":"ls"}'
curl http://localhost:8080/cmd -d '{"cmd":"rm somefile"}'
curl http://localhost:8080/cmd -d '{"cmd":"ls"}'
