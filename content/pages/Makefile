# requires package recutils on Ubuntu
all:
	recsel books.rec

five:
	recsel -e 'Rating=5' books.rec

finished:
	recsel -e 'Status="finished" && Date_Read >> "2020/01/01"' -p Title,Date_Read books.rec

