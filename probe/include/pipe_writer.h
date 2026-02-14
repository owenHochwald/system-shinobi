#ifndef PIPE_WRITER_H
#define PIPE_WRITER_H

// Open or create a named pipe for writing
int pipe_open(const char *path);

// Write CPU percentage as JSON to the pipe
int pipe_write_cpu(int fd, double cpu_percent);

// Close pipe and clean up
void pipe_close(int fd, const char *path);

#endif // PIPE_WRITER_H
