#ifndef CPU_H
#define CPU_H

#include <stdint.h>

typedef struct {
    uint64_t user;
    uint64_t system;
    uint64_t idle;
    uint64_t nice;
} CpuSample;

// Sample current CPU ticks from Mach kernel
int cpu_sample(CpuSample *out);

// Calculate CPU usage percentage between two samples (0.0 - 100.0)
double cpu_delta(const CpuSample *prev, const CpuSample *cur);

#endif // CPU_H
