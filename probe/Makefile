CC = clang
CFLAGS = -Wall -Wextra -Werror -std=c11 -Iinclude
LDFLAGS =

SRC_DIR = src
TEST_DIR = test
BUILD_DIR = build

SRCS = $(SRC_DIR)/main.c $(SRC_DIR)/cpu.c $(SRC_DIR)/pipe_writer.c
TEST_SRCS = $(TEST_DIR)/test_cpu.c $(SRC_DIR)/cpu.c
OBJS = $(SRCS:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o)
TEST_OBJS = $(TEST_SRCS:.c=.o)

TARGET = probe
TEST_TARGET = test_cpu

.PHONY: all clean test test_cpu test_pipe

all: $(TARGET)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(TARGET): $(BUILD_DIR) $(OBJS)
	$(CC) $(CFLAGS) $(OBJS) -o $(TARGET) $(LDFLAGS)

$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c | $(BUILD_DIR)
	$(CC) $(CFLAGS) -c $< -o $@

# Unit test for CPU delta calculation
test_cpu: $(TEST_DIR)/test_cpu.c $(SRC_DIR)/cpu.c
	$(CC) $(CFLAGS) $^ -o $(TEST_TARGET) -lm
	./$(TEST_TARGET)

# Integration test for named pipe
test_pipe: $(TARGET)
	cd $(TEST_DIR) && ./test_pipe.sh

test: test_cpu test_pipe

clean:
	rm -rf $(BUILD_DIR) $(TARGET) $(TEST_TARGET)
	rm -f /tmp/shinobi.pipe
