#include "../include/pipe_writer.h"
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <sys/stat.h>
#include <time.h>
#include <unistd.h>

int pipe_open(const char *path) {
  // Remove old pipe if it exists
  unlink(path);

  // Create named pipe (FIFO)
  if (mkfifo(path, 0666) == -1) {
    return -1;
  }

  // Open for writing (blocks until reader connects)
  int fd = open(path, O_WRONLY);
  return fd;
}

int pipe_write_cpu(int fd, double cpu_percent) {
  char buffer[128];
  time_t timestamp = time(NULL);

  int len = snprintf(buffer, sizeof(buffer),
                     "{\"cpu_percent\":%.1f,\"timestamp\":%ld}\n", cpu_percent,
                     timestamp);

  if (len < 0 || len >= (int)sizeof(buffer)) {
    return -1;
  }

  ssize_t written = write(fd, buffer, strlen(buffer));
  return (written > 0) ? 0 : -1;
}

void pipe_close(int fd, const char *path) {
  if (fd >= 0) {
    close(fd);
  }
  unlink(path);
}
