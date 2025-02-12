#include "CUnit/CUnit.h"
#include "CUnit/Basic.h"

#include <stdlib.h>

#include "chunk.h"


int init_suite_chunk(void){
	return 0;
}

int clean_suite_chunk(void){
	/*free(chunk);*/
	return 0;
}

// Test for chunks
void testInitChunk(Chunk* chunk){
	initChunk(chunk);
	CU_ASSERT_EQUAL( chunk->capacity,0);
	CU_ASSERT_EQUAL( chunk->count,0);
	CU_ASSERT_EQUAL( chunk->code,NULL);
	CU_ASSERT_EQUAL( chunk->lines,NULL);
}

void testWriteChunk(Chunk* chunk, uint8_t byte, int line){
	writeChunk(chunk, byte, line);
	CU_ASSERT_EQUAL( chunk->capacity,8);
	CU_ASSERT_EQUAL( chunk->count,1);
	CU_ASSERT_EQUAL( chunk->code[chunk->count-1],byte);
	CU_ASSERT_EQUAL( chunk->lines[chunk->count-1],line);
}

void testChunk(void){
	// Test for chunk
	Chunk* chunk;
	testInitChunk(chunk);
	testWriteChunk(chunk, OP_RETURN, 1);
}

int main() {
	if (CUE_SUCCESS != CU_initialize_registry()){
		return CU_get_error();
	}
	CU_pSuite pSuite = CU_add_suite("Chunk test", init_suite_chunk, clean_suite_chunk);

	if (pSuite == NULL){
		CU_cleanup_registry();
		return CU_get_error();
	}
	
	if (NULL == CU_add_test(pSuite,"Chunk Tests", testChunk)){
		CU_cleanup_registry();
		return CU_get_error();
	}

	// Run
	CU_basic_set_mode(CU_BRM_VERBOSE);
	CU_basic_run_tests();
	CU_cleanup_registry();
	return 0;
}
