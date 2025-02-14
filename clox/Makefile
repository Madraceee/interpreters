clean:
	@rm clox*
build: 
	@gcc main.c chunk.c memory.c debug.c value.c vm.c compiler.c scanner.c table.c object.c -o clox
test:
	@gcc  test.c chunk.c memory.c debug.c value.c vm.c compiler.c scanner.c -o test -lcunit && ./test
