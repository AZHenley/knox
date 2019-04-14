#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <stdbool.h>

char* concat(const char *s1, const char *s2)
{
    const size_t len1 = strlen(s1);
    const size_t len2 = strlen(s2);
    char *result = malloc(len1 + len2 + 1); 
    // TODO: Check for malloc errors.
    memcpy(result, s1, len1);
    memcpy(result + len1, s2, len2 + 1);
    return result;
}

char* copy(const char *s)
{
    const size_t len = strlen(s);
    char *result = malloc(len + 1); 
    // TODO: Check for malloc errors.
    memcpy(result, s, len + 1);
    return result;
}