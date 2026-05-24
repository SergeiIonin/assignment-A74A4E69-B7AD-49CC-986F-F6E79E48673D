LOG_FILE="$1"

if [ -z "$LOG_FILE" ]; then
  echo "Usage: source calculateAvgResponse.sh <log_file>"
  return 1
fi

if [ ! -f "$LOG_FILE" ]; then
  echo "Log file not found: $LOG_FILE"
  return 1
fi

awk '
/took/ && !/\/dashboard\/0 / {
    val = $NF
    if (sub(/ms$/, "", val)) {
        ms = val + 0
    } else if (sub(/s$/, "", val)) {
        ms = (val + 0) * 1000
    }
    sum += ms
    count++
}
END {
    if (count > 0)
        printf "Average response time: %.3f ms (over %d requests)\n", sum / count, count
    else
        print "No response time entries found in log"
}
' "$LOG_FILE"
