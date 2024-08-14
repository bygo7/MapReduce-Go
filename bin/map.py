import hashlib
import json
import os
import re
import sys
import uuid
from typing import List, Iterator, Tuple


def main():
	# Inputs:
	# 	input_file_paths: str[]
	# 	reduce_tasks: int
	# Outputs:
	# 	output_file_paths: str[]
	# 	output_file_bytes: int[]
    # Read in inputs
    input_file_paths = json.loads(sys.argv[1])
    reduce_tasks = int(sys.argv[2])
    all_lines = []

    # Read in input files
    for input_file_path in input_file_paths:
        with open(input_file_path, 'r') as f:
            lines = f.readlines()
            all_lines.extend(lines)

    # Write to R output files based on hash
    output_file_paths = []
    output_file_bytes = []
    rand_hash = str(uuid.uuid4())
    for i in range(reduce_tasks):
        output_file_path = "output_" + rand_hash + "_" + str(i) + '.txt'
        output_file_paths.append(output_file_path)
    for word, count in word_count(all_lines):
        # Hash word
        word_hash = consistent_hash(word, reduce_tasks)
        # Write word to file
        with open(output_file_paths[word_hash], 'a') as f:
            f.write("%s\t%d\n" % (word, count))
    for output_file_path in output_file_paths:
        output_file_bytes.append(os.path.getsize(output_file_path))
    
    # Print outputs to stdout in JSON format
    print(json.dumps({
        'output_file_paths': output_file_paths,
        'output_file_bytes': output_file_bytes
    }))

def word_count(file: List[str]) -> Iterator[Tuple[str, int]]:
    output = {}

    for line in file:
        # split the line into words
        words = re.findall(r'[a-z](?:[a-z\'‘’]*[a-z])?', line.lower())
        # increase counters
        for word in words:
            output[word] = output.get(word, 0) + 1

    # Sort alphabetically
    output = dict(sorted(output.items(), key=lambda item: item[0]))
    
    yield from output.items()

def consistent_hash(word: str, reduce_tasks: int) -> int:
    # Hash word
    hash = hashlib.sha256(word.encode('utf-8')).hexdigest()
    # Convert hash to int
    hash_int = int(hash, 16)
    # Mod by number of reduce tasks
    return hash_int % reduce_tasks

if __name__ == '__main__':
    main()

# python3 cmd/worker/mapper/map.py "John Bunyan___The Works of John Bunyan.txt" 2
