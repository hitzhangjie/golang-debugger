english:
	unlink ./README.md
	unlink ./SUMMARY.md
	ln -s ./1-introduction.md ./README.md
	ln -s ./SUMMARY.en_US.md ./SUMMARY.md
	gitbook serve

chinese:
	unlink ./README.md
	unlink ./SUMMARY.md
	ln -s ./1-introduction.zh_CN.md ./README.md
	ln -s ./SUMMARY.zh_CN.md ./SUMMARY.md
	gitbook serve

.PHONY: english
.PHONY: chinese

clean:
	rm -rf ./_book
	#rm -rf ./node_modules
	unlink ./README.md
	unlink ./SUMMARY.md

install:
	gitbook install
