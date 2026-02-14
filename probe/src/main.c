#include "../include/cpu.h"
#include "../include/pipe_writer.h"
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#define PIPE_PATH "/tmp/shinobi.pipe"
#define SAMPLE_INTERVAL 1

static volatile int running = 1;
static int pipe_fd = -1;

void signal_handler(int sig) {
  (void)sig;
  running = 0;
}

int main(void) {
  signal(SIGINT, signal_handler);
  signal(SIGTERM, signal_handler);

  pipe_fd = pipe_open(PIPE_PATH);
  if (pipe_fd < 0) {
    fprintf(stderr, "Failed to open pipe at %s\n", PIPE_PATH);
    return 1;
  }

  CpuSample prev, cur;
  if (cpu_sample(&prev) != 0) {
    fprintf(stderr, "Failed to get initial CPU sample\n");
    pipe_close(pipe_fd, PIPE_PATH);
    return 1;
  }

  while (running) {
    sleep(SAMPLE_INTERVAL);

    if (cpu_sample(&cur) != 0) {
      fprintf(stderr, "Failed to sample CPU\n");
      continue;
    }

    double cpu_percent = cpu_delta(&prev, &cur);

    if (pipe_write_cpu(pipe_fd, cpu_percent) != 0) {
      fprintf(stderr, "Failed to write to pipe\n");
    }

    prev = cur;
  }

  pipe_close(pipe_fd, PIPE_PATH);
  return 0;
}
