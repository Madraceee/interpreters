#ifndef clox_value_h
#define clox_value_h

#include "common.h"

typedef double Value;

// ValueArray Used to store literal values
// Since constant cannot be directly stored in OPCODE
// We use a constants Pool to store the values
// These values are refered from the chunk as and when required.
typedef struct {
	int capacity;
	int count;
	Value* values;
} ValueArray;

void initValueArray(ValueArray* array);
void writeValueArray(ValueArray* array, Value value);
void freeValueArray(ValueArray* array);

#endif
