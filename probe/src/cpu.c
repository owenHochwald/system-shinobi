#include "../include/cpu.h"
#include <mach/mach.h>
#include <mach/mach_host.h>
#include <mach/processor_info.h>

int cpu_sample(CpuSample *out) {
  mach_msg_type_number_t count;
  processor_cpu_load_info_t cpu_load;
  natural_t processor_count;

  kern_return_t kr = host_processor_info(
      mach_host_self(), PROCESSOR_CPU_LOAD_INFO, &processor_count,
      (processor_info_array_t *)&cpu_load, &count);

  if (kr != KERN_SUCCESS) {
    return -1;
  }

  // Aggregate across all cores
  out->user = 0;
  out->system = 0;
  out->idle = 0;
  out->nice = 0;

  for (natural_t i = 0; i < processor_count; i++) {
    out->user += cpu_load[i].cpu_ticks[CPU_STATE_USER];
    out->system += cpu_load[i].cpu_ticks[CPU_STATE_SYSTEM];
    out->idle += cpu_load[i].cpu_ticks[CPU_STATE_IDLE];
    out->nice += cpu_load[i].cpu_ticks[CPU_STATE_NICE];
  }

  vm_deallocate(mach_task_self(), (vm_address_t)cpu_load, count);
  return 0;
}

double cpu_delta(const CpuSample *prev, const CpuSample *cur) {
  uint64_t delta_user = cur->user - prev->user;
  uint64_t delta_system = cur->system - prev->system;
  uint64_t delta_idle = cur->idle - prev->idle;
  uint64_t delta_nice = cur->nice - prev->nice;

  uint64_t delta_total = delta_user + delta_system + delta_idle + delta_nice;
  if (delta_total == 0) {
    return 0.0;
  }

  uint64_t delta_active = delta_user + delta_system + delta_nice;
  return (100.0 * delta_active) / delta_total;
}
