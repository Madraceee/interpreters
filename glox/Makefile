clean:
	@rm -rf ./bin
build: clean
	@go build -o ./bin/glox
run:
	@./bin/glox
generateAst:
	@go build  -o ./parser/tool/ ./parser/tool/ && ./parser/tool/tool && go fmt  ./parser/ast.go
