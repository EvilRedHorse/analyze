make all:
	go build

count = 1 

pkgs = \
  ./lockcheck \
  ./jsontag \
  ./responsewritercheck \

run = .

# fmt calls go fmt on all packges
fmt:
	go fmt $(pkgs)

# vet calls go vet on all packages
vet:
	go vet $(pkgs)

# test runs short tests 
test:
	go test -short -tags='debug testing netgo' -timeout=5s $(pkgs) -run=$(run) -count=$(count)

# test-long runs long tests 
test-long: fmt vet
	GORACE='$(racevars)' go test -race -v -failfast -tags='testing debug netgo' -timeout=3600s $(pkgs) -run=$(run) -count=$(count)


