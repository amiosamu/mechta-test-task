#!/bin/bash

# Output file
output_file="numbers.json"

# Number of entries
num_entries=1000000

# Function to generate a random number between -10 and 10
random_number() {
  echo $((RANDOM % 21 - 10))
}

# Start JSON array
echo "[" > "$output_file"

# Generate JSON entries
for ((i=1; i<=num_entries; i++))
do
  a=$(random_number)
  b=$(random_number)
  
  if [[ $i -eq $num_entries ]]
  then
    # Last entry without trailing comma
    echo "  {\"a\": $a, \"b\": $b}" >> "$output_file"
  else
    # Entries with trailing comma
    echo "  {\"a\": $a, \"b\": $b}," >> "$output_file"
  fi
done

# End JSON array
echo "]" >> "$output_file"

echo "Generated $num_entries entries in $output_file"
