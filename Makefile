APP_NAME := bnnchatbot
GO_FILES := $(wildcard *.go)

.PHONY: all clean

all: $(APP_NAME).exe

$(APP_NAME).exe: $(GO_FILES)
	go build -ldflags -H=windowsgui -o ./target/$(APP_NAME).exe .

clean:
	rm -f $(APP_NAME).exe