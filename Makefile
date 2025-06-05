# Makefile

BPF_CLANG ?= clang
BPF_LLC   ?= llc
BPF_DIR   := bpf
OUT_DIR   := internal/tracer/bpf/compiled

CFLAGS := -O2 -g -target bpf -D__TARGET_ARCH_x86 -I$(BPF_DIR)

BPF_SRCS := $(wildcard $(BPF_DIR)/*.bpf.c)
BPF_OBJS := $(patsubst $(BPF_DIR)/%.bpf.c, $(OUT_DIR)/%.bpf.o, $(BPF_SRCS))

.PHONY: all clean bpf

all: bpf

bpf: $(BPF_OBJS)

$(OUT_DIR)/%.bpf.o: $(BPF_DIR)/%.bpf.c
	@mkdir -p $(OUT_DIR)
	$(BPF_CLANG) $(CFLAGS) -c $< -o $@

clean:
	rm -rf $(OUT_DIR)
