#ifndef MARKABLE_REFERENCE_H
#define MARKABLE_REFERENCE_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef void*      markable_reference;
markable_reference markable_reference_init(uintptr_t reference, bool mark);
void               markable_reference_free(markable_reference instance);
bool               markable_reference_compare_and_set(markable_reference instance, uintptr_t expected_reference, uintptr_t new_reference,
                                                      bool expected_mark, bool new_mark);
uintptr_t          markable_reference_get_reference(markable_reference instance);
bool               markable_reference_is_marked();

#ifdef __cplusplus
}
#endif

#endif