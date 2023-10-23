#!/bin/bash

url="http://127.0.0.1:8090/primes"
total_requests=100
total_time=0

for i in $(seq 1 $total_requests); do
    # Capture only the real time, discard output and only show the time taken
    time_taken=$( { time curl -so /dev/null $url; } 2>&1 | grep real | awk '{print $2}' | sed 's/s//')
    total_time=$(echo "$total_time + $time_taken" | bc)
done

avg_time=$(echo "scale=3; $total_time / $total_requests" | bc)
echo "Average latency for $url over $total_requests requests: $avg_time seconds"
