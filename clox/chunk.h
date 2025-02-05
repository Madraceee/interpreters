#ifndef clox_chunk_h
#define clox_chunk_h

#include "value.h"
#include <stdint.h>

typedef enum{
	OP_CONSTANT,
	OP_RETURN
}OpCode;

typedef struct {
	int capacity;
	int count;
	uint8_t* code;
	int* lines;
	ValueArray constants;
} Chunk;

void initChunk(Chunk* chunk);
void freeChunk(Chunk* chunk);
void writeChunk(Chunk* chunk, uint8_t byte, int line);
int addConstants(Chunk* chunk, Value value);

#endif
