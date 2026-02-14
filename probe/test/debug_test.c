#include <stdio.h>
#include <stdint.h>
#include "../include/cpu.h"

int main(void) {
    // This is the failing test case
    CpuSample prev = {100, 50, 850, 0};
    CpuSample cur = {850, 100, 50, 0};

    printf("prev: user=%llu, sys=%llu, idle=%llu, nice=%llu\n",
           prev.user, prev.system, prev.idle, prev.nice);
    printf("cur:  user=%llu, sys=%llu, idle=%llu, nice=%llu\n",
           cur.user, cur.system, cur.idle, cur.nice);

    uint64_t delta_user = cur.user - prev.user;
    uint64_t delta_system = cur.system - prev.system;
    uint64_t delta_idle = cur.idle - prev.idle;  // This will underflow!
    uint64_t delta_nice = cur.nice - prev.nice;

    printf("\nDeltas (as uint64_t):\n");
    printf("delta_user=%llu, delta_sys=%llu, delta_idle=%llu, delta_nice=%llu\n",
           delta_user, delta_system, delta_idle, delta_nice);

    uint64_t delta_total = delta_user + delta_system + delta_idle + delta_nice;
    printf("delta_total=%llu\n", delta_total);

    double result = cpu_delta(&prev, &cur);
    printf("\nCPU result: %.2f%%\n", result);
    printf("Expected: 80.00%%\n");

    return 0;
}
