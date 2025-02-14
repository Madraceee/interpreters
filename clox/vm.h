#ifndef clox_vm_h
#define clox_vm_h

#include "chunk.h"
#include "table.h"
#include "value.h"
#include <stdint.h>

#define STACK_MAX 256

// VM
// chunk contains all the bytecode and other data
// ip is a pointer which points to the current chunk/instruction
// stack used to store intermediate values for calculation
typedef struct{
	Chunk* chunk;
	uint8_t* ip;
	Value stack[STACK_MAX];
	Value* stack_top;
	Table strings;
	Obj* objects;
} VM;

typedef enum {
	INTERPRET_OK,
	INTERPRET_COMPILE_ERROR,
	INTERPRET_RUNTIME_ERROR,
}InterpretResult;

extern VM vm;

void initVM();
void freeVM();
InterpretResult interpret(const char* source);
void push(Value value);
Value pop();

#endif
