#!/bin/bash

echo "Benchmarking PyTorch LLM inference..."

START=$(date +%s)

python3 -c "
import torch
from transformers import AutoModelForCausalLM, AutoTokenizer
model = AutoModelForCausalLM.from_pretrained('gpt2')
tokenizer = AutoTokenizer.from_pretrained('gpt2')
inputs = tokenizer('Hello, world!', return_tensors='pt')
outputs = model.generate(**inputs)
print(tokenizer.decode(outputs[0]))
"

END=$(date +%s)
echo "Inference completed in $(($END - $START)) seconds"