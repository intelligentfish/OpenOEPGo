target=openOEP
$(target):	main.go
	go build
clean:
	-rm $(target)
