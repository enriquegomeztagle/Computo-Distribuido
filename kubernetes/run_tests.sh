#!/bin/bash

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[1;34m'
NC='\033[0m'

total_tests=0
passed_tests=0
failed_tests=0

for dir in $(go list -f '{{.Dir}}' ./...); do
  if find "$dir" -type f -name '*_test.go' | grep -q .; then
    echo -e "\t${GREEN}\nRunning tests in ${YELLOW}$dir${NC}"
    
    if go test -race "$dir"; then
      echo -e "\t${GREEN}Tests passed in $dir${NC}"
      passed_tests=$((passed_tests + 1))
    else
      echo -e "\t${RED}Tests failed in $dir${NC}"
      failed_tests=$((failed_tests + 1))
    fi
    total_tests=$((total_tests + 1))
  else
    echo -e "\t${YELLOW}Skipping $dir (no test files)${NC}"
  fi
done

if [ $total_tests -gt 0 ]; then
  pass_percentage=$((100 * passed_tests / total_tests))
  fail_percentage=$((100 * failed_tests / total_tests))
else
  pass_percentage=0
  fail_percentage=0
fi

echo -e "\n---------------------------------"
echo -e "           ${BLUE}Test Summary${NC}           "
echo -e "---------------------------------"
echo -e "Total Packages Tested: ${YELLOW}$total_tests${NC}"
echo -e "Passed Packages: ${GREEN}$passed_tests${NC} (${GREEN}$pass_percentage%${NC})"
echo -e "Failed Packages: ${RED}$failed_tests${NC} (${RED}$fail_percentage%${NC})"
echo -e "---------------------------------"

if [ $pass_percentage -eq 100 ]; then
  echo -e "\n${GREEN}🎉 All tests passed! 😊${NC}"
elif [ $fail_percentage -eq 100 ]; then
  echo -e "\n${RED}💔 All tests failed 😢${NC}"
else
  echo -e "\n${YELLOW}⚠️  Some tests failed. Please review the results.${NC}"
fi

if [ $failed_tests -gt 0 ]; then
  exit 1
else
  exit 0
fi
