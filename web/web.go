package web

//go:generate npm run build
//go:generate togo http -package web --input dist/** -output web_gen.go --trim-prefix dist
