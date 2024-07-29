import json
import argparse
import os

# Set up argument parser
parser = argparse.ArgumentParser(description='Process frame metadata JSON file.')
parser.add_argument('json_file', type=str, help='Path to the frame_metadata.json file')

# Parse the arguments
args = parser.parse_args()

# Extract the directory path from the input file path
input_dir = os.path.dirname(args.json_file)

# Load the JSON array from the file
with open(args.json_file, 'r') as f:
    data = json.load(f)

# Replace each JSON object in the file, one per line
with open(args.json_file, 'w') as f:
    for item in data:
        f.write(json.dumps(item) + '\n')