default:
	cd v1 && go test ./...
install:
	# go test ./...
	go generate ./...
	go install ./cmd/mdcoach
# *.gen.go files are hidden in .gitignore
default1: assets/assets.gen.go assets/img.gen.go
	cd ./cmd/mdcoach && go build -trimpath -o ~/.local/bin/mdcoach && chmod +x ~/.local/bin/mdcoach
	# mdcoach --output test/test-demo.html assets/test-demo.md assets/test-another.md
	# firefox test/test-demo.notes.html
assets/assets.gen.go:
	cd assets && zassets js sass --output assets.gen.go --refine --embed --var Cache --package assets
assets/img.gen.go:
	cd assets && zassets img --output img.gen.go --embed --var Img --package assets
frontend:
	npm run build
tests: clean compile-assets
	# go test -run YamlVisualizer
	mkdir -p /tmp/mdcoachtest/.cache/
	rsync -a assets/bundle/ /tmp/mdcoachtest/.cache/
	go test .
	# firefox /tmp/mdcoachtest/handout.html
	# devenv:
	# 	go get -u github.com/chai2010/webp
	# 	go get -u github.com/alecthomas/chroma
	# 	# sudo npm install --global webpack webpack-cli mini-css-extract-plugin
	# 	# cd assets && npm install --save-dev webpack webpack-cli
	# 	# cd assets && npm install style-loader css-loader sass-loader node-sass file-loader
clean:
	rm -f assets/assets.gen.go
	rm -f assets/img.gen.go
