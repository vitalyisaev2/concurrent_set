#include "markable_reference.h"
#include "markable_reference.hpp"

extern "C" {
markable_reference markable_reference_init(uintptr_t reference, bool mark)
{
    auto result = new MarkableReference<uintptr_t>(reference, mark);
}

void      markable_reference_free(markable_reference instance);
bool      markable_reference_compare_and_set(markable_reference instance, uintptr_t expected_reference, uintptr_t new_reference,
                                             bool expected_mark, bool new_mark);
uintptr_t markable_reference_get_reference(markable_reference instance);
bool      markable_reference_is_marked();
}