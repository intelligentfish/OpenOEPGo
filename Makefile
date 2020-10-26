target=agent
$(target):	main.go
	go build
clean:
	-rm $(target)