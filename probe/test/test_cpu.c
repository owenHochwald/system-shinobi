#include <assert.h>
#include <stdio.h>
#include <math.h>
#include "../include/cpu.h"

void test_cpu_delta_with_usage(void) {
    // prev: 100 user, 50 sys, 850 idle, 0 nice = 1000 total ticks
    // cur:  250 user, 100 sys, 1650 idle, 0 nice = 2000 total ticks
    // delta: 150 user, 50 sys, 800 idle = 1000 new ticks
    // active: 150 user + 50 sys = 200 active ticks
    // cpu_percent = (200 / 1000) * 100 = 20%
    CpuSample prev = {100, 50, 850, 0};
    CpuSample cur = {250, 100, 1650, 0};

    double result = cpu_delta(&prev, &cur);
    assert(fabs(result - 20.0) < 0.01);
    printf("✓ cpu_delta with 20%% usage\n");
}

void test_cpu_delta_high_usage(void) {
    // prev: 100, 50, 850, 0 = 1000
    // cur:  850, 100, 50, 0 = 1000
    // delta_user = 750, delta_sys = 50, delta_idle = -800, delta_total = 1000
    // delta_active = 750 + 50 = 800
    // cpu_percent = (800 / 1000) * 100 = 80%
    CpuSample prev2 = {100, 50, 850, 0};
    CpuSample cur2 = {850, 100, 50, 0};

    double result2 = cpu_delta(&prev2, &cur2);
    assert(fabs(result2 - 80.0) < 0.01);
    printf("✓ cpu_delta with 80%% usage\n");
}

void test_cpu_delta_zero_usage(void) {
    // Identical samples = no CPU usage
    CpuSample prev = {100, 50, 850, 0};
    CpuSample cur = {100, 50, 850, 0};

    double result = cpu_delta(&prev, &cur);
    assert(fabs(result - 0.0) < 0.01);
    printf("✓ cpu_delta with 0%% usage (identical samples)\n");
}

int main(void) {
    test_cpu_delta_with_usage();
    test_cpu_delta_high_usage();
    test_cpu_delta_zero_usage();
    printf("\n✓ All CPU delta tests passed!\n");
    return 0;
}
