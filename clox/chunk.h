#ifndef clox_chunk_h
#define clox_chunk_h

#include "value.h"
#include <stdint.h>

typedef enum{
	OP_CONSTANT,
	OP_NIL,
	OP_TRUE,
	OP_FALSE,
	OP_RETURN,
	OP_NEGATE,
	OP_ADD,
	OP_SUBTRACT,
	OP_MULTIPLY,
	OP_DIVIDE,
	OP_NOT,
}OpCode;

// Chunk Used to store instrcutions as bytecode
typedef struct {
	int capacity;
	int count;
	uint8_t* code;
	int* lines;
	ValueArray constants;
} Chunk;

// Functions to increase or decrease the size of chunks
void initChunk(Chunk* chunk);
void freeChunk(Chunk* chunk);
void writeChunk(Chunk* chunk, uint8_t byte, int line);
int addConstant(Chunk* chunk, Value value);

#endif
