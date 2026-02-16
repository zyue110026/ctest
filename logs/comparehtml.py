import re
from pathlib import Path

# Paths
file_a = Path("/Users/yuezhang/research/k8s-config-test/kubernetes-1.34.2/test/ctest/logs/ctest_unit_logs_20260209T183706.html")
file_b = Path("/Users/yuezhang/research/k8s-config-test/kubernetes-1.34.2/test/ctest/logs/ctest_unit_logs_20260216T153623.html")
output_file = file_b.parent / f"{file_b.stem}_compare.html"

# Regex to match any test line with duration: --- PASS/FAIL/SKIP: TestName (123.45s)
test_pattern = re.compile(r'^--- (PASS|FAIL|SKIP): (\S+) \([^\)]+\)', re.MULTILINE)

# Read file A and collect all test names
with open(file_a, "r") as f:
    content_a = f.read()

tests_a = test_pattern.findall(content_a)
# Only names
test_names_a = set(name for _, name in tests_a)

# Read file B
with open(file_b, "r") as f:
    content_b = f.read()

# Regex to match FAIL lines only with duration
fail_pattern = re.compile(r'^--- FAIL: (\S+) \([^\)]+\)', re.MULTILINE)

def mark_new_fail(match):
    test_name = match.group(1)
    if test_name not in test_names_a:
        return match.group(0) + " [NEW]"
    return match.group(0)

new_content = fail_pattern.sub(mark_new_fail, content_b)

# Save new file
with open(output_file, "w") as f:
    f.write(new_content)

print(f"Comparison done! Output saved to {output_file}")
