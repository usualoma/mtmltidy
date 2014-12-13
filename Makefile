update-data:
	rm -f data/tags.json
	$(MAKE) bindata.go

bindata.go: data/tags.json
	go-bindata data

data/tags.json:
	curl -s http://www.movabletype.jp/documentation/appendices/tags/ | hxselect 'li.xfolkentry' | perl -MJSON -ne 'BEGIN { %tags = (functions => [], blocks => []) }; END { print JSON->new->encode(\%tags) } push $$tags{$$2 . "s"}, $$1 if m{.*>MT(.*?)<.*tag=.(function|block)}' > $@
