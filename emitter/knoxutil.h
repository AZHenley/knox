#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <stdbool.h>

// Concatensate two strings.
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

// Copy the contents of one string to a new string.
char* copy(const char *s)
{
    const size_t len = strlen(s);
    char *result = malloc(len + 1); 
    // TODO: Check for malloc errors.
    memcpy(result, s, len + 1);
    return result;
}

int random(int min, int max) {
    return (rand() % (max - min + 1)) + min; 
}

float randomf(float min, float max) {
    return min + ((float)rand() / RAND_MAX) * (max - min);
}

double randomd(double min, double max) {
    return min + ((double)rand() / RAND_MAX) * (max - min);
}