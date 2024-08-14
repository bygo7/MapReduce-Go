import json
import os
import sys
import uuid
from itertools import groupby
from operator import itemgetter
from typing import Iterator, List


def main(separator='\t'):
    # Inputs:
	# 	input_file_paths: str[]
	# Outputs:
	# 	output_file_path: str
	# 	output_file_bytes: int

    # Read in inputs
    input_file_paths = json.loads(sys.argv[1])
    all_lines = []

    # Read in input files
    for input_file_path in input_file_paths:
        with open(input_file_path, 'r') as f:
            lines = f.readlines()
            all_lines.extend(lines)
    
    # Count words
    word_counts = parse_word_count(all_lines, separator=separator)

    # Write to output file in JSON
    output_file_path = "output_" + str(uuid.uuid4()) + '.txt'
    with open(output_file_path, 'w') as f:
        for current_word, group in groupby(word_counts, itemgetter(0)):
            try:
                total_count = sum(int(count) for _, count in group)
                f.write("%s%s%d\n" % (current_word, separator, total_count))
            except ValueError:
                # count was not a number, so silently discard this item
                pass
    output_file_bytes = os.path.getsize(output_file_path)
    
    # Print outputs to stdout in JSON format
    print(json.dumps({
        'output_file_path': output_file_path,
        'output_file_bytes': output_file_bytes
    }))

def parse_word_count(lines: List[str], separator='\t') -> Iterator[List[str]]:
    # Sort alphabetically by word (format: word\tcount)
    lines.sort()
    # Group by word
    for line in lines:
        yield line.rstrip().split(separator, 1)

if __name__ == '__main__':
    main()

# python3 cmd/worker/reducer/reduce.py '["output_cf05aeee-06e5-4544-845b-9cbc111ff259_0.txt","output_cf05aeee-06e5-4544-845b-9cbc111ff259_1.txt"]'
