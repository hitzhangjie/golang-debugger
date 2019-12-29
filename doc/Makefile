english:
	rm -rf ./_book
	ln -sf ./1-introduction.md ./README.md
	ln -sf ./SUMMARY.en_US.md ./SUMMARY.md
	gitbook serve

chinese:
	rm -rf ./_book
	ln -sf ./1-introduction.zh_CN.md ./README.md
	ln -sf ./SUMMARY.zh_CN.md ./SUMMARY.md
	gitbook serve

.PHONY: english
.PHONY: chinese

clean:
	rm -rf _book
	#rm -rf ./node_modules
	rm -f ./README.md
	rm -f ./SUMMARY.md

install:
	gitbook install
