deploy: 
	VERSION=$(shell git describe --tags --abbrev=0 | awk -F. -v OFS=. '{$$NF = $$NF + 1;} 1') && \
	git commit -am "deploy" && \
	git tag $$VERSION && \
	git push origin $$VERSION && \
	go list -m github.com/jordanlumley/gobalena@$$VERSION && \
	git pull origin main && \
	git merge $$VERSION && \
	git push