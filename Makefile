1:
	cat words.txt | sed -E 's/[^[:alnum:]]/ /g' | tr -s ' ' '\n' | sort | uniq --count | sed -E 's/^\s+//g' | awk '{print $$2 " " $$1; }' | sort > results/1.txt

2:
	go run main.go single | sort > results/2.txt

diff: 1 2
	diff results/1.txt results/2.txt

gen10:
	for i in {1..10}; do cat orig.txt >> words.txt; done

gen50:
	for i in {1..50}; do cat orig.txt >> words.txt; done

gen100:
	for i in {1..100}; do cat orig.txt >> words.txt; done

3:
	go run main.go mrd | sort -n -r > results/3.txt
